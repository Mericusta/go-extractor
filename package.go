package extractor

import (
	"fmt"
	"path/filepath"
	"regexp"
)

// 以文件为单位提取
// 以包为单位整合

var (
	GO_PACKAGE_EXPRESSION            string = `package\s+(?P<NAME>[[:alpha:]][_\w]+)`
	GoPackageRegexp                         = regexp.MustCompile(GO_PACKAGE_EXPRESSION)
	GoPackageRegexpSubmatchNameIndex        = GoPackageRegexp.SubexpIndex("NAME")
	GO_IMPORT_EXPRESSION             string = `((?P<ALIAS>\w+)\s+)?"(?P<PATH>[/_\.\w-]+)"`
	GoImportRegexp                          = regexp.MustCompile(GO_IMPORT_EXPRESSION)
	GoImportRegexpSubmatchAliasIndex        = GoImportRegexp.SubexpIndex("ALIAS")
	GoImportRegexpSubmatchPathIndex         = GoImportRegexp.SubexpIndex("PATH")
)

type GoPackageInfo struct {
	Name    string
	Path    string
	FileMap map[string]*GoFileInfo
}

func (gpi *GoPackageInfo) ImportPath(project string) string {
	return filepath.ToSlash(fmt.Sprintf("%v/%v", project, gpi.Path))
}

func (gpi *GoPackageInfo) MakeUp() string {
	return fmt.Sprintf("package %v", gpi.Name)
}

func ExtractGoFilePackage(fileContent []byte) string {
	submatch := GoPackageRegexp.FindSubmatch(fileContent)
	if GoPackageRegexpSubmatchNameIndex >= len(submatch) {
		panic("can not match package")
	}
	return string(submatch[GoPackageRegexpSubmatchNameIndex])
}
