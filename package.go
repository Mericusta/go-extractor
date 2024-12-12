package extractor

import (
	"fmt"
	"go/ast"
	"os"
	"path/filepath"
)

type GoPackageMetaTypeConstraints interface {
	ast.Node
}

// GoPackageMeta go package 的 meta 数据
type GoPackageMeta struct {
	// 空组合
	*meta[ast.Node]

	// package 标识
	ident string

	// absolutePath 绝对路径
	absolutePath string

	// package 导入路径
	importPath string // import path

	// package 内所有文件的 meta 数据
	// - key: 文件名称
	fileMetaMap map[string]*GoFileMeta[*ast.File]

	// package 内所有 var 的 meta 数据
	// - key: var 标识
	varMetaMap map[string]*GoVarMeta[*ast.ValueSpec]

	// package 内所有 func 的 meta 数据
	// - key: func 标识
	funcMetaMap map[string]*GoFuncMeta[*ast.FuncDecl]

	// package 内所有 struct 的 meta 数据
	// - key: struct 标识
	structMetaMap map[string]*GoStructMeta[*ast.TypeSpec]

	// package 内所有 interface 的 meta 数据
	// - key: interface 标识
	interfaceMetaMap map[string]*GoInterfaceMeta[*ast.TypeSpec]

	// package 内所有 类型约束 的 meta 数据
	// - key: 类型约束 标识
	typeConstraintsMetaMap map[string]*GoInterfaceMeta[*ast.TypeSpec]
}

// newGoPackageMeta 通过 ast 构造 package 的 meta
func newGoPackageMeta(ident, absolutePath, importPath string) *GoPackageMeta {
	return &GoPackageMeta{
		ident:                  ident,
		absolutePath:           absolutePath,
		importPath:             importPath,
		fileMetaMap:            make(map[string]*GoFileMeta[*ast.File]),
		varMetaMap:             make(map[string]*GoVarMeta[*ast.ValueSpec]),
		funcMetaMap:            make(map[string]*GoFuncMeta[*ast.FuncDecl]),
		structMetaMap:          make(map[string]*GoStructMeta[*ast.TypeSpec]),
		interfaceMetaMap:       make(map[string]*GoInterfaceMeta[*ast.TypeSpec]),
		typeConstraintsMetaMap: make(map[string]*GoInterfaceMeta[*ast.TypeSpec]),
	}
}

// -------------------------------- extractor --------------------------------

// ExtractGoPackageMeta 通过 package 的绝对路径提取 package 的 meta 数据
// - 指定忽略文件
// - 递归提取
// - 无法获得 package 的导入路径
func ExtractGoPackageMeta[T GoPackageMetaTypeConstraints](packageRelativePath string, ignoreFiles map[string]struct{}) (*GoPackageMeta, error) {
	return extractGoPackageMeta[T](packageRelativePath, ignoreFiles, false)
}

// ExtractGoPackageMetaWithSpecPaths 通过 package 的绝对路径提取 package 的 meta 数据
// - 指定特定文件
// - 递归提取
// - 无法获得 package 的导入路径
func ExtractGoPackageMetaWithSpecPaths[T GoPackageMetaTypeConstraints](packageRelativePath string, specFiles map[string]struct{}) (*GoPackageMeta, error) {
	return extractGoPackageMeta[T](packageRelativePath, specFiles, true)
}

