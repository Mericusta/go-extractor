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
	PackageMap map[string]*goPackageMeta
}

func ExtractGoProjectMeta(projectPath string, ignorePaths map[string]struct{}) (*GoProjectMeta, error) {
	return extractGoProjectMeta(projectPath, ignorePaths, false)
}

func ExtractGoProjectMetaWithSpecPaths(projectPath string, specPaths map[string]struct{}) (*GoProjectMeta, error) {
	return extractGoProjectMeta(projectPath, specPaths, true)
}

func extractGoProjectMeta(projectPath string, paths map[string]struct{}, spec bool) (*GoProjectMeta, error) {
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
	}

	// get ignore abs path
	pathsAbs := make(map[string]struct{})
	for path := range paths {
		pathAbs, err := filepath.Abs(path)
		if err != nil {
			return nil, err
		}
		pathsAbs[pathAbs] = struct{}{}
	}
	if spec {
		pathsAbs[filepath.Join(projectPathAbs, "go.mod")] = struct{}{}
	}

	if projectDirInfo.IsDir() {
		hasGoMod := false
		err = filepath.WalkDir(projectPathAbs, func(path string, d fs.DirEntry, err error) error {
			if path == projectPathAbs {
				return nil
			}
			if !d.IsDir() {
				if (!spec && isInPaths(pathsAbs, path)) || (spec && !isInPaths(pathsAbs, path)) {
					return nil
				}
				if d.Name() == "go.mod" {
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
		if (!spec && isInPaths(pathsAbs, projectPathAbs)) || (spec && !isInPaths(pathsAbs, projectPathAbs)) {
			return nil, fmt.Errorf("project path not in handle list")
		}

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
