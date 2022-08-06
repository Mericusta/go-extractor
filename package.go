package extractor

import (
	"fmt"
	"go/ast"
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

type goPackageMeta struct {
	Name          string                   // pkg name
	PkgPath       string                   // pkg path
	ImportPath    string                   // import path
	pkgFileMap    map[string]*goFileMeta   // each file meta
	pkgStructDecl map[string]*goStructMeta // each struct meta
}

func (gpm *goPackageMeta) searchDeclaration(objectKind ast.ObjKind) ast.Spec {

	return nil
}

func (gpm *goPackageMeta) SearchStructMeta(structName string) *goStructMeta {
	if len(gpm.pkgStructDecl) > 0 {
		if _, has := gpm.pkgStructDecl[structName]; has {
			return gpm.pkgStructDecl[structName]
		}
	}

	var gsm *goStructMeta
	for _, gfm := range gpm.pkgFileMap {
		if gfm.fileAST != nil && gfm.fileAST.Scope != nil && gsm == nil {
			gsm = searchGoStructMeta(gfm.fileAST, structName)
		}
	}

	if gpm.pkgStructDecl == nil {
		gpm.pkgStructDecl = make(map[string]*goStructMeta)
	}
	gpm.pkgStructDecl[gsm.StructName()] = gsm

	return gpm.pkgStructDecl[gsm.StructName()]
}

func SearchFunctionDeclaration() {

}