// extractGoPackageMeta 通过 package 的结对路径提取 package 的 meta 数据
// - 递归提取
// - 无法获得 package 的导入路径
func extractGoPackageMeta[T GoPackageMetaTypeConstraints](packageRelativePath string, files map[string]struct{}, spec bool) (*GoPackageMeta, error) {
	// 查找绝对路径
	packagePathAbs, err := filepath.Abs(packageRelativePath)
	if err != nil {
		return nil, err
	}
	// 查找文件夹
	packageDirStat, err := os.Stat(packageRelativePath)
	if err != nil {
		return nil, err.(*os.PathError)
	}

	// 包内所有文件
	pathsAbsMap := make(map[string]struct{})
	for fileName := range files {
		pathAbs, err := filepath.Abs(filepath.Join(packageRelativePath, fileName))
		if err != nil {
			return nil, err
		}
		pathsAbsMap[pathAbs] = struct{}{}
	}

	// 构造 meta 数据
	packageMeta := &GoPackageMeta{
		absolutePath:     packagePathAbs,
		fileMetaMap:      make(map[string]*GoFileMeta[*ast.File]),
		varMetaMap:       make(map[string]*GoVarMeta[*ast.ValueSpec]),
		structMetaMap:    make(map[string]*GoStructMeta[*ast.TypeSpec]),
		interfaceMetaMap: make(map[string]*GoInterfaceMeta[*ast.TypeSpec]),
		funcMetaMap:      make(map[string]*GoFuncMeta[*ast.FuncDecl]),
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

			filePkg := fileMeta.PackageName()
			if len(packageMeta.ident) > 0 && packageMeta.ident != filePkg {
				fmt.Printf("difference package name %v - %v in file %v", packageMeta.ident, filePkg, filePathAbs)
				continue
			}
			packageMeta.ident = filePkg
			packageMeta.fileMetaMap[fileInfo.Name()] = fileMeta
		}
	} else {
		return nil, fmt.Errorf("package path '%v' is not a folder", packagePathAbs)
	}

	return packageMeta, nil
}

// ExtractAll 提取 package 内所有 var，func，struct，interface 的 meta 数据
func (gpm *GoPackageMeta) ExtractAll() {
	// 提取 var
	gpm.extractVar()

	// 提取 func
	gpm.extractFunc()

	// 提取 struct
	gpm.extractStruct()

	// 提取 method
	gpm.extractMethod()

	// 提取 interface
	gpm.extractInterface()
}

// extractVar 提取 var 的 meta 数据
func (gpm *GoPackageMeta) extractVar() {
	for _, gfm := range gpm.fileMetaMap {
		if gfm.node != nil {
			ast.Inspect(gfm.node, func(n ast.Node) bool {
				switch {
				case IsVarNode(n):
					for _, specNode := range n.(*ast.GenDecl).Specs {
						valueSpec, ok := specNode.(*ast.ValueSpec)
						if valueSpec != nil && ok {
							for _, ident := range valueSpec.Names {
								gpm.varMetaMap[ident.Name] = newGoVarMeta(newMeta(valueSpec, gfm.path), ident.Name)
							}
						}
					}
					return false // 只查找顶层为 var 的节点
				case IsImportNode(n) || IsFuncNode(n) || IsTypeNode(n) || IsMethodNode(n):
					return false // 顶层为其他节点直接跳过
				}
				return true
			})
		}
	}
}

// extractFunc 提取 func 的 meta 数据
func (gpm *GoPackageMeta) extractFunc() {
	for _, gfm := range gpm.fileMetaMap {
		if gfm.node != nil {
			ast.Inspect(gfm.node, func(n ast.Node) bool {
				switch {
				case IsFuncNode(n):
					funcDecl := n.(*ast.FuncDecl)
					funcIdent := funcDecl.Name.String()
					gpm.funcMetaMap[funcIdent] = newGoFuncMeta(newMeta(funcDecl, gfm.path), funcIdent)
					return false // 只查找顶层为 func 的节点
				case IsImportNode(n) || IsVarNode(n) || IsTypeNode(n) || IsMethodNode(n):
					return false // 顶层为其他节点直接跳过
				}
				return true
			})
		}
	}
}

// extractStruct 提取 struct 的 meta 数据
func (gpm *GoPackageMeta) extractStruct() {
	for _, gfm := range gpm.fileMetaMap {
		if gfm.node != nil {
			ast.Inspect(gfm.node, func(n ast.Node) bool {
				if IsTypeNode(n) {
					for _, specNode := range n.(*ast.GenDecl).Specs {
						if IsStructNode(specNode) {
							typeSpec := specNode.(*ast.TypeSpec)
							structIdent := typeSpec.Name.String()
							gpm.structMetaMap[structIdent] = newGoStructMeta(newMeta(typeSpec, gfm.path), structIdent)
						}
					}
					return false // 只查找顶层为 struct 的节点
				}
				return true
			})
		}
	}
}

