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
	CallMeta     []*compareCallMeta
}

type compareGoMemberMeta struct {
	Name    string
	Tag     string
	Doc     []string
	Comment string
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
	ArgType    int
	Arg        string
	From       string
	Value      interface{}
	CallMeta   *compareCallMeta
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
					},
					"NoDocExampleFunc": {
						FunctionName: "NoDocExampleFunc",
					},
					"OneLineDocExampleFunc": {
						FunctionName: "OneLineDocExampleFunc",
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
								Name: "P",
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
								Name: "ParentStruct",
							},
							"v": {
								Name: "v",
								Tag:  "`ast:init,default=1`",
								Doc: []string{
									"// v this is member doc line1",
									"// v this is member doc line2",
								},
								Comment: "// this is member single comment line",
							},
							"sub": {
								Name: "sub",
							},
						},
						StructMethodMeta: map[string]*compareGoMethodMeta{
							"ExampleFunc": {
								compareGoFunctionMeta: &compareGoFunctionMeta{
									FunctionName: "ExampleFunc",
									CallMeta: []*compareCallMeta{
										{
											Expression: `NewExampleStruct(v)`,
											Call:       "NewExampleStruct",
											Args: []*compareArgMeta{
												{
													Arg: "v",
												},
											},
										},
										{
											Expression: `fmt.Println("module.ExampleStruct.ExampleFunc Hello go-extractor", es, es.v, es.V(), nes, nes.v, nes.V(), nes.sub.v, es.sub.V(), globalExampleStruct)`,
											Call:       "Println",
											From:       "fmt",
											Args: []*compareArgMeta{
												{Expression: `"module.ExampleStruct.ExampleFunc Hello go-extractor"`},
												{Expression: `"es"`},
												{Expression: `"es.v"`},
												{Expression: `"es.V()"`},
												{Expression: `"nes"`},
												{Expression: `"nes.v"`},
												{Expression: `"nes.V()"`},
												{Expression: `"nes.sub.v"`},
												{Expression: `"es.sub.V()"`},
												{Expression: `"globalExampleStruct"`},
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
									CallMeta: []*compareCallMeta{
										{
											Expression: `fmt.Println("module.ExampleStruct.ExampleFuncWithPointerReceiver Hello go-extractor")`,
											Call:       "Println",
											From:       "fmt",
											Args: []*compareArgMeta{
												{Expression: `"module.ExampleStruct.ExampleFuncWithPointerReceiver Hello go-extractor"`},
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
					},
					"ExampleFunc": {
						FunctionName: "ExampleFunc",
						CallMeta: []*compareCallMeta{
							{
								Expression: `s.ExampleFunc(s.v)`,
								From:       "s",
								Call:       "ExampleFunc",
								Args: []*compareArgMeta{
									{Expression: `s.v`},
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

			// for memberName, _gmm := range _gsm.StructMemberMeta {
			// 	gsm.SearchMemberMeta(memberName)
			// 	gmm, has := gsm.memberDecl[memberName]
			// 	if gmm == nil || !has {
			// 		panic(memberName)
			// 	}
			// 	checkMemberMeta(gmm, _gmm)
			// }

			// for methodName, _gmm := range _gsm.StructMethodMeta {
			// 	gpm.SearchMethodMeta(_structName, methodName)
			// 	gmm, has := gsm.methodDecl[methodName]
			// 	if gmm == nil || !has {
			// 		panic(methodName)
			// 	}
			// 	checkMethodMeta(gmm, _gmm)
			// }
		}

		// 	for interfaceName, _gim := range _gpm.pkgInterfaceMeta {
		// 		gpm.SearchInterfaceMeta(interfaceName)
		// 		gim, has := gpm.pkgInterfaceDecl[interfaceName]
		// 		if gim == nil || !has {
		// 			panic(interfaceName)
		// 		}
		// 		checkInterfaceMeta(gim, _gim)
		// 	}

		// 	for funcName, _gfm := range _gpm.pkgFunctionMeta {
		// 		gpm.SearchFunctionMeta(funcName)
		// 		gfm, has := gpm.pkgFunctionDecl[funcName]
		// 		if gfm == nil || !has {
		// 			panic(funcName)
		// 		}
		// 		checkFunctionMeta(gfm, _gfm)
		// 	}
	}
}

func Panic(v, c any) {
	panic(fmt.Sprintf("%v != %v", v, c))
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

	// file
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
		panic(gim.InterfaceName())
	}
}

func checkFunctionMeta(gfm *GoFunctionMeta, _gfm *compareGoFunctionMeta) {
	if gfm.FunctionName() != _gfm.FunctionName {
		panic(gfm.FunctionName())
	}
	stpslice.Compare(gfm.Doc(), _gfm.Doc)
	// for _, _gcm := range _gfm.CallMeta {
	// 	for _, gcm := range gfm.SearchCallMeta(_gcm.Call, _gcm.From) {
	// 		checkCallMeta(gcm, _gcm)
	// 	}
	// }
}

func checkMemberMeta(gmm *GoMemberMeta, _gmm *compareGoMemberMeta) {
	if gmm.MemberName() != _gmm.Name {
		panic(gmm.MemberName())
	}
	if gmm.Tag() != _gmm.Tag {
		panic(gmm.Tag())
	}
	stpslice.Compare(gmm.Doc(), _gmm.Doc)
	if gmm.Comment() != _gmm.Comment {
		panic(gmm.Comment())
	}
}

func checkMethodMeta(gmm *GoMethodMeta, _gmm *compareGoMethodMeta) {
	if gmm.FunctionName() != _gmm.FunctionName {
		panic(gmm.FunctionName())
	}
	if recvStruct, pointerReceiver := gmm.RecvStruct(); recvStruct != _gmm.RecvStruct || pointerReceiver != _gmm.PointerReceiver {
		panic(fmt.Sprintf("%v, %v", recvStruct, pointerReceiver))
	}
	// for _, _gcm := range _gmm.CallMeta {
	// 	for _, gcm := range gmm.SearchCallMeta(_gcm.Call, _gcm.From) {
	// 		checkCallMeta(gcm, _gcm)
	// 	}
	// }
}

// func checkCallMeta(gcm *GoCallMeta, _gcm *compareCallMeta) {
// 	if (gcm != nil) != (_gcm != nil) {
// 		panic(gcm != nil)
// 	}
// 	if gcm == nil {
// 		return
// 	}
// 	if gcm.Expression() != _gcm.Expression {
// 		panic(gcm.Expression())
// 	}
// 	if gcm.Call() != _gcm.Call {
// 		panic(gcm.Call())
// 	}
// 	if gcm.From() != _gcm.From {
// 		panic(gcm.From())
// 	}
// 	if len(gcm.Args()) != len(_gcm.Args) {
// 		panic(len(gcm.Args()))
// 	}
// 	for _, _arg := range _gcm.Args {
// 		for _, arg := range gcm.Args() {
// 			fmt.Printf("compare %v with %v\n", arg.Expression(), _arg.Expression)
// 			if arg.Expression() == _arg.Expression {
// 				// if arg.ArgType() != _arg.ArgType {
// 				// 	panic(arg.ArgType())
// 				// }
// 				// if arg.Arg() != _arg.Arg {
// 				// 	panic(arg.Arg())
// 				// }
// 				// if arg.From() != _arg.From {
// 				// 	panic(arg.From())
// 				// }
// 				// if !reflect.DeepEqual(arg.Value(), _arg.Value) {
// 				// 	panic(arg.Value())
// 				// }
// 				// checkCallMeta(arg.CallMeta(), _arg.CallMeta)
// 				goto NEXT_ARG
// 			}
// 		}
// 		panic(_arg)
// 	NEXT_ARG:
// 	}
// }

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
		{
			Expression: `HaveReadGP(1)`,
			Call:       "HaveReadGP",
			Args: []*compareArgMeta{
				{
					Value: int32(1),
				},
			},
		},
		{
			Expression: `GetPlayerLevel()`,
			Call:       "GetPlayerLevel",
			Args:       nil,
		},
		{
			Expression: `HaveReadGP("gamephone")`,
			Call:       "HaveReadGP",
			Args: []*compareArgMeta{
				{
					Value: `"gamephone"`},
			},
		},
		{
			Expression: `HaveReadGP("gamephone",1,"remove")`,
			Call:       "HaveReadGP",
			Args: []*compareArgMeta{
				{
					Value: `"gamephone"`,
				},
				{
					Value: int32(1),
				},
				{
					Value: `"remove"`,
				},
			},
		},
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
