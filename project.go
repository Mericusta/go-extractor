package extractor

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
)

// 以文件为单位提取
// 以包为单位整合
// 以包为单位查找

// NOTE: 不用 parser.ParseDir 是因为我需要自行组织结构

// GoProjectMeta go 项目 meta 数据
type GoProjectMeta struct {
	// 项目绝对路径
	absolutePath string

	// 项目名称，提取自 go.mod 文件
	// 如果没有，则无项目名称，也无法导出
	moduleName string

	// 项目内所有 package 的 meta 数据
	// - key: package 的导入路径
	// - value: package 的 meta 数据
	packageMap map[string]*GoPackageMeta
}

// newGoProjectMeta 通过 ast 构造 项目 的 meta
func newGoProjectMeta(absolutePath, moduleName string) *GoProjectMeta {
	return &GoProjectMeta{
		absolutePath: absolutePath,
		moduleName:   moduleName,
		packageMap:   make(map[string]*GoPackageMeta),
	}
}

// -------------------------------- extractor --------------------------------

// extractGoProjectMetaByDir 通过指定目录提取项目 meta 数据，针对特定目录执行特定操作
func extractGoProjectMetaByDir(projectPath string, toHandlePaths map[string]struct{}, spec bool) (*GoProjectMeta, error) {
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
		absolutePath: projectAbsPath,
		packageMap:   make(map[string]*GoPackageMeta),
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

// ExtractGoProjectMeta 通过指定目录提取项目 meta 数据
// - 指定忽略路径
// - 递归提取
func ExtractGoProjectMeta(projectPath string, ignorePaths map[string]struct{}) (*GoProjectMeta, error) {
	return extractGoProjectMeta(projectPath, ignorePaths, false)
}

// ExtractGoProjectMeta 通过指定目录提取项目 meta 数据
// - 指定特定路径
// - 递归提取
func ExtractGoProjectMetaWithSpecPaths(projectPath string, specPaths map[string]struct{}) (*GoProjectMeta, error) {
	return extractGoProjectMeta(projectPath, specPaths, true)
}

// extractGoProjectMeta 通过指定目录提取项目 meta 数据
// - 递归提取
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
		absolutePath: projectAbsPath,
		packageMap:   make(map[string]*GoPackageMeta),
	}

	if projectDirStat.IsDir() {
		// project 是目录则必须存在 go.mod
		hasGoMod := false
		err = filepath.WalkDir(projectAbsPath, func(walkPath string, d fs.DirEntry, err error) error {
			if walkPath == projectAbsPath {
				// 跳过根目录
				return nil
			}
			if !d.IsDir() {
				if (!spec && isInPaths(toHandleAbsPaths, walkPath)) || (spec && !isInPaths(toHandleAbsPaths, walkPath)) {
					// 跳过 非指定情况下的忽略路径 以及 指定情况下的非指定路径
					return nil
				}
				// 处理 go.mod 文件
				if d.Name() == "go.mod" {
					if projectAbsPath != filepath.Dir(walkPath) {
						return fmt.Errorf("go.mod not exists project root path")
					}
					hasGoMod = true
					moduleName, err := ExtractGoModuleName(walkPath)
					if err != nil {
						return err
					}
					projectMeta.moduleName = moduleName
				} else if filepath.Ext(walkPath) == ".go" {
					fileMeta, err := ExtractGoFileMeta(walkPath)
					if err != nil {
						return err
					}
					filePkg := fileMeta.PackageName()
					if filePkg == "main" {
						if _, has := projectMeta.packageMap[filePkg]; !has {
							projectMeta.packageMap[filePkg] = newGoPackageMeta(filePkg, filepath.Dir(walkPath), "")
						}
						projectMeta.packageMap[filePkg].fileMetaMap[fileMeta.ident] = fileMeta
					} else {
						relPath := "."
						fileDir := filepath.Dir(walkPath)
						if fileDir != projectAbsPath {
							relPath, err = filepath.Rel(projectAbsPath, filepath.Dir(fileDir))
							if err != nil {
								return err
							}
						}
						pkgImportPath := FormatFilePathWithOS(filepath.Clean(fmt.Sprintf("%v/%v/%v", projectMeta.moduleName, relPath, filePkg)), "linux")
						if _, has := projectMeta.packageMap[pkgImportPath]; !has {
							projectMeta.packageMap[pkgImportPath] = newGoPackageMeta(filePkg, filepath.Dir(walkPath), pkgImportPath)
						}
						projectMeta.packageMap[pkgImportPath].fileMetaMap[fileMeta.ident] = fileMeta
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
		// project 不是目录则按照单个文件处理
		if (!spec && isInPaths(toHandleAbsPaths, projectAbsPath)) || (spec && !isInPaths(toHandleAbsPaths, projectAbsPath)) {
			return nil, fmt.Errorf("project path not in handle list")
		}

		gfm, err := ExtractGoFileMeta(projectAbsPath)
		if err != nil {
			return nil, err
		}
		projectMeta.packageMap[gfm.PackageName()] = newGoPackageMeta(gfm.PackageName(), filepath.Dir(projectAbsPath), "")
		projectMeta.packageMap[gfm.PackageName()].fileMetaMap[filepath.Base(projectAbsPath)] = gfm
	}

	// 提取所有 package 的 meta 数据后，递归提取
	for _, gpm := range projectMeta.packageMap {
		// 提取 package 的所有子 meta 数据
		gpm.ExtractAll()
	}

	return projectMeta, nil
}

// -------------------------------- extractor --------------------------------

// SearchPackageMeta 根据 package 的 导入路径 搜索 package 的 meta 数据
func (gpm *GoProjectMeta) SearchPackageMeta(packageImportPath string) *GoPackageMeta {
	return gpm.packageMap[packageImportPath]
}

// -------------------------------- unit test --------------------------------

func (gpm *GoProjectMeta) AbsolutePath() string                  { return gpm.absolutePath }
func (gpm *GoProjectMeta) ModuleName() string                    { return gpm.moduleName }
func (gpm *GoProjectMeta) PackageMap() map[string]*GoPackageMeta { return gpm.packageMap }

// -------------------------------- unit test --------------------------------