// extractMethod 提取 method 的 meta 数据
func (gpm *GoPackageMeta) extractMethod() {
	for _, gfm := range gpm.fileMetaMap {
		if gfm.node != nil {
			ast.Inspect(gfm.node, func(n ast.Node) bool {
				switch {
				case IsMethodNode(n):
					funcDecl := n.(*ast.FuncDecl)
					funcIdent := funcDecl.Name.String()
					gmm := newGoMethodMeta(newMeta(funcDecl, gfm.path), funcIdent)
					gsm, has := gpm.structMetaMap[gmm.Receiver().TypeIdent()]
					if gsm != nil && has {
						gsm.methodMetaMap[funcIdent] = gmm
					}
					return false // 只查找顶层为 func 的节点
				case IsImportNode(n) || IsVarNode(n) || IsTypeNode(n) || IsFuncNode(n):
					return false // 顶层为其他节点直接跳过
				}
				return true
			})
		}
	}
}

// extractVar 提取 interface 的 meta 数据
func (gpm *GoPackageMeta) extractInterface() {
	for _, gfm := range gpm.fileMetaMap {
		if gfm.node != nil {
			ast.Inspect(gfm.node, func(n ast.Node) bool {
				if IsTypeNode(n) {
					for _, specNode := range n.(*ast.GenDecl).Specs {
						if IsInterfaceNode(specNode) && !IsTypeConstraintsNode(specNode) {
							typeSpec := specNode.(*ast.TypeSpec)
							interfaceIdent := typeSpec.Name.String()
							gpm.interfaceMetaMap[interfaceIdent] = newGoInterfaceMeta(newMeta(typeSpec, gfm.path), interfaceIdent)
						}
					}
					return false // 只查找顶层为 interface 的节点
				}
				return true
			})
		}
	}
}

// -------------------------------- extractor --------------------------------

// SearchFileMeta 根据 文件名 搜索 文件 的 meta 数据
func (gpm *GoPackageMeta) SearchFileMeta(fileName string) *GoFileMeta[*ast.File] {
	return gpm.fileMetaMap[fileName]
}

// SearchVarMeta 根据 var 名称 搜索 var 的 meta 数据
func (gpm *GoPackageMeta) SearchVarMeta(varIdent string) *GoVarMeta[*ast.ValueSpec] {
	return gpm.varMetaMap[varIdent]
}

func (gpm *GoPackageMeta) SearchFuncMeta(funcIdent string) *GoFuncMeta[*ast.FuncDecl] {
	return gpm.funcMetaMap[funcIdent]
}

func (gpm *GoPackageMeta) SearchStructMeta(structName string) *GoStructMeta[*ast.TypeSpec] {
	return gpm.structMetaMap[structName]
}

func (gpm *GoPackageMeta) SearchInterfaceMeta(interfaceIdent string) *GoInterfaceMeta[*ast.TypeSpec] {
	return gpm.interfaceMetaMap[interfaceIdent]
}

func (gpm *GoPackageMeta) StructNames() []string {
	structNames := make([]string, 0)
	for _, gfm := range gpm.fileMetaMap {
		ast.Inspect(gfm.node, func(n ast.Node) bool {
			if IsTypeNode(n) {
				structNames = append(structNames, n.(*ast.TypeSpec).Name.String())
				return false
			}
			return true
		})
	}
	return structNames
}

func (gpm *GoPackageMeta) InterfaceNames() []string {
	interfaceNames := make([]string, 0)
	for _, gfm := range gpm.fileMetaMap {
		ast.Inspect(gfm.node, func(n ast.Node) bool {
			if IsTypeNode(n) {
				interfaceNames = append(interfaceNames, n.(*ast.TypeSpec).Name.String())
				return false
			}
			return true
		})
	}
	return interfaceNames
}

