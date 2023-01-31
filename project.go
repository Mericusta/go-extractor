package extractor

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	stpmap "github.com/Mericusta/go-stp/map"
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
	return stpmap.Key(gpm.packageMap)
}

func (gpm *GoProjectMeta) SearchPackageMeta(pkgImportPath string) *GoPackageMeta {
	return gpm.packageMap[pkgImportPath]
}

func (gpm *GoProjectMeta) SearchArgType(gam *GoArgMeta) string {
	var goPackageMeta *GoPackageMeta
	var importPkgPath string
	headMeta := gam.Head()
	importMeta, ok := headMeta.typeMeta.(*GoImportMeta)
	if ok {
		importPkgPath = importMeta.ImportPath()
	} else {
		filePkg, err := extractGoFilePkgName(gam.meta.path)
		if err != nil {
			fmt.Printf("extract go file %v pkg name occurs error: %v\n", gam.meta.path, err)
			return ""
		}
		if filePkg == "main" {
			importPkgPath = filePkg
		} else {
			// relPath := "."
			// fileDir := filepath.Dir(path)
			// if fileDir != projectAbsPath {
			// 	relPath, err = filepath.Rel(projectAbsPath, filepath.Dir(fileDir))
			// 	if err != nil {
			// 		return err
			// 	}
			// }
			// pkgImportPath := FormatFilePathWithOS(filepath.Clean(fmt.Sprintf("%v/%v/%v", projectMeta.moduleName, relPath, filePkg)), "linux")
		}
	}
	goPackageMeta = gpm.SearchPackageMeta(importPkgPath)
	fmt.Printf("goPackageMeta = %v\n", goPackageMeta)

	return ""
}
