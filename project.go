package extractor

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// GoProjectMeta
// go project which uses go module
// go.mod must exists project root path
type GoProjectMeta struct {
	ProjectPath string // projects absolute path
	ModuleName  string // projects module name
	PackageMap  map[string]*goPackageMeta
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
	if !projectDirInfo.IsDir() {
		return nil, fmt.Errorf("project path is not dir")
	}

	projectMeta := &GoProjectMeta{
		ProjectPath: projectPathAbs,
		PackageMap:  make(map[string]*goPackageMeta),
		ignorePaths: make(map[string]struct{}),
	}
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
		if path == projectPathAbs {
			return nil
		}
		if _, isIgnore := ignorePaths[filepath.Dir(path)]; isIgnore {
			return nil
		}
		if info.IsDir() {
			fmt.Printf("dir = %v\n", path)
		} else {
			if info.Name() == "go.mod" {
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
				fmt.Printf("path = %v\n", path)
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
					pkgImportPath := fmt.Sprintf("%v/%v", projectMeta.ModuleName, filePkg)
					if _, has := projectMeta.PackageMap[pkgImportPath]; !has {
						projectMeta.PackageMap[pkgImportPath] = &goPackageMeta{
							Name:       filePkg,
							ImportPath: pkgImportPath,
							PkgPath:    filepath.Dir(path),
							pkgFileMap: map[string]*goFileMeta{
								fileMeta.Name: fileMeta,
							},
						}
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

	return projectMeta, nil
}
