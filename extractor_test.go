package extractor

import (
	"fmt"
	"sort"
	"testing"

	stpmap "github.com/Mericusta/go-stp/map"
	stpslice "github.com/Mericusta/go-stp/slice"
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
	Name    string
	Path    string
	PkgName string
}

type compareGoStructMeta struct {
	Expression       string
	StructName       string
	Doc              []string
	StructMemberMeta map[string]*compareGoMemberMeta
	StructMethodMeta map[string]*compareGoMethodMeta
}

type compareGoInterfaceMeta struct {
	InterfaceName string
}

type compareGoFunctionMeta struct {
	FunctionName string
	Doc          []string
	CallMeta     map[string][]*compareGoCallMeta
	// VarMeta map[string]
}

type compareGoMemberMeta struct {
	Expression string
	MemberName string
	Tag        string
	Comment    string
	Doc        []string
}

type compareGoMethodMeta struct {
	*compareGoFunctionMeta
	RecvStruct      string
	PointerReceiver bool
}

type compareGoCallMeta struct {
	Expression string
	Call       string
	From       string
	Args       []*compareGoArgMeta
}

type compareGoArgMeta struct {
	Expression string
	Head       *compareGoVariableMeta
}

type compareGoVariableMeta struct {
	Expression string
	Name       string
	Type       string
}

