package extractor

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

type GoMethodMeta struct {
	methodDecl *ast.FuncDecl
}

func ExtractGoMethodMeta(extractFilepath string, structName, methodName string) (*GoMethodMeta, error) {
	fileAST, err := parser.ParseFile(token.NewFileSet(), extractFilepath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	gmm := SearchGoMethodMeta(fileAST, structName, methodName)
	if gmm.methodDecl == nil {
		return nil, fmt.Errorf("can not find method decl")
	}

	return gmm, nil
}

func SearchGoMethodMeta(fileAST *ast.File, structName, methodName string) *GoMethodMeta {
	var methodDecl *ast.FuncDecl
	ast.Inspect(fileAST, func(n ast.Node) bool {
		if IsMethodNode(n) {
			decl := n.(*ast.FuncDecl)
			if decl.Name.String() == methodName {
				recvStructName, _ := extractMethodRecvStruct(decl)
				if recvStructName == structName {
					methodDecl = decl
				}
			}
		}
		return methodDecl == nil
	})
	if methodDecl == nil {
		return nil
	}
	return &GoMethodMeta{
		methodDecl: methodDecl,
	}
}

func IsMethodNode(n ast.Node) bool {
	decl, ok := n.(*ast.FuncDecl)
	return ok && decl.Recv != nil && len(decl.Recv.List) > 0
}

func (gmm *GoMethodMeta) MethodName() string {
	return gmm.methodDecl.Name.String()
}

func (gmm *GoMethodMeta) RecvStruct() (string, bool) {
	return extractMethodRecvStruct(gmm.methodDecl)
}

func extractMethodRecvStruct(methodDecl *ast.FuncDecl) (string, bool) {
	if len(methodDecl.Recv.List) < 1 {
		return "", false
	}

	var pointerReceiver bool
	var recvTypeIdentNode ast.Node
	switch methodDecl.Recv.List[0].Type.(type) {
	case *ast.Ident:
		pointerReceiver = false
		recvTypeIdentNode = methodDecl.Recv.List[0].Type
	case *ast.StarExpr:
		pointerReceiver = true
		recvTypeIdentNode = methodDecl.Recv.List[0].Type.(*ast.StarExpr).X
	}

	recvTypeIdent, ok := recvTypeIdentNode.(*ast.Ident)
	if !ok {
		return "", false
	}
	return recvTypeIdent.Name, pointerReceiver
}
