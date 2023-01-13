package extractor

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

type GoInterfaceMeta struct {
	typeSpec     *ast.TypeSpec
	commentGroup *ast.CommentGroup
	methodMeta   map[string]*GoInterfaceMethodMeta
}

type GoInterfaceMethodMeta struct {
	methodField *ast.Field
}

func ExtractGoInterfaceMeta(extractFilepath string, interfaceName string) (*GoInterfaceMeta, error) {
	fileAST, err := parser.ParseFile(token.NewFileSet(), extractFilepath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	gim := SearchGoInterfaceMeta(fileAST, interfaceName)
	if gim.typeSpec == nil {
		return nil, fmt.Errorf("can not find interface decl")
	}

	return gim, nil
}

func SearchGoInterfaceMeta(fileAST *ast.File, interfaceName string) *GoInterfaceMeta {
	var interfaceDecl *ast.TypeSpec
	var commentDecl *ast.CommentGroup
	ast.Inspect(fileAST, func(n ast.Node) bool {
		if genDecl, ok := n.(*ast.GenDecl); ok {
			ast.Inspect(genDecl, func(n ast.Node) bool {
				if IsInterfaceNode(n) {
					typeSpec := n.(*ast.TypeSpec)
					if typeSpec.Name.String() == interfaceName {
						interfaceDecl = typeSpec
						commentDecl = genDecl.Doc
						return false
					}
				}
				return true
			})
			return false // genDecl traverse done
		}
		return interfaceDecl == nil // already found
	})
	if interfaceDecl == nil {
		return nil
	}
	return &GoInterfaceMeta{
		typeSpec:     interfaceDecl,
		commentGroup: commentDecl,
		methodMeta:   make(map[string]*GoInterfaceMethodMeta),
	}
}

func IsInterfaceNode(n ast.Node) bool {
	typeSpec, ok := n.(*ast.TypeSpec)
	if !ok {
		return false
	}
	if typeSpec.Type == nil {
		return false
	}
	_, ok = typeSpec.Type.(*ast.InterfaceType)
	return ok
}

func (gim *GoInterfaceMeta) PrintAST() {
	ast.Print(token.NewFileSet(), gim.typeSpec)
}

func (gim *GoInterfaceMeta) InterfaceName() string {
	return gim.typeSpec.Name.String()
}

// SearchMethodDecl search method decl from node.(*ast.InterfaceType)
func (gim *GoInterfaceMeta) SearchMethodDecl(methodName string) *GoInterfaceMethodMeta {
	gim.ForeachMethodDecl(func(f *ast.Field) bool {
		if f.Names[0].Name == methodName {
			gim.methodMeta[methodName] = &GoInterfaceMethodMeta{methodField: f}
			return false
		}
		return true
	})
	return gim.methodMeta[methodName]
}

func (gim *GoInterfaceMeta) ForeachMethodDecl(f func(*ast.Field) bool) {
	interfaceType := gim.typeSpec.Type.(*ast.InterfaceType)
	if interfaceType.Methods == nil {
		return
	}
	for _, methodField := range interfaceType.Methods.List {
		_, ok := methodField.Type.(*ast.FuncType)
		if ok {
			if !f(methodField) {
				break
			}
		}
	}
}
