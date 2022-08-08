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

type goPackageMeta struct {
	Name             string                      // pkg name
	PkgPath          string                      // pkg path
	ImportPath       string                      // import path
	pkgFileMap       map[string]*goFileMeta      // all file meta
	pkgStructDecl    map[string]*goStructMeta    // all struct meta
	pkgInterfaceDecl map[string]*goInterfaceMeta // all interface meta
	pkgFunctionDecl  map[string]*goFunctionMeta  // all function meta
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

func (gpm *goPackageMeta) SearchInterfaceMeta(interfaceName string) *goInterfaceMeta {
	if len(gpm.pkgInterfaceDecl) > 0 {
		if _, has := gpm.pkgInterfaceDecl[interfaceName]; has {
			return gpm.pkgInterfaceDecl[interfaceName]
		}
	}

	var gim *goInterfaceMeta
	for _, gfm := range gpm.pkgFileMap {
		if gfm.fileAST != nil && gfm.fileAST.Scope != nil && gim == nil {
			gim = searchGoInterfaceMeta(gfm.fileAST, interfaceName)
		}
	}

	if gpm.pkgInterfaceDecl == nil {
		gpm.pkgInterfaceDecl = make(map[string]*goInterfaceMeta)
	}
	gpm.pkgInterfaceDecl[gim.InterfaceName()] = gim

	return gpm.pkgInterfaceDecl[gim.InterfaceName()]
}

func (gpm *goPackageMeta) SearchFunctionMeta(functionName string) *goFunctionMeta {
	if len(gpm.pkgFunctionDecl) > 0 {
		if _, has := gpm.pkgFunctionDecl[functionName]; has {
			return gpm.pkgFunctionDecl[functionName]
		}
	}

	var gsm *goFunctionMeta
	for _, gfm := range gpm.pkgFileMap {
		if gfm.fileAST != nil && gfm.fileAST.Scope != nil && gsm == nil {
			gsm = searchGoFunctionMeta(gfm.fileAST, functionName)
		}
	}

	if gpm.pkgFunctionDecl == nil {
		gpm.pkgFunctionDecl = make(map[string]*goFunctionMeta)
	}
	gpm.pkgFunctionDecl[gsm.FunctionName()] = gsm

	return gpm.pkgFunctionDecl[gsm.FunctionName()]
}

// SearchMethodMeta search method implement from node.(*ast.File)
func (gpm *goPackageMeta) SearchMethodMeta(structName, methodName string) *goMethodMeta {
	for _, gfm := range gpm.pkgFileMap {
		gmm := searchGoMethodMeta(gfm.fileAST, structName, methodName)
		if gmm != nil && gmm.structMethodDecl != nil {
			return gmm
		}
	}
	return nil
}
