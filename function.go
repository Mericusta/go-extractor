package extractor

import (
	"fmt"
	"go/ast"
	"os"
	"strings"
)

type GoFunctionMeta struct {
	// funcDecl            *ast.FuncDecl
	*meta
	callMeta map[string][]*GoCallMeta
	// nonSelectorCallMeta map[string][]*GoCallMeta
	// selectorCallMeta    map[string]map[string][]*GoCallMeta
}

func NewGoFunctionMeta(m *meta) *GoFunctionMeta {
	return &GoFunctionMeta{
		meta:     m,
		callMeta: make(map[string][]*GoCallMeta),
	}
}

func ExtractGoFunctionMeta(extractFilepath string, functionName string) (*GoFunctionMeta, error) {
	goFileMeta, err := ExtractGoFileMeta(extractFilepath)
	if err != nil {
		return nil, err
	}

	gfm := SearchGoFunctionMeta(goFileMeta.meta, functionName)
	if gfm.node == nil {
		return nil, fmt.Errorf("can not find function node")
	}

	return gfm, nil
}

func SearchGoFunctionMeta(m *meta, functionName string) *GoFunctionMeta {
	var funcDecl *ast.FuncDecl
	ast.Inspect(m.node, func(n ast.Node) bool {
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
		meta:     m.newMeta(funcDecl),
		callMeta: make(map[string][]*GoCallMeta),
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

func (gfm *GoFunctionMeta) ReturnTypes() []string {
	if gfm.node.(*ast.FuncDecl).Type == nil || gfm.node.(*ast.FuncDecl).Type.Results == nil || len(gfm.node.(*ast.FuncDecl).Type.Results.List) == 0 {
		return nil
	}

	rLen := len(gfm.node.(*ast.FuncDecl).Type.Results.List)
	returnTypes := make([]string, 0, rLen)
	for _, field := range gfm.node.(*ast.FuncDecl).Type.Results.List {
		ident, ok := field.Type.(*ast.Ident)
		if ident == nil || !ok {
			continue
		}
		returnTypes = append(returnTypes, ident.String())
	}
	return returnTypes
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

func (gfm *GoFunctionMeta) SearchCallMeta(call string) []*GoCallMeta {
	if gcm, has := gfm.callMeta[call]; gcm != nil && has {
		return gcm
	}

	if gfm.node.(*ast.FuncDecl) == nil {
		return nil
	}

	gfm.callMeta[call] = SearchGoCallMeta(gfm.meta, call)

	return gfm.callMeta[call]

	// if len(from) == 0 {
	// 	if callMetaSlice, has := gfm.nonSelectorCallMeta[call]; has && len(callMetaSlice) > 0 {
	// 		return callMetaSlice
	// 	}
	// } else {
	// 	if selector, has := gfm.selectorCallMeta[from]; has && len(selector) > 0 {
	// 		if callMetaSlice, has := selector[call]; has && len(callMetaSlice) > 0 {
	// 			return callMetaSlice
	// 		}
	// 	}
	// }

	// if gfm.node.(*ast.FuncDecl) == nil {
	// 	return nil
	// }

	// gcm := SearchGoCallMeta(gfm.meta, gfm.node.(*ast.FuncDecl), call)
	// if gcm != nil {
	// 	if len(from) == 0 {
	// 		gfm.nonSelectorCallMeta[call] = append(gfm.nonSelectorCallMeta[call], gcm)
	// 		return gfm.nonSelectorCallMeta[call]
	// 	} else {
	// 		if gfm.selectorCallMeta[from] == nil {
	// 			gfm.selectorCallMeta[from] = make(map[string][]*GoCallMeta)
	// 		}
	// 		gfm.selectorCallMeta[from][call] = append(gfm.selectorCallMeta[from][call], gcm)
	// 		return gfm.selectorCallMeta[from][call]
	// 	}
	// }

	return nil
}
