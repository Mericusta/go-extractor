package extractor

import (
	"go/ast"
)

type GoTypeConstraintsMetaTypeConstraints interface {
	*ast.TypeSpec

	ast.Node
}

type GoTypeConstraintsMeta[T GoTypeConstraintsMetaTypeConstraints] struct {
	// typeSpec     *ast.TypeSpec
	*meta[T]
	commentGroup *ast.CommentGroup
	// methodMeta   map[string]*GoTypeConstraintsMethodMeta
}

// func ExtractGoTypeConstraintsMeta[T GoTypeConstraintsMetaTypeConstraints](extractFilepath, interfaceName string) (*GoTypeConstraintsMeta[T], error) {
// 	gfm, err := ExtractGoFileMeta(extractFilepath)
// 	if err != nil {
// 		return nil, err
// 	}

// 	gtcm := SearchGoTypeConstraintsMeta(gfm, interfaceName)
// 	if gtcm == nil {
// 		return nil, fmt.Errorf("can not find interface node")
// 	}

// 	return gtcm, nil
// }

// func SearchGoTypeConstraintsMeta[T GoTypeConstraintsMetaTypeConstraints](gfm *GoFileMeta[*ast.File], interfaceName string) *GoTypeConstraintsMeta[T] {
// 	var interfaceDecl *ast.TypeSpec
// 	var commentDecl *ast.CommentGroup
// 	ast.Inspect(gfm.node, func(n ast.Node) bool {
// 		if genDecl, ok := n.(*ast.GenDecl); ok {
// 			ast.Inspect(genDecl, func(n ast.Node) bool {
// 				if IsTypeNode(n) {
// 					typeSpec := n.(*ast.TypeSpec)
// 					if typeSpec.Name.String() == interfaceName {
// 						interfaceDecl = typeSpec
// 						commentDecl = genDecl.Doc
// 						return false
// 					}
// 				}
// 				return true
// 			})
// 			return false // genDecl traverse done
// 		}
// 		return interfaceDecl == nil // already found
// 	})
// 	if interfaceDecl == nil {
// 		return nil
// 	}
// 	return &GoTypeConstraintsMeta[T]{
// 		meta:         newMeta(interfaceDecl, gfm.path),
// 		commentGroup: commentDecl,
// 		methodMeta:   make(map[string]*GoTypeConstraintsMethodMeta),
// 	}
// }

// func (gtcm *GoTypeConstraintsMeta) InterfaceName() string {
// 	return gtcm.node.(*ast.TypeSpec).Name.String()
// }

// func (gtcm *GoTypeConstraintsMeta) Doc() []string {
// 	if gtcm.node == nil || gtcm.commentGroup == nil || len(gtcm.commentGroup.List) == 0 {
// 		return nil
// 	}
// 	commentSlice := make([]string, 0, len(gtcm.commentGroup.List))
// 	for _, comment := range gtcm.commentGroup.List {
// 		commentSlice = append(commentSlice, comment.Text)
// 	}
// 	return commentSlice
// }

// // SearchMethodDecl search method decl from node.(*ast.InterfaceType)
// func (gtcm *GoTypeConstraintsMeta) SearchMethodDecl(methodName string) *GoTypeConstraintsMethodMeta {
// 	gtcm.ForeachMethodDecl(func(f *ast.Field) bool {
// 		if f.Names[0].String() == methodName && IsInterfaceMethodNode(f) {
// 			gtcm.methodMeta[methodName] = NewGoTypeConstraintsMethodMeta(
// 				gtcm.copyMeta(f), gtcm,
// 			)
// 			return false
// 		}
// 		return true
// 	})
// 	return gtcm.methodMeta[methodName]
// }

// func (gtcm *GoTypeConstraintsMeta) ForeachMethodDecl(f func(*ast.Field) bool) {
// 	interfaceType := gtcm.node.(*ast.TypeSpec).Type.(*ast.InterfaceType)
// 	if interfaceType.Methods == nil {
// 		return
// 	}
// 	for _, methodField := range interfaceType.Methods.List {
// 		_, ok := methodField.Type.(*ast.FuncType)
// 		if ok {
// 			if !f(methodField) {
// 				break
// 			}
// 		}
// 	}
// }

// func (gtcm *GoTypeConstraintsMeta) TypeParams() []*GoVarMeta {
// 	if gtcm.node == nil || gtcm.node.(*ast.TypeSpec).TypeParams == nil || len(gtcm.node.(*ast.TypeSpec).TypeParams.List) == 0 {
// 		return nil
// 	}

// 	tParamLen := len(gtcm.node.(*ast.TypeSpec).TypeParams.List)
// 	tParams := make([]*GoVarMeta, 0, tParamLen)
// 	for _, field := range gtcm.node.(*ast.TypeSpec).TypeParams.List {
// 		for _, name := range field.Names {
// 			tParams = append(tParams, &GoVarMeta{
// 				meta:     gtcm.copyMeta(field),
// 				ident:    name.String(),
// 				typeMeta: gtcm.copyMeta(field.Type),
// 			})
// 		}
// 	}
// 	return tParams
// }

