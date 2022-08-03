package extractor

import (
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

func ExtractGoFileInterfaceDeclaration(extractFilePath string, parseComments bool) (map[string]*GoInterfaceInfo, error) {

	// for objectName, object := range fileAST.Scope.Objects {

	// }

	// return fileInterfaceDeclarationMap
	return nil, nil
}
