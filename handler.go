package extractor

// import "go/ast"

// var (
// 	nodeHandler         func(n ast.Node, post ...func(ast.Node) bool) bool
// 	starExprHandler     func(n ast.Node) bool
// 	selectorExprHandler func(n ast.Node) bool
// 	indexExprHandler    func(n ast.Node) bool
// )

// func init() {
// 	nodeHandler = func(n ast.Node, posts ...func(ast.Node) bool) bool {
// 		ident, ok := n.(*ast.Ident)
// 		if ident != nil && ok {
// 			typeIdent = ident.String()
// 			return false
// 		} else {
// 			for _, post := range posts {
// 				if post != nil && !post(n) {
// 					return false
// 				}
// 			}
// 		}
// 		return true
// 	}
// 	starExprHandler = func(n ast.Node) bool {
// 		// 遇到 StarExpr 取 X
// 		starExpr, ok := n.(*ast.StarExpr)
// 		if starExpr == nil || !ok {
// 			return true
// 		}
// 		return nodeHandler(starExpr.X, starExprHandler, selectorExprHandler, indexExprHandler)
// 	}
// 	selectorExprHandler = func(n ast.Node) bool {
// 		// 遇到 SelectorExpr 取 Sel
// 		selectorExpr, ok := n.(*ast.SelectorExpr)
// 		if selectorExpr == nil || !ok {
// 			return true
// 		}
// 		return nodeHandler(selectorExpr.Sel, starExprHandler, selectorExprHandler, indexExprHandler)
// 	}
// 	indexExprHandler = func(n ast.Node) bool {
// 		// 遇到 IndexExp 取 X
// 		indexExpr, ok := n.(*ast.IndexExpr)
// 		if indexExpr == nil || !ok {
// 			return true
// 		}
// 		return nodeHandler(indexExpr.X, starExprHandler, selectorExprHandler, indexExprHandler)
// 	}
// }
