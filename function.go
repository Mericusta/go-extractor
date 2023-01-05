package extractor

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"strings"
)

type GoFunctionDeclaration struct {
	Content           []byte
	FunctionSignature string
	This              *GoVariableDefinition
	ParamsList        []*GoVariableDefinition
	ReturnList        []*GoTypeDeclaration // not support named return
	BodyContent       []byte
}

func (d *GoFunctionDeclaration) Traversal(deep int) {
	fmt.Printf("%v- Function Signature: %v\n", strings.Repeat("\t", deep), d.FunctionSignature)
	fmt.Printf("%v- Function Param List: ", strings.Repeat("\t", deep))
	if len(d.ParamsList) > 0 {
		fmt.Println()
		for index, paramDeclaration := range d.ParamsList {
			fmt.Printf("%v- Param.%v: %v %v\n", strings.Repeat("\t", deep+1), index, paramDeclaration.VariableSignature, paramDeclaration.TypeDeclaration.MakeUp())
		}
	} else {
		fmt.Printf("None\n")
	}
	fmt.Printf("%v- Function Return List: ", strings.Repeat("\t", deep))
	if len(d.ReturnList) > 0 {
		fmt.Println()
		for index, returnDeclaration := range d.ReturnList {
			fmt.Printf("%v- Return.%v: %v\n", strings.Repeat("\t", deep+1), index, returnDeclaration.MakeUp())
		}
	} else {
		fmt.Printf("None\n")
	}
	fmt.Printf("%v- MakeUp: |%v|\n", strings.Repeat("\t", deep), d.MakeUp())
}

func (d *GoFunctionDeclaration) MakeUp() string {
	makeUpTemplate := `func[THIS_SCOPE] [SIGNATURE][PARAM_SCOPE][RETURN_SCOPE] [BODY_SCOPE]`
	thisScopeTemplate := `(THIS_VARIABLE)`
	paramScopeTemplate := `(EACH_PARAM)`
	singleReturnScopeTemplate := `EACH_RETURN`
	multiReturnScopeTemplate := `(EACH_RETURN)`
	bodyScopeTemplate := `{BODY_CONTENT}`

	makeUpReplaceKeywordThisScope := `[THIS_SCOPE]`
	makeUpReplaceKeywordSignature := `[SIGNATURE]`
	makeUpReplaceKeywordParamScope := `[PARAM_SCOPE]`
	makeUpReplaceKeywordReturnScope := `[RETURN_SCOPE]`
	makeUpReplaceKeywordBodyScope := `[BODY_SCOPE]`
	makeUpReplaceKeywordThisVariable := `THIS_VARIABLE`
	makeUpReplaceKeywordEachParam := `EACH_PARAM`
	makeUpReplaceKeywordEachReturn := `EACH_RETURN`
	makeUpReplaceKeywordBodyContent := `BODY_CONTENT`

	// signature
	makeUpContent := strings.Replace(makeUpTemplate, makeUpReplaceKeywordSignature, d.FunctionSignature, -1)

	// this scope
	var thisScopeContent string
	if d.This != nil {
		thisScopeContent = strings.Replace(thisScopeTemplate, makeUpReplaceKeywordThisVariable, d.This.MakeUp(), -1)
	}
	makeUpContent = strings.Replace(makeUpContent, makeUpReplaceKeywordThisScope, thisScopeContent, -1)

	// param scope
	builder := strings.Builder{}
	for index, eachParam := range d.ParamsList {
		builder.WriteString(eachParam.MakeUp())
		if index > 0 {
			builder.WriteRune(',')
			builder.WriteRune(' ')
		}
	}
	paramScopeContent := strings.Replace(paramScopeTemplate, makeUpReplaceKeywordEachParam, builder.String(), -1)
	makeUpContent = strings.Replace(makeUpContent, makeUpReplaceKeywordParamScope, paramScopeContent, -1)

	// return scope
	builder.Reset()
	multiReturn := false
	for index, eachReturn := range d.ReturnList {
		builder.WriteString(eachReturn.MakeUp())
		if index > 0 {
			builder.WriteRune(',')
			builder.WriteRune(' ')
			multiReturn = true
		}
	}
	var returnScopeContent string
	if multiReturn {
		returnScopeContent = strings.Replace(multiReturnScopeTemplate, makeUpReplaceKeywordEachReturn, builder.String(), -1)
	} else {
		returnScopeContent = strings.Replace(singleReturnScopeTemplate, makeUpReplaceKeywordEachReturn, builder.String(), -1)
	}
	if len(returnScopeContent) > 0 {
		returnScopeContent = fmt.Sprintf(" %v", returnScopeContent)
	}
	makeUpContent = strings.Replace(makeUpContent, makeUpReplaceKeywordReturnScope, returnScopeContent, -1)

	// body scope
	bodyContent := strings.Replace(bodyScopeTemplate, makeUpReplaceKeywordBodyContent, string(d.BodyContent), -1)
	makeUpContent = strings.Replace(makeUpContent, makeUpReplaceKeywordBodyScope, bodyContent, -1)

	// add new line
	makeUpContent = "\n" + makeUpContent + "\n"

	// remove all CR
	makeUpContent = strings.ReplaceAll(makeUpContent, "\r", "")

	return makeUpContent
}

