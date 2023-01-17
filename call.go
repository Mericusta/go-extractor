package extractor

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type GoCallMeta struct {
	*meta
	from *GoFromMeta
	args []*GoArgMeta
}

func ParseGoCallMeta(expression string) *GoCallMeta {
	callAST, err := parser.ParseExpr(expression)
	if err != nil {
		panic(err)
	}
	callExpr, ok := callAST.(*ast.CallExpr)
	if !ok {
		ast.Print(token.NewFileSet(), callAST)
		panic(callAST)
	}
	return &GoCallMeta{
		meta: &meta{node: callExpr},
		args: make([]*GoArgMeta, 0),
	}
}

func SearchGoCallMeta(m *meta, call string) []*GoCallMeta {
	lastSelectorIndex := strings.LastIndex(call, ".")
	searchSelector := lastSelectorIndex >= 0 && lastSelectorIndex < len(call)
	var fromMeta *GoFromMeta
	callMetaSlice := make([]*GoCallMeta, 0)
	ast.Inspect(m.node, func(n ast.Node) bool {
		if searchSelector && IsSelectorCallNode(n) {
			fromMeta = &GoFromMeta{meta: m.newMeta(n.(*ast.CallExpr).Fun.(*ast.SelectorExpr).X)}
			if fromMeta.Expression() == call[:lastSelectorIndex] && n.(*ast.CallExpr).Fun.(*ast.SelectorExpr).Sel.String() == call[lastSelectorIndex+1:] {
				callMetaSlice = append(callMetaSlice, &GoCallMeta{
					meta: m.newMeta(n.(*ast.CallExpr)),
					from: fromMeta,
				})
			}
		} else if !searchSelector && IsNonSelectorCallNode(n) {
			if n.(*ast.CallExpr).Fun.(*ast.Ident).String() == call {
				callMetaSlice = append(callMetaSlice, &GoCallMeta{
					meta: m.newMeta(n.(*ast.CallExpr)),
					from: fromMeta,
				})
			}
		}
		return true
	})
	return callMetaSlice
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

func (gcm *GoCallMeta) From() string {
	if gcm.from == nil {
		return ""
	}
	return gcm.from.Expression()
}

func (gcm *GoCallMeta) Call() string {
	if IsNonSelectorCallNode(gcm.node) {
		return gcm.node.(*ast.CallExpr).Fun.(*ast.Ident).String()
	} else if IsSelectorCallNode(gcm.node) {
		return gcm.node.(*ast.CallExpr).Fun.(*ast.SelectorExpr).Sel.String()
	}
	return ""
}

func (gcm *GoCallMeta) Args() []*GoArgMeta {
	callExpr := gcm.node.(*ast.CallExpr)
	if len(callExpr.Args) == 0 {
		return nil
	}
	args := make([]*GoArgMeta, 0, len(callExpr.Args))
	for _, argExpr := range callExpr.Args {
		args = append(args, &GoArgMeta{meta: gcm.meta.newMeta(argExpr)})
		// switch argExpr := argExpr.(type) {
		// case *ast.BasicLit: // 内建类型
		// 	switch argExpr.Kind {
		// 	case token.INT:
		// 		argValue, err := strconv.ParseInt(argExpr.Value, 10, 32)
		// 		if err != nil {
		// 			panic(err)
		// 		}
		// 		args = append(args, &GoArgMeta{
		// 			node:  argExpr,
		// 			value: argValue,
		// 		})
		// 	case token.STRING:
		// 		args = append(args, &GoArgMeta{
		// 			node:  argExpr,
		// 			value: argExpr.Value,
		// 		})
		// 	}
		// case *ast.Ident: // 变量
		// 	args = append(args, &GoArgMeta{
		// 		node: argExpr,
		// 		arg:  argExpr.String(),
		// 	})
		// case *ast.SelectorExpr: // 变量及其 from
		// 	args = append(args, &GoArgMeta{
		// 		node: argExpr,
		// 		from: argExpr.X.(*ast.Ident).String(),
		// 		arg:  argExpr.Sel.String(),
		// 	})
		// case *ast.CallExpr: // 函数调用
		// 	args = append(args, &GoArgMeta{
		// 		node: argExpr,
		// 		callMeta: &GoCallMeta{
		// 			fileMeta: gcm.fileMeta,
		// 			callExpr: argExpr,
		// 		},
		// 	})
		// }
	}
	return args
}
