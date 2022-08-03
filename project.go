package extractor

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// GoProjectMeta
// go project which uses go module
type GoProjectMeta struct {
	ProjectPath string // projects absolute path
	ModuleName  string // projects module name
	pkgPath     string // go.mod relative path to ProjectPath
	cmdPath     string // main.go relative path to ProjectPath if exists
	PackageMap  map[string]*goPackageMeta
	ignorePaths map[string]interface{}
}

func ExtractGoProjectMeta(projectPath string, ignorePaths map[string]interface{}) (*GoProjectMeta, error) {
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
		ignorePaths: ignorePaths,
	}

	// TODO:
	// search main.go and go.mod
	err = filepath.Walk(projectPathAbs, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			switch info.Name() {
			case "go.mod":
				moduleName, err := extractGoModuleName(path)
				if err != nil {
					return err
				}
				projectMeta.ModuleName = moduleName

				if len(projectMeta.pkgPath) > 0 {
					return fmt.Errorf("there are more than one go.mod in go projects")
				}

				dirPath := filepath.Dir(path)
				relPath, err := filepath.Rel(projectPathAbs, dirPath)
				if err != nil {
					return err
				}
				projectMeta.pkgPath = relPath
			case "main.go":
				dirPath := filepath.Dir(path)

				relPath, err := filepath.Rel(projectPathAbs, dirPath)
				if err != nil {
					return err
				}
				projectMeta.cmdPath = relPath

				mainPkgFileMetaMap := make(map[string]*goFileMeta)
				cmdDirEntrySlice, err := os.ReadDir(dirPath)
				if err != nil {
					return err
				}
				for _, cmdDirEntry := range cmdDirEntrySlice {
					if cmdDirEntry.IsDir() {
						continue
					}
					if filepath.Ext(cmdDirEntry.Name()) == ".go" {
						mainPkgFileMetaMap[cmdDirEntry.Name()] = nil
					}
				}

				projectMeta.PackageMap["main"] = &goPackageMeta{
					Name:       "main",
					PkgPath:    dirPath,
					pkgFileMap: mainPkgFileMetaMap,
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return projectMeta, nil
}
