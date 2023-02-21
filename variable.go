package extractor

import "go/ast"

// type VariableTypeEnum int

// const (
// 	TYPE_UNKNOWN = iota
// 	TYPE_PKG_ALIAS
// 	TYPE_ASSIGNMENT
// 	TYPE_FUNC_CALL
// 	TYPE_VAR_FIELD
// 	TYPE_CONSTANTS
// )

// // 包别名
// // - name 别名
// // - typeMeta *GoImportMeta
// // 变量赋值
// // - name 变量名
// // - typeMeta 赋值表达式 *meta // TODO: 区分出来
// // 变量声明
// // - name 变量名
// // - typeMeta 类型 *meta
// // 函数参数表
// // - name 变量名
// // - typeMeta 类型 *meta
// // 函数调用
// // - name 函数名
// // - typeMeta 返回值类型 *meta
// // 常量
// // - name 值
// // - typeMeta 常量 *meta
type GoVariableMeta struct {
	*meta
	name     string
	typeMeta Meta
	// typeEnum VariableTypeEnum
}

// func ExtractGoVariableMeta(extractFilepath string, variableName string) (*GoVariableMeta, error) {
// 	goFileMeta, err := ExtractGoFileMeta(extractFilepath)
// 	if err != nil {
// 		return nil, err
// 	}

// 	gam := SearchGoVariableMeta(goFileMeta.meta, variableName)
// 	if gam.node == nil {
// 		return nil, fmt.Errorf("can not find function node")
// 	}

// 	return gam, nil
// }

// func SearchGoVariableMeta(m *meta, variableName string) *GoVariableMeta {
// 	var variableDecl *ast.ValueSpec
// 	var typeMeta Meta
// 	ast.Inspect(m.node, func(n ast.Node) bool {
// 		if IsVarNode(n) {
// 			decl := n.(*ast.ValueSpec)
// 			if decl.Names[0].String() == variableName {
// 				variableDecl = decl
// 				// TODO: typeMeta
// 			}
// 		}
// 		return variableDecl == nil
// 	})
// 	if variableDecl == nil {
// 		return nil
// 	}
// 	return &GoVariableMeta{
// 		meta:     m.newMeta(variableDecl),
// 		name:     variableName,
// 		typeMeta: typeMeta,
// 	}
// }

// // TODO:
// func IsVarNode(n ast.Node) bool {
// 	valueSpec, ok := n.(*ast.ValueSpec)
// 	if !ok {
// 		return false
// 	}
// 	return valueSpec.Names[0].Obj.Kind == ast.Var
// }

func (gvm *GoVariableMeta) Name() string {
	return gvm.name
}

func (gvm *GoVariableMeta) Type() string {
	return gvm.typeMeta.Expression()
}

func (gvm *GoVariableMeta) typeNode() ast.Expr {
	return gvm.typeMeta.(*meta).node.(ast.Expr)
}

func (gvm *GoVariableMeta) IsPointer() bool {
	_, ok := gvm.meta.node.(*ast.Field).Type.(*ast.StarExpr)
	return ok
}
