package extractor

// import (
// 	"fmt"
// 	"go/ast"
// 	"go/parser"
// 	"go/token"
// 	"strings"
// )

// type GoCallMeta struct {
// 	*meta
// 	from *GoFromMeta // TODO: 改成 GoVariableMeta
// 	args []*GoArgMeta
// }

// func ParseGoCallMeta(expression string) *GoCallMeta {
// 	callAST, err := parser.ParseExpr(expression)
// 	if err != nil {
// 		panic(err)
// 	}
// 	callExpr, ok := callAST.(*ast.CallExpr)
// 	if !ok {
// 		ast.Print(token.NewFileSet(), callAST)
// 		panic(callAST)
// 	}
// 	return &GoCallMeta{
// 		meta: &meta{node: callExpr},
// 		args: make([]*GoArgMeta, 0),
// 	}
// }

// func SearchGoCallMeta(m *meta, call string) []*GoCallMeta {
// 	lastSelectorIndex := strings.LastIndex(call, ".")
// 	searchSelector := lastSelectorIndex >= 0 && lastSelectorIndex < len(call)
// 	var fromMeta *GoFromMeta
// 	callMetaSlice := make([]*GoCallMeta, 0)
// 	ast.Inspect(m.node, func(n ast.Node) bool {
// 		if searchSelector && IsSelectorCallNode(n) {
// 			fromMeta = &GoFromMeta{meta: m.newMeta(n.(*ast.CallExpr).Fun.(*ast.SelectorExpr).X)}
// 			if fromMeta.Expression() == call[:lastSelectorIndex] && n.(*ast.CallExpr).Fun.(*ast.SelectorExpr).Sel.String() == call[lastSelectorIndex+1:] {
// 				callMetaSlice = append(callMetaSlice, &GoCallMeta{
// 					meta: m.newMeta(n.(*ast.CallExpr)),
// 					from: fromMeta,
// 				})
// 			}
// 		} else if !searchSelector && IsNonSelectorCallNode(n) {
// 			if n.(*ast.CallExpr).Fun.(*ast.Ident).String() == call {
// 				callMetaSlice = append(callMetaSlice, &GoCallMeta{
// 					meta: m.newMeta(n.(*ast.CallExpr)),
// 					from: fromMeta,
// 				})
// 			}
// 		}
// 		return true
// 	})
// 	return callMetaSlice
// }

// func ExtractGoCallMeta(m *meta) map[string][]*GoCallMeta {
// 	callMetaMap := make(map[string][]*GoCallMeta, 0)
// 	ast.Inspect(m.node, func(n ast.Node) bool {
// 		if IsSelectorCallNode(n) {
// 			fromMeta := &GoFromMeta{meta: m.newMeta(n.(*ast.CallExpr).Fun.(*ast.SelectorExpr).X)}
// 			callExpression := fmt.Sprintf("%v.%v", fromMeta.Expression(), n.(*ast.CallExpr).Fun.(*ast.SelectorExpr).Sel.String())
// 			callMetaMap[callExpression] = append(callMetaMap[callExpression], &GoCallMeta{
// 				meta: m.newMeta(n.(*ast.CallExpr)),
// 				from: fromMeta,
// 			})
// 		} else if IsNonSelectorCallNode(n) {
// 			callExpression := n.(*ast.CallExpr).Fun.(*ast.Ident).String()
// 			callMetaMap[callExpression] = append(callMetaMap[callExpression], &GoCallMeta{
// 				meta: m.newMeta(n.(*ast.CallExpr)),
// 			})
// 		}
// 		return true
// 	})
// 	return callMetaMap
// }

// func (gcm *GoCallMeta) From() string {
// 	if gcm.from == nil {
// 		return ""
// 	}
// 	return gcm.from.Expression()
// }

// func (gcm *GoCallMeta) Call() string {
// 	if IsNonSelectorCallNode(gcm.node) {
// 		return gcm.node.(*ast.CallExpr).Fun.(*ast.Ident).String()
// 	} else if IsSelectorCallNode(gcm.node) {
// 		return gcm.node.(*ast.CallExpr).Fun.(*ast.SelectorExpr).Sel.String()
// 	}
// 	return ""
// }

// func (gcm *GoCallMeta) Args() []*GoArgMeta {
// 	callExpr := gcm.node.(*ast.CallExpr)
// 	if len(callExpr.Args) == 0 {
// 		return nil
// 	}
// 	args := make([]*GoArgMeta, 0, len(callExpr.Args))
// 	for _, argExpr := range callExpr.Args {
// 		args = append(args, &GoArgMeta{meta: gcm.meta.newMeta(argExpr)})
// 	}
// 	return args
// }
