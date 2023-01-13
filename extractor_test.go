package extractor

import (
	"fmt"
	"reflect"
	"testing"

	stpmap "github.com/Mericusta/go-stp/map"
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
}

type compareGoMemberMeta struct {
	Name    string
	Tag     string
	Doc     []string
	Comment string
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
						Doc: []string{
							"// ExampleFunc this is example function",
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
						Doc: []string{
							"// ExampleStruct this is an example struct",
							"// this is struct comment",
							"// this is another struct comment",
						},
						StructMemberMeta: map[string]*compareGoMemberMeta{
							"v": {
								Name: "v",
								Tag:  "`ast:init,default=1`",
								Doc: []string{
									"// v this is member doc line1",
									"// v this is member doc line2",
								},
								Comment: "// this is member single comment line",
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
						Doc: []string{
							"// NewExampleStruct this is new example struct",
							"// @param           value",
							"// @return          pointer to ExampleStruct",
						},
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
	if len(gpm.FileNames()) != len(stpmap.Key(_gpm.pkgFileMap)) {
		panic(len(gpm.FileNames()))
	}
	for _, _fileName := range stpmap.Key(_gpm.pkgFileMap) {
		for _, fileName := range gpm.FileNames() {
			if fileName == _fileName {
				goto NEXT_FILE
			}
		}
		panic(_fileName)
	NEXT_FILE:
	}
	if len(gpm.StructNames()) != len(stpmap.Key(_gpm.pkgStructMeta)) {
		panic(len(gpm.StructNames()))
	}
	for _, _structName := range stpmap.Key(_gpm.pkgStructMeta) {
		for _, structName := range gpm.StructNames() {
			if structName == _structName {
				goto NEXT_STRUCT
			}
		}
		panic(_structName)
	NEXT_STRUCT:
	}
	if len(gpm.InterfaceNames()) != len(stpmap.Key(_gpm.pkgInterfaceMeta)) {
		panic(len(gpm.InterfaceNames()))
	}
	for _, _structName := range stpmap.Key(_gpm.pkgInterfaceMeta) {
		for _, structName := range gpm.InterfaceNames() {
			if structName == _structName {
				goto NEXT_INTERFACE
			}
		}
		panic(_structName)
	NEXT_INTERFACE:
	}
	if len(gpm.MethodNames()) != len(stpmap.Key(_gpm.pkgStructMeta)) {
		panic(len(gpm.MethodNames()))
	}
	for _structName, _gsm := range _gpm.pkgStructMeta {
		structMethodMap := gpm.MethodNames()
		methods, has := structMethodMap[_structName]
		if !has {
			panic(_structName)
		}
		if len(methods) != len(_gsm.StructMethodMeta) {
			panic(len(methods))
		}
		for _methodName := range _gsm.StructMethodMeta {
			for _, methodName := range methods {
				if _methodName == methodName {
					goto NEXT_METHOD
				}
			}
			panic(_methodName)
		NEXT_METHOD:
		}
	}
	// for _, _structName := range stpmap.Key(_gpm.pkgInterfaceMeta) {
	// 	for _, structName := range gpm.MethodNames() {
	// 		if structName == _structName {
	// 			goto NEXT_INTERFACE
	// 		}
	// 	}
	// 	panic(_structName)
	// NEXT_INTERFACE:
	// }
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
	checkDoc(gsm.Doc(), _gsm.Doc)
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
	checkDoc(gfm.Doc(), _gfm.Doc)
}

func checkMemberMeta(gmm *GoMemberMeta, _gmm *compareGoMemberMeta) {
	if gmm.MemberName() != _gmm.Name {
		panic(gmm.MemberName())
	}
	if gmm.Tag() != _gmm.Tag {
		panic(gmm.Tag())
	}
	checkDoc(gmm.Doc(), _gmm.Doc)
	if gmm.Comment() != _gmm.Comment {
		panic(gmm.Comment())
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

func checkDoc(doc, _doc []string) {
	if len(doc) != len(_doc) {
		panic(len(doc))
	}
	for _, _comment := range _doc {
		for _, comment := range doc {
			if comment == _comment {
				goto NEXT_COMMENT
			}
		}
		panic(_comment)
	NEXT_COMMENT:
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
	goProjectMeta, err := ExtractGoProjectMeta(standardProjectRelPath, standardProjectIgnorePathMap)
	if err != nil {
		panic(err)
	}

	for pkgName, replaceFunctionDoc := range replaceDoc {
		gpm, has := goProjectMeta.PackageMap[pkgName]
		if gpm == nil || !has {
			panic(pkgName)
		}
		for funcName, _replace := range replaceFunctionDoc {
			gpm.SearchFunctionMeta(funcName)
			gfm, has := gpm.pkgFunctionDecl[funcName]
			if gfm == nil || !has {
				panic(funcName)
			}
			checkDoc(gfm.Doc(), _replace.originDoc)
			originContent, replaceContent, err := gfm.ReplaceFunctionDoc(_replace.replaceDoc)
			if err != nil {
				panic(err)
			}
			if originContent != _replace.originContent {
				panic(originContent)
			}
			if replaceContent != _replace.replaceContent {
				panic(replaceContent)
			}
		}
	}
}

type compareCallMeta struct {
	Expression string
	Call       string
	Args       []interface{}
}

var (
	compareGoCallMetaSlice = []*compareCallMeta{
		{
			Expression: "HaveReadGP(1)",
			Call:       "HaveReadGP",
			Args:       []interface{}{int32(1)},
		},
		{
			Expression: "GetPlayerLevel()",
			Call:       "GetPlayerLevel",
			Args:       nil,
		},
		{
			Expression: `HaveReadGP("gamephone")`,
			Call:       "HaveReadGP",
			Args:       []interface{}{`"gamephone"`},
		},
		{
			Expression: `HaveReadGP("gamephone",1,"remove")`,
			Call:       "HaveReadGP",
			Args: []interface{}{
				`"gamephone"`, int32(1), `"remove"`,
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
		// 	Expression: "func() { HaveReadGP(1,HaveReadGP(1));HaveReadGP(1,HaveReadGP(1)) }",
		// 	Call:       "HaveReadGP",
		// 	Args:       []interface{}{int32(1)},
		// },
		// TODO: not support
		// {
		// 	Expression: "HaveReadGP(1)HaveReadGP(2)HaveReadGP(3)",
		// 	Call:       "HaveReadGP",
		// 	Args:       []interface{}{int32(1)},
		// },
		// {
		// 	Expression: "HaveReadGP(1) and HaveReadGP(2) and HaveReadGP(3)",
		// 	Call:       "HaveReadGP",
		// 	Args:       []interface{}{int32(1)},
		// },
	}
)

func TestParseGoCallMeta(t *testing.T) {
	for _, _gcm := range compareGoCallMetaSlice {
		gcm := ExtractGoCallMeta(_gcm.Expression)
		gcm.PrintAST()

		if gcm.Expression() != _gcm.Expression {
			panic(gcm.Expression())
		}
		if gcm.Call() != _gcm.Call {
			panic(gcm.Call())
		}

		if len(gcm.Args()) != len(_gcm.Args) {
			panic(len(gcm.Args()))
		}
		for _, _arg := range _gcm.Args {
			for _, arg := range gcm.Args() {
				if reflect.DeepEqual(arg, _arg) {
					goto NEXT_PARAM
				}
			}
			panic(_arg)
		NEXT_PARAM:
		}
	}
}
