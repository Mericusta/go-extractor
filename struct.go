package extractor

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"regexp"
	"strings"
)

type GoStructInfo struct {
	Name                 string
	MemberDeclarationMap map[string]*GoVariableDefinition
}

var (
	CONSTRUCT_STRUCT_TEMPLATE       = `[RP_STRUCT_NAME]{}`
	REPLACE_KEYWORD_STRUCT_KEY_TYPE = "[RP_STRUCT_NAME]"
)

// Construct 生成该结构体的空构造方法
func (g *GoStructInfo) Construct() string {
	return strings.Replace(CONSTRUCT_STRUCT_TEMPLATE, REPLACE_KEYWORD_STRUCT_KEY_TYPE, g.Name, -1)
}

var (
	GO_STRUCT_DECLARATION_SCOPE_BEGIN_EXPRESSION string = `type\s+(?P<NAME>\w+)\s+struct\s+\{`
	GoStructDeclarationScopeBeginRegexp                 = regexp.MustCompile(GO_STRUCT_DECLARATION_SCOPE_BEGIN_EXPRESSION)
	GoStructRegexpSubmatchNameIndex                     = GoStructDeclarationScopeBeginRegexp.SubexpIndex("NAME")
	GoStructDeclarationScopeBeginRune                   = '{'
)

// ExtractGoFileStructDeclaration 从 go 文件中提取结构体声明，不支持匿名结构体
func ExtractGoFileStructDeclaration(r io.Reader) map[string]*GoStructInfo {
	fileContent := CleanFileComment(r)

	// fmt.Println()
	// fmt.Printf("fileContent = %v", string(fileContent))
	// fmt.Println()
	// return

	fileStructDeclarationMap := make(map[string]*GoStructInfo)
	for _, structDeclarationScopeBeginIndexSlice := range GoStructDeclarationScopeBeginRegexp.FindAllStringIndex(string(fileContent), -1) {
		submatchSlice := GoStructDeclarationScopeBeginRegexp.FindStringSubmatch(string(fileContent[structDeclarationScopeBeginIndexSlice[0]:structDeclarationScopeBeginIndexSlice[1]]))
		structName := submatchSlice[GoStructRegexpSubmatchNameIndex]
		fileStructDeclarationMap[structName] = &GoStructInfo{
			Name:                 structName,
			MemberDeclarationMap: make(map[string]*GoVariableDefinition),
		}

		// fmt.Println()
		// fmt.Printf("structDeclarationScopeBeginIndexSlice = |%v|\n", structDeclarationScopeBeginIndexSlice)
		// fmt.Printf("structDeclarationScope = |%v|\n", string(fileContent[structDeclarationScopeBeginIndexSlice[0]:structDeclarationScopeBeginIndexSlice[1]]))
		// fmt.Printf("structName = %v\n", structName)
		// fmt.Println()
		// return

		structDeclarationScopeBeginRune := rune(fileContent[structDeclarationScopeBeginIndexSlice[1]-1])
		structDeclarationScopeEndRune := GetAnotherPunctuationMark(structDeclarationScopeBeginRune)
		structDeclarationLength := CalculatePunctuationMarksContentLength(
			string(fileContent[structDeclarationScopeBeginIndexSlice[1]+1:]),
			structDeclarationScopeBeginRune,
			structDeclarationScopeEndRune,
			InvalidScopePunctuationMarkMap,
		)
		if structDeclarationLength < 0 {
			fmt.Printf("Error: struct %v can not find struct end index\n", structName)
			continue
		}

		// fmt.Println()
		// fmt.Printf("struct content = |%v|", string(fileContent[structDeclarationScopeBeginIndexSlice[1]:structDeclarationScopeBeginIndexSlice[1]+structDeclarationLength]))
		// fmt.Println()
		// return

		for _, lineContent := range strings.Split(string(fileContent[structDeclarationScopeBeginIndexSlice[1]:structDeclarationScopeBeginIndexSlice[1]+structDeclarationLength]), "\n") {
			trimSpaceString := strings.TrimSpace(lineContent)
			if len(trimSpaceString) == 0 {
				continue
			}

			// fmt.Println()
			// fmt.Printf("|%v|", trimSpaceString)
			// fmt.Println()

			submatchSlice := GoVariableDeclarationRegexp.FindStringSubmatch(trimSpaceString)
			if len(submatchSlice) == 0 {
				continue
			}
			memberName := submatchSlice[GoVariableDeclarationRegexpSubmatchNameIndex]
			fileStructDeclarationMap[structName].MemberDeclarationMap[memberName] = &GoVariableDefinition{
				VariableSignature: memberName,
				TypeDeclaration:   ExtractGoVariableTypeDeclaration(submatchSlice[GoVariableDeclarationRegexpSubmatchTypeIndex]),
			}

			// fmt.Println()
			// fileStructDeclarationMap[structName].MemberDeclarationMap[memberName].TypeDeclaration.Traversal(0)
			// fmt.Println()
		}
	}
	return fileStructDeclarationMap
}

type goStructMeta struct {
	typeSpec *ast.TypeSpec
}

func extractGoStructMeta(extractFilepath string, structName string) (*goStructMeta, error) {
	astFile, err := parser.ParseFile(token.NewFileSet(), extractFilepath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	gsm := searchGoStructMeta(astFile, structName)
	if gsm.typeSpec == nil {
		return nil, fmt.Errorf("can not find struct decl")
	}

	return gsm, nil
}

func searchGoStructMeta(fileAST *ast.File, structName string) *goStructMeta {
	var structDecl *ast.TypeSpec
	ast.Inspect(fileAST, func(n ast.Node) bool {
		if n == nil || structDecl != nil {
			return false
		}
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}
		if typeSpec.Name.String() == structName {
			structDecl = typeSpec
			return false
		}
		return true
	})
	return &goStructMeta{
		typeSpec: structDecl,
	}
}

func (gsm *goStructMeta) StructName() string {
	return gsm.typeSpec.Name.String()
}
