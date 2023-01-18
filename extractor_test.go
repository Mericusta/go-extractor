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
	CallMeta     map[string][]*compareCallMeta
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

type compareCallMeta struct {
	Expression string
	Call       string
	From       string
	Args       []*compareArgMeta
}

type compareArgMeta struct {
	Expression string
	Head       string
	// ArgType    int
	// Arg        string
	// From       string
	// Value      interface{}
	// CallMeta   *compareCallMeta
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
						CallMeta: map[string][]*compareCallMeta{
							"fmt.Println": {
								{
									From: "fmt",
									Call: "Println",
									Args: []*compareArgMeta{
										{
											Expression: `"pkg.ExampleFunc, Hello go-extractor"`,
											Head:       `"pkg.ExampleFunc, Hello go-extractor"`,
										},
										{
											Expression: `s.V()`,
											Head:       "s",
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
						CallMeta: map[string][]*compareCallMeta{
							"fmt.Println": {
								{
									From: "fmt",
									Call: "Println",
									Args: []*compareArgMeta{
										{
											Expression: `"pkg.NoDocExampleFunc, Hello go-extractor"`,
											Head:       `"pkg.NoDocExampleFunc, Hello go-extractor"`,
										},
										{
											Expression: `s.V()`,
											Head:       "s",
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
						CallMeta: map[string][]*compareCallMeta{
							"fmt.Println": {
								{
									From: "fmt",
									Call: "Println",
									Args: []*compareArgMeta{
										{
											Expression: `"pkg.OneLineDocExampleFunc, Hello go-extractor"`,
											Head:       `"pkg.OneLineDocExampleFunc, Hello go-extractor"`,
										},
										{
											Expression: `s.V()`,
											Head:       "s",
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
						CallMeta: map[string][]*compareCallMeta{
							"fmt.Println": {
								&compareCallMeta{
									From: "fmt",
									Call: "Println",
									Args: []*compareArgMeta{
										{
											Expression: `"pkg.ImportSelectorFunc, Hello go-extractor"`,
											Head:       `"pkg.ImportSelectorFunc, Hello go-extractor"`,
										},
										{
											Expression: `module.NewExampleStruct(s.V()).Sub().ParentStruct.P`,
											Head:       "module",
										},
									},
								},
							},
							"module.NewExampleStruct": {
								{
									From: "module",
									Call: "NewExampleStruct",
									Args: []*compareArgMeta{
										{
											Expression: `s.V()`,
											Head:       "s",
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
							"P": {
								MemberName: "P",
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
									CallMeta: map[string][]*compareCallMeta{
										"NewExampleStruct": {
											{
												Call: "NewExampleStruct",
												Args: []*compareArgMeta{
													{
														Expression: `v`,
														Head:       "v",
													},
												},
											},
											{
												Call: "NewExampleStruct",
												Args: []*compareArgMeta{
													{
														Expression: `nes.sub.V()`,
														Head:       "nes",
													},
												},
											},
										},
										"fmt.Println": {
											{
												From: "fmt",
												Call: "Println",
												Args: []*compareArgMeta{
													{
														Expression: `"module.ExampleStruct.ExampleFunc Hello go-extractor"`,
														Head:       `"module.ExampleStruct.ExampleFunc Hello go-extractor"`,
													},
													{
														Expression: `es`,
														Head:       "es",
													},
													{
														Expression: `es.v`,
														Head:       "es",
													},
													{
														Expression: `es.V()`,
														Head:       "es",
													},
													{
														Expression: `nes`,
														Head:       "nes",
													},
													{
														Expression: `nes.v`,
														Head:       "nes",
													},
													{
														Expression: `nes.V()`,
														Head:       "nes",
													},
													{
														Expression: `nes.sub.v`,
														Head:       "nes",
													},
													{
														Expression: `es.sub.V()`,
														Head:       "es",
													},
													{
														Expression: `globalExampleStruct`,
														Head:       "globalExampleStruct",
													},
													{
														Expression: `NewExampleStruct(nes.sub.V())`,
														Head:       "NewExampleStruct",
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
										"es.sub.V": {
											{
												From: "es.sub",
												Call: "V",
											},
										},
										"nes.sub.V": {
											{
												From: "nes.sub",
												Call: "V",
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
									CallMeta: map[string][]*compareCallMeta{
										"fmt.Println": {
											{
												From: "fmt",
												Call: "Println",
												Args: []*compareArgMeta{
													{
														Expression: `"module.ExampleStruct.ExampleFuncWithPointerReceiver Hello go-extractor"`,
														Head:       `"module.ExampleStruct.ExampleFuncWithPointerReceiver Hello go-extractor"`,
													},
												},
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
						CallMeta: map[string][]*compareCallMeta{
							"rand.Intn": {
								{
									From: "rand",
									Call: "Intn",
									Args: []*compareArgMeta{
										{
											Expression: "v",
											Head:       "v",
										},
									},
								},
							},
						},
					},
					"ExampleFunc": {
						FunctionName: "ExampleFunc",
						CallMeta: map[string][]*compareCallMeta{
							"s.ExampleFunc": {
								{
									From: "s",
									Call: "ExampleFunc",
									Args: []*compareArgMeta{
										{
											Expression: `s.v`,
											Head:       "s",
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
	if gim.InterfaceName() != _gim.InterfaceName {
		Panic(gim.InterfaceName(), _gim.InterfaceName)
	}
}

func checkFunctionMeta(gfm *GoFunctionMeta, _gfm *compareGoFunctionMeta) {
	if gfm.FunctionName() != _gfm.FunctionName {
		Panic(gfm.FunctionName(), _gfm.FunctionName)
	}
	stpslice.Compare(gfm.Doc(), _gfm.Doc)
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
	recvStruct, pointerReceiver := gmm.RecvStruct()
	if recvStruct != _gmm.RecvStruct {
		Panic(recvStruct, _gmm.RecvStruct)
	}
	if pointerReceiver != _gmm.PointerReceiver {
		Panic(pointerReceiver, _gmm.PointerReceiver)
	}
	checkFunctionMeta(gmm.GoFunctionMeta, _gmm.compareGoFunctionMeta)
}

func checkCallMeta(gcm *GoCallMeta, _gcm *compareCallMeta) {
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
	if len(gcm.Args()) != len(_gcm.Args) {
		Panic(gcm.Args(), _gcm.Args)
	}
	args := gcm.Args()
	for index, _arg := range _gcm.Args {
		arg := args[index]
		if arg.Expression() != _arg.Expression {
			Panic(arg.Expression(), _arg.Expression)
		}
		if arg.Head() != _arg.Head {
			Panic(arg.Head(), _arg.Head)
		}
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
	compareGoCallMetaSlice = []*compareCallMeta{
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
