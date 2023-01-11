package extractor

import (
	"fmt"
	"testing"
)

type compareGoProjectMeta struct {
	ProjectPath string
	ModuleName  string
	PackageMap  map[string]*compareGoPackageMeta
}

type compareGoPackageMeta struct {
	Name             string
	PkgPath          string
	ImportPath       string
	pkgFileMap       map[string]*compareGoFileMeta
	pkgStructMeta    map[string]*compareGoStructMeta
	pkgInterfaceMeta map[string]*compareGoInterfaceMeta
	pkgFunctionMeta  map[string]*compareGoFunctionMeta
}

type compareGoFileMeta struct {
	Name string
	Path string
}

type compareGoStructMeta struct {
	StructName       string
	Comments         []string
	StructMemberMeta map[string]*compareGoMemberMeta
	StructMethodMeta map[string]*compareGoMethodMeta
}

type compareGoInterfaceMeta struct {
	InterfaceName string
}

type compareGoFunctionMeta struct {
	FunctionName string
}

type compareGoMemberMeta struct {
	Name string
}

type compareGoMethodMeta struct {
	MethodName      string
	RecvStruct      string
	PointerReceiver bool
}

var (
	standardProjectRelPath       = "./testdata/standardProject"
	standardProjectIgnorePathMap = map[string]struct{}{
		standardProjectRelPath + "/vendor": {},
	}
	standardProjectAbsPath    = "d:\\Projects\\go-extractor\\testdata\\standardProject"
	standardProjectModuleName = "standardProject"
	standardProjectMeta       = &compareGoProjectMeta{
		ProjectPath: standardProjectAbsPath,
		ModuleName:  standardProjectModuleName,
		PackageMap: map[string]*compareGoPackageMeta{
			"main": {
				Name:    "main",
				PkgPath: standardProjectAbsPath + "\\cmd",
				pkgFileMap: map[string]*compareGoFileMeta{
					"main.go": {
						Name: "main.go",
						Path: standardProjectAbsPath + "\\cmd\\main.go",
					},
					"init.go": {
						Name: "init.go",
						Path: standardProjectAbsPath + "\\cmd\\init.go",
					},
				},
				pkgFunctionMeta: map[string]*compareGoFunctionMeta{
					"main": {
						FunctionName: "main",
					},
				},
			},
			standardProjectModuleName + "/pkg": {
				Name:       "pkg",
				PkgPath:    standardProjectAbsPath + "\\pkg",
				ImportPath: standardProjectModuleName + "/pkg",
				pkgFileMap: map[string]*compareGoFileMeta{
					"pkg.go": {
						Name: "pkg.go",
						Path: standardProjectAbsPath + "\\pkg\\pkg.go",
					},
				},
				pkgFunctionMeta: map[string]*compareGoFunctionMeta{
					"ExampleFunc": {
						FunctionName: "ExampleFunc",
					},
				},
			},
			standardProjectModuleName + "/pkg/pkgInterface": {
				Name:       "pkgInterface",
				PkgPath:    standardProjectAbsPath + "\\pkg\\interface",
				ImportPath: standardProjectModuleName + "/pkg/pkgInterface",
				pkgFileMap: map[string]*compareGoFileMeta{
					"interface.go": {
						Name: "interface.go",
						Path: standardProjectAbsPath + "\\pkg\\interface\\interface.go",
					},
				},
				pkgInterfaceMeta: map[string]*compareGoInterfaceMeta{
					"ExampleInterface": {
						InterfaceName: "ExampleInterface",
					},
				},
			},
			standardProjectModuleName + "/pkg/module": {
				Name:       "module",
				PkgPath:    standardProjectAbsPath + "\\pkg\\module",
				ImportPath: standardProjectModuleName + "/pkg/module",
				pkgFileMap: map[string]*compareGoFileMeta{
					"module.go": {
						Name: "module.go",
						Path: standardProjectAbsPath + "\\pkg\\module\\module.go",
					},
				},
				pkgStructMeta: map[string]*compareGoStructMeta{
					"ExampleStruct": {
						StructName: "ExampleStruct",
						Comments: []string{
							"// ExampleStruct this is an example struct",
							"// this is struct comment",
							"// this is another struct comment",
						},
						StructMemberMeta: map[string]*compareGoMemberMeta{
							"v": {
								Name: "v",
							},
						},
						StructMethodMeta: map[string]*compareGoMethodMeta{
							"ExampleFunc": {
								MethodName:      "ExampleFunc",
								RecvStruct:      "ExampleStruct",
								PointerReceiver: false,
							},
							"AnotherExampleFunc": {
								MethodName:      "AnotherExampleFunc",
								RecvStruct:      "ExampleStruct",
								PointerReceiver: true,
							},
							"V": {
								MethodName:      "V",
								RecvStruct:      "ExampleStruct",
								PointerReceiver: false,
							},
						},
					},
				},
				pkgFunctionMeta: map[string]*compareGoFunctionMeta{
					"NewExampleStruct": {
						FunctionName: "NewExampleStruct",
					},
					"ExampleFunc": {
						FunctionName: "ExampleFunc",
					},
				},
			},
		},
	}
)

