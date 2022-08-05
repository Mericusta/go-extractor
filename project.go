package extractor

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// GoProjectMeta
type GoProjectMeta struct {
	// projects absolute path
	ProjectPath string

	// projects module name, extract from go.mod
	// if there is no module name
	// that means project doesn't have go.mod file
	// and also can not export
	ModuleName string

	// projects all package meta
	// package import path as key
	PackageMap map[string]*goPackageMeta

	// ignore file or folder
	ignorePaths map[string]struct{}
}

func ExtractGoProjectMeta(projectPath string, ignorePaths map[string]struct{}) (*GoProjectMeta, error) {
	// get project abs path
	projectPathAbs, err := filepath.Abs(projectPath)
	if err != nil {
		return nil, err
	}
	projectDirInfo, err := os.Stat(projectPath)
	if err != nil {
		return nil, err.(*os.PathError)
	}

	projectMeta := &GoProjectMeta{
		ProjectPath: projectPathAbs,
		PackageMap:  make(map[string]*goPackageMeta),
		ignorePaths: make(map[string]struct{}),
	}

	if projectDirInfo.IsDir() {
		// get ignore abs path
		for ignorePath := range ignorePaths {
			ignorePathAbs, err := filepath.Abs(ignorePath)
			if err != nil {
				return nil, err
			}
			projectMeta.ignorePaths[ignorePathAbs] = struct{}{}
		}

		hasGoMod := false
		err = filepath.Walk(projectPathAbs, func(path string, info fs.FileInfo, err error) error {
			fmt.Printf("path = %v\n", path)
			if path == projectPathAbs {
				fmt.Printf("path is project path\n")
				return nil
			}
			if projectMeta.isSubpath(path) {
				fmt.Printf("path is ignore path\n")
				return nil
			}
			if info.IsDir() {
				fmt.Printf("path is dir\n")
			} else {
				if info.Name() == "go.mod" {
					fmt.Printf("path is go.mod\n")
					if projectPathAbs != filepath.Dir(path) {
						return fmt.Errorf("go.mod not exists project root path")
					}
					hasGoMod = true
					moduleName, err := extractGoModuleName(path)
					if err != nil {
						return err
					}
					projectMeta.ModuleName = moduleName
				} else if filepath.Ext(path) == ".go" {
					fmt.Printf("path is .go file\n")
					fileMeta, err := extractGoFileMeta(path)
					if err != nil {
						return err
					}
					filePkg := fileMeta.PkgName()
					if filePkg == "main" {
						if _, has := projectMeta.PackageMap[filePkg]; !has {
							projectMeta.PackageMap[filePkg] = &goPackageMeta{
								Name:    filePkg,
								PkgPath: filepath.Dir(path),
								pkgFileMap: map[string]*goFileMeta{
									fileMeta.Name: fileMeta,
								},
							}
						} else {
							projectMeta.PackageMap[filePkg].pkgFileMap[fileMeta.Name] = fileMeta
						}
					} else {
						relPath := "."
						fileDir := filepath.Dir(path)
						if fileDir != projectPathAbs {
							relPath, err = filepath.Rel(projectPathAbs, filepath.Dir(fileDir))
							if err != nil {
								return err
							}
						}
						pkgImportPath := FormatFilePathWithOS(filepath.Clean(fmt.Sprintf("%v/%v/%v", projectMeta.ModuleName, relPath, filePkg)), "linux")
						if _, has := projectMeta.PackageMap[pkgImportPath]; !has {
							projectMeta.PackageMap[pkgImportPath] = &goPackageMeta{
								Name:       filePkg,
								ImportPath: pkgImportPath,
								PkgPath:    filepath.Dir(path),
								pkgFileMap: map[string]*goFileMeta{
									fileMeta.Name: fileMeta,
								},
							}
						} else {
							projectMeta.PackageMap[pkgImportPath].pkgFileMap[filepath.Base(path)] = fileMeta
						}
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
		gfm, err := extractGoFileMeta(projectPathAbs)
		if err != nil {
			return nil, err
		}
		projectMeta.PackageMap[gfm.PkgName()] = &goPackageMeta{
			Name:    gfm.PkgName(),
			PkgPath: filepath.Dir(projectPathAbs),
			pkgFileMap: map[string]*goFileMeta{
				filepath.Base(projectPathAbs): gfm,
			},
		}
	}
	return projectMeta, nil
}

func (gpm *GoProjectMeta) isSubpath(path string) bool {
	for ignorePath := range gpm.ignorePaths {
		if strings.Contains(path, ignorePath) {
			return true
		}
	}
	return false
}
