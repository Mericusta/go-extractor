package extractor

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

type GoFunctionMeta struct {
	// fileMeta *GoFileMeta
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
		if n == fileMeta.fileAST {
			return true
		}
		if n == nil || funcDecl != nil {
			return false
		}
		decl, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}
		if decl.Recv == nil && decl.Name.String() == functionName {
			funcDecl = decl
			return false
		}
		return true
	})
	if funcDecl == nil {
		return nil
	}
	return &GoFunctionMeta{
		// fileMeta: fileMeta,
		funcDecl: funcDecl,
	}
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

func (gfm *GoFunctionMeta) Comments() []string {
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

// // TODO:
// func (gfm *GoFunctionMeta) UpdateComments(comments []string) {
// 	if gfm.funcDecl.Doc == nil {
// 		return
// 	}

// 	if len(gfm.funcDecl.Doc.List) != len(comments) {
// 		return
// 	}

// 	for index, comment := range gfm.funcDecl.Doc.List {
// 		comment.Text = comments[index]
// 	}

// 	outputFile, err := os.OpenFile(gfm.fileMeta.Path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer outputFile.Close()

// 	buffer := &bytes.Buffer{}
// 	if gfm.fileMeta.fileAST.Scope != nil {
// 		for name, object := range gfm.fileMeta.fileAST.Scope.Objects {
// 			if object.Kind == ast.Fun && name == gfm.FunctionName() {
// 				decl := object.Decl.(*ast.FuncDecl)
// 				if declLen := decl.End() - decl.Pos(); buffer.Cap() < int(declLen) {
// 					buffer.Grow(int(declLen))
// 				}
// 				err = format.Node(buffer, gfm.fileMeta.fileSet, decl)
// 				if err != nil {
// 					panic(err)
// 				}
// 				outputFile.WriteAt(buffer.Bytes(), int64(decl.Pos()))
// 				// outputFile.Write(buffer.Bytes())
// 				buffer.Reset()
// 				break
// 			}
// 		}
// 	}
// }
