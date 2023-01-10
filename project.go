package extractor

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
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
	PackageMap map[string]*GoPackageMeta
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
		ProjectPath: projectAbsPath,
		PackageMap:  make(map[string]*GoPackageMeta),
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
					projectMeta.ModuleName = moduleName
				} else if filepath.Ext(path) == ".go" {
					fileMeta, err := ExtractGoFileMeta(path)
					if err != nil {
						return err
					}
					filePkg := fileMeta.PkgName()
					if filePkg == "main" {
						if _, has := projectMeta.PackageMap[filePkg]; !has {
							projectMeta.PackageMap[filePkg] = &GoPackageMeta{
								Name:    filePkg,
								PkgPath: filepath.Dir(path),
								pkgFileMap: map[string]*GoFileMeta{
									fileMeta.Name: fileMeta,
								},
							}
						} else {
							projectMeta.PackageMap[filePkg].pkgFileMap[fileMeta.Name] = fileMeta
						}
					} else {
						relPath := "."
						fileDir := filepath.Dir(path)
						if fileDir != projectAbsPath {
							relPath, err = filepath.Rel(projectAbsPath, filepath.Dir(fileDir))
							if err != nil {
								return err
							}
						}
						pkgImportPath := FormatFilePathWithOS(filepath.Clean(fmt.Sprintf("%v/%v/%v", projectMeta.ModuleName, relPath, filePkg)), "linux")
						if _, has := projectMeta.PackageMap[pkgImportPath]; !has {
							projectMeta.PackageMap[pkgImportPath] = &GoPackageMeta{
								Name:       filePkg,
								ImportPath: pkgImportPath,
								PkgPath:    filepath.Dir(path),
								pkgFileMap: map[string]*GoFileMeta{
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
		if (!spec && isInPaths(toHandleAbsPaths, projectAbsPath)) || (spec && !isInPaths(toHandleAbsPaths, projectAbsPath)) {
			return nil, fmt.Errorf("project path not in handle list")
		}

		gfm, err := ExtractGoFileMeta(projectAbsPath)
		if err != nil {
			return nil, err
		}
		projectMeta.PackageMap[gfm.PkgName()] = &GoPackageMeta{
			Name:    gfm.PkgName(),
			PkgPath: filepath.Dir(projectAbsPath),
			pkgFileMap: map[string]*GoFileMeta{
				filepath.Base(projectAbsPath): gfm,
			},
		}
	}

	return projectMeta, nil
}