var (
	GO_FUNCTION_DECLARATION_SCOPE_BEGIN_EXPRESSION                     string = `\nfunc\s+(\((?P<THIS>\w+)\s+(?P<THIS_TYPE>(\*)?\w+)\))?\s*(?P<NAME>\w+)\s*(?P<PARAMS_SCOPE_BEGIN>\()`
	GoFunctionDeclarationScopeBeginRegexp                                     = regexp.MustCompile(GO_FUNCTION_DECLARATION_SCOPE_BEGIN_EXPRESSION)
	GoFunctionDeclarationScopeBeginRegexpSubmatchThisIndex                    = GoFunctionDeclarationScopeBeginRegexp.SubexpIndex("THIS")
	GoFunctionDeclarationScopeBeginRegexpSubmatchThisTypeIndex                = GoFunctionDeclarationScopeBeginRegexp.SubexpIndex("THIS_TYPE")
	GoFunctionDeclarationScopeBeginRegexpSubmatchNameIndex                    = GoFunctionDeclarationScopeBeginRegexp.SubexpIndex("NAME")
	GoFunctionDeclarationScopeBeginRegexpSubmatchParamsScopeBeginIndex        = GoFunctionDeclarationScopeBeginRegexp.SubexpIndex("PARAMS_SCOPE_BEGIN")
)

