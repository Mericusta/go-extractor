package extractor

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

type GoImportDeclaration struct {
	ImportAliasPathMap map[string]string
}

func (d *GoImportDeclaration) Traversal(deep int) {
	fmt.Printf("%v- Import Package List: ", strings.Repeat("\t", deep))
	if len(d.ImportAliasPathMap) > 0 {
		fmt.Println()
		for alias, path := range d.ImportAliasPathMap {
			fmt.Printf("%v- Alias: %v\n", strings.Repeat("\t", deep+1), alias)
			fmt.Printf("%v- Path: %v\n", strings.Repeat("\t", deep+1), path)
		}
	} else {
		fmt.Printf("None\n")
	}
	fmt.Printf("%v- MakeUp: |%v|\n", strings.Repeat("\t", deep), d.MakeUp(false))
}

func (d *GoImportDeclaration) MakeUp(withBracket bool) string {
	makeUpTemplate := `import [IMPORT_SCOPE]`
	singleImportTemplate := `EACH_IMPORT`
	multiImportTemplate := `(EACH_IMPORT)`
	eachImportTemplate := `[ALIAS] "[PATH]"`

	makeUpReplaceKeywordImportScope := `[IMPORT_SCOPE]`
	makeUpReplaceKeywordEachImport := `EACH_IMPORT`
	makeUpReplaceKeywordAlias := `[ALIAS]`
	makeUpReplaceKeywordPath := `[PATH]`

	// import scope
	builder := strings.Builder{}
	for alias, path := range d.ImportAliasPathMap {
		replaceAlias := alias
		if filepath.Base(path) == alias {
			replaceAlias = ""
		}
		builder.WriteString("\n\t")
		eachImportContent := strings.Replace(eachImportTemplate, makeUpReplaceKeywordAlias, replaceAlias, -1)
		eachImportContent = strings.Replace(eachImportContent, makeUpReplaceKeywordPath, path, -1)
		builder.WriteString(strings.TrimSpace(eachImportContent))
	}

	var makeUpContent string
	l := len(d.ImportAliasPathMap)
	switch {
	case l == 1 && !withBracket:
		importScopeContent := strings.Replace(singleImportTemplate, makeUpReplaceKeywordEachImport, strings.TrimSpace(builder.String()), -1)
		makeUpContent = strings.Replace(makeUpTemplate, makeUpReplaceKeywordImportScope, importScopeContent, -1)
	case (l == 1 && withBracket) || l > 1:
		builder.WriteString("\n")
		importScopeContent := strings.Replace(multiImportTemplate, makeUpReplaceKeywordEachImport, builder.String(), -1)
		makeUpContent = strings.Replace(makeUpTemplate, makeUpReplaceKeywordImportScope, importScopeContent, -1)
	}

	// add new line
	makeUpContent = "\n" + makeUpContent + "\n"

	// remove all CR
	makeUpContent = strings.ReplaceAll(makeUpContent, "\r", "")

	return makeUpContent
}

var (
	GO_IMPORT_SCOPE_BEGIN_EXPRESSION              = `import\s*\(`
	GoImportScopeBeginRegexp                      = regexp.MustCompile(GO_IMPORT_SCOPE_BEGIN_EXPRESSION)
	GO_SINGLE_IMPORT_SCOPE_EXPRESSION             = `import\s+(?P<CONTENT>(\w+\s+)?"\S+")`
	GoSingleImportScopeRegexp                     = regexp.MustCompile(GO_SINGLE_IMPORT_SCOPE_EXPRESSION)
	GoSingleImportScopeRegexpSubmatchContentIndex = GoSingleImportScopeRegexp.SubexpIndex("CONTENT")
)

func ExtractGoFileImportDeclaration(fileContent []byte) *GoImportDeclaration {
	var importScopeContent []byte
	importScopeBeginIndexSlice := GoImportScopeBeginRegexp.FindIndex(fileContent)
	if len(importScopeBeginIndexSlice) == 0 {
		subMatchSlice := GoSingleImportScopeRegexp.FindSubmatch(fileContent)
		if len(subMatchSlice) > 0 {
			importScopeContent = subMatchSlice[GoSingleImportScopeRegexpSubmatchContentIndex]
		} else {
			return nil
		}
	} else {
		importScopeContent = GetScopeContentBetweenPunctuationMarks(fileContent, importScopeBeginIndexSlice[1]-1)
	}

	// fmt.Printf("import scope content: |%v|\n", string(fileContent[importScopeBeginIndexSlice[1]+1:importScopeBeginIndexSlice[1]+1+importScopeLength]))
	goImportDeclaration := &GoImportDeclaration{
		ImportAliasPathMap: make(map[string]string),
	}
	for _, eachImportString := range strings.Split(string(importScopeContent), "\n") {
		for _, submatchSlice := range GoImportRegexp.FindAllStringSubmatch(strings.TrimSpace(eachImportString), -1) {
			if GoImportRegexpSubmatchAliasIndex == -1 || len(submatchSlice[GoImportRegexpSubmatchAliasIndex]) == 0 {
				goImportDeclaration.ImportAliasPathMap[filepath.Base(submatchSlice[GoImportRegexpSubmatchPathIndex])] = submatchSlice[GoImportRegexpSubmatchPathIndex]
			} else {
				goImportDeclaration.ImportAliasPathMap[submatchSlice[GoImportRegexpSubmatchAliasIndex]] = submatchSlice[GoImportRegexpSubmatchPathIndex]
			}
		}
	}
	return goImportDeclaration
}
