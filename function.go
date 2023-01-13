package extractor

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

type GoFunctionMeta struct {
	fileMeta *GoFileMeta
	funcDecl *ast.FuncDecl
}

func ExtractGoFunctionMeta(extractFilepath string, functionName string) (*GoFunctionMeta, error) {
	fileSet := token.NewFileSet()
	fileAST, err := parser.ParseFile(fileSet, extractFilepath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	gfm := SearchGoFunctionMeta(&GoFileMeta{fileSet: fileSet, fileAST: fileAST}, functionName)
	if gfm.funcDecl == nil {
		return nil, fmt.Errorf("can not find function decl")
	}

	return gfm, nil
}

func SearchGoFunctionMeta(fileMeta *GoFileMeta, functionName string) *GoFunctionMeta {
	var funcDecl *ast.FuncDecl
	ast.Inspect(fileMeta.fileAST, func(n ast.Node) bool {
		if IsFuncNode(n) {
			decl := n.(*ast.FuncDecl)
			if decl.Name.String() == functionName {
				funcDecl = decl
			}
		}
		return funcDecl == nil
	})
	if funcDecl == nil {
		return nil
	}
	return &GoFunctionMeta{
		fileMeta: fileMeta,
		funcDecl: funcDecl,
	}
}

func IsFuncNode(n ast.Node) bool {
	funcDecl, ok := n.(*ast.FuncDecl)
	return ok && funcDecl.Recv == nil
}

func (gfm *GoFunctionMeta) PrintAST() {
	ast.Print(token.NewFileSet(), gfm.funcDecl)
}

func (gfm *GoFunctionMeta) FunctionName() string {
	return gfm.funcDecl.Name.String()
}

func (gfm *GoFunctionMeta) CallMap() map[string][]*ast.CallExpr {
	// ast.Print(token.NewFileSet(), gfm.funcDecl.Body)
	callMap := make(map[string][]*ast.CallExpr)
	for _, e := range gfm.funcDecl.Body.List {
		exprStmt, ok := e.(*ast.ExprStmt)
		if exprStmt == nil || !ok || exprStmt.X == nil {
			continue
		}
		callExpr, ok := exprStmt.X.(*ast.CallExpr)
		if callExpr == nil || !ok {
			continue
		}

		ident, ok := callExpr.Fun.(*ast.Ident)
		if ident == nil || !ok {
			continue
		}
		callMap[ident.Name] = append(callMap[ident.Name], callExpr)
	}
	return callMap
}

func (gfm *GoFunctionMeta) Doc() []string {
	if gfm.funcDecl == nil || gfm.funcDecl.Doc == nil || len(gfm.funcDecl.Doc.List) == 0 {
		return nil
	}
	commentSlice := make([]string, 0, len(gfm.funcDecl.Doc.List))
	for _, comment := range gfm.funcDecl.Doc.List {
		commentSlice = append(commentSlice, comment.Text)
	}
	return commentSlice
}

func (gfm *GoFunctionMeta) ReturnTypes() []string {
	if gfm.funcDecl.Type == nil || gfm.funcDecl.Type.Results == nil || len(gfm.funcDecl.Type.Results.List) == 0 {
		return nil
	}

	rLen := len(gfm.funcDecl.Type.Results.List)
	returnTypes := make([]string, 0, rLen)
	for _, field := range gfm.funcDecl.Type.Results.List {
		ident, ok := field.Type.(*ast.Ident)
		if ident == nil || !ok {
			continue
		}
		returnTypes = append(returnTypes, ident.Name)
	}
	return returnTypes
}

func (gfm *GoFunctionMeta) ReplaceFunctionDoc(newDoc []string) (string, string, error) {
	fileContent, err := os.ReadFile(gfm.fileMeta.Path)
	if err != nil {
		return "", "", err
	}

	originPos := gfm.funcDecl.Pos()
	originEnd := gfm.funcDecl.End()
	if gfm.funcDecl.Doc != nil {
		originPos = gfm.funcDecl.Doc.Pos()
	}
	originContent := string(fileContent[originPos-1 : originEnd])

	gfm.funcDecl.Doc = &ast.CommentGroup{
		List: make([]*ast.Comment, 0, len(newDoc)),
	}
	for _, comment := range newDoc {
		gfm.funcDecl.Doc.List = append(gfm.funcDecl.Doc.List, &ast.Comment{
			Text: comment,
		})
	}

	buffer := &bytes.Buffer{}
	err = format.Node(buffer, gfm.fileMeta.fileSet, gfm.funcDecl)
	if err != nil {
		panic(err)
	}

	return strings.ReplaceAll(originContent, "\r", ""), strings.ReplaceAll(buffer.String(), "\r", ""), nil
}
