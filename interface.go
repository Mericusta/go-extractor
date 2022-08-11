package extractor

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
)

type GoInterfaceInfo struct {
	Name                   string
	FunctionDeclarationMap map[string]*GoFunctionDeclaration
}

var (
	GO_INTERFACE_DECLARATION_SCOPE_BEGIN_EXPRESSION                             string = `type\s+(?P<NAME>\w+)\s+interface\s+\{`
	GoInterfaceDeclarationScopeBeginRegexp                                             = regexp.MustCompile(GO_INTERFACE_DECLARATION_SCOPE_BEGIN_EXPRESSION)
	GoInterfaceRegexpSubmatchNameIndex                                                 = GoInterfaceDeclarationScopeBeginRegexp.SubexpIndex("NAME")
	GoInterfaceDeclarationScopeBeginRune                                               = '{'
	GO_INTERFACE_FUNCTION_DECLARATION_SCOPE_BEGIN_EXPRESSION                    string = `(?P<NAME>\w+)\s*(?P<PARAMS_SCOPE_BEGIN>\()`
	GoInterfaceFunctionDeclarationScopeBeginRegexp                                     = regexp.MustCompile(GO_INTERFACE_FUNCTION_DECLARATION_SCOPE_BEGIN_EXPRESSION)
	GoInterfaceFunctionDeclarationScopeBeginRegexpSubmatchNameIndex                    = GoInterfaceFunctionDeclarationScopeBeginRegexp.SubexpIndex("NAME")
	GoInterfaceFunctionDeclarationScopeBeginRegexpSubmatchParamsScopeBeginIndex        = GoInterfaceFunctionDeclarationScopeBeginRegexp.SubexpIndex("PARAMS_SCOPE_BEGIN")
)

func ExtractGoFileInterfaceDeclaration(extractFilepath string, parseComments bool) (map[string]*GoInterfaceInfo, error) {
	return nil, nil
}

type goInterfaceMeta struct {
	typeSpec   *ast.TypeSpec
	methodMeta map[string]*goInterfaceMethodMeta
}

type goInterfaceMethodMeta struct {
	methodField *ast.Field
}

func extractGoInterfaceMeta(extractFilepath string, interfaceName string) (*goInterfaceMeta, error) {
	fileAST, err := parser.ParseFile(token.NewFileSet(), extractFilepath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	gim := searchGoInterfaceMeta(fileAST, interfaceName)
	if gim.typeSpec == nil {
		return nil, fmt.Errorf("can not find interface decl")
	}

	return gim, nil
}

func searchGoInterfaceMeta(fileAST *ast.File, interfaceName string) *goInterfaceMeta {
	var interfaceDecl *ast.TypeSpec
	ast.Inspect(fileAST, func(n ast.Node) bool {
		if n == fileAST {
			return true
		}
		if n == nil || interfaceDecl != nil {
			return false
		}
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}
		if typeSpec.Type == nil {
			return false
		}
		_, ok = typeSpec.Type.(*ast.InterfaceType)
		if !ok {
			return true
		}
		if typeSpec.Name.String() == interfaceName {
			interfaceDecl = typeSpec
			return false
		}
		return true
	})
	return &goInterfaceMeta{
		typeSpec: interfaceDecl,
	}
}

func (gim *goInterfaceMeta) PrintAST() {
	ast.Print(token.NewFileSet(), gim.typeSpec)
}

func (gim *goInterfaceMeta) InterfaceName() string {
	return gim.typeSpec.Name.String()
}

// SearchMethodDecl search method decl from node.(*ast.InterfaceType)
func (gim *goInterfaceMeta) SearchMethodDecl(methodName string) *goInterfaceMethodMeta {
	gim.ForeachMethodDecl(func(f *ast.Field) bool {
		if f.Names[0].Name == methodName {
			gim.methodMeta[methodName] = &goInterfaceMethodMeta{methodField: f}
			return false
		}
		return true
	})
	return gim.methodMeta[methodName]
}

func (gim *goInterfaceMeta) ForeachMethodDecl(f func(*ast.Field) bool) {
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