// type GoTypeConstraintsMethodMeta struct {
// 	*meta
// 	interfaceMeta *GoTypeConstraintsMeta
// 	receiverMeta  *GoVarMeta
// }

// func NewGoTypeConstraintsMethodMeta(m *meta, gtcm *GoTypeConstraintsMeta) *GoTypeConstraintsMethodMeta {
// 	return &GoTypeConstraintsMethodMeta{meta: m, interfaceMeta: gtcm}
// }

// func (gtcmm *GoTypeConstraintsMethodMeta) FunctionName() string {
// 	return gtcmm.node.(*ast.Field).Names[0].String()
// }

// func (gtcmm *GoTypeConstraintsMethodMeta) Doc() []string {
// 	if gtcmm.node.(*ast.Field) == nil || gtcmm.node.(*ast.Field).Doc == nil || len(gtcmm.node.(*ast.Field).Doc.List) == 0 {
// 		return nil
// 	}
// 	commentSlice := make([]string, 0, len(gtcmm.node.(*ast.Field).Doc.List))
// 	for _, comment := range gtcmm.node.(*ast.Field).Doc.List {
// 		commentSlice = append(commentSlice, comment.Text)
// 	}
// 	return commentSlice
// }

// func (gtcmm *GoTypeConstraintsMethodMeta) TypeParams() []*GoVarMeta {
// 	return gtcmm.interfaceMeta.TypeParams()
// }

// func (gtcmm *GoTypeConstraintsMethodMeta) Params() []*GoVarMeta {
// 	if gtcmm.node.(*ast.Field).Type == nil || gtcmm.node.(*ast.Field).Type.(*ast.FuncType).Params == nil || len(gtcmm.node.(*ast.Field).Type.(*ast.FuncType).Params.List) == 0 {
// 		return nil
// 	}

// 	pLen := len(gtcmm.node.(*ast.Field).Type.(*ast.FuncType).Params.List)
// 	params := make([]*GoVarMeta, 0, pLen)
// 	for index, field := range gtcmm.node.(*ast.Field).Type.(*ast.FuncType).Params.List {
// 		params = append(params, &GoVarMeta{
// 			meta:     gtcmm.copyMeta(field),
// 			ident:    fmt.Sprintf("p%v", index),
// 			typeMeta: gtcmm.copyMeta(field.Type),
// 		})

// 	}
// 	return params
// }

// func (gtcmm *GoTypeConstraintsMethodMeta) ReturnTypes() []*GoVarMeta {
// 	if gtcmm.node.(*ast.Field).Type == nil || gtcmm.node.(*ast.Field).Type.(*ast.FuncType).Results == nil || len(gtcmm.node.(*ast.Field).Type.(*ast.FuncType).Results.List) == 0 {
// 		return nil
// 	}

// 	rLen := len(gtcmm.node.(*ast.Field).Type.(*ast.FuncType).Results.List)
// 	returns := make([]*GoVarMeta, 0, rLen)
// 	for _, field := range gtcmm.node.(*ast.Field).Type.(*ast.FuncType).Results.List {
// 		// TODO: not support named return value
// 		returns = append(returns, &GoVarMeta{
// 			meta:     gtcmm.copyMeta(field),
// 			ident:    "",
// 			typeMeta: gtcmm.copyMeta(field.Type),
// 		})
// 	}
// 	return returns
// }

// func (gtcmm *GoTypeConstraintsMethodMeta) RecvInterface() (string, bool) {
// 	return gtcmm.interfaceMeta.InterfaceName(), true
// }

// func (gtcmm *GoTypeConstraintsMethodMeta) Recv() *GoVarMeta {
// 	if gtcmm.receiverMeta != nil {
// 		return gtcmm.receiverMeta
// 	}

// 	var receiverTypeExpr ast.Expr = ast.NewIdent(gtcmm.interfaceMeta.InterfaceName())
// 	typeParams := gtcmm.TypeParams()
// 	if l := len(typeParams); l > 0 {
// 		typeParamsExpr := make([]ast.Expr, 0, l)
// 		for _, typeParam := range typeParams {
// 			typeParamsExpr = append(typeParamsExpr, ast.NewIdent(typeParam.Ident()))
// 		}
// 		if l == 1 {
// 			receiverTypeExpr = &ast.IndexExpr{
// 				X:     receiverTypeExpr,
// 				Index: typeParamsExpr[0],
// 			}
// 		} else {
// 			receiverTypeExpr = &ast.IndexListExpr{
// 				X:       receiverTypeExpr,
// 				Indices: typeParamsExpr,
// 			}
// 		}
// 	}
// 	gtcmm.receiverMeta = &GoVarMeta{
// 		meta:     gtcmm.copyMeta(gtcmm.interfaceMeta.node),
// 		ident:    "i",
// 		typeMeta: gtcmm.copyMeta(receiverTypeExpr),
// 	}

// 	return gtcmm.receiverMeta
// }