func (gpm *GoPackageMeta) FunctionNames() []string {
	functionNames := make([]string, 0)
	for _, gfm := range gpm.fileMetaMap {
		ast.Inspect(gfm.node, func(n ast.Node) bool {
			if IsFuncNode(n) {
				functionNames = append(functionNames, n.(*ast.FuncDecl).Name.String())
				return false
			}
			return true
		})
	}
	return functionNames
}

// func (gpm *GoPackageMeta) SearchFunctionMeta(functionName string) *GoFuncMeta {
// 	if gfm, has := gpm.funcMetaMap[functionName]; gfm != nil && has {
// 		return gfm
// 	}

// 	for _, gfm := range gpm.fileMetaMap {
// 		if gfm.node != nil {
// 			gfm := SearchGoFunctionMeta(gfm, functionName)
// 			if gfm != nil {
// 				gpm.funcMetaMap[gfm.FunctionName()] = gfm
// 				break
// 			}
// 		}
// 	}

// 	return gpm.funcMetaMap[functionName]
// }

// func (gpm *GoPackageMeta) MethodNames() map[string][]string {
// 	methodNames := make(map[string][]string, 0)

// 	for _, gfm := range gpm.fileMetaMap {
// 		ast.Inspect(gfm.node, func(n ast.Node) bool {
// 			if IsMethodNode(n) {
// 				decl := n.(*ast.FuncDecl)
// 				recvStructName, _ := extractMethodRecvStruct(decl)
// 				methodNames[recvStructName] = append(methodNames[recvStructName], decl.Name.String())
// 			}
// 			return true
// 		})
// 	}
// 	return methodNames
// }

// func (gpm *GoPackageMeta) SearchMethodMeta(structName, methodName string) *GoMethodMeta {
// 	if gsm, hasGSM := gpm.structMetaMap[structName]; gsm != nil && hasGSM {
// 		if gmm, hasGMM := gsm.methodDecl[methodName]; gmm != nil && hasGMM {
// 			return gmm
// 		}
// 	} else {
// 		gpm.SearchStructMeta(structName)
// 	}

// 	for _, gfm := range gpm.fileMetaMap {
// 		if gfm.node != nil {
// 			gmm := SearchGoMethodMeta(gfm, structName, methodName)
// 			if gmm != nil {
// 				gpm.structMetaMap[structName].methodDecl[methodName] = gmm
// 				break
// 			}
// 		}
// 	}

// 	return gpm.structMetaMap[structName].methodDecl[methodName]
// }

// -------------------------------- unit test --------------------------------

func (gpm *GoPackageMeta) Ident() string                                  { return gpm.ident }
func (gpm *GoPackageMeta) AbsolutePath() string                           { return gpm.absolutePath }
func (gpm *GoPackageMeta) ImportPath() string                             { return gpm.importPath }
func (gpm *GoPackageMeta) FileMetaMap() map[string]*GoFileMeta[*ast.File] { return gpm.fileMetaMap }
func (gpm *GoPackageMeta) VariableMetaMap() map[string]*GoVarMeta[*ast.ValueSpec] {
	return gpm.varMetaMap
}
func (gpm *GoPackageMeta) FuncMetaMap() map[string]*GoFuncMeta[*ast.FuncDecl] {
	return gpm.funcMetaMap
}
func (gpm *GoPackageMeta) StructMetaMap() map[string]*GoStructMeta[*ast.TypeSpec] {
	return gpm.structMetaMap
}
func (gpm *GoPackageMeta) InterfaceMetaMap() map[string]*GoInterfaceMeta[*ast.TypeSpec] {
	return gpm.interfaceMetaMap
}
func (gpm *GoPackageMeta) TypeConstraintsMetaMap() map[string]*GoInterfaceMeta[*ast.TypeSpec] {
	return gpm.typeConstraintsMetaMap
}

// -------------------------------- unit test --------------------------------
