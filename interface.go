package extractor

import (
	"fmt"
	"go/ast"
)

type GoInterfaceMeta struct {
	// typeSpec     *ast.TypeSpec
	*meta
	commentGroup *ast.CommentGroup
	methodMeta   map[string]*GoInterfaceMethodMeta
}

func ExtractGoInterfaceMeta(extractFilepath, interfaceName string) (*GoInterfaceMeta, error) {
	gfm, err := ExtractGoFileMeta(extractFilepath)
	if err != nil {
		return nil, err
	}

	gim := SearchGoInterfaceMeta(gfm, interfaceName)
	if gim == nil {
		return nil, fmt.Errorf("can not find interface node")
	}

	return gim, nil
}

func SearchGoInterfaceMeta(gfm *GoFileMeta, interfaceName string) *GoInterfaceMeta {
	var interfaceDecl *ast.TypeSpec
	var commentDecl *ast.CommentGroup
	ast.Inspect(gfm.node, func(n ast.Node) bool {
		if genDecl, ok := n.(*ast.GenDecl); ok {
			ast.Inspect(genDecl, func(n ast.Node) bool {
				if IsInterfaceNode(n) {
					typeSpec := n.(*ast.TypeSpec)
					if typeSpec.Name.String() == interfaceName {
						interfaceDecl = typeSpec
						commentDecl = genDecl.Doc
						return false
					}
				}
				return true
			})
			return false // genDecl traverse done
		}
		return interfaceDecl == nil // already found
	})
	if interfaceDecl == nil {
		return nil
	}
	return &GoInterfaceMeta{
		meta:         gfm.newMeta(interfaceDecl),
		commentGroup: commentDecl,
		methodMeta:   make(map[string]*GoInterfaceMethodMeta),
	}
}

func IsInterfaceNode(n ast.Node) bool {
	typeSpec, ok := n.(*ast.TypeSpec)
	if !ok {
		return false
	}
	if typeSpec.Type == nil {
		return false
	}
	_, ok = typeSpec.Type.(*ast.InterfaceType)
	return ok
}

func (gim *GoInterfaceMeta) InterfaceName() string {
	return gim.node.(*ast.TypeSpec).Name.String()
}

func (gim *GoInterfaceMeta) Doc() []string {
	if gim.node == nil || gim.commentGroup == nil || len(gim.commentGroup.List) == 0 {
		return nil
	}
	commentSlice := make([]string, 0, len(gim.commentGroup.List))
	for _, comment := range gim.commentGroup.List {
		commentSlice = append(commentSlice, comment.Text)
	}
	return commentSlice
}

// SearchMethodDecl search method decl from node.(*ast.InterfaceType)
func (gim *GoInterfaceMeta) SearchMethodDecl(methodName string) *GoInterfaceMethodMeta {
	gim.ForeachMethodDecl(func(f *ast.Field) bool {
		if f.Names[0].String() == methodName && IsInterfaceMethodNode(f) {
			gim.methodMeta[methodName] = NewGoInterfaceMethodMeta(
				gim.newMeta(f), gim,
			)
			return false
		}
		return true
	})
	return gim.methodMeta[methodName]
}

func (gim *GoInterfaceMeta) ForeachMethodDecl(f func(*ast.Field) bool) {
	interfaceType := gim.node.(*ast.TypeSpec).Type.(*ast.InterfaceType)
	if interfaceType.Methods == nil {
		return
	}
	for _, methodField := range interfaceType.Methods.List {
		_, ok := methodField.Type.(*ast.FuncType)
		if ok {
			if !f(methodField) {
				break
			}
		}
	}
}

func (gim *GoInterfaceMeta) TypeParams() []*GoVariableMeta {
	if gim.node == nil || gim.node.(*ast.TypeSpec).TypeParams == nil || len(gim.node.(*ast.TypeSpec).TypeParams.List) == 0 {
		return nil
	}

	tParamLen := len(gim.node.(*ast.TypeSpec).TypeParams.List)
	tParams := make([]*GoVariableMeta, 0, tParamLen)
	for _, field := range gim.node.(*ast.TypeSpec).TypeParams.List {
		for _, name := range field.Names {
			tParams = append(tParams, &GoVariableMeta{
				meta:     gim.newMeta(field),
				name:     name.String(),
				typeMeta: gim.newMeta(field.Type),
			})
		}
	}
	return tParams
}

type GoInterfaceMethodMeta struct {
	*meta
	interfaceMeta *GoInterfaceMeta
	receiverMeta  *GoVariableMeta
}

