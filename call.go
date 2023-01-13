package extractor

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
)

type GoCallMeta struct {
	fileMeta *GoFileMeta
	callExpr *ast.CallExpr
	args     []*GoArgMeta
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
		callExpr: callExpr,
		args:     make([]*GoArgMeta, 0),
	}
}

func SearchGoCallMeta(fileMeta *GoFileMeta, node ast.Node, call, from string) *GoCallMeta {
	searchSelector := len(from) > 0
	var callDecl *ast.CallExpr
	ast.Inspect(node, func(n ast.Node) bool {
		if searchSelector && IsSelectorCallNode(n) {
			fromIdent := n.(*ast.CallExpr).Fun.(*ast.SelectorExpr).X.(*ast.Ident)
			callIdent := n.(*ast.CallExpr).Fun.(*ast.SelectorExpr).Sel
			if fromIdent.Name == from && callIdent.Name == call {
				callDecl = n.(*ast.CallExpr)
			}
		} else if !searchSelector && IsNonSelectorCallNode(n) {
			if n.(*ast.CallExpr).Fun.(*ast.Ident).Name == call {
				callDecl = n.(*ast.CallExpr)
			}
		}
		return callDecl == nil
	})
	return &GoCallMeta{
		fileMeta: fileMeta,
		callExpr: callDecl,
	}
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

func (gcm *GoCallMeta) PrintAST() {
	ast.Print(token.NewFileSet(), gcm.callExpr)
}

// TODO: 非必要不存储额外信息
func (gcm *GoCallMeta) Expression() string {
	return gcm.fileMeta.Expression(gcm.callExpr.Pos(), gcm.callExpr.End())
}

func (gcm *GoCallMeta) Call() string {
	if IsNonSelectorCallNode(gcm.callExpr) {
		return gcm.callExpr.Fun.(*ast.Ident).Name
	} else if IsSelectorCallNode(gcm.callExpr) {
		return gcm.callExpr.Fun.(*ast.SelectorExpr).Sel.Name
	}
	return ""
}

func (gcm *GoCallMeta) From() string {
	if IsSelectorCallNode(gcm.callExpr) {
		return gcm.callExpr.Fun.(*ast.SelectorExpr).X.(*ast.Ident).Name
	}
	return ""
}

func (gcm *GoCallMeta) Args() []*GoArgMeta {
	if gcm.callExpr == nil || len(gcm.callExpr.Args) == 0 {
		return nil
	}
	args := make([]*GoArgMeta, 0, len(gcm.callExpr.Args))
	for _, argExpr := range gcm.callExpr.Args {
		switch argExpr := argExpr.(type) {
		case *ast.BasicLit:
			switch argExpr.Kind {
			case token.INT:
				argValue, err := strconv.ParseInt(argExpr.Value, 10, 32)
				if err != nil {
					panic(err)
				}
				args = append(args, &GoArgMeta{
					node:  argExpr,
					arg:   argExpr.Value,
					value: argValue,
				})
			case token.STRING:
				args = append(args, &GoArgMeta{
					node:  argExpr,
					arg:   argExpr.Value,
					value: argExpr.Value,
				})
			}
		case *ast.Ident:
			args = append(args, &GoArgMeta{
				node: argExpr,
				arg:  argExpr.Name,
			})
		case *ast.SelectorExpr:
			args = append(args, &GoArgMeta{
				node: argExpr,
				from: argExpr.X.(*ast.Ident).Name,
				arg:  argExpr.Sel.Name,
			})
		case *ast.CallExpr:
			args = append(args, &GoArgMeta{
				node: argExpr,
				callMeta: &GoCallMeta{
					fileMeta: gcm.fileMeta,
					callExpr: argExpr,
				},
			})
		}
	}
	return args
}