func TestExtractGoProjectMeta(t *testing.T) {
	goProjectMeta, err := ExtractGoProjectMeta(standardProjectRelPath, standardProjectIgnorePathMap)
	if err != nil {
		panic(err)
	}

	checkProjectMeta(goProjectMeta, standardProjectMeta)

	for pkgName, _gpm := range standardProjectMeta.PackageMap {
		gpm, has := goProjectMeta.PackageMap[pkgName]
		if gpm == nil || !has {
			panic(pkgName)
		}
		checkPackageMeta(gpm, _gpm)

		for fileName, _gfm := range _gpm.pkgFileMap {
			gfm, has := gpm.pkgFileMap[fileName]
			if gfm == nil || !has {
				panic(fileName)
			}
			checkFileMeta(gfm, _gfm)
			gfm.OutputAST()
		}

		for structName, _gsm := range _gpm.pkgStructMeta {
			gpm.SearchStructMeta(structName)
			gsm, has := gpm.pkgStructDecl[structName]
			if gsm == nil || !has {
				panic(structName)
			}
			checkStructMeta(gsm, _gsm)

			for memberName, _gmm := range _gsm.StructMemberMeta {
				gsm.SearchMemberMeta(memberName)
				gmm, has := gsm.memberDecl[memberName]
				if gmm == nil || !has {
					panic(memberName)
				}
				checkMemberMeta(gmm, _gmm)
			}

			for methodName, _gmm := range _gsm.StructMethodMeta {
				gpm.SearchMethodMeta(structName, methodName)
				gmm, has := gsm.methodDecl[methodName]
				if gmm == nil || !has {
					panic(methodName)
				}
				checkMethodMeta(gmm, _gmm)
			}
		}

		for interfaceName, _gim := range _gpm.pkgInterfaceMeta {
			gpm.SearchInterfaceMeta(interfaceName)
			gim, has := gpm.pkgInterfaceDecl[interfaceName]
			if gim == nil || !has {
				panic(interfaceName)
			}
			checkInterfaceMeta(gim, _gim)
		}

		for funcName, _gfm := range _gpm.pkgFunctionMeta {
			gpm.SearchFunctionMeta(funcName)
			gfm, has := gpm.pkgFunctionDecl[funcName]
			if gfm == nil || !has {
				panic(funcName)
			}
			checkFunctionMeta(gfm, _gfm)
		}
	}
}

func checkProjectMeta(gpm *GoProjectMeta, _gpm *compareGoProjectMeta) {
	if gpm.ModuleName != _gpm.ModuleName {
		panic(gpm.ModuleName)
	}
	if gpm.ProjectPath != _gpm.ProjectPath {
		panic(gpm.ProjectPath)
	}
}

func checkPackageMeta(gpm *GoPackageMeta, _gpm *compareGoPackageMeta) {
	if gpm.Name != _gpm.Name {
		panic(gpm.Name)
	}
	if gpm.PkgPath != _gpm.PkgPath {
		panic(gpm.PkgPath)
	}
	if gpm.ImportPath != _gpm.ImportPath {
		panic(gpm.ImportPath)
	}
}

func checkFileMeta(gfm *GoFileMeta, _gfm *compareGoFileMeta) {
	if gfm.Name != _gfm.Name {
		panic(gfm.Name)
	}
	if gfm.Path != _gfm.Path {
		panic(gfm.Path)
	}
}

func checkStructMeta(gsm *GoStructMeta, _gsm *compareGoStructMeta) {
	if gsm.StructName() != _gsm.StructName {
		panic(gsm.StructName())
	}
	for _, _comment := range _gsm.Comments {
		for _, comment := range gsm.Comments() {
			if comment == _comment {
				goto NEXT_COMMENT
			}
		}
		panic(_comment)
	NEXT_COMMENT:
	}
}

func checkInterfaceMeta(gim *GoInterfaceMeta, _gim *compareGoInterfaceMeta) {
	if gim.InterfaceName() != _gim.InterfaceName {
		panic(gim.InterfaceName())
	}
}

func checkFunctionMeta(gfm *GoFunctionMeta, _gfm *compareGoFunctionMeta) {
	if gfm.FunctionName() != _gfm.FunctionName {
		panic(gfm.FunctionName())
	}
}

func checkMemberMeta(gmm *GoMemberMeta, _gmm *compareGoMemberMeta) {
	if gmm.MemberName() != _gmm.Name {
		panic(gmm.MemberName())
	}
}

func checkMethodMeta(gmm *GoMethodMeta, _gmm *compareGoMethodMeta) {
	if gmm.MethodName() != _gmm.MethodName {
		panic(gmm.MethodName())
	}
	if recvStruct, pointerReceiver := gmm.RecvStruct(); recvStruct != _gmm.RecvStruct || pointerReceiver != _gmm.PointerReceiver {
		panic(fmt.Sprintf("%v, %v", recvStruct, pointerReceiver))
	}
}
