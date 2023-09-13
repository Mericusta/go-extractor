package extractor

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"

	stp "github.com/Mericusta/go-stp"
)

// GoProjectMeta
type GoProjectMeta struct {
	// projects absolute path
	projectPath string

	// projects module name, extract from go.mod
	// if there is no module name
	// that means project doesn't have go.mod file
	// and also can not export
	moduleName string

	// projects all package meta
	// package import path as key
	packageMap map[string]*GoPackageMeta
}

// NOTE: 不用 parser.ParseDir 是因为我需要自行组织结构

func ExtractGoProjectMetaByDir(projectPath string, toHandlePaths map[string]struct{}, spec bool) (*GoProjectMeta, error) {
	projectAbsPath, err := filepath.Abs(projectPath)
	if err != nil {
		return nil, err
	}
	projectDirStat, err := os.Stat(projectPath)
	if err != nil {
		return nil, err.(*os.PathError)
	}

	toHandleAbsPaths := make(map[string]struct{})
	for toHandleRelPath := range toHandlePaths {
		pathAbs, err := filepath.Abs(toHandleRelPath)
		if err != nil {
			return nil, err
		}
		toHandleAbsPaths[pathAbs] = struct{}{}
	}
	if spec {
		toHandleAbsPaths[filepath.Join(projectAbsPath, "go.mod")] = struct{}{}
	}

	projectMeta := &GoProjectMeta{
		projectPath: projectAbsPath,
		packageMap:  make(map[string]*GoPackageMeta),
	}

	tmpParseDir := filepath.Join(projectAbsPath, "tmp")
	if projectDirStat.IsDir() {
		err := CreateDir(tmpParseDir)

		if err != nil {
			panic(err)
		}
		// filepath.WalkDir()
		TraverseDirectorySpecificFileWithFunction(
			projectAbsPath, "go",
			func(s string, de fs.DirEntry) error {
				fmt.Printf("s = %v\n", s)
				fmt.Printf("de.Name = %v\n", de.Name())
				f1, err := ioutil.ReadFile(s)
				if err != nil {
					panic(err)
				}
				err = ioutil.WriteFile(filepath.Join(tmpParseDir, de.Name()), f1, 0644)
				if err != nil {
					panic(err)
				}
				return nil
			},
		)
	} else {
		panic("not dir")
	}

	projectFileSet := token.NewFileSet()
	pkgMap, err := parser.ParseDir(projectFileSet, tmpParseDir, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	for pkgPath := range pkgMap {
		fmt.Printf("pkgPath %v\n", pkgPath)
	}

	err = os.RemoveAll(tmpParseDir)
	if err != nil {
		panic(err)
	}

	return projectMeta, nil
}

func ExtractGoProjectMeta(projectPath string, ignorePaths map[string]struct{}) (*GoProjectMeta, error) {
	return extractGoProjectMeta(projectPath, ignorePaths, false)
}

func ExtractGoProjectMetaWithSpecPaths(projectPath string, specPaths map[string]struct{}) (*GoProjectMeta, error) {
	return extractGoProjectMeta(projectPath, specPaths, true)
}

func extractGoProjectMeta(projectPath string, toHandlePaths map[string]struct{}, spec bool) (*GoProjectMeta, error) {
	projectAbsPath, err := filepath.Abs(projectPath)
	if err != nil {
		return nil, err
	}
	projectDirStat, err := os.Stat(projectPath)
	if err != nil {
		return nil, err.(*os.PathError)
	}

	toHandleAbsPaths := make(map[string]struct{})
	for toHandleRelPath := range toHandlePaths {
		pathAbs, err := filepath.Abs(toHandleRelPath)
		if err != nil {
			return nil, err
		}
		toHandleAbsPaths[pathAbs] = struct{}{}
	}
	if spec {
		toHandleAbsPaths[filepath.Join(projectAbsPath, "go.mod")] = struct{}{}
	}

	projectMeta := &GoProjectMeta{
		projectPath: projectAbsPath,
		packageMap:  make(map[string]*GoPackageMeta),
	}

	if projectDirStat.IsDir() {
		hasGoMod := false
		err = filepath.WalkDir(projectAbsPath, func(path string, d fs.DirEntry, err error) error {
			if path == projectAbsPath {
				return nil
			}
			if !d.IsDir() {
				if (!spec && isInPaths(toHandleAbsPaths, path)) || (spec && !isInPaths(toHandleAbsPaths, path)) {
					return nil
				}
				if d.Name() == "go.mod" {
					if projectAbsPath != filepath.Dir(path) {
						return fmt.Errorf("go.mod not exists project root path")
					}
					hasGoMod = true
					moduleName, err := ExtractGoModuleName(path)
					if err != nil {
						return err
					}
					projectMeta.moduleName = moduleName
				} else if filepath.Ext(path) == ".go" {
					fileMeta, err := ExtractGoFileMeta(path)
					if err != nil {
						return err
					}
					filePkg := fileMeta.PkgName()
					if filePkg == "main" {
						if _, has := projectMeta.packageMap[filePkg]; !has {
							projectMeta.packageMap[filePkg] = NewGoPackageMeta(filePkg, filepath.Dir(path), "")
						}
						projectMeta.packageMap[filePkg].pkgFileMap[fileMeta.name] = fileMeta
					} else {
						relPath := "."
						fileDir := filepath.Dir(path)
						if fileDir != projectAbsPath {
							relPath, err = filepath.Rel(projectAbsPath, filepath.Dir(fileDir))
							if err != nil {
								return err
							}
						}
						pkgImportPath := FormatFilePathWithOS(filepath.Clean(fmt.Sprintf("%v/%v/%v", projectMeta.moduleName, relPath, filePkg)), "linux")
						if _, has := projectMeta.packageMap[pkgImportPath]; !has {
							projectMeta.packageMap[pkgImportPath] = NewGoPackageMeta(filePkg, filepath.Dir(path), pkgImportPath)
						}
						projectMeta.packageMap[pkgImportPath].pkgFileMap[fileMeta.name] = fileMeta
					}
				}
			}
			return nil
		})

		if err != nil {
			return nil, err
		}

		if !hasGoMod {
			return nil, fmt.Errorf("go.mod not exists in project path")
		}
	} else {
		if (!spec && isInPaths(toHandleAbsPaths, projectAbsPath)) || (spec && !isInPaths(toHandleAbsPaths, projectAbsPath)) {
			return nil, fmt.Errorf("project path not in handle list")
		}

		gfm, err := ExtractGoFileMeta(projectAbsPath)
		if err != nil {
			return nil, err
		}
		projectMeta.packageMap[gfm.PkgName()] = NewGoPackageMeta(gfm.PkgName(), filepath.Dir(projectAbsPath), "")
		projectMeta.packageMap[gfm.PkgName()].pkgFileMap[filepath.Base(projectAbsPath)] = gfm
	}

	return projectMeta, nil
}

func (gpm *GoProjectMeta) ProjectPath() string {
	return gpm.projectPath
}

func (gpm *GoProjectMeta) ModuleName() string {
	return gpm.moduleName
}

func (gpm *GoProjectMeta) Packages() []string {
	return stp.Key(gpm.packageMap)
}

func (gpm *GoProjectMeta) SearchPackageMeta(pkgImportPath string) *GoPackageMeta {
	return gpm.packageMap[pkgImportPath]
}

// func (gpm *GoProjectMeta) SearchArgType(gam *GoArgMeta) string {
// 	var goPackageMeta *GoPackageMeta
// 	var pkgImportPath string
// 	headMeta := gam.Head()
// 	importMeta, ok := headMeta.typeMeta.(*GoImportMeta)
// 	if ok {
// 		pkgImportPath = importMeta.ImportPath()
// 	} else {
// 		filePkg, err := extractGoFilePkgName(gam.meta.path)
// 		if err != nil {
// 			fmt.Printf("extract go file %v pkg name occurs error: %v\n", gam.meta.path, err)
// 			return ""
// 		}
// 		if filePkg == "main" {
// 			pkgImportPath = filePkg
// 		} else {
// 			relPath := "."
// 			fileDir := filepath.Dir(gam.meta.path)
// 			if fileDir != gpm.projectPath {
// 				relPath, err = filepath.Rel(gpm.projectPath, filepath.Dir(fileDir))
// 				if err != nil {
// 					fmt.Printf("get rel path from project path %v to arg file path %v occurs error: %v", gpm.projectPath, gam.meta.path, err)
// 					return ""
// 				}
// 			}
// 			pkgImportPath = FormatFilePathWithOS(filepath.Clean(fmt.Sprintf("%v/%v/%v", gpm.moduleName, relPath, filePkg)), "linux")
// 		}
// 	}
// 	goPackageMeta = gpm.SearchPackageMeta(strings.Trim(pkgImportPath, "\""))
// 	fmt.Printf("arg %v, importPkgPath = %v\n", gam.Expression(), pkgImportPath)
// 	fmt.Printf("goPackageMeta = %v\n", goPackageMeta)

// 	nextMeta := headMeta
// 	lastMetaType := nextMeta.typeEnum
// SWITCH:
// 	switch nextMeta.typeEnum {
// 	case TYPE_PKG_ALIAS:
// 		fmt.Printf("next meta %v is TYPE_PKG_ALIAS\n", nextMeta.Name())
// 		lastMetaType = TYPE_PKG_ALIAS
// 		nextMeta = gam.next(nextMeta.node)
// 		goto SWITCH
// 	case TYPE_ASSIGNMENT:
// 		fmt.Printf("next meta is TYPE_ASSIGNMENT\n")
// 	case TYPE_FUNC_CALL:
// 		fmt.Printf("next meta %v is TYPE_FUNC_CALL\n", nextMeta.Name())

// 		nextMeta = gam.next(nextMeta.node)
// 		// TODO: func 还是 method 取决于上一个 meta
// 		if lastMetaType == TYPE_PKG_ALIAS {
// 			gfm := goPackageMeta.SearchFunctionMeta(nextMeta.Name())
// 			fmt.Printf("search gfm %v from meta %v\n", gfm, nextMeta.Name())

// 			// nextMeta.typeMeta =
// 		}
// 		goto SWITCH
// 	case TYPE_VAR_FIELD:
// 		fmt.Printf("next meta is TYPE_VAR_FIELD\n")
// 	case TYPE_CONSTANTS:
// 		fmt.Printf("next meta is TYPE_CONSTANTS\n")
// 	default:
// 		fmt.Printf("next meta is TYPE_UNKNOWN\n")
// 	}

// 	// for next != nil {
// 	// 	fmt.Printf("next = %v\n", next.Expression())
// 	// 	fmt.Printf("next = %v\n", next.Name())

// 	// 	gfm := goPackageMeta.SearchFunctionMeta(next.Name())

// 	// 	next = gam.next(next.node)
// 	// }

// 	// nodeSlice := make([]ast.Node, 0, 64)
// 	// var selectorExpr *ast.SelectorExpr
// 	// ast.Inspect(gam.node, func(n ast.Node) bool {
// 	// 	nodeSlice = append(nodeSlice, n)
// 	// 	selectorExpr, ok = n.(*ast.SelectorExpr)
// 	// 	if !ok {
// 	// 		return true
// 	// 	}
// 	// 	if
// 	// 	return false
// 	// })

// 	// for _, n := range nodeSlice {
// 	// 	ast.Print(token.NewFileSet(), n)
// 	// 	fmt.Println()
// 	// }

// 	return ""
// }
