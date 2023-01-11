package extractor

import (
	"fmt"
	"os"
	"path/filepath"
)

// 以文件为单位提取
// 以包为单位整合
// 以包为单位查找

type GoPackageMeta struct {
	Name             string                      // pkg name
	PkgPath          string                      // pkg path
	ImportPath       string                      // import path
	pkgFileMap       map[string]*GoFileMeta      // all file meta
	pkgStructDecl    map[string]*GoStructMeta    // all struct meta
	pkgInterfaceDecl map[string]*GoInterfaceMeta // all interface meta
	pkgFunctionDecl  map[string]*GoFunctionMeta  // all function meta
}

func NewGoPackageMeta(name, pkgPath, importPath string) *GoPackageMeta {
	return &GoPackageMeta{
		Name:             name,
		PkgPath:          pkgPath,
		ImportPath:       importPath,
		pkgFileMap:       make(map[string]*GoFileMeta),
		pkgStructDecl:    make(map[string]*GoStructMeta),
		pkgInterfaceDecl: make(map[string]*GoInterfaceMeta),
		pkgFunctionDecl:  make(map[string]*GoFunctionMeta),
	}
}

func ExtractGoPackageMeta(projectPath string, ignoreFiles map[string]struct{}) (*GoPackageMeta, error) {
	return extractGoPackageMeta(projectPath, ignoreFiles, false)
}

func ExtractGoPackageMetaWithSpecPaths(projectPath string, specFiles map[string]struct{}) (*GoPackageMeta, error) {
	return extractGoPackageMeta(projectPath, specFiles, true)
}

func extractGoPackageMeta(packagePath string, filePaths map[string]struct{}, spec bool) (*GoPackageMeta, error) {
	packagePathAbs, err := filepath.Abs(packagePath)
	if err != nil {
		return nil, err
	}
	packageDirStat, err := os.Stat(packagePath)
	if err != nil {
		return nil, err.(*os.PathError)
	}

	pathsAbsMap := make(map[string]struct{})
	for path := range filePaths {
		pathAbs, err := filepath.Abs(path)
		if err != nil {
			return nil, err
		}
		pathsAbsMap[pathAbs] = struct{}{}
	}

	packageMeta := &GoPackageMeta{
		PkgPath:          packagePathAbs,
		pkgFileMap:       make(map[string]*GoFileMeta),
		pkgStructDecl:    make(map[string]*GoStructMeta),
		pkgInterfaceDecl: make(map[string]*GoInterfaceMeta),
		pkgFunctionDecl:  make(map[string]*GoFunctionMeta),
	}

	if packageDirStat.IsDir() {
		fileSlice, err := os.ReadDir(packagePathAbs)
		if err != nil {
			return nil, err
		}

		for _, fileInfo := range fileSlice {
			filePathAbs := filepath.Join(packagePathAbs, fileInfo.Name())
			fileStat, err := os.Stat(filePathAbs)
			if err != nil {
				fmt.Printf("get file '%v' state occurs error: %v", filePathAbs, err)
				continue
			}

			if fileStat.IsDir() || (!spec && isInPaths(pathsAbsMap, filePathAbs)) || (spec && !isInPaths(pathsAbsMap, filePathAbs)) || filepath.Ext(fileInfo.Name()) != ".go" {
				continue
			}

			fileMeta, err := ExtractGoFileMeta(filePathAbs)
			if err != nil {
				fmt.Printf("extract go file meta from file path '%v' occurs error: %v", filePathAbs, err)
				continue
			}

			filePkg := fileMeta.PkgName()
			if len(packageMeta.Name) > 0 && packageMeta.Name != filePkg {
				fmt.Printf("difference package name %v - %v in file %v", packageMeta.Name, filePkg, filePathAbs)
				continue
			}
			packageMeta.Name = filePkg
			packageMeta.pkgFileMap[fileInfo.Name()] = fileMeta
		}
	} else {
		return nil, fmt.Errorf("package path '%v' is not a folder", packagePathAbs)
	}

	return packageMeta, nil
}

func (gpm *GoPackageMeta) SearchStructMeta(structName string) *GoStructMeta {
	if gsm, has := gpm.pkgStructDecl[structName]; gsm != nil && has {
		return gsm
	}

	for _, gfm := range gpm.pkgFileMap {
		if gfm.fileAST != nil && gfm.fileAST.Scope != nil {
			gsm := SearchGoStructMeta(gfm.fileAST, structName)
			if gsm != nil {
				gpm.pkgStructDecl[gsm.StructName()] = gsm
				break
			}
		}
	}

	return gpm.pkgStructDecl[structName]
}

func (gpm *GoPackageMeta) SearchInterfaceMeta(interfaceName string) *GoInterfaceMeta {
	if gim, has := gpm.pkgInterfaceDecl[interfaceName]; gim != nil && has {
		return gim
	}

	for _, gfm := range gpm.pkgFileMap {
		if gfm.fileAST != nil && gfm.fileAST.Scope != nil {
			gim := SearchGoInterfaceMeta(gfm.fileAST, interfaceName)
			if gim != nil {
				gpm.pkgInterfaceDecl[gim.InterfaceName()] = gim
				break
			}
		}
	}

	return gpm.pkgInterfaceDecl[interfaceName]
}

func (gpm *GoPackageMeta) SearchFunctionMeta(functionName string) *GoFunctionMeta {
	if gfm, has := gpm.pkgFunctionDecl[functionName]; gfm != nil && has {
		return gfm
	}

	for _, gfm := range gpm.pkgFileMap {
		if gfm.fileAST != nil && gfm.fileAST.Scope != nil {
			gfm := SearchGoFunctionMeta(gfm, functionName)
			if gfm != nil {
				gpm.pkgFunctionDecl[gfm.FunctionName()] = gfm
				break
			}
		}
	}

	return gpm.pkgFunctionDecl[functionName]
}

func (gpm *GoPackageMeta) SearchMethodMeta(structName, methodName string) *GoMethodMeta {
	if gsm, hasGSM := gpm.pkgStructDecl[structName]; gsm != nil && hasGSM {
		if gmm, hasGMM := gsm.methodDecl[methodName]; gmm != nil && hasGMM {
			return gmm
		}
	} else {
		gpm.SearchStructMeta(structName)
	}

	for _, gfm := range gpm.pkgFileMap {
		if gfm.fileAST != nil && gfm.fileAST.Scope != nil {
			gmm := SearchGoMethodMeta(gfm.fileAST, structName, methodName)
			if gmm != nil {
				gpm.pkgStructDecl[structName].methodDecl[methodName] = gmm
				break
			}
		}
	}

	return gpm.pkgStructDecl[structName].methodDecl[methodName]
}