func ExtractGoFileFunctionDeclaration(content []byte) map[string]*GoFunctionDeclaration {
	if len(content) == 0 {
		return nil
	}

	if GoFunctionDeclarationScopeBeginRegexpSubmatchThisIndex == -1 || GoFunctionDeclarationScopeBeginRegexpSubmatchThisTypeIndex == -1 || GoFunctionDeclarationScopeBeginRegexpSubmatchNameIndex == -1 || GoFunctionDeclarationScopeBeginRegexpSubmatchParamsScopeBeginIndex == -1 {
		panic("sub match index is -1")
	}

	functionDeclarationMap := make(map[string]*GoFunctionDeclaration)
	for _, functionDeclarationScopeBeginSubmatchIndexSlice := range GoFunctionDeclarationScopeBeginRegexp.FindAllSubmatchIndex(content, -1) {
		contentAfterBeginScopeLength := 0

		// signature
		functionName := strings.TrimSpace(string(content[functionDeclarationScopeBeginSubmatchIndexSlice[GoFunctionDeclarationScopeBeginRegexpSubmatchNameIndex*2]:functionDeclarationScopeBeginSubmatchIndexSlice[GoFunctionDeclarationScopeBeginRegexpSubmatchNameIndex*2+1]]))

		// this scope
		var thisDeclaration *GoVariableDefinition
		if functionDeclarationScopeBeginSubmatchIndexSlice[GoFunctionDeclarationScopeBeginRegexpSubmatchThisIndex*2] != -1 &&
			functionDeclarationScopeBeginSubmatchIndexSlice[GoFunctionDeclarationScopeBeginRegexpSubmatchThisIndex*2+1] != -1 &&
			functionDeclarationScopeBeginSubmatchIndexSlice[GoFunctionDeclarationScopeBeginRegexpSubmatchThisTypeIndex*2] != -1 &&
			functionDeclarationScopeBeginSubmatchIndexSlice[GoFunctionDeclarationScopeBeginRegexpSubmatchThisTypeIndex*2+1] != -1 {
			thisSignature := strings.TrimSpace(string(content[functionDeclarationScopeBeginSubmatchIndexSlice[GoFunctionDeclarationScopeBeginRegexpSubmatchThisIndex*2]:functionDeclarationScopeBeginSubmatchIndexSlice[GoFunctionDeclarationScopeBeginRegexpSubmatchThisIndex*2+1]]))
			thisTypeContent := strings.TrimSpace(string(content[functionDeclarationScopeBeginSubmatchIndexSlice[GoFunctionDeclarationScopeBeginRegexpSubmatchThisTypeIndex*2]:functionDeclarationScopeBeginSubmatchIndexSlice[GoFunctionDeclarationScopeBeginRegexpSubmatchThisTypeIndex*2+1]]))
			thisDeclaration = &GoVariableDefinition{
				VariableSignature: thisSignature,
				TypeDeclaration:   ExtractGoVariableTypeDeclaration(thisTypeContent),
			}
		}

		// params scope
		paramsScopeBeginRuneIndex := functionDeclarationScopeBeginSubmatchIndexSlice[GoFunctionDeclarationScopeBeginRegexpSubmatchParamsScopeBeginIndex*2] // '(' index
		paramsScopeBeginRune := rune(content[paramsScopeBeginRuneIndex])                                                                                   // '('
		paramsScopeEndRune := GetAnotherPunctuationMark(paramsScopeBeginRune)                                                                              // ')'
		paramsScopeLength := CalculatePunctuationMarksContentLength(
			string(content[paramsScopeBeginRuneIndex+1:]),
			paramsScopeBeginRune, paramsScopeEndRune, InvalidScopePunctuationMarkMap,
		)
		paramsScopeEndRuneIndex := paramsScopeBeginRuneIndex + paramsScopeLength + 1        // ')' index
		paramsListContent := content[paramsScopeBeginRuneIndex+1 : paramsScopeEndRuneIndex] // between '(' and ')'
		paramsList := ExtractorFunctionParamsList(paramsListContent)
		contentAfterBeginScopeLength += paramsScopeLength + 1

		// returns scope
		returnsScopeBeginRuneIndex := paramsScopeEndRuneIndex + 1 // after params scope end rune index
		returnsScopeEndRuneIndex := returnsScopeBeginRuneIndex    // before body scope begin rune index
		bodyScopeBeginRuneIndex := -1                             // '{' index
		keywordStack := []string{Keyword_func}
		word := make([]byte, 0, 16)
		for rIndex := 0; rIndex != len(content[paramsScopeEndRuneIndex:]); rIndex++ {
			contentIndex := paramsScopeEndRuneIndex + rIndex
			r := rune(content[contentIndex])
			switch {
			case IsCharacter(r):
				word = append(word, byte(r))
			default:
				if len(word) > 0 {
					if IsGolangScopeKeyword(string(word)) {
						keywordStack = append(keywordStack, string(word))
					}
					word = make([]byte, 0, 16)
				}
				stackLength := len(keywordStack)
				switch {
				case IsSpaceRune(r):
				case r == '(':
					scopeLength := CalculatePunctuationMarksContentLength(
						string(content[contentIndex+1:]),
						'(', ')', InvalidScopePunctuationMarkMap,
					)
					if stackLength == 1 && keywordStack[0] == Keyword_func {
						returnsScopeBeginRuneIndex = contentIndex
					} else {
						if keywordStack[stackLength-1] == Keyword_func {
							if stackLength-2 >= 0 {
								keywordStack = keywordStack[0 : stackLength-1]
							} else {
								panic("stack length error")
							}
						}
					}
					// func (int) -> func (int)
					//      |                 |
					rIndex += (1 + scopeLength)
				case r == '[':
					scopeLength := CalculatePunctuationMarksContentLength(
						string(content[contentIndex+1:]),
						'[', ']', InvalidScopePunctuationMarkMap,
					)
					// map[int]int -> map[int]int
					//    |                  |
					rIndex += (1 + scopeLength)
				case r == '{':
					if stackLength == 1 && keywordStack[0] == Keyword_func {
						if content[returnsScopeBeginRuneIndex] == '(' {
							returnsScopeEndRuneIndex = contentIndex - 1
							for content[returnsScopeEndRuneIndex] != ')' {
								returnsScopeEndRuneIndex--
							}
						} else {
							returnsScopeEndRuneIndex = contentIndex
						}
						bodyScopeBeginRuneIndex = contentIndex
						goto SEARCH_END
					} else {
						scopeLength := CalculatePunctuationMarksContentLength(
							string(content[contentIndex+1:]),
							'{', '}', InvalidScopePunctuationMarkMap,
						)
						if keywordStack[stackLength-1] == Keyword_func {
							panic("syntax error")
						}
						// interface{} -> interface{}
						//          |               |
						// struct{ v interface{} } -> struct{ v interface{} }
						//       |                                          |
						rIndex += (1 + scopeLength)
						if stackLength-2 >= 0 {
							keywordStack = keywordStack[0 : stackLength-1]
						} else {
							panic("stack length error")
						}
					}
				}
			}
		}
	SEARCH_END:
		returnList := ExtractorFunctionReturnList(content[returnsScopeBeginRuneIndex+1 : returnsScopeEndRuneIndex])
		contentAfterBeginScopeLength += returnsScopeEndRuneIndex - returnsScopeBeginRuneIndex + 1

		// body scope
		if bodyScopeBeginRuneIndex < 0 {
			panic("body scope begin rune index is -1")
		}
		bodyLength := CalculatePunctuationMarksContentLength(
			string(content[bodyScopeBeginRuneIndex+1:]),
			'{', '}', InvalidScopePunctuationMarkMap,
		)
		if bodyLength < 0 {
			panic("function body length is -1")
		}
		bodyScopeEndRuneIndex := bodyScopeBeginRuneIndex + 1 + bodyLength
		contentAfterBeginScopeLength += bodyLength + 1

		functionDeclarationMap[functionName] = &GoFunctionDeclaration{
			Content:           content[functionDeclarationScopeBeginSubmatchIndexSlice[0] : functionDeclarationScopeBeginSubmatchIndexSlice[1]+1+contentAfterBeginScopeLength],
			FunctionSignature: functionName,
			This:              thisDeclaration,
			ParamsList:        paramsList,
			ReturnList:        returnList,
			BodyContent:       content[bodyScopeBeginRuneIndex+1 : bodyScopeEndRuneIndex],
		}
	}

	return functionDeclarationMap
}

