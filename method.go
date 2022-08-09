package extractor

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

type goMethodMeta struct {
	interfaceMethodDecl *ast.Field    // method declared in an interface
	structMethodDecl    *ast.FuncDecl // methods implemented in struct
}

func extractGoMethodMeta(extractFilepath, structInterfaceName, methodName string) (*goMethodMeta, error) {
	fileAST, err := parser.ParseFile(token.NewFileSet(), extractFilepath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	gmm := searchGoMethodMeta(fileAST, structInterfaceName, methodName)
	if gmm != nil && gmm.interfaceMethodDecl == nil && gmm.structMethodDecl == nil {
		return nil, fmt.Errorf("can not find struct/interface %v method %v decl", structInterfaceName, methodName)
	}

	return gmm, nil
}

// searchGoMethodMeta search method from node.(*ast.File)
// TODO: split search interface and method impl
func searchGoMethodMeta(fileAST *ast.File, structInterfaceName, methodName string) *goMethodMeta {
	var interfaceMethodDecl *ast.Field
	var structMethodDecl *ast.FuncDecl
	ast.Inspect(fileAST, func(n ast.Node) bool {
		if n == fileAST {
			return true
		}
		if n == nil || interfaceMethodDecl != nil || structMethodDecl != nil {
			return false
		}

		switch n.(type) {
		case *ast.FuncDecl:
			// struct method
			funcDecl := n.(*ast.FuncDecl)
			if funcDecl.Name.String() == methodName {
				if funcDecl.Recv != nil && len(funcDecl.Recv.List) > 0 {
					var recvTypeIdentNode ast.Node
					switch funcDecl.Recv.List[0].Type.(type) {
					case *ast.Ident:
						recvTypeIdentNode = funcDecl.Recv.List[0].Type
					case *ast.StarExpr:
						recvTypeIdentNode = funcDecl.Recv.List[0].Type.(*ast.StarExpr).X
					}
					recvTypeIdent, ok := recvTypeIdentNode.(*ast.Ident)
					if ok && recvTypeIdent.Name == structInterfaceName {
						structMethodDecl = funcDecl
						return false
					}
				}
				return false
			}
		case *ast.GenDecl: // interface top decl
			return true
		case *ast.TypeSpec:
			// interface method
			typeSpec := n.(*ast.TypeSpec)
			if typeSpec.Type == nil {
				return false
			}
			interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
			if !ok {
				return false
			}
			if typeSpec.Name.String() == structInterfaceName && interfaceType.Methods != nil && len(interfaceType.Methods.List) > 0 {
				for _, methodField := range interfaceType.Methods.List {
					if ok && methodField.Names[0].Name == methodName {
						interfaceMethodDecl = methodField
						return false
					}
				}
				return true
			}
		}

		return false
	})

	if interfaceMethodDecl == nil && structMethodDecl == nil {
		return nil
	}

	return &goMethodMeta{
		interfaceMethodDecl: interfaceMethodDecl,
		structMethodDecl:    structMethodDecl,
	}
}

func (gmm *goMethodMeta) PrintAST() {
	if gmm.interfaceMethodDecl != nil {
		ast.Print(token.NewFileSet(), gmm.interfaceMethodDecl)
	} else if gmm.structMethodDecl != nil {
		ast.Print(token.NewFileSet(), gmm.structMethodDecl)
	}
}

func (gmm *goMethodMeta) MethodName() string {
	if gmm.interfaceMethodDecl != nil {
		return gmm.structMethodDecl.Name.String()
	} else if gmm.structMethodDecl != nil {
		return gmm.structMethodDecl.Name.String()
	}
	return "<nil>"
}
