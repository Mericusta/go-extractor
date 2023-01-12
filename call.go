package extractor

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
)

type GoCallMeta struct {
	expression string
	callExpr   *ast.CallExpr
}

func ExtractGoCallMeta(expression string) *GoCallMeta {
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
		expression: expression,
		callExpr:   callExpr,
	}
}

func (gcm *GoCallMeta) PrintAST() {
	ast.Print(token.NewFileSet(), gcm.callExpr)
}

func (gcm *GoCallMeta) Expression() string {
	return gcm.expression
}

func (gcm *GoCallMeta) Call() string {
	if gcm.callExpr == nil || gcm.callExpr.Fun == nil {
		return ""
	}
	ident, ok := gcm.callExpr.Fun.(*ast.Ident)
	if !ok {
		return ""
	}
	return ident.Name
}

func (gcm *GoCallMeta) Args() []interface{} {
	if gcm.callExpr == nil || len(gcm.callExpr.Args) == 0 {
		return nil
	}
	args := make([]interface{}, 0, len(gcm.callExpr.Args))
	for _, argExpr := range gcm.callExpr.Args {
		switch argExpr.(type) {
		case *ast.BasicLit:
			basicLit := argExpr.(*ast.BasicLit)
			switch basicLit.Kind {
			case token.INT:
				arg, err := strconv.ParseInt(basicLit.Value, 10, 32)
				if err != nil {
					panic(err)
				}
				args = append(args, int32(arg))
			case token.STRING:
				args = append(args, basicLit.Value)
			}
			// switch basicLit
		}
		// switch basicLit, ok := argExpr.(*ast.BasicLit)

		// args = append(args, arg)
	}
	return args
}