func ExtractorFunctionParamsList(content []byte) []*GoVariableDefinition {
	splitContent := RecursiveSplitUnderSameDeepPunctuationMarksContent(string(content), GetLeftPunctuationMarkList(), ",")
	var sameTypeParamSlice, paramsSlice []*GoVariableDefinition
	for _, content := range splitContent.ContentList {
		// fmt.Printf("param content = |%v|\n", strings.TrimSpace(content))
		if len(content) == 0 {
			panic("param content is empty")
		}
		paramDeclaration := &GoVariableDefinition{}
		paramContentSubmatchSlice := GoVariableDeclarationRegexp.FindStringSubmatch(strings.TrimSpace(content))
		if len(paramContentSubmatchSlice) == 0 {
			paramDeclaration.VariableSignature = strings.TrimSpace(content)
			sameTypeParamSlice = append(sameTypeParamSlice, paramDeclaration)
		} else {
			paramDeclaration.VariableSignature = paramContentSubmatchSlice[GoVariableDeclarationRegexpSubmatchNameIndex]
			paramDeclaration.TypeDeclaration = ExtractGoVariableTypeDeclaration(paramContentSubmatchSlice[GoVariableDeclarationRegexpSubmatchTypeIndex])
			for _, sameTypeParam := range sameTypeParamSlice {
				sameTypeParam.TypeDeclaration = paramDeclaration.TypeDeclaration
			}
			sameTypeParamSlice = nil
		}
		paramsSlice = append(paramsSlice, paramDeclaration)
	}

	// {
	// 	for index, paramDeclaration := range paramsSlice {
	// 		fmt.Printf("%v param: %v\n", index, paramDeclaration.VariableSignature)
	// 		fmt.Printf("%v param type: %v\n", index, paramDeclaration.TypeDeclaration.MakeUp())
	// 	}
	// }
	return paramsSlice
}

