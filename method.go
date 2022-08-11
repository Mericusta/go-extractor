package extractor

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

type goMethodMeta struct {
	methodDecl *ast.FuncDecl
}

func extractGoMethodMeta(extractFilepath string, structName, methodName string) (*goMethodMeta, error) {
	fileAST, err := parser.ParseFile(token.NewFileSet(), extractFilepath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	gmm := searchGoMethodMeta(fileAST, structName, methodName)
	if gmm.methodDecl == nil {
		return nil, fmt.Errorf("can not find method decl")
	}

	return gmm, nil
}

func searchGoMethodMeta(fileAST *ast.File, structName, methodName string) *goMethodMeta {
	var methodDecl *ast.FuncDecl
	ast.Inspect(fileAST, func(n ast.Node) bool {
		if n == fileAST {
			return true
		}
		if n == nil || methodDecl != nil {
			return false
		}
		decl, ok := n.(*ast.FuncDecl)
		if !ok {
			return false
		}
		if decl.Name.String() == methodName {
			if decl.Recv != nil && len(decl.Recv.List) > 0 {
				var recvTypeIdentNode ast.Node
				switch decl.Recv.List[0].Type.(type) {
				case *ast.Ident:
					recvTypeIdentNode = decl.Recv.List[0].Type
				case *ast.StarExpr:
					recvTypeIdentNode = decl.Recv.List[0].Type.(*ast.StarExpr).X
				}
				recvTypeIdent, ok := recvTypeIdentNode.(*ast.Ident)
				if ok && recvTypeIdent.Name == structName {
					methodDecl = decl
					return false
				}
			}
			return false
		}
		return true
	})
	return &goMethodMeta{
		methodDecl: methodDecl,
	}
}

func (gmm *goMethodMeta) MethodName() string {
	return gmm.methodDecl.Name.String()
}
