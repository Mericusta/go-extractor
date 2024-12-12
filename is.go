package extractor

import (
	"go/ast"
	"go/token"
)

func IsImportNode(n ast.Node) bool {
	_, ok := n.(*ast.ImportSpec)
	return ok
}

// 判断是否是 ast.Node 是否是 var 关键字定义域
func IsVarNode(n ast.Node) bool {
	genDecl, ok := n.(*ast.GenDecl)
	return ok && genDecl.Tok == token.VAR
}

func IsFuncNode(n ast.Node) bool {
	funcDecl, ok := n.(*ast.FuncDecl)
	return ok && funcDecl.Recv == nil
}

// IsTypeNode 判断 ast.Node 是否是 type 关键字定义域
func IsTypeNode(n ast.Node) bool {
	genDecl, ok := n.(*ast.GenDecl)
	return genDecl != nil && ok
}

// IsStructNode 判断 ast.Spec 是否是 struct 节点
func IsStructNode(n ast.Spec) bool {
	typeSpec, ok := n.(*ast.TypeSpec)
	if typeSpec == nil || !ok {
		return false
	}
	structType, ok := typeSpec.Type.(*ast.StructType)
	return structType != nil && ok
}

func IsMethodNode(n ast.Node) bool {
	decl, ok := n.(*ast.FuncDecl)
	return ok && decl.Recv != nil && len(decl.Recv.List) == 1
}

// IsInterfaceNode 判断 ast.Spec 是否是 interface 节点
func IsInterfaceNode(n ast.Spec) bool {
	typeSpec, ok := n.(*ast.TypeSpec)
	if typeSpec == nil || !ok {
		return false
	}

	interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
	return interfaceType != nil && ok
}

// IsInterfaceMethodNode
func IsInterfaceMethodNode(n *ast.Field) bool {
	typeNode := n.Type
	if typeNode == nil {
		return false
	}
	funcType, ok := typeNode.(*ast.FuncType)
	return funcType != nil && ok
}

func IsTypeConstraintsNode(n ast.Spec) bool {
	typeSpec, ok := n.(*ast.TypeSpec)
	if !ok {
		return false
	}
	if typeSpec.Type == nil {
		return false
	}
	interfaceTypeNode, ok := typeSpec.Type.(*ast.InterfaceType)
	if !ok {
		return false
	}
	if interfaceTypeNode == nil {
		return false
	}
	if interfaceTypeNode.Methods == nil || len(interfaceTypeNode.Methods.List) == 0 {
		return true
	}
	// 存在 *ast.BinaryExpr 不存在 *ast.FuncType 即视为 类型约束
	hasBinaryExpr, hasFuncType := false, false
	for _, method := range interfaceTypeNode.Methods.List {
		_, isBinaryExpr := method.Type.(*ast.BinaryExpr)
		if isBinaryExpr {
			hasBinaryExpr = true
		}
		_, isFuncType := method.Type.(*ast.FuncType)
		if isFuncType {
			hasFuncType = true
		}
	}
	return hasBinaryExpr && !hasFuncType
}

func IsTypeConstraintsMethodNode(n ast.Node) bool {
	typeNode := n.(*ast.Field).Type
	if typeNode == nil {
		return false
	}
	funcType, ok := typeNode.(*ast.FuncType)
	return funcType != nil && ok
}

func IsSelectorNode(n ast.Node) bool {
	_, ok := n.(*ast.SelectorExpr)
	return ok
}

func IsCallNode(n ast.Node) bool {
	_, ok := n.(*ast.CallExpr)
	return ok
}

func IsNonSelectorCallNode(node ast.Node) bool {
	callExpr, ok := node.(*ast.CallExpr)
	if !ok {
		return false
	}
	if callExpr.Fun == nil {
		return false
	}
	_, ok = callExpr.Fun.(*ast.Ident)
	return ok
}

func IsSelectorCallNode(node ast.Node) bool {
	callExpr, ok := node.(*ast.CallExpr)
	if !ok {
		return false
	}
	if callExpr.Fun == nil {
		return false
	}
	_, ok = callExpr.Fun.(*ast.SelectorExpr)
	return ok
}