func ExtractorFunctionReturnList(content []byte) []*GoTypeDeclaration {
	contentWithoutSpace := strings.TrimSpace(string(content))
	if len(contentWithoutSpace) == 0 {
		return nil
	}
	splitContent := RecursiveSplitUnderSameDeepPunctuationMarksContent(string(content), GetLeftPunctuationMarkList(), ",")
	returnTypeDeclarationSlice := make([]*GoTypeDeclaration, 0)
	for _, content := range splitContent.ContentList {
		if len(content) == 0 {
			panic("return content is empty")
		}
		typeDeclaration := ExtractGoVariableTypeDeclaration(strings.TrimSpace(content))
		returnTypeDeclarationSlice = append(returnTypeDeclarationSlice, typeDeclaration)
	}
	return returnTypeDeclarationSlice
}

// ----------------------------------------------------------------

type goFunctionMeta struct {
	funcDecl *ast.FuncDecl
}

func extractGoFunctionMeta(extractFilepath string, functionName string) (*goFunctionMeta, error) {
	fileAST, err := parser.ParseFile(token.NewFileSet(), extractFilepath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	gfm := searchGoFunctionMeta(fileAST, functionName)
	if gfm.funcDecl == nil {
		return nil, fmt.Errorf("can not find function decl")
	}

	return gfm, nil
}

func searchGoFunctionMeta(fileAST *ast.File, functionName string) *goFunctionMeta {
	var funcDecl *ast.FuncDecl
	ast.Inspect(fileAST, func(n ast.Node) bool {
		if n == fileAST {
			return true
		}
		if n == nil || funcDecl != nil {
			return false
		}
		decl, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}
		if decl.Recv == nil && decl.Name.String() == functionName {
			funcDecl = decl
			return false
		}
		return true
	})
	if funcDecl == nil {
		return nil
	}
	return &goFunctionMeta{
		funcDecl: funcDecl,
	}
}

func (gfm *goFunctionMeta) PrintAST() {
	ast.Print(token.NewFileSet(), gfm.funcDecl)
}

func (gfm *goFunctionMeta) FunctionName() string {
	return gfm.funcDecl.Name.String()
}

func (gfm *goFunctionMeta) IsMethod() bool {
	return gfm.funcDecl.Recv != nil
}

func (gfm *goFunctionMeta) RecvStruct() string {
	if !gfm.IsMethod() || len(gfm.funcDecl.Recv.List) < 1 {
		return ""
	}

	var recvTypeIdentNode ast.Node
	switch gfm.funcDecl.Recv.List[0].Type.(type) {
	case *ast.Ident:
		recvTypeIdentNode = gfm.funcDecl.Recv.List[0].Type
	case *ast.StarExpr:
		recvTypeIdentNode = gfm.funcDecl.Recv.List[0].Type.(*ast.StarExpr).X
	}

	recvTypeIdent, ok := recvTypeIdentNode.(*ast.Ident)
	if !ok {
		return ""
	}
	return recvTypeIdent.Name
}

func (gfm *goFunctionMeta) CallMap() map[string]*ast.CallExpr {
	ast.Print(token.NewFileSet(), gfm.funcDecl.Body)
	callMap := make(map[string]*ast.CallExpr)
	for _, e := range gfm.funcDecl.Body.List {
		exprStmt, ok := e.(*ast.ExprStmt)
		if exprStmt == nil || !ok || exprStmt.X == nil {
			continue
		}
		callExpr, ok := exprStmt.X.(*ast.CallExpr)
		if callExpr == nil || !ok {
			continue
		}

		ident, ok := callExpr.Fun.(*ast.Ident)
		if ident == nil || !ok {
			continue
		}
		callMap[ident.Name] = callExpr
	}
	return callMap
}

func (gfm *goFunctionMeta) Comments() []string {
	if gfm.funcDecl.Doc == nil {
		return nil
	}

	commentSlice := make([]string, 0, len(gfm.funcDecl.Doc.List))
	for _, comment := range gfm.funcDecl.Doc.List {
		commentSlice = append(commentSlice, comment.Text)
	}
	return commentSlice
}

func (gfm *goFunctionMeta) UpdateComments(comments []string) {
	if gfm.funcDecl.Doc == nil {
		return
	}

	if len(gfm.funcDecl.Doc.List) != len(comments) {
		return
	}

	for index, comment := range gfm.funcDecl.Doc.List {
		comment.Text = comments[index]
	}
}
