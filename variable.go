package extractor

import (
	"fmt"
	"go/ast"
	"go/parser"
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

type GoVarMetaTypeConstraints interface {
	*ast.ValueSpec | *ast.Field | *ast.InterfaceType

	ast.Node
}

// GoVarMeta go var 的 meta 数据
type GoVarMeta struct {
	// 组合基本 meta 数据
	// ast 节点，要求为 *ast.ValueSpec 或 *ast.Field 或 *ast.InterfaceType
	// 以 ast 节点 为单位执行 AST/PrintAST/Expression/Format
	*meta

	// var 标识
	ident string

	// 基类型: *ExampleStruct[T] -> ExampleStruct
	typeIdent string

	// 类型表达式: *ExampleStruct[T]
	typeExpression string

	// 是否是接口
	isInterface bool

	// 是否是指针
	isPointer bool
}

// newGoVarMeta 通过 ast 构造 var 的 meta 数据
func newGoVarMeta(m *meta, ident string, stopExtract ...bool) *GoVarMeta {
	gvm := &GoVarMeta{meta: m, ident: ident}
	if len(stopExtract) == 0 {
		gvm.ExtractAll()
	}
	return gvm
}

// -------------------------------- unit test --------------------------------

func (gvm *GoVarMeta) Ident() string          { return gvm.ident }
func (gvm *GoVarMeta) TypeIdent() string      { return gvm.typeIdent }
func (gvm *GoVarMeta) TypeExpression() string { return gvm.typeExpression }

// -------------------------------- unit test --------------------------------

// -------------------------------- extractor --------------------------------

// ExtractGoVarMeta 通过文件的绝对路径和 var 的 标识 提取文件中 var 的 meta 数据
func ExtractGoVarMeta(extractFilepath, varIdent string) (*GoVarMeta, error) {
	// 提取 package
	gpm, err := ExtractGoPackageMeta(extractFilepath, nil)
	if err != nil {
		return nil, err
	}

	// 提取 var
	gpm.extractVar()

	// 搜索 var
	gvm := gpm.SearchVarMeta(varIdent)
	if gvm == nil {
		return nil, fmt.Errorf("can not find var node")
	}

	return gvm, nil
}

func (gvm *GoVarMeta) ExtractAll() {
	// 提取 类型
	gvm.extractType()
}

func (gvm *GoVarMeta) extractType() {
	var (
		n        ast.Node = gvm.node
		typeExpr ast.Expr
	)
	switch node := n.(type) {
	case *ast.ValueSpec:
		typeExpr = node.Type
	case *ast.Field:
		typeExpr = node.Type
	}
	if typeExpr == nil {
		return
	}
	// 取 type expression
	gvm.typeExpression = newMeta(typeExpr, gvm.path).Expression()
	// 取 type ident
	var (
		nodeHandler          func(n ast.Node, post ...func(ast.Node) bool) bool
		starExprHandler      func(n ast.Node) bool
		selectorExprHandler  func(n ast.Node) bool
		indexExprHandler     func(n ast.Node) bool
		indexListExprHandler func(n ast.Node) bool
		interfaceTypeHandler func(n ast.Node) bool
		arrayTypeHandler     func(n ast.Node) bool
		mapTypeHandler       func(n ast.Node) bool
	)
	nodeHandler = func(n ast.Node, posts ...func(ast.Node) bool) bool {
		ident, ok := n.(*ast.Ident)
		if ident != nil && ok {
			gvm.typeIdent = ident.String()
			return false
		} else {
			for _, post := range posts {
				if post != nil && !post(n) {
					return false
				}
			}
		}
		return true
	}
	starExprHandler = func(n ast.Node) bool {
		// 遇到 StarExpr 取 X
		starExpr, ok := n.(*ast.StarExpr)
		if starExpr == nil || !ok {
			return true
		}
		return nodeHandler(starExpr.X, starExprHandler, selectorExprHandler, indexExprHandler, indexListExprHandler, interfaceTypeHandler, arrayTypeHandler, mapTypeHandler)
	}
	selectorExprHandler = func(n ast.Node) bool {
		// 遇到 SelectorExpr 取 Sel
		selectorExpr, ok := n.(*ast.SelectorExpr)
		if selectorExpr == nil || !ok {
			return true
		}
		return nodeHandler(selectorExpr.Sel, starExprHandler, selectorExprHandler, indexExprHandler, indexListExprHandler, interfaceTypeHandler, arrayTypeHandler, mapTypeHandler)
	}
	indexExprHandler = func(n ast.Node) bool {
		// 遇到 IndexExp 取 X
		indexExpr, ok := n.(*ast.IndexExpr)
		if indexExpr == nil || !ok {
			return true
		}
		return nodeHandler(indexExpr.X, starExprHandler, selectorExprHandler, indexExprHandler, indexListExprHandler, interfaceTypeHandler, arrayTypeHandler, mapTypeHandler)
	}
	indexListExprHandler = func(n ast.Node) bool {
		// 遇到 IndexExp 取 X
		indexExpr, ok := n.(*ast.IndexListExpr)
		if indexExpr == nil || !ok {
			return true
		}
		return nodeHandler(indexExpr.X, starExprHandler, selectorExprHandler, indexExprHandler, indexListExprHandler, interfaceTypeHandler, arrayTypeHandler, mapTypeHandler)
	}
	interfaceTypeHandler = func(n ast.Node) bool {
		// 遇到 InterfaceType 直接取整个 expression -> 最好是 any
		// TODO: 可能需要处理特殊语法，带有 method 的 interface
		interfaceType, ok := n.(*ast.InterfaceType)
		if interfaceType == nil || !ok {
			return true
		}
		gvm.typeIdent = gvm.typeExpression
		return false
	}
	arrayTypeHandler = func(n ast.Node) bool {
		// 遇到 ArrayType 直接取整个 expression
		arrayType, ok := n.(*ast.ArrayType)
		if arrayType == nil || !ok {
			return true
		}
		gvm.typeIdent = gvm.typeExpression
		return false
	}
	mapTypeHandler = func(n ast.Node) bool {
		// 遇到 MapType 直接取整个 expression
		mapType, ok := n.(*ast.MapType)
		if mapType == nil || !ok {
			return true
		}
		gvm.typeIdent = gvm.typeExpression
		return false
	}
	ast.Inspect(typeExpr, func(n ast.Node) bool {
		return n != nil && nodeHandler(n, starExprHandler, selectorExprHandler, indexExprHandler, indexListExprHandler, interfaceTypeHandler, arrayTypeHandler, mapTypeHandler)
	})
}

// -------------------------------- extractor --------------------------------

// -------------------------------- maker --------------------------------

// MakeUpVarMeta 通过参数生成 var 的 meta 数据
func MakeUpVarMeta(ident, typeExpression string) *GoVarMeta {
	typeExpr, err := parser.ParseExpr(typeExpression)
	if typeExpr == nil || err != nil {
		return nil
	}
	astField := &ast.Field{
		Names: []*ast.Ident{ast.NewIdent(ident)},
		Type:  typeExpr,
	}

	gvm := newGoVarMeta(newMeta(astField, ""), ident, true)

	return gvm
}

func (gvm *GoVarMeta) make() *ast.Field {
	var (
		n        ast.Node = gvm.node
		typeExpr ast.Expr
	)
	switch node := n.(type) {
	case *ast.ValueSpec:
		typeExpr = node.Type
	case *ast.Field:
		typeExpr = node.Type
	}
	return &ast.Field{
		Names: []*ast.Ident{ast.NewIdent(gvm.ident)},
		Type:  typeExpr,
	}
}

// -------------------------------- maker --------------------------------

// func (gvm *GoVarMeta[T]) Tag() string {
// 	if gvm.node.(*ast.Field).Tag == nil {
// 		return ""
// 	}
// 	return gvm.node.(*ast.Field).Tag.Value
// }

// func (gvm *GoVarMeta[T]) Doc() []string {
// 	if gvm.node.(*ast.Field).Doc == nil {
// 		return nil
// 	}
// 	commentSlice := make([]string, 0, len(gvm.node.(*ast.Field).Doc.List))
// 	for _, comment := range gvm.node.(*ast.Field).Doc.List {
// 		commentSlice = append(commentSlice, comment.Text)
// 	}
// 	return commentSlice
// }

// func (gvm *GoVarMeta[T]) Comment() string {
// 	if gvm.node.(*ast.Field).Comment == nil || len(gvm.node.(*ast.Field).Comment.List) == 0 {
// 		return ""
// 	}
// 	return gvm.node.(*ast.Field).Comment.List[0].Text
// }

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

// func (gvm *GoVarMeta[T]) Type() (string, string, UnderlyingType) {
// 	underlyingString, underlyingEnum := "", UnderlyingType(0)
// 	switch gvm.typeMeta.(*meta).node.(type) {
// 	case *ast.ArrayType:
// 		underlyingString, underlyingEnum = "array", UNDERLYING_TYPE_ARRAY
// 	case *ast.StructType:
// 		underlyingString, underlyingEnum = "struct", UNDERLYING_TYPE_STRUCT
// 	case *ast.StarExpr:
// 		underlyingString, underlyingEnum = "pointer", UNDERLYING_TYPE_POINTER
// 	case *ast.FuncType:
// 		underlyingString, underlyingEnum = "func", UNDERLYING_TYPE_FUNC
// 	case *ast.InterfaceType:
// 		underlyingString, underlyingEnum = "interface", UNDERLYING_TYPE_INTERFACE
// 	case *ast.MapType:
// 		underlyingString, underlyingEnum = "map", UNDERLYING_TYPE_MAP
// 	case *ast.ChanType:
// 		underlyingString, underlyingEnum = "chan", UNDERLYING_TYPE_CHAN
// 	default:
// 		underlyingString, underlyingEnum = gvm.typeMeta.Expression(), UNDERLYING_TYPE_IDENT
// 	}
// 	return gvm.typeMeta.Expression(), underlyingString, underlyingEnum
// }

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

// func (gvm *GoVarMeta[T]) typeNode() ast.Expr {
// 	return gvm.typeMeta.(*meta).node.(ast.Expr)
// }

// func (gvm *GoVarMeta[T]) IsPointer() bool {
// 	_, ok := gvm.meta.node.(*ast.Field).Type.(*ast.StarExpr)
// 	return ok
// }