func IsInterfaceMethodNode(n ast.Node) bool {
	typeNode := n.(*ast.Field).Type
	if typeNode == nil {
		return false
	}
	funcType, ok := typeNode.(*ast.FuncType)
	return funcType != nil && ok
}

func NewGoInterfaceMethodMeta(m *meta, gim *GoInterfaceMeta) *GoInterfaceMethodMeta {
	return &GoInterfaceMethodMeta{meta: m, interfaceMeta: gim}
}

func (gimm *GoInterfaceMethodMeta) FunctionName() string {
	return gimm.node.(*ast.Field).Names[0].String()
}

func (gimm *GoInterfaceMethodMeta) Doc() []string {
	if gimm.node.(*ast.Field) == nil || gimm.node.(*ast.Field).Doc == nil || len(gimm.node.(*ast.Field).Doc.List) == 0 {
		return nil
	}
	commentSlice := make([]string, 0, len(gimm.node.(*ast.Field).Doc.List))
	for _, comment := range gimm.node.(*ast.Field).Doc.List {
		commentSlice = append(commentSlice, comment.Text)
	}
	return commentSlice
}

func (gimm *GoInterfaceMethodMeta) TypeParams() []*GoVariableMeta {
	return gimm.interfaceMeta.TypeParams()
}

func (gimm *GoInterfaceMethodMeta) Params() []*GoVariableMeta {
	if gimm.node.(*ast.Field).Type == nil || gimm.node.(*ast.Field).Type.(*ast.FuncType).Params == nil || len(gimm.node.(*ast.Field).Type.(*ast.FuncType).Params.List) == 0 {
		return nil
	}

	pLen := len(gimm.node.(*ast.Field).Type.(*ast.FuncType).Params.List)
	params := make([]*GoVariableMeta, 0, pLen)
	for index, field := range gimm.node.(*ast.Field).Type.(*ast.FuncType).Params.List {
		params = append(params, &GoVariableMeta{
			meta:     gimm.newMeta(field),
			name:     fmt.Sprintf("p%v", index),
			typeMeta: gimm.newMeta(field.Type),
		})

	}
	return params
}

func (gimm *GoInterfaceMethodMeta) ReturnTypes() []*GoVariableMeta {
	if gimm.node.(*ast.Field).Type == nil || gimm.node.(*ast.Field).Type.(*ast.FuncType).Results == nil || len(gimm.node.(*ast.Field).Type.(*ast.FuncType).Results.List) == 0 {
		return nil
	}

	rLen := len(gimm.node.(*ast.Field).Type.(*ast.FuncType).Results.List)
	returns := make([]*GoVariableMeta, 0, rLen)
	for _, field := range gimm.node.(*ast.Field).Type.(*ast.FuncType).Results.List {
		// TODO: not support named return value
		returns = append(returns, &GoVariableMeta{
			meta:     gimm.newMeta(field),
			name:     "",
			typeMeta: gimm.newMeta(field.Type),
		})
	}
	return returns
}

func (gimm *GoInterfaceMethodMeta) RecvInterface() (string, bool) {
	return gimm.interfaceMeta.InterfaceName(), true
}

func (gimm *GoInterfaceMethodMeta) Recv() *GoVariableMeta {
	if gimm.receiverMeta != nil {
		return gimm.receiverMeta
	}

	var receiverTypeExpr ast.Expr = ast.NewIdent(gimm.interfaceMeta.InterfaceName())
	typeParams := gimm.TypeParams()
	if l := len(typeParams); l > 0 {
		typeParamsExpr := make([]ast.Expr, 0, l)
		for _, typeParam := range typeParams {
			typeParamsExpr = append(typeParamsExpr, ast.NewIdent(typeParam.Name()))
		}
		if l == 1 {
			receiverTypeExpr = &ast.IndexExpr{
				X:     receiverTypeExpr,
				Index: typeParamsExpr[0],
			}
		} else {
			receiverTypeExpr = &ast.IndexListExpr{
				X:       receiverTypeExpr,
				Indices: typeParamsExpr,
			}
		}
	}
	gimm.receiverMeta = &GoVariableMeta{
		meta:     gimm.newMeta(gimm.interfaceMeta.node),
		name:     "i",
		typeMeta: gimm.newMeta(receiverTypeExpr),
	}

	return gimm.receiverMeta
}

func (gimm *GoInterfaceMethodMeta) MakeUnitTest(typeArgs []string) (string, []byte) {
	return makeTest(unittestMaker, gimm, "", typeArgs)
}

func (gimm *GoInterfaceMethodMeta) MakeBenchmark(typeArgs []string) (string, []byte) {
	return makeTest(benchmarkMaker, gimm, "", typeArgs)
}
