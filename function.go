package extractor

import (
	"fmt"
	"go/ast"
	"os"
	"strings"
)

type GoFunctionMeta struct {
	*meta // *ast.FuncDecl
	// callMeta map[string][]*GoCallMeta
	// nonSelectorCallMeta map[string][]*GoCallMeta
	// selectorCallMeta    map[string]map[string][]*GoCallMeta
}

func NewGoFunctionMeta(m *meta) *GoFunctionMeta {
	return &GoFunctionMeta{
		meta: m,
		// callMeta: make(map[string][]*GoCallMeta),
	}
}

func ExtractGoFunctionMeta(extractFilepath string, functionName string) (*GoFunctionMeta, error) {
	goFileMeta, err := ExtractGoFileMeta(extractFilepath)
	if err != nil {
		return nil, err
	}

	gfm := SearchGoFunctionMeta(goFileMeta, functionName)
	if gfm == nil {
		return nil, fmt.Errorf("can not find function node")
	}

	return gfm, nil
}

func SearchGoFunctionMeta(gfm *GoFileMeta, functionName string) *GoFunctionMeta {
	var funcDecl *ast.FuncDecl
	ast.Inspect(gfm.node, func(n ast.Node) bool {
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
		meta: gfm.newMeta(funcDecl),
		// callMeta: make(map[string][]*GoCallMeta),
		// nonSelectorCallMeta: make(map[string][]*GoCallMeta),
		// selectorCallMeta:    make(map[string]map[string][]*GoCallMeta),
	}
}

func IsFuncNode(n ast.Node) bool {
	funcDecl, ok := n.(*ast.FuncDecl)
	return ok && funcDecl.Recv == nil
}

func (gfm *GoFunctionMeta) FunctionName() string {
	return gfm.node.(*ast.FuncDecl).Name.String()
}

func (gfm *GoFunctionMeta) Doc() []string {
	if gfm.node.(*ast.FuncDecl) == nil || gfm.node.(*ast.FuncDecl).Doc == nil || len(gfm.node.(*ast.FuncDecl).Doc.List) == 0 {
		return nil
	}
	commentSlice := make([]string, 0, len(gfm.node.(*ast.FuncDecl).Doc.List))
	for _, comment := range gfm.node.(*ast.FuncDecl).Doc.List {
		commentSlice = append(commentSlice, comment.Text)
	}
	return commentSlice
}

func (gfm *GoFunctionMeta) TypeParams() []*GoVariableMeta {
	if gfm.node.(*ast.FuncDecl).Type == nil || gfm.node.(*ast.FuncDecl).Type.TypeParams == nil || len(gfm.node.(*ast.FuncDecl).Type.TypeParams.List) == 0 {
		return nil
	}

	tParamLen := len(gfm.node.(*ast.FuncDecl).Type.TypeParams.List)
	tParams := make([]*GoVariableMeta, 0, tParamLen)
	for _, field := range gfm.node.(*ast.FuncDecl).Type.TypeParams.List {
		for _, name := range field.Names {
			tParams = append(tParams, &GoVariableMeta{
				meta:     gfm.newMeta(field),
				name:     name.String(),
				typeMeta: gfm.newMeta(field.Type),
			})
		}
	}
	return tParams
}

func (gfm *GoFunctionMeta) Params() []*GoVariableMeta {
	if gfm.node.(*ast.FuncDecl).Type == nil || gfm.node.(*ast.FuncDecl).Type.Params == nil || len(gfm.node.(*ast.FuncDecl).Type.Params.List) == 0 {
		return nil
	}

	pLen := len(gfm.node.(*ast.FuncDecl).Type.Params.List)
	params := make([]*GoVariableMeta, 0, pLen)
	for _, field := range gfm.node.(*ast.FuncDecl).Type.Params.List {
		for _, name := range field.Names {
			params = append(params, &GoVariableMeta{
				meta:     gfm.newMeta(field),
				name:     name.String(),
				typeMeta: gfm.newMeta(field.Type),
			})
		}
	}
	return params
}

func (gfm *GoFunctionMeta) ReturnTypes() []*GoVariableMeta {
	if gfm.node.(*ast.FuncDecl).Type == nil || gfm.node.(*ast.FuncDecl).Type.Results == nil || len(gfm.node.(*ast.FuncDecl).Type.Results.List) == 0 {
		return nil
	}

	rLen := len(gfm.node.(*ast.FuncDecl).Type.Results.List)
	returns := make([]*GoVariableMeta, 0, rLen)
	for _, field := range gfm.node.(*ast.FuncDecl).Type.Results.List {
		// TODO: not support named return value
		returns = append(returns, &GoVariableMeta{
			meta:     gfm.newMeta(field),
			name:     "",
			typeMeta: gfm.newMeta(field.Type),
		})
	}
	return returns
}

func (gfm *GoFunctionMeta) Expression() string {
	originPos := gfm.node.(*ast.FuncDecl).Pos()
	originEnd := gfm.node.(*ast.FuncDecl).End()
	if gfm.node.(*ast.FuncDecl).Doc != nil {
		originPos = gfm.node.(*ast.FuncDecl).Doc.Pos()
	}

	fileContent, err := os.ReadFile(gfm.path)
	if err != nil {
		return ""
	}
	fileContentLen := len(fileContent)
	if originPos > originEnd || int(originPos) >= fileContentLen || int(originEnd) >= fileContentLen {
		return ""
	}
	return strings.TrimSpace(string(fileContent[originPos-1 : originEnd-1]))
}

// func (gfm *GoFunctionMeta) ReplaceFunctionDoc(newDoc []string) (string, string, error) {
// 	originContent := gfm.Expression()

// 	gfm.node.(*ast.FuncDecl).Doc = &ast.CommentGroup{
// 		List: make([]*ast.Comment, 0, len(newDoc)),
// 	}
// 	for _, comment := range newDoc {
// 		gfm.node.(*ast.FuncDecl).Doc.List = append(gfm.node.(*ast.FuncDecl).Doc.List, &ast.Comment{
// 			Text: comment,
// 		})
// 	}

// 	buffer := &bytes.Buffer{}
// 	err := format.Node(buffer, gfm.fileMeta.fileSet, gfm.node.(*ast.FuncDecl))
// 	if err != nil {
// 		panic(err)
// 	}

// 	return strings.ReplaceAll(originContent, "\r", ""), strings.ReplaceAll(buffer.String(), "\r", ""), nil
// }

// func (gfm *GoFunctionMeta) SearchCallMeta(call string) []*GoCallMeta {
// 	if gcm, has := gfm.callMeta[call]; gcm != nil && has {
// 		return gcm
// 	}

// 	if gfm.node.(*ast.FuncDecl) == nil {
// 		return nil
// 	}

// 	gfm.callMeta[call] = SearchGoCallMeta(gfm.meta, call)

// 	return gfm.callMeta[call]
// }

// func (gfm *GoFunctionMeta) Calls() map[string][]*GoCallMeta {
// 	return ExtractGoCallMeta(gfm.meta)
// }

// func (gfm *GoFunctionMeta) Search