var (
	standardProjectRelPath       = "./testdata/standardProject"
	standardProjectIgnorePathMap = map[string]struct{}{
		standardProjectRelPath + "/vendor": {},
	}
	// standardProjectAbsPath    = "d:\\Projects\\go-extractor\\testdata\\standardProject"
	standardProjectAbsPath    = "d:\\Projects\\SGAME\\server-dev\\gameServer\\game_server\\pkg\\github.com\\Mericusta\\go-extractor\\testdata\\standardProject"
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
						Name:    "main.go",
						Path:    standardProjectAbsPath + "\\cmd\\main.go",
						PkgName: "main",
					},
					"init.go": {
						Name:    "init.go",
						Path:    standardProjectAbsPath + "\\cmd\\init.go",
						PkgName: "main",
					},
				},
				pkgFunctionMeta: map[string]*compareGoFunctionMeta{
					"main": {
						FunctionName: "main",
						CallMeta: map[string][]*compareGoCallMeta{
							"pkg.ExampleFunc": {
								{
									Call: "ExampleFunc",
									From: "pkg",
									Args: []*compareGoArgMeta{
										{
											Expression: `module.NewExampleStruct(10)`,
											Head: &compareGoVariableMeta{
												Expression: `module`,
												Name:       "module",
												Type:       `"standardProject/pkg/module"`,
											},
										},
									}},
							},
							"module.ExampleFunc": {
								{
									Call: "ExampleFunc",
									From: "module",
									Args: []*compareGoArgMeta{
										{
											Expression: `module.NewExampleStruct(11)`,
											Head: &compareGoVariableMeta{
												Expression: `module`,
												Name:       "module",
												Type:       `"standardProject/pkg/module"`,
											},
										},
									},
								},
							},
							"module.NewExampleStruct": {
								{
									Call: "NewExampleStruct",
									From: "module",
									Args: []*compareGoArgMeta{
										{
											Expression: `10`,
											Head: &compareGoVariableMeta{
												Expression: `10`,
												Name:       "10",
												Type:       `10`,
											},
										},
									},
								},
								{
									Call: "NewExampleStruct",
									From: "module",
									Args: []*compareGoArgMeta{
										{
											Expression: `11`,
											Head: &compareGoVariableMeta{
												Expression: `11`,
												Name:       "11",
												Type:       `11`,
											},
										},
									},
								},
							},
							"Init": {
								{
									Call: "Init",
								},
							},
						},
					},
					"Init": {
						FunctionName: "Init",
					},
				},
			},
			standardProjectModuleName + "/pkg": {
				Name:       "pkg",
				PkgPath:    standardProjectAbsPath + "\\pkg",
				ImportPath: standardProjectModuleName + "/pkg",
				pkgFileMap: map[string]*compareGoFileMeta{
					"pkg.go": {
						Name:    "pkg.go",
						Path:    standardProjectAbsPath + "\\pkg\\pkg.go",
						PkgName: "pkg",
					},
				},
				pkgFunctionMeta: map[string]*compareGoFunctionMeta{
					"ExampleFunc": {
						FunctionName: "ExampleFunc",
						Doc: []string{
							"// ExampleFunc this is example function",
						},
						CallMeta: map[string][]*compareGoCallMeta{
							"fmt.Println": {
								{
									From: "fmt",
									Call: "Println",
									Args: []*compareGoArgMeta{
										{
											Expression: `"pkg.ExampleFunc, Hello go-extractor"`,
											Head: &compareGoVariableMeta{
												Expression: `"pkg.ExampleFunc, Hello go-extractor"`,
												Name:       `"pkg.ExampleFunc, Hello go-extractor"`,
												Type:       `"pkg.ExampleFunc, Hello go-extractor"`,
											},
										},
										{
											Expression: `s.V()`,
											Head: &compareGoVariableMeta{
												Expression: `s *module.ExampleStruct`,
												Name:       "s",
												Type:       `*module.ExampleStruct`,
											},
										},
									},
								},
							},
							"s.V": {
								{
									From: "s",
									Call: "V",
								},
							},
						},
					},
					"NoDocExampleFunc": {
						FunctionName: "NoDocExampleFunc",
						CallMeta: map[string][]*compareGoCallMeta{
							"fmt.Println": {
								{
									From: "fmt",
									Call: "Println",
									Args: []*compareGoArgMeta{
										{
											Expression: `"pkg.NoDocExampleFunc, Hello go-extractor"`,
											Head: &compareGoVariableMeta{
												Expression: `"pkg.NoDocExampleFunc, Hello go-extractor"`,
												Name:       `"pkg.NoDocExampleFunc, Hello go-extractor"`,
												Type:       `"pkg.NoDocExampleFunc, Hello go-extractor"`,
											},
										},
										{
											Expression: `s.V()`,
											Head: &compareGoVariableMeta{
												Expression: `s *module.ExampleStruct`,
												Name:       "s",
												Type:       `*module.ExampleStruct`,
											},
										},
									},
								},
							},
							"s.V": {
								{
									From: "s",
									Call: "V",
								},
							},
						},
					},
					"OneLineDocExampleFunc": {
						FunctionName: "OneLineDocExampleFunc",
						CallMeta: map[string][]*compareGoCallMeta{
							"fmt.Println": {
								{
									From: "fmt",
									Call: "Println",
									Args: []*compareGoArgMeta{
										{
											Expression: `"pkg.OneLineDocExampleFunc, Hello go-extractor"`,
											Head: &compareGoVariableMeta{
												Expression: `"pkg.OneLineDocExampleFunc, Hello go-extractor"`,
												Name:       `"pkg.OneLineDocExampleFunc, Hello go-extractor"`,
												Type:       `"pkg.OneLineDocExampleFunc, Hello go-extractor"`,
											},
										},
										{
											Expression: `s.V()`,
											Head: &compareGoVariableMeta{
												Expression: `s *module.ExampleStruct`,
												Name:       "s",
												Type:       `*module.ExampleStruct`,
											},
										},
									},
								},
							},
							"s.V": {
								{
									From: "s",
									Call: "V",
								},
							},
						},
					},
					"ImportSelectorFunc": {
						FunctionName: "ImportSelectorFunc",
						CallMeta: map[string][]*compareGoCallMeta{
							"fmt.Println": {
								&compareGoCallMeta{
									From: "fmt",
									Call: "Println",
									Args: []*compareGoArgMeta{
										{
											Expression: `"pkg.ImportSelectorFunc, Hello go-extractor"`,
											Head: &compareGoVariableMeta{
												Expression: `"pkg.ImportSelectorFunc, Hello go-extractor"`,
												Name:       `"pkg.ImportSelectorFunc, Hello go-extractor"`,
												Type:       `"pkg.ImportSelectorFunc, Hello go-extractor"`,
											},
										},
										{
											Expression: `module.NewExampleStruct(s.V()).Sub().ParentStruct.P`,
											Head: &compareGoVariableMeta{
												Expression: `module`,
												Name:       "module",
												Type:       `"standardProject/pkg/module"`,
											},
										},
									},
								},
							},
							"module.NewExampleStruct": {
								{
									From: "module",
									Call: "NewExampleStruct",
									Args: []*compareGoArgMeta{
										{
											Expression: `s.V()`,
											Head: &compareGoVariableMeta{
												Expression: `s *module.ExampleStruct`,
												Name:       "s",
												Type:       `*module.ExampleStruct`,
											},
										},
									},
								},
							},
							"module.NewExampleStruct(s.V()).Sub": {
								{
									From: "module.NewExampleStruct(s.V())",
									Call: "Sub",
								},
							},
							"s.V": {
								{
									From: "s",
									Call: "V",
								},
							},
						},
					},
				},
			},
			standardProjectModuleName + "/pkg/pkgInterface": {
				Name:       "pkgInterface",
				PkgPath:    standardProjectAbsPath + "\\pkg\\interface",
				ImportPath: standardProjectModuleName + "/pkg/pkgInterface",
				pkgFileMap: map[string]*compareGoFileMeta{
					"interface.go": {
						Name:    "interface.go",
						Path:    standardProjectAbsPath + "\\pkg\\interface\\interface.go",
						PkgName: "pkgInterface",
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
						Name:    "module.go",
						Path:    standardProjectAbsPath + "\\pkg\\module\\module.go",
						PkgName: "module",
					},
				},
				pkgStructMeta: map[string]*compareGoStructMeta{
					"ParentStruct": {
						StructName: "ParentStruct",
						StructMemberMeta: map[string]*compareGoMemberMeta{
							"p": {
								MemberName: "p",
								Comment:    "// parent value",
							},
						},
					},
					"ExampleStruct": {
						StructName: "ExampleStruct",
						Doc: []string{
							"// ExampleStruct this is an example struct",
							"// this is struct comment",
							"// this is another struct comment",
						},
						StructMemberMeta: map[string]*compareGoMemberMeta{
							"ParentStruct": {
								MemberName: "ParentStruct",
								Comment:    "// parent struct",
							},
							"v": {
								MemberName: "v",
								Tag:        "`ast:init,default=1`",
								Doc: []string{
									"// v this is member doc line1",
									"// v this is member doc line2",
								},
								Comment: "// this is member single comment line",
							},
							"sub": {
								MemberName: "sub",
							},
						},
						StructMethodMeta: map[string]*compareGoMethodMeta{
							"ExampleFunc": {
								compareGoFunctionMeta: &compareGoFunctionMeta{
									FunctionName: "ExampleFunc",
									CallMeta: map[string][]*compareGoCallMeta{
										"NewExampleStruct": {
											{
												Call: "NewExampleStruct",
												Args: []*compareGoArgMeta{
													{
														Expression: `v`,
														Head: &compareGoVariableMeta{
															Expression: `v int`,
															Name:       "v",
															Type:       `int`,
														},
													},
												},
											},
											{
												Call: "NewExampleStruct",
												Args: []*compareGoArgMeta{
													{
														Expression: `nes.Sub().ParentStruct.P()`,
														Head: &compareGoVariableMeta{
															Expression: `nes := NewExampleStruct(v)`,
															Name:       "nes",
															Type:       `NewExampleStruct(v)`,
														},
													},
												},
											},
										},
										"es.DoubleReturnFunc": {
											{
												From: "es",
												Call: "DoubleReturnFunc",
											},
										},
										"nes.DoubleReturnFunc": {
											{
												From: "nes",
												Call: "DoubleReturnFunc",
											},
										},
										"fmt.Println": {
											{
												From: "fmt",
												Call: "Println",
												Args: []*compareGoArgMeta{
													{
														Expression: `"module.ExampleStruct.ExampleFunc Hello go-extractor"`,
														Head: &compareGoVariableMeta{
															Expression: `"module.ExampleStruct.ExampleFunc Hello go-extractor"`,
															Name:       `"module.ExampleStruct.ExampleFunc Hello go-extractor"`,
															Type:       `"module.ExampleStruct.ExampleFunc Hello go-extractor"`,
														},
													},
													{
														Expression: `es`,
														Head: &compareGoVariableMeta{
															Expression: `es ExampleStruct`,
															Name:       "es",
															Type:       `ExampleStruct`,
														},
													},
													{
														Expression: `es.v`,
														Head: &compareGoVariableMeta{
															Expression: `es ExampleStruct`,
															Name:       "es",
															Type:       `ExampleStruct`,
														},
													},
													{
														Expression: `es.V()`,
														Head: &compareGoVariableMeta{
															Expression: `es ExampleStruct`,
															Name:       "es",
															Type:       `ExampleStruct`,
														},
													},
													{
														Expression: `esP`,
														Head: &compareGoVariableMeta{
															Expression: `esP, esSubV := es.DoubleReturnFunc()`,
															Name:       "esP",
															Type:       `es.DoubleReturnFunc()`,
														},
													},
													{
														Expression: `esSubV`,
														Head: &compareGoVariableMeta{
															Expression: `esP, esSubV := es.DoubleReturnFunc()`,
															Name:       "esSubV",
															Type:       `es.DoubleReturnFunc()`,
														},
													},
													{
														Expression: `nes`,
														Head: &compareGoVariableMeta{
															Expression: `nes := NewExampleStruct(v)`,
															Name:       "nes",
															Type:       `NewExampleStruct(v)`,
														},
													},
													{
														Expression: `nes.v`,
														Head: &compareGoVariableMeta{
															Expression: `nes := NewExampleStruct(v)`,
															Name:       "nes",
															Type:       `NewExampleStruct(v)`,
														},
													},
													{
														Expression: `nes.V()`,
														Head: &compareGoVariableMeta{
															Expression: `nes := NewExampleStruct(v)`,
															Name:       "nes",
															Type:       `NewExampleStruct(v)`,
														},
													},
													{
														Expression: `nesP`,
														Head: &compareGoVariableMeta{
															Expression: `nesP, nesSubV := nes.DoubleReturnFunc()`,
															Name:       "nesP",
															Type:       `nes.DoubleReturnFunc()`,
														},
													},
													{
														Expression: `nesSubV`,
														Head: &compareGoVariableMeta{
															Expression: `nesP, nesSubV := nes.DoubleReturnFunc()`,
															Name:       "nesSubV",
															Type:       `nes.DoubleReturnFunc()`,
														},
													},
													{
														Expression: `globalExampleStruct`,
														Head: &compareGoVariableMeta{
															Expression: `var globalExampleStruct *ExampleStruct`,
															Name:       "globalExampleStruct",
															Type:       `*ExampleStruct`,
														},
													},
													{
														Expression: `NewExampleStruct(nes.Sub().ParentStruct.P())`,
														Head: &compareGoVariableMeta{
															Expression: `NewExampleStruct`,
															Name:       "NewExampleStruct",
															Type:       `*ExampleStruct`,
														},
													},
												},
											},
										},
										"es.V": {
											{
												From: "es",
												Call: "V",
											},
										},
										"nes.V": {
											{
												From: "nes",
												Call: "V",
											},
										},
										"nes.Sub": {
											{
												From: "nes",
												Call: "Sub",
											},
										},
										"nes.Sub().ParentStruct.P": {
											{
												From: "nes.Sub().ParentStruct",
												Call: "P",
											},
										},
									},
								},
								RecvStruct:      "ExampleStruct",
								PointerReceiver: false,
							},
							"ExampleFuncWithPointerReceiver": {
								compareGoFunctionMeta: &compareGoFunctionMeta{
									FunctionName: "ExampleFuncWithPointerReceiver",
									CallMeta: map[string][]*compareGoCallMeta{
										"fmt.Println": {
											{
												From: "fmt",
												Call: "Println",
												Args: []*compareGoArgMeta{
													{
														Expression: `"module.ExampleStruct.ExampleFuncWithPointerReceiver Hello go-extractor"`,
														Head: &compareGoVariableMeta{
															Expression: `"module.ExampleStruct.ExampleFuncWithPointerReceiver Hello go-extractor"`,
															Name:       `"module.ExampleStruct.ExampleFuncWithPointerReceiver Hello go-extractor"`,
															Type:       `"module.ExampleStruct.ExampleFuncWithPointerReceiver Hello go-extractor"`,
														},
													},
												},
											},
										},
									},
								},
								RecvStruct:      "ExampleStruct",
								PointerReceiver: true,
							},
							"DoubleReturnFunc": {
								compareGoFunctionMeta: &compareGoFunctionMeta{
									FunctionName: "DoubleReturnFunc",
									CallMeta: map[string][]*compareGoCallMeta{
										"es.P": {
											{
												From: "es",
												Call: "P",
											},
										},
										"es.sub.V": {
											{
												From: "es.sub",
												Call: "V",
											},
										},
									},
								},
								RecvStruct:      "ExampleStruct",
								PointerReceiver: true,
							},
							"V": {
								compareGoFunctionMeta: &compareGoFunctionMeta{
									FunctionName: "V",
								},
								RecvStruct:      "ExampleStruct",
								PointerReceiver: false,
							},
						},
					},
				},
				pkgFunctionMeta: map[string]*compareGoFunctionMeta{
					"NewExampleStruct": {
						FunctionName: "NewExampleStruct",
						Doc: []string{
							"// NewExampleStruct this is new example struct",
							"// @param           value",
							"// @return          pointer to ExampleStruct",
						},
						CallMeta: map[string][]*compareGoCallMeta{
							"random.Intn": {
								{
									From: "random",
									Call: "Intn",
									Args: []*compareGoArgMeta{
										{
											Expression: `v`,
											Head: &compareGoVariableMeta{
												Expression: `v int`,
												Name:       "v",
												Type:       `int`,
											},
										},
									},
								},
							},
						},
					},
					"ExampleFunc": {
						FunctionName: "ExampleFunc",
						CallMeta: map[string][]*compareGoCallMeta{
							"s.ExampleFunc": {
								{
									From: "s",
									Call: "ExampleFunc",
									Args: []*compareGoArgMeta{
										{
											Expression: `s.v`,
											Head: &compareGoVariableMeta{
												Expression: `s *ExampleStruct`,
												Name:       "s",
												Type:       `*ExampleStruct`,
											},
										},
									},
								},
							},
						},
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
		gpm := goProjectMeta.SearchPackageMeta(pkgName)
		if gpm == nil {
			Panic(gpm, _gpm)
		}
		checkPackageMeta(gpm, _gpm)

		for _fileName, _gfm := range _gpm.pkgFileMap {
			gfm := gpm.SearchFileMeta(_fileName)
			if gfm == nil {
				Panic(gfm, _gfm)
			}
			checkFileMeta(gfm, _gfm)
		}

		for _structName, _gsm := range _gpm.pkgStructMeta {
			gsm := gpm.SearchStructMeta(_structName)
			if gsm == nil {
				Panic(gsm, _gsm)
			}
			checkStructMeta(gsm, _gsm)

			for memberName, _gmm := range _gsm.StructMemberMeta {
				gmm := gsm.SearchMemberMeta(memberName)
				if gmm == nil {
					Panic(gmm, _gmm)
				}
				checkMemberMeta(gmm, _gmm)
			}

			for methodName, _gmm := range _gsm.StructMethodMeta {
				gmm := gpm.SearchMethodMeta(_structName, methodName)
				if gmm == nil {
					Panic(gmm, _gmm)
				}
				checkMethodMeta(gmm, _gmm)
			}
		}

		for interfaceName, _gim := range _gpm.pkgInterfaceMeta {
			gim := gpm.SearchInterfaceMeta(interfaceName)
			if gim == nil {
				Panic(gim, _gim)
			}
			checkInterfaceMeta(gim, _gim)
		}

		for funcName, _gfm := range _gpm.pkgFunctionMeta {
			gfm := gpm.SearchFunctionMeta(funcName)
			if gfm == nil {
				Panic(gfm, _gfm)
			}
			checkFunctionMeta(gfm, _gfm)
		}
	}
}

func Panic(v, c any) {
	panic(fmt.Sprintf("%+v != %+v", v, c))
}

func checkProjectMeta(gpm *GoProjectMeta, _gpm *compareGoProjectMeta) {
	// basic
	if gpm.ProjectPath() != _gpm.ProjectPath {
		Panic(gpm.ProjectPath(), _gpm.ProjectPath)
	}
	if gpm.ModuleName() != _gpm.ModuleName {
		Panic(gpm.ModuleName(), _gpm.ModuleName)
	}

	// packages
	packages := gpm.Packages()
	sort.Strings(packages)
	_packages := stpmap.Key(_gpm.PackageMap)
	sort.Strings(_packages)
	if !stpslice.Compare(packages, _packages) {
		Panic(packages, _packages)
	}
}

func checkPackageMeta(gpm *GoPackageMeta, _gpm *compareGoPackageMeta) {
	// basic
	if gpm.Name() != _gpm.Name {
		Panic(gpm.Name(), _gpm.Name)
	}
	if gpm.PkgPath() != _gpm.PkgPath {
		Panic(gpm.PkgPath(), _gpm.PkgPath)
	}
	if gpm.ImportPath() != _gpm.ImportPath {
		Panic(gpm.ImportPath(), _gpm.ImportPath)
	}

	// file
	fileNames := gpm.FileNames()
	sort.Strings(fileNames)
	_fileNames := stpmap.Key(_gpm.pkgFileMap)
	sort.Strings(_fileNames)
	if !stpslice.Compare(fileNames, _fileNames) {
		Panic(fileNames, _fileNames)
	}

	// struct
	structNames := gpm.StructNames()
	sort.Strings(structNames)
	_structNames := stpmap.Key(_gpm.pkgStructMeta)
	sort.Strings(_structNames)
	if !stpslice.Compare(structNames, _structNames) {
		Panic(structNames, _structNames)
	}

	// interface
	interfaceNames := gpm.InterfaceNames()
	sort.Strings(interfaceNames)
	_interfaceNames := stpmap.Key(_gpm.pkgInterfaceMeta)
	sort.Strings(_interfaceNames)
	if !stpslice.Compare(interfaceNames, _interfaceNames) {
		Panic(interfaceNames, _interfaceNames)
	}

	// function
	functionNames := gpm.FunctionNames()
	sort.Strings(functionNames)
	_functionNames := stpmap.Key(_gpm.pkgFunctionMeta)
	sort.Strings(_functionNames)
	if !stpslice.Compare(functionNames, _functionNames) {
		Panic(functionNames, _functionNames)
	}
}

func checkFileMeta(gfm *GoFileMeta, _gfm *compareGoFileMeta) {
	if gfm.Name() != _gfm.Name {
		Panic(gfm.Name(), _gfm.Name)
	}
	if gfm.Path() != _gfm.Path {
		Panic(gfm.Path(), _gfm.Path)
	}
	if gfm.PkgName() != _gfm.PkgName {
		Panic(gfm.PkgName(), _gfm.PkgName)
	}
	gfm.OutputAST()
}

func checkStructMeta(gsm *GoStructMeta, _gsm *compareGoStructMeta) {
	// basic
	if gsm.StructName() != _gsm.StructName {
		Panic(gsm.StructName(), _gsm.StructName)
	}
	stpslice.Compare(gsm.Doc(), _gsm.Doc)

	// member
	memberNames := gsm.Members()
	sort.Strings(memberNames)
	_memberNames := stpmap.Key(_gsm.StructMemberMeta)
	sort.Strings(_memberNames)
	if !stpslice.Compare(memberNames, _memberNames) {
		Panic(memberNames, _memberNames)
	}
}

func checkInterfaceMeta(gim *GoInterfaceMeta, _gim *compareGoInterfaceMeta) {
	// basic
	if gim.InterfaceName() != _gim.InterfaceName {
		Panic(gim.InterfaceName(), _gim.InterfaceName)
	}
}

func checkFunctionMeta(gfm *GoFunctionMeta, _gfm *compareGoFunctionMeta) {
	// basic
	if gfm.FunctionName() != _gfm.FunctionName {
		Panic(gfm.FunctionName(), _gfm.FunctionName)
	}
	stpslice.Compare(gfm.Doc(), _gfm.Doc)

	// call
	calls := stpmap.Key(gfm.Calls())
	sort.Strings(calls)
	_calls := stpmap.Key(_gfm.CallMeta)
	sort.Strings(_calls)
	if !stpslice.Compare(calls, _calls) {
		gfm.Calls()
		Panic(calls, _calls)
	}

	for _call, _gcmSlice := range _gfm.CallMeta {
		gcmSlice := gfm.SearchCallMeta(_call)
		if len(gcmSlice) != len(_gcmSlice) {
			Panic(gcmSlice, _gcmSlice)
		}
		for index, _gcm := range _gcmSlice {
			gcm := gcmSlice[index]
			checkCallMeta(gcm, _gcm)
		}
	}
}

func checkMemberMeta(gmm *GoMemberMeta, _gmm *compareGoMemberMeta) {
	// basic
	if gmm.MemberName() != _gmm.MemberName {
		Panic(gmm.MemberName(), _gmm.MemberName)
	}
	if gmm.Tag() != _gmm.Tag {
		Panic(gmm.Tag(), _gmm.Tag)
	}
	if gmm.Comment() != _gmm.Comment {
		Panic(gmm.Comment(), _gmm.Comment)
	}
	stpslice.Compare(gmm.Doc(), _gmm.Doc)
}

func checkMethodMeta(gmm *GoMethodMeta, _gmm *compareGoMethodMeta) {
	// basic
	recvStruct, pointerReceiver := gmm.RecvStruct()
	if recvStruct != _gmm.RecvStruct {
		Panic(recvStruct, _gmm.RecvStruct)
	}
	if pointerReceiver != _gmm.PointerReceiver {
		Panic(pointerReceiver, _gmm.PointerReceiver)
	}

	// function
	checkFunctionMeta(gmm.GoFunctionMeta, _gmm.compareGoFunctionMeta)
}

func checkCallMeta(gcm *GoCallMeta, _gcm *compareGoCallMeta) {
	// basic
	if (gcm != nil) != (_gcm != nil) {
		Panic(gcm, _gcm)
	}
	if gcm == nil {
		return
	}
	if gcm.From() != _gcm.From {
		Panic(gcm.From(), _gcm.From)
	}
	if gcm.Call() != _gcm.Call {
		Panic(gcm.Call(), _gcm.Call)
	}

	// args
	if len(gcm.Args()) != len(_gcm.Args) {
		Panic(gcm.Args(), _gcm.Args)
	}
	args := gcm.Args()
	for index, _arg := range _gcm.Args {
		arg := args[index]
		if arg.Expression() != _arg.Expression {
			Panic(arg.Expression(), _arg.Expression)
		}
		checkVariableMeta(arg.Head(), _arg.Head)
	}
}

func checkVariableMeta(gvm *GoVariableMeta, _gvm *compareGoVariableMeta) {
	if gvm.Name() != _gvm.Name {
		Panic(gvm.Name(), _gvm.Name)
	}
	if gvm.Type() != _gvm.Type {
		Panic(gvm.Type(), _gvm.Type)
	}
}

type replaceFunctionDoc struct {
	originDoc      []string
	replaceDoc     []string
	originContent  string
	replaceContent string
}

var (
	replaceDoc = map[string]map[string]*replaceFunctionDoc{
		standardProjectModuleName + "/pkg": {
			"ExampleFunc": {
				originDoc: []string{
					"// ExampleFunc this is example function",
				},
				replaceDoc: []string{
					"// ExampleFunc this is example function doc after replace, line 1",
				},
				originContent: `// ExampleFunc this is example function
func ExampleFunc(s *module.ExampleStruct) {
	fmt.Println("pkg.ExampleFunc, Hello go-extractor,", s.V())
}`,
				replaceContent: `// ExampleFunc this is example function doc after replace, line 1
func ExampleFunc(s *module.ExampleStruct) {
	fmt.Println("pkg.ExampleFunc, Hello go-extractor,", s.V())
}`,
			},
			"NoDocExampleFunc": {
				originDoc: nil,
				replaceDoc: []string{
					"// NoDocExampleFunc this is no-doc example function doc after replace, line 1",
				},
				originContent: `func NoDocExampleFunc(s *module.ExampleStruct) {
	fmt.Println("pkg.ExampleFunc, Hello go-extractor,", s.V())
}`,
				replaceContent: `// NoDocExampleFunc this is no-doc example function doc after replace, line 1
func NoDocExampleFunc(s *module.ExampleStruct) {
	fmt.Println("pkg.ExampleFunc, Hello go-extractor,", s.V())
}`,
			},
			"OneLineDocExampleFunc": {
				originDoc: []string{
					"// OneLineDocExampleFunc this is one-line-doc example function",
				},
				replaceDoc: []string{
					"// OneLineDocExampleFunc this is one-line-doc example function doc after replace, line 1",
					"// OneLineDocExampleFunc this is one-line-doc example function doc after replace, line 2",
				},
				originContent: `// OneLineDocExampleFunc this is one-line-doc example function
func OneLineDocExampleFunc(s *module.ExampleStruct) {
	fmt.Println("pkg.ExampleFunc, Hello go-extractor,", s.V())
}`,
				replaceContent: `// OneLineDocExampleFunc this is one-line-doc example function doc after replace, line 1
// OneLineDocExampleFunc this is one-line-doc example function doc after replace, line 2
func OneLineDocExampleFunc(s *module.ExampleStruct) {
	fmt.Println("pkg.ExampleFunc, Hello go-extractor,", s.V())
}`,
			},
		},
	}
)

func TestReplaceGoProjectMeta(t *testing.T) {
	// goProjectMeta, err := ExtractGoProjectMeta(standardProjectRelPath, standardProjectIgnorePathMap)
	// if err != nil {
	// 	panic(err)
	// }

	// for pkgName, replaceFunctionDoc := range replaceDoc {
	// 	gpm, has := goProjectMeta.PackageMap[pkgName]
	// 	if gpm == nil || !has {
	// 		panic(pkgName)
	// 	}
	// 	for funcName, _replace := range replaceFunctionDoc {
	// 		gpm.SearchFunctionMeta(funcName)
	// 		gfm, has := gpm.pkgFunctionDecl[funcName]
	// 		if gfm == nil || !has {
	// 			panic(funcName)
	// 		}
	// 		checkDoc(gfm.Doc(), _replace.originDoc)
	// 		originContent, replaceContent, err := gfm.ReplaceFunctionDoc(_replace.replaceDoc)
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		if originContent != _replace.originContent {
	// 			panic(originContent)
	// 		}
	// 		if replaceContent != _replace.replaceContent {
	// 			panic(replaceContent)
	// 		}
	// 	}
	// }
}

var (
	compareGoCallMetaSlice = []*compareGoCallMeta{
		// {
		// 	Expression: `HaveReadGP(1)`,
		// 	Call:       "HaveReadGP",
		// 	Args: []*compareArgMeta{
		// 		{
		// 			Value: int32(1),
		// 		},
		// 	},
		// },
		// {
		// 	Expression: `GetPlayerLevel()`,
		// 	Call:       "GetPlayerLevel",
		// 	Args:       nil,
		// },
		// {
		// 	Expression: `HaveReadGP("gamephone")`,
		// 	Call:       "HaveReadGP",
		// 	Args: []*compareArgMeta{
		// 		{
		// 			Value: `"gamephone"`},
		// 	},
		// },
		// {
		// 	Expression: `HaveReadGP("gamephone",1,"remove")`,
		// 	Call:       "HaveReadGP",
		// 	Args: []*compareArgMeta{
		// 		{
		// 			Value: `"gamephone"`,
		// 		},
		// 		{
		// 			Value: int32(1),
		// 		},
		// 		{
		// 			Value: `"remove"`,
		// 		},
		// 	},
		// },
		// TODO: syntax tree
		// {
		// 	Expression: `HaveReadGP(1) && HaveReadGP(2)`,
		// 	Call:       "HaveReadGP",
		// 	Args: []interface{}{
		// 		int32(1),
		// 	},
		// },
		// {
		// 	Expression: `HaveReadGP(1, HaveReadGP(2)) && HaveReadGP(3)`,
		// 	Call:       "HaveReadGP",
		// 	Args: []interface{}{
		// 		int32(1),
		// 	},
		// },
		// TODO: func wrapper
		// {
		// 	Expression: `func() { HaveReadGP(1,HaveReadGP(1));HaveReadGP(1,HaveReadGP(1)) }`,
		// 	Call:       "HaveReadGP",
		// 	Args:       []interface{}{int32(1)},
		// },
		// TODO: not support
		// {
		// 	Expression: `HaveReadGP(1)HaveReadGP(2)HaveReadGP(3)`,
		// 	Call:       "HaveReadGP",
		// 	Args:       []interface{}{int32(1)},
		// },
		// {
		// 	Expression: `HaveReadGP(1) and HaveReadGP(2) and HaveReadGP(3)``,
		// 	Call:       "HaveReadGP",
		// 	Args:       []interface{}{int32(1)},
		// },
	}
)

func TestParseGoCallMeta(t *testing.T) {
	// for _, _gcm := range compareGoCallMetaSlice {
	// 	gcm := ParseGoCallMeta(_gcm.Expression)
	// 	gcm.PrintAST()

	// 	if gcm.Expression() != _gcm.Expression {
	// 		panic(gcm.Expression())
	// 	}
	// 	if gcm.Call() != _gcm.Call {
	// 		panic(gcm.Call())
	// 	}

	// 	if len(gcm.Args()) != len(_gcm.Args) {
	// 		panic(len(gcm.Args()))
	// 	}
	// 	for _, _arg := range _gcm.Args {
	// 		for _, arg := range gcm.Args() {
	// 			if reflect.DeepEqual(arg, _arg) {
	// 				goto NEXT_PARAM
	// 			}
	// 		}
	// 		panic(_arg)
	// 	NEXT_PARAM:
	// 	}
	// }
}
