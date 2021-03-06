package extractor

import (
	"io"
	"regexp"
	"strings"
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

func ExtractGoFileInterfaceDeclaration(r io.Reader) map[string]*GoInterfaceInfo {
	fileContent := CleanFileComment(r)

	fileInterfaceDeclarationMap := make(map[string]*GoInterfaceInfo)
	for _, interfaceDeclarationScopeBeginIndexSlice := range GoInterfaceDeclarationScopeBeginRegexp.FindAllStringIndex(string(fileContent), -1) {
		submatchSlice := GoInterfaceDeclarationScopeBeginRegexp.FindStringSubmatch(string(fileContent[interfaceDeclarationScopeBeginIndexSlice[0]:interfaceDeclarationScopeBeginIndexSlice[1]]))
		interfaceName := submatchSlice[GoInterfaceRegexpSubmatchNameIndex]
		// if interfaceName != "GoInterfaceDeclaration" {
		// 	continue
		// }
		fileInterfaceDeclarationMap[interfaceName] = &GoInterfaceInfo{
			Name:                   interfaceName,
			FunctionDeclarationMap: make(map[string]*GoFunctionDeclaration),
		}

		// {
		// 	fmt.Println()
		// 	fmt.Printf("interfaceDeclarationScopeBeginIndexSlice = |%v|\n", interfaceDeclarationScopeBeginIndexSlice)
		// 	fmt.Printf("interfaceDeclarationScope = |%v|\n", string(fileContent[interfaceDeclarationScopeBeginIndexSlice[0]:interfaceDeclarationScopeBeginIndexSlice[1]]))
		// 	fmt.Printf("interfaceName = %v\n", interfaceName)
		// 	fmt.Println()
		// 	return nil
		// }

		interfaceDeclarationScopeBeginRune := rune(fileContent[interfaceDeclarationScopeBeginIndexSlice[1]-1])
		interfaceDeclarationScopeEndRune := GetAnotherPunctuationMark(interfaceDeclarationScopeBeginRune)
		interfaceDeclarationLength := CalculatePunctuationMarksContentLength(
			string(fileContent[interfaceDeclarationScopeBeginIndexSlice[1]+1:]),
			interfaceDeclarationScopeBeginRune,
			interfaceDeclarationScopeEndRune,
			InvalidScopePunctuationMarkMap,
		)
		if interfaceDeclarationLength < 0 {
			continue
		}

		// {
		// 	fmt.Println()
		// 	fmt.Printf("interface content = |%v|", string(fileContent[interfaceDeclarationScopeBeginIndexSlice[1]:interfaceDeclarationScopeBeginIndexSlice[1]+interfaceDeclarationLength]))
		// 	fmt.Println()
		// 	return nil
		// }

		for _, lineContent := range strings.Split(string(fileContent[interfaceDeclarationScopeBeginIndexSlice[1]:interfaceDeclarationScopeBeginIndexSlice[1]+interfaceDeclarationLength]), "\n") {
			trimSpaceString := strings.TrimSpace(lineContent)
			if len(trimSpaceString) == 0 {
				continue
			}

			submatchSlice := GoInterfaceFunctionDeclarationScopeBeginRegexp.FindStringSubmatch(trimSpaceString)
			if len(submatchSlice) == 0 {
				continue
			}
			// fmt.Printf("trimSpaceString = |%v|\n", trimSpaceString)

			// signature
			functionName := submatchSlice[GoInterfaceFunctionDeclarationScopeBeginRegexpSubmatchNameIndex]

			// // params
			// // fmt.Printf("submatchSlice = %v\n", submatchSlice)
			// // fmt.Printf("submatchSlice[%v] = %v\n", GoStructRegexpSubmatchNameIndex, submatchSlice[GoStructRegexpSubmatchNameIndex])
			// // fmt.Printf("submatchSlice[%v] = %v\n", GoFunctionDeclarationScopeBeginRegexpSubmatchParamsScopeBeginIndex, submatchSlice[GoFunctionDeclarationScopeBeginRegexpSubmatchParamsScopeBeginIndex])
			// submatchIndexSlice := GoFunctionDeclarationScopeBeginRegexp.FindStringSubmatchIndex(trimSpaceString)
			// // fmt.Printf("submatchIndexSlice = %v\n", submatchIndexSlice)
			// // fmt.Printf("submatchIndexSlice[%v] = %v, %v\n", GoFunctionDeclarationScopeBeginRegexpSubmatchParamsScopeBeginIndex*2, submatchIndexSlice[GoFunctionDeclarationScopeBeginRegexpSubmatchParamsScopeBeginIndex*2], string(trimSpaceString[submatchIndexSlice[GoFunctionDeclarationScopeBeginRegexpSubmatchParamsScopeBeginIndex*2]]))
			// paramsScopeBeginRune := rune(trimSpaceString[submatchIndexSlice[GoFunctionDeclarationScopeBeginRegexpSubmatchParamsScopeBeginIndex*2]])
			// paramsScopeEndRune := GetAnotherPunctuationMark(paramsScopeBeginRune)
			// paramsScopeBeginIndex := submatchIndexSlice[GoFunctionDeclarationScopeBeginRegexpSubmatchParamsScopeBeginIndex*2+1]
			// paramsScopeLength := CalculatePunctuationMarksContentLength(
			// 	trimSpaceString[submatchIndexSlice[GoFunctionDeclarationScopeBeginRegexpSubmatchParamsScopeBeginIndex*2+1]:],
			// 	paramsScopeBeginRune,
			// 	paramsScopeEndRune,
			// 	InvalidScopePunctuationMarkMap,
			// )
			// paramsScopeEndIndex := paramsScopeBeginIndex + paramsScopeLength + 1
			// var paramsScopeContent string
			// if paramsScopeLength > 0 {
			// 	paramsScopeContent = trimSpaceString[paramsScopeBeginIndex:paramsScopeEndIndex]
			// }
			// fmt.Printf("paramsScopeContent = |%v|\n", paramsScopeContent)

			// // returns
			// var returnsScopeContent string
			// if paramsScopeEndIndex+1 < len(trimSpaceString) && len(trimSpaceString[paramsScopeEndIndex+1:]) > 0 {
			// 	returnsScopeContent = strings.TrimSpace(trimSpaceString[paramsScopeEndIndex+1:])
			// }
			// fmt.Printf("returnsScopeContent = |%v|\n", returnsScopeContent)

			fileInterfaceDeclarationMap[interfaceName].FunctionDeclarationMap[functionName] = &GoFunctionDeclaration{
				FunctionSignature: functionName,
				// TypeDeclaration:   ExtractGoVariableTypeDeclaration(submatchSlice[GoVariableDeclarationRegexpSubmatchTypeIndex]),
			}
		}
	}

	return fileInterfaceDeclarationMap
}
