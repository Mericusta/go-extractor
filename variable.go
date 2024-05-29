package extractor

import (
	"go/ast"
	"strings"
)

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
	*meta    // *ast.Field
	name     string
	typeMeta Meta
	// typeEnum VariableTypeEnum
}

func (gvm *GoVariableMeta) Tag() string {
	if gvm.node.(*ast.Field).Tag == nil {
		return ""
	}
	return gvm.node.(*ast.Field).Tag.Value
}

func (gvm *GoVariableMeta) Doc() []string {
	if gvm.node.(*ast.Field).Doc == nil {
		return nil
	}
	commentSlice := make([]string, 0, len(gvm.node.(*ast.Field).Doc.List))
	for _, comment := range gvm.node.(*ast.Field).Doc.List {
		commentSlice = append(commentSlice, comment.Text)
	}
	return commentSlice
}

func (gvm *GoVariableMeta) Comment() string {
	if gvm.node.(*ast.Field).Comment == nil || len(gvm.node.(*ast.Field).Comment.List) == 0 {
		return ""
	}
	return gvm.node.(*ast.Field).Comment.List[0].Text
}

type UnderlyingType int

const (
	UNDERLYING_TYPE_IDENT     = iota + 1 // any others *ast.Ident
	UNDERLYING_TYPE_ARRAY                // *ast.ArrayType
	UNDERLYING_TYPE_STRUCT               // *ast.StructType
	UNDERLYING_TYPE_POINTER              // *ast.StarExpr
	UNDERLYING_TYPE_FUNC                 // *ast.FuncType
	UNDERLYING_TYPE_INTERFACE            // *ast.InterfaceType
	UNDERLYING_TYPE_MAP                  // *ast.MapType
	UNDERLYING_TYPE_CHAN                 // *ast.ChanType
)

func (gvm *GoVariableMeta) Type() (string, string, UnderlyingType) {
	underlyingString, underlyingEnum := "", UnderlyingType(0)
	switch gvm.typeMeta.(*meta).node.(type) {
	case *ast.ArrayType:
		underlyingString, underlyingEnum = "array", UNDERLYING_TYPE_ARRAY
	case *ast.StructType:
		underlyingString, underlyingEnum = "struct", UNDERLYING_TYPE_STRUCT
	case *ast.StarExpr:
		underlyingString, underlyingEnum = "pointer", UNDERLYING_TYPE_POINTER
	case *ast.FuncType:
		underlyingString, underlyingEnum = "func", UNDERLYING_TYPE_FUNC
	case *ast.InterfaceType:
		underlyingString, underlyingEnum = "interface", UNDERLYING_TYPE_INTERFACE
	case *ast.MapType:
		underlyingString, underlyingEnum = "map", UNDERLYING_TYPE_MAP
	case *ast.ChanType:
		underlyingString, underlyingEnum = "chan", UNDERLYING_TYPE_CHAN
	default:
		underlyingString, underlyingEnum = gvm.typeMeta.Expression(), UNDERLYING_TYPE_IDENT
	}
	return gvm.typeMeta.Expression(), underlyingString, underlyingEnum
}

func traitFrom(expression string, isPointer bool) string {
	if isPointer {
		expression = expression[1:]
	}
	expressionSlice := strings.Split(expression, ".")
	l := len(expressionSlice)
	if l > 1 {
		return strings.Join(expressionSlice[:l-1], ".")
	}
	return ""
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

func (gvm *GoVariableMeta) typeNode() ast.Expr {
	return gvm.typeMeta.(*meta).node.(ast.Expr)
}

func (gvm *GoVariableMeta) IsPointer() bool {
	_, ok := gvm.meta.node.(*ast.Field).Type.(*ast.StarExpr)
	return ok
}
