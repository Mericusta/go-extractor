package extractor

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	stp "github.com/Mericusta/go-stp"
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
	TypeParams       []*compareGoVariableMeta
	StructMemberMeta map[string]*compareGoVariableMeta
	StructMethodMeta map[string]*compareGoMethodMeta
}

type compareGoInterfaceMeta struct {
	Expression          string
	InterfaceName       string
	Doc                 []string
	TypeParams          []*compareGoVariableMeta
	InterfaceMethodMeta map[string]*compareGoInterfaceMethodMeta
}

type compareGoFunctionDeclMeta struct {
	FunctionName string
	Doc          []string
	TypeParams   []*compareGoVariableMeta
	Params       []*compareGoVariableMeta
	ReturnTypes  []*compareGoVariableMeta
}

type compareGoInterfaceMethodMeta struct {
	compareGoFunctionDeclMeta
}

type compareGoFunctionMeta struct {
	compareGoFunctionDeclMeta
	CallMeta map[string][]*compareGoCallMeta
	// VarMeta map[string]
}

// type compareGoMemberMeta struct {
// 	Expression string
// 	MemberName string
// 	Tag        string
// 	Comment    string
// 	Doc        []string
// }

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
	// Head       *compareGoVariableMeta
}

type compareGoVariableMeta struct {
	Expression           string
	Name                 string
	TypeExpression       string
	TypeUnderlyingString string
	TypeUnderlyingEnum   UnderlyingType
	Tag                  string
	Comment              string
	Doc                  []string
}

type compareGoImportMeta struct {
	Expression string
	Alias      string
	ImportPath string
}

var (
	standardProjectRelPath       = "./testdata/standardProject"
	standardProjectIgnorePathMap = map[string]struct{}{
		standardProjectRelPath + "/vendor": {},
	}
	// standardProjectAbsPath = "d:\\Projects\\go-extractor\\testdata\\standardProject"
	standardProjectAbsPath = "/Users/dragonplus/Projects/github.com/Mericustar/go-extractor/testdata/standardProject"
	// standardProjectAbsPath    = "d:\\Projects\\SGAME\\server-dev\\gameServer\\game_server\\pkg\\github.com\\Mericusta\\go-extractor\\testdata\\standardProject"
	standardProjectModuleName = "standardProject"
	standardProjectMeta       = &compareGoProjectMeta{
		ProjectPath: standardProjectAbsPath,
		ModuleName:  standardProjectModuleName,
		PackageMap: map[string]*compareGoPackageMeta{
			"main": {
				Name:    "main",
				PkgPath: stp.FormatFilePathWithOS(standardProjectAbsPath + "\\cmd"),
				pkgFileMap: map[string]*compareGoFileMeta{
					"main.go": {
						Name:    "main.go",
						Path:    stp.FormatFilePathWithOS(standardProjectAbsPath + "\\cmd\\main.go"),
						PkgName: "main",
					},
					"init.go": {
						Name:    "init.go",
						Path:    stp.FormatFilePathWithOS(standardProjectAbsPath + "\\cmd\\init.go"),
						PkgName: "main",
					},
				},
				pkgFunctionMeta: map[string]*compareGoFunctionMeta{
					"main": {
						compareGoFunctionDeclMeta: compareGoFunctionDeclMeta{
							FunctionName: "main",
						},
						CallMeta: map[string][]*compareGoCallMeta{
							"pkg.ExampleFunc": {
								{
									Call: "ExampleFunc",
									From: "pkg",
									Args: []*compareGoArgMeta{
										{
											Expression: `module.NewExampleStruct(10)`,
											// Head: &compareGoVariableMeta{
											// 	Expression: `module`,
											// 	Name:       "module",
											// 	Type: &compareGoImportMeta{
											// 		Expression: `"standardProject/pkg/module"`,
											// 		Alias:      "module",
											// 		ImportPath: `"standardProject/pkg/module"`,
											// 	},
											// },
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
											// Head: &compareGoVariableMeta{
											// 	Expression: `module`,
											// 	Name:       "module",
											// 	Type: &compareGoImportMeta{
											// 		Expression: `"standardProject/pkg/module"`,
											// 		Alias:      "module",
											// 		ImportPath: `"standardProject/pkg/module"`,
											// 	},
											// },
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
											// Head: &compareGoVariableMeta{
											// 	Expression: `10`,
											// 	Name:       "10",
											// 	Type:       `10`,
											// },
										},
									},
								},
								{
									Call: "NewExampleStruct",
									From: "module",
									Args: []*compareGoArgMeta{
										{
											Expression: `11`,
											// Head: &compareGoVariableMeta{
											// 	Expression: `11`,
											// 	Name:       "11",
											// 	Type:       `11`,
											// },
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
						compareGoFunctionDeclMeta: compareGoFunctionDeclMeta{
							FunctionName: "Init",
						},
					},
				},
			},
			standardProjectModuleName + "/pkg": {
				Name:       "pkg",
				PkgPath:    stp.FormatFilePathWithOS(standardProjectAbsPath + "\\pkg"),
				ImportPath: standardProjectModuleName + "/pkg",
				pkgFileMap: map[string]*compareGoFileMeta{
					"pkg.go": {
						Name:    "pkg.go",
						Path:    stp.FormatFilePathWithOS(standardProjectAbsPath + "\\pkg\\pkg.go"),
						PkgName: "pkg",
					},
				},
				pkgFunctionMeta: map[string]*compareGoFunctionMeta{
					"ExampleFunc": {
						compareGoFunctionDeclMeta: compareGoFunctionDeclMeta{
							FunctionName: "ExampleFunc",
							Doc: []string{
								"// ExampleFunc this is example function",
							},
							Params: []*compareGoVariableMeta{
								{
									Expression:           `s *module.ExampleStruct`,
									Name:                 "s",
									TypeExpression:       "*module.ExampleStruct",
									TypeUnderlyingString: "pointer",
									TypeUnderlyingEnum:   UNDERLYING_TYPE_POINTER,
								},
							},
						},
						CallMeta: map[string][]*compareGoCallMeta{
							"fmt.Println": {
								{
									From: "fmt",
									Call: "Println",
									Args: []*compareGoArgMeta{
										{
											Expression: `"pkg.ExampleFunc, Hello go-extractor"`,
											// Head: &compareGoVariableMeta{
											// 	Expression: `"pkg.ExampleFunc, Hello go-extractor"`,
											// 	Name:       `"pkg.ExampleFunc, Hello go-extractor"`,
											// 	Type:       `"pkg.ExampleFunc, Hello go-extractor"`,
											// },
										},
										{
											Expression: `s.V()`,
											// Head: &compareGoVariableMeta{
											// 	Expression: `s *module.ExampleStruct`,
											// 	Name:       "s",
											// 	Type:       `*module.ExampleStruct`,
											// },
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
						compareGoFunctionDeclMeta: compareGoFunctionDeclMeta{
							FunctionName: "NoDocExampleFunc",
							Params: []*compareGoVariableMeta{
								{
									Expression:           `s *module.ExampleStruct`,
									Name:                 "s",
									TypeExpression:       "*module.ExampleStruct",
									TypeUnderlyingString: "pointer",
									TypeUnderlyingEnum:   UNDERLYING_TYPE_POINTER,
								},
							},
						},
						CallMeta: map[string][]*compareGoCallMeta{
							"fmt.Println": {
								{
									From: "fmt",
									Call: "Println",
									Args: []*compareGoArgMeta{
										{
											Expression: `"pkg.NoDocExampleFunc, Hello go-extractor"`,
											// Head: &compareGoVariableMeta{
											// 	Expression: `"pkg.NoDocExampleFunc, Hello go-extractor"`,
											// 	Name:       `"pkg.NoDocExampleFunc, Hello go-extractor"`,
											// 	Type:       `"pkg.NoDocExampleFunc, Hello go-extractor"`,
											// },
										},
										{
											Expression: `s.V()`,
											// Head: &compareGoVariableMeta{
											// 	Expression: `s *module.ExampleStruct`,
											// 	Name:       "s",
											// 	Type:       `*module.ExampleStruct`,
											// },
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
						compareGoFunctionDeclMeta: compareGoFunctionDeclMeta{
							FunctionName: "OneLineDocExampleFunc",
							Doc: []string{
								"// OneLineDocExampleFunc this is one-line-doc example function",
							},
							Params: []*compareGoVariableMeta{
								{
									Expression:           `s *module.ExampleStruct`,
									Name:                 "s",
									TypeExpression:       "*module.ExampleStruct",
									TypeUnderlyingString: "pointer",
									TypeUnderlyingEnum:   UNDERLYING_TYPE_POINTER,
								},
							},
						},
						CallMeta: map[string][]*compareGoCallMeta{
							"fmt.Println": {
								{
									From: "fmt",
									Call: "Println",
									Args: []*compareGoArgMeta{
										{
											Expression: `"pkg.OneLineDocExampleFunc, Hello go-extractor"`,
											// Head: &compareGoVariableMeta{
											// 	Expression: `"pkg.OneLineDocExampleFunc, Hello go-extractor"`,
											// 	Name:       `"pkg.OneLineDocExampleFunc, Hello go-extractor"`,
											// 	Type:       `"pkg.OneLineDocExampleFunc, Hello go-extractor"`,
											// },
										},
										{
											Expression: `s.V()`,
											// Head: &compareGoVariableMeta{
											// 	Expression: `s *module.ExampleStruct`,
											// 	Name:       "s",
											// 	Type:       `*module.ExampleStruct`,
											// },
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
						compareGoFunctionDeclMeta: compareGoFunctionDeclMeta{
							FunctionName: "ImportSelectorFunc",
							Params: []*compareGoVariableMeta{
								{
									Expression:           `s *module.ExampleStruct`,
									Name:                 "s",
									TypeExpression:       "*module.ExampleStruct",
									TypeUnderlyingString: "pointer",
									TypeUnderlyingEnum:   UNDERLYING_TYPE_POINTER,
								},
							},
						},
						CallMeta: map[string][]*compareGoCallMeta{
							"fmt.Println": {
								&compareGoCallMeta{
									From: "fmt",
									Call: "Println",
									Args: []*compareGoArgMeta{
										{
											Expression: `"pkg.ImportSelectorFunc, Hello go-extractor"`,
											// Head: &compareGoVariableMeta{
											// 	Expression: `"pkg.ImportSelectorFunc, Hello go-extractor"`,
											// 	Name:       `"pkg.ImportSelectorFunc, Hello go-extractor"`,
											// 	Type:       `"pkg.ImportSelectorFunc, Hello go-extractor"`,
											// },
										},
										{
											Expression: `module.NewExampleStruct(s.V()).Sub().ParentStruct.P`,
											// Head: &compareGoVariableMeta{
											// 	Expression: `module`,
											// 	Name:       "module",
											// 	Type: &compareGoImportMeta{
											// 		Expression: `"standardProject/pkg/module"`,
											// 		Alias:      "module",
											// 		ImportPath: `"standardProject/pkg/module"`,
											// 	},
											// },
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
											// Head: &compareGoVariableMeta{
											// 	Expression: `s *module.ExampleStruct`,
											// 	Name:       "s",
											// 	Type:       `*module.ExampleStruct`,
											// },
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
				PkgPath:    stp.FormatFilePathWithOS(standardProjectAbsPath + "\\pkg\\interface"),
				ImportPath: standardProjectModuleName + "/pkg/pkgInterface",
				pkgFileMap: map[string]*compareGoFileMeta{
					"interface.go": {
						Name:    "interface.go",
						Path:    stp.FormatFilePathWithOS(standardProjectAbsPath + "\\pkg\\interface\\interface.go"),
						PkgName: "pkgInterface",
					},
				},
				pkgInterfaceMeta: map[string]*compareGoInterfaceMeta{
					"ExampleInterface": {
						InterfaceName: "ExampleInterface",
						InterfaceMethodMeta: map[string]*compareGoInterfaceMethodMeta{
							"ExampleFunc": {
								compareGoFunctionDeclMeta{
									FunctionName: "ExampleFunc",
									Doc:          []string{"// This is ExampleFunc Doc"},
									Params: []*compareGoVariableMeta{
										{
											Expression:           `int`,
											Name:                 "p0",
											TypeExpression:       `int`,
											TypeUnderlyingString: `int`,
											TypeUnderlyingEnum:   UNDERLYING_TYPE_IDENT,
										},
									},
								},
							},
							"AnotherExampleFunc": {
								compareGoFunctionDeclMeta{
									FunctionName: "AnotherExampleFunc",
									Doc:          []string{"// This is AnotherExampleFunc Doc"},
									Params: []*compareGoVariableMeta{
										{
											Expression:           `int`,
											Name:                 "p0",
											TypeExpression:       `int`,
											TypeUnderlyingString: `int`,
											TypeUnderlyingEnum:   UNDERLYING_TYPE_IDENT,
										},
										{
											Expression:           `[]int`,
											Name:                 "p1",
											TypeExpression:       `[]int`,
											TypeUnderlyingString: `array`,
											TypeUnderlyingEnum:   UNDERLYING_TYPE_ARRAY,
										},
									},
									ReturnTypes: []*compareGoVariableMeta{
										{
											Expression:           `int`,
											TypeExpression:       `int`,
											TypeUnderlyingString: `int`,
											TypeUnderlyingEnum:   UNDERLYING_TYPE_IDENT,
										},
										{
											Expression:           `[]int`,
											TypeExpression:       `[]int`,
											TypeUnderlyingString: `array`,
											TypeUnderlyingEnum:   UNDERLYING_TYPE_ARRAY,
										},
									},
								},
							},
						},
					},
					"ExampleTemplateInterface": {
						InterfaceName: "ExampleTemplateInterface",
						TypeParams: []*compareGoVariableMeta{
							{
								Expression:           `T any`,
								Name:                 "T",
								TypeExpression:       "any",
								TypeUnderlyingString: "any",
								TypeUnderlyingEnum:   UNDERLYING_TYPE_IDENT,
							},
						},
						InterfaceMethodMeta: map[string]*compareGoInterfaceMethodMeta{
							"ExampleFunc": {
								compareGoFunctionDeclMeta{
									FunctionName: "ExampleFunc",
									Doc:          []string{"// This is ExampleFunc Doc"},
									TypeParams: []*compareGoVariableMeta{
										{
											Expression:           `T any`,
											Name:                 "T",
											TypeExpression:       `any`,
											TypeUnderlyingString: `any`,
											TypeUnderlyingEnum:   UNDERLYING_TYPE_IDENT,
										},
									},
									Params: []*compareGoVariableMeta{
										{
											Expression:           `T`,
											Name:                 "p0",
											TypeExpression:       `T`,
											TypeUnderlyingString: `T`,
											TypeUnderlyingEnum:   UNDERLYING_TYPE_IDENT,
										},
									},
								},
							},
							"AnotherExampleFunc": {
								compareGoFunctionDeclMeta{
									FunctionName: "AnotherExampleFunc",
									Doc:          []string{"// This is AnotherExampleFunc Doc"},
									TypeParams: []*compareGoVariableMeta{
										{
											Expression:           `T any`,
											Name:                 "T",
											TypeExpression:       `any`,
											TypeUnderlyingString: `any`,
											TypeUnderlyingEnum:   UNDERLYING_TYPE_IDENT,
										},
									},
									Params: []*compareGoVariableMeta{
										{
											Expression:           `T`,
											Name:                 "p0",
											TypeExpression:       `T`,
											TypeUnderlyingString: `T`,
											TypeUnderlyingEnum:   UNDERLYING_TYPE_IDENT,
										},
										{
											Expression:           `[]T`,
											Name:                 "p1",
											TypeExpression:       `[]T`,
											TypeUnderlyingString: `array`,
											TypeUnderlyingEnum:   UNDERLYING_TYPE_ARRAY,
										},
									},
									ReturnTypes: []*compareGoVariableMeta{
										{
											Expression:           `T`,
											TypeExpression:       `T`,
											TypeUnderlyingString: `T`,
											TypeUnderlyingEnum:   UNDERLYING_TYPE_IDENT,
										},
										{
											Expression:           `[]T`,
											TypeExpression:       `[]T`,
											TypeUnderlyingString: `array`,
											TypeUnderlyingEnum:   UNDERLYING_TYPE_ARRAY,
										},
									},
								},
							},
						},
					},
				},
			},
			standardProjectModuleName + "/pkg/module": {
				Name:       "module",
				PkgPath:    stp.FormatFilePathWithOS(standardProjectAbsPath + "\\pkg\\module"),
				ImportPath: standardProjectModuleName + "/pkg/module",
				pkgFileMap: map[string]*compareGoFileMeta{
					"module.go": {
						Name:    "module.go",
						Path:    stp.FormatFilePathWithOS(standardProjectAbsPath + "\\pkg\\module\\module.go"),
						PkgName: "module",
					},
				},
				pkgStructMeta: map[string]*compareGoStructMeta{
					"ParentStruct": {
						StructName: "ParentStruct",
						StructMemberMeta: map[string]*compareGoVariableMeta{
							"p": {
								Expression:           `p int`,
								Name:                 "p",
								TypeExpression:       `int`,
								TypeUnderlyingString: "int",
								TypeUnderlyingEnum:   UNDERLYING_TYPE_IDENT,
								Comment:              "// parent value",
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
						StructMemberMeta: map[string]*compareGoVariableMeta{
							"ParentStruct": {
								Expression:           `*ParentStruct`,
								Name:                 "ParentStruct",
								TypeExpression:       `*ParentStruct`,
								TypeUnderlyingString: "pointer",
								TypeUnderlyingEnum:   UNDERLYING_TYPE_POINTER,
								Comment:              "// parent struct",
							},
							"v": {
								Expression:           "v   int `ast:init,default=1`",
								Name:                 "v",
								TypeExpression:       `int`,
								TypeUnderlyingString: "int",
								TypeUnderlyingEnum:   UNDERLYING_TYPE_IDENT,
								Tag:                  "`ast:init,default=1`",
								Doc: []string{
									"// v this is member doc line1",
									"// v this is member doc line2",
								},
								Comment: "// this is member single comment line",
							},
							"sub": {
								Expression:           `sub *ExampleStruct`,
								Name:                 "sub",
								TypeExpression:       `*ExampleStruct`,
								TypeUnderlyingString: "pointer",
								TypeUnderlyingEnum:   UNDERLYING_TYPE_POINTER,
							},
						},
						StructMethodMeta: map[string]*compareGoMethodMeta{
							"ExampleFunc": {
								compareGoFunctionMeta: &compareGoFunctionMeta{
									compareGoFunctionDeclMeta: compareGoFunctionDeclMeta{
										FunctionName: "ExampleFunc",
										Params: []*compareGoVariableMeta{
											{
												Expression:           `v int`,
												Name:                 "v",
												TypeExpression:       "int",
												TypeUnderlyingString: "int",
												TypeUnderlyingEnum:   UNDERLYING_TYPE_IDENT,
											},
										},
									},
									CallMeta: map[string][]*compareGoCallMeta{
										"NewExampleStruct": {
											{
												Call: "NewExampleStruct",
												Args: []*compareGoArgMeta{
													{
														Expression: `v`,
														// Head: &compareGoVariableMeta{
														// 	Expression: `v int`,
														// 	Name:       "v",
														// 	Type:       `int`,
														// },
													},
												},
											},
											{
												Call: "NewExampleStruct",
												Args: []*compareGoArgMeta{
													{
														Expression: `nes.Sub().ParentStruct.P()`,
														// Head: &compareGoVariableMeta{
														// 	Expression: `nes *ExampleStruct`,
														// 	Name:       "nes",
														// 	Type:       `*ExampleStruct`,
														// },
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
														// Head: &compareGoVariableMeta{
														// 	Expression: `"module.ExampleStruct.ExampleFunc Hello go-extractor"`,
														// 	Name:       `"module.ExampleStruct.ExampleFunc Hello go-extractor"`,
														// 	Type:       `"module.ExampleStruct.ExampleFunc Hello go-extractor"`,
														// },
													},
													{
														Expression: `es`,
														// Head: &compareGoVariableMeta{
														// 	Expression: `es ExampleStruct`,
														// 	Name:       "es",
														// 	Type:       `ExampleStruct`,
														// },
													},
													{
														Expression: `es.v`,
														// Head: &compareGoVariableMeta{
														// 	Expression: `es ExampleStruct`,
														// 	Name:       "es",
														// 	Type:       `ExampleStruct`,
														// },
													},
													{
														Expression: `es.V()`,
														// Head: &compareGoVariableMeta{
														// 	Expression: `es ExampleStruct`,
														// 	Name:       "es",
														// 	Type:       `ExampleStruct`,
														// },
													},
													{
														Expression: `esP`,
														// Head: &compareGoVariableMeta{
														// 	Expression: `esP, esSubV := es.DoubleReturnFunc()`,
														// 	Name:       "esP",
														// 	Type:       `es.DoubleReturnFunc()`,
														// },
													},
													{
														Expression: `esSubV`,
														// Head: &compareGoVariableMeta{
														// 	Expression: `esP, esSubV := es.DoubleReturnFunc()`,
														// 	Name:       "esSubV",
														// 	Type:       `es.DoubleReturnFunc()`,
														// },
													},
													{
														Expression: `nes`,
														// Head: &compareGoVariableMeta{
														// 	Expression: `nes *ExampleStruct`,
														// 	Name:       "nes",
														// 	Type:       `*ExampleStruct`,
														// },
													},
													{
														Expression: `nes.v`,
														// Head: &compareGoVariableMeta{
														// 	Expression: `nes *ExampleStruct`,
														// 	Name:       "nes",
														// 	Type:       `*ExampleStruct`,
														// },
													},
													{
														Expression: `nes.V()`,
														// Head: &compareGoVariableMeta{
														// 	Expression: `nes *ExampleStruct`,
														// 	Name:       "nes",
														// 	Type:       `*ExampleStruct`,
														// },
													},
													{
														Expression: `nesP`,
														// Head: &compareGoVariableMeta{
														// 	Expression: `nesP, nesSubV := nes.DoubleReturnFunc()`,
														// 	Name:       "nesP",
														// 	Type:       `nes.DoubleReturnFunc()`,
														// },
													},
													{
														Expression: `nesSubV`,
														// Head: &compareGoVariableMeta{
														// 	Expression: `nesP, nesSubV := nes.DoubleReturnFunc()`,
														// 	Name:       "nesSubV",
														// 	Type:       `nes.DoubleReturnFunc()`,
														// },
													},
													{
														Expression: `globalExampleStruct`,
														// Head: &compareGoVariableMeta{
														// 	Expression: `globalExampleStruct *ExampleStruct`,
														// 	Name:       "globalExampleStruct",
														// 	Type:       `*ExampleStruct`,
														// },
													},
													{
														Expression: `NewExampleStruct(nes.Sub().ParentStruct.P())`,
														// Head: &compareGoVariableMeta{
														// 	Expression: `<IGNORE>`,
														// 	Name:       "NewExampleStruct",
														// 	Type:       `*ExampleStruct`,
														// },
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
									compareGoFunctionDeclMeta: compareGoFunctionDeclMeta{
										FunctionName: "ExampleFuncWithPointerReceiver",
										Params: []*compareGoVariableMeta{
											{
												Expression:           `v int`,
												Name:                 "v",
												TypeExpression:       "int",
												TypeUnderlyingString: "int",
												TypeUnderlyingEnum:   UNDERLYING_TYPE_IDENT,
											},
										},
									},
									CallMeta: map[string][]*compareGoCallMeta{
										"fmt.Println": {
											{
												From: "fmt",
												Call: "Println",
												Args: []*compareGoArgMeta{
													{
														Expression: `"module.ExampleStruct.ExampleFuncWithPointerReceiver Hello go-extractor"`,
														// Head: &compareGoVariableMeta{
														// 	Expression: `"module.ExampleStruct.ExampleFuncWithPointerReceiver Hello go-extractor"`,
														// 	Name:       `"module.ExampleStruct.ExampleFuncWithPointerReceiver Hello go-extractor"`,
														// 	Type:       `"module.ExampleStruct.ExampleFuncWithPointerReceiver Hello go-extractor"`,
														// },
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
									compareGoFunctionDeclMeta: compareGoFunctionDeclMeta{
										FunctionName: "DoubleReturnFunc",
										ReturnTypes: []*compareGoVariableMeta{
											{
												Expression:           `int`,
												TypeExpression:       "int",
												TypeUnderlyingString: "int",
												TypeUnderlyingEnum:   UNDERLYING_TYPE_IDENT,
											},
											{
												Expression:           `int`,
												TypeExpression:       "int",
												TypeUnderlyingString: "int",
												TypeUnderlyingEnum:   UNDERLYING_TYPE_IDENT,
											},
										},
									},
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
									compareGoFunctionDeclMeta: compareGoFunctionDeclMeta{
										FunctionName: "V",
										ReturnTypes: []*compareGoVariableMeta{
											{
												Expression:           `int`,
												TypeExpression:       "int",
												TypeUnderlyingString: "int",
												TypeUnderlyingEnum:   UNDERLYING_TYPE_IDENT,
											},
										},
									},
								},
								RecvStruct:      "ExampleStruct",
								PointerReceiver: false,
							},
							"Sub": {
								compareGoFunctionMeta: &compareGoFunctionMeta{
									compareGoFunctionDeclMeta: compareGoFunctionDeclMeta{
										FunctionName: "Sub",
										ReturnTypes: []*compareGoVariableMeta{
											{
												Expression:           `*ExampleStruct`,
												TypeExpression:       "*ExampleStruct",
												TypeUnderlyingString: "pointer",
												TypeUnderlyingEnum:   UNDERLYING_TYPE_POINTER,
											},
										},
									},
								},
								RecvStruct:      "ExampleStruct",
								PointerReceiver: true,
							},
						},
					},
				},
				pkgFunctionMeta: map[string]*compareGoFunctionMeta{
					"NewExampleStruct": {
						compareGoFunctionDeclMeta: compareGoFunctionDeclMeta{
							FunctionName: "NewExampleStruct",
							Doc: []string{
								"// NewExampleStruct this is new example struct",
								"// @param           value",
								"// @return          pointer to ExampleStruct",
							},
							Params: []*compareGoVariableMeta{
								{
									Expression:           `v int`,
									Name:                 "v",
									TypeExpression:       "int",
									TypeUnderlyingString: "int",
									TypeUnderlyingEnum:   UNDERLYING_TYPE_IDENT,
								},
							},
							ReturnTypes: []*compareGoVariableMeta{
								{
									Expression:           `*ExampleStruct`,
									TypeExpression:       "*ExampleStruct",
									TypeUnderlyingString: "pointer",
									TypeUnderlyingEnum:   UNDERLYING_TYPE_POINTER,
								},
							},
						},
						CallMeta: map[string][]*compareGoCallMeta{
							"random.Intn": {
								{
									From: "random",
									Call: "Intn",
									Args: []*compareGoArgMeta{
										{
											Expression: `v`,
											// Head: &compareGoVariableMeta{
											// 	Expression: `v int`,
											// 	Name:       "v",
											// 	Type:       `int`,
											// },
										},
									},
								},
							},
						},
					},
					"ExampleFunc": {
						compareGoFunctionDeclMeta: compareGoFunctionDeclMeta{
							FunctionName: "ExampleFunc",
							Params: []*compareGoVariableMeta{
								{
									Expression:           `s *ExampleStruct`,
									Name:                 "s",
									TypeExpression:       "*ExampleStruct",
									TypeUnderlyingString: "pointer",
									TypeUnderlyingEnum:   UNDERLYING_TYPE_POINTER,
								},
							},
						},
						CallMeta: map[string][]*compareGoCallMeta{
							"s.ExampleFunc": {
								{
									From: "s",
									Call: "ExampleFunc",
									Args: []*compareGoArgMeta{
										{
											Expression: `s.v`,
											// Head: &compareGoVariableMeta{
											// 	Expression: `s *ExampleStruct`,
											// 	Name:       "s",
											// 	Type:       `*ExampleStruct`,
											// },
										},
									},
								},
							},
						},
					},
				},
			},
			standardProjectModuleName + "/pkg/template": {
				Name:       "template",
				PkgPath:    stp.FormatFilePathWithOS(standardProjectAbsPath + "\\pkg\\template"),
				ImportPath: standardProjectModuleName + "/pkg/template",
				pkgFileMap: map[string]*compareGoFileMeta{
					"template.go": {
						Name:    "template.go",
						Path:    stp.FormatFilePathWithOS(standardProjectAbsPath + "\\pkg\\template\\template.go"),
						PkgName: "template",
					},
				},
				pkgFunctionMeta: map[string]*compareGoFunctionMeta{
					"OneTemplateFunc": {
						compareGoFunctionDeclMeta: compareGoFunctionDeclMeta{
							FunctionName: "OneTemplateFunc",
							TypeParams: []*compareGoVariableMeta{
								{
									Expression:           `T any`,
									Name:                 "T",
									TypeExpression:       "any",
									TypeUnderlyingString: "any",
									TypeUnderlyingEnum:   UNDERLYING_TYPE_IDENT,
								},
							},
							Params: []*compareGoVariableMeta{
								{
									Expression:           `tv *T`,
									Name:                 "tv",
									TypeExpression:       "*T",
									TypeUnderlyingString: "pointer",
									TypeUnderlyingEnum:   UNDERLYING_TYPE_POINTER,
								},
							},
							ReturnTypes: []*compareGoVariableMeta{
								{
									Expression:           `*T`,
									TypeExpression:       "*T",
									TypeUnderlyingString: "pointer",
									TypeUnderlyingEnum:   UNDERLYING_TYPE_POINTER,
								},
							},
						},
					},
					"DoubleSameTemplateFunc": {
						compareGoFunctionDeclMeta: compareGoFunctionDeclMeta{
							FunctionName: "DoubleSameTemplateFunc",
							TypeParams: []*compareGoVariableMeta{
								{
									Expression:           `T1, T2 any`,
									Name:                 "T1",
									TypeExpression:       "any",
									TypeUnderlyingString: "any",
									TypeUnderlyingEnum:   UNDERLYING_TYPE_IDENT,
								},
								{
									Expression:           `T1, T2 any`,
									Name:                 "T2",
									TypeExpression:       "any",
									TypeUnderlyingString: "any",
									TypeUnderlyingEnum:   UNDERLYING_TYPE_IDENT,
								},
							},
							Params: []*compareGoVariableMeta{
								{
									Expression:           `tv1 T1`,
									Name:                 "tv1",
									TypeExpression:       "T1",
									TypeUnderlyingString: "T1",
									TypeUnderlyingEnum:   UNDERLYING_TYPE_IDENT,
								},
								{
									Expression:           `tv2 T2`,
									Name:                 "tv2",
									TypeExpression:       "T2",
									TypeUnderlyingString: "T2",
									TypeUnderlyingEnum:   UNDERLYING_TYPE_IDENT,
								},
							},
							ReturnTypes: []*compareGoVariableMeta{
								{
									Expression:           `*T1`,
									TypeExpression:       "*T1",
									TypeUnderlyingString: "pointer",
									TypeUnderlyingEnum:   UNDERLYING_TYPE_POINTER,
								},
								{
									Expression:           `*T2`,
									TypeExpression:       "*T2",
									TypeUnderlyingString: "pointer",
									TypeUnderlyingEnum:   UNDERLYING_TYPE_POINTER,
								},
							},
						},
					},
					"DoubleDifferenceTemplateFunc": {
						compareGoFunctionDeclMeta: compareGoFunctionDeclMeta{
							FunctionName: "DoubleDifferenceTemplateFunc",
							TypeParams: []*compareGoVariableMeta{
								{
									Expression:           `T1 any`,
									Name:                 "T1",
									TypeExpression:       "any",
									TypeUnderlyingString: "any",
									TypeUnderlyingEnum:   UNDERLYING_TYPE_IDENT,
								},
								{
									Expression:           `T2 comparable`,
									Name:                 "T2",
									TypeExpression:       "comparable",
									TypeUnderlyingString: "comparable",
									TypeUnderlyingEnum:   UNDERLYING_TYPE_IDENT,
								},
							},
							Params: []*compareGoVariableMeta{
								{
									Expression:           `tv1 T1`,
									Name:                 "tv1",
									TypeExpression:       "T1",
									TypeUnderlyingString: "T1",
									TypeUnderlyingEnum:   UNDERLYING_TYPE_IDENT,
								},
								{
									Expression:           `tv2 T2`,
									Name:                 "tv2",
									TypeExpression:       "T2",
									TypeUnderlyingString: "T2",
									TypeUnderlyingEnum:   UNDERLYING_TYPE_IDENT,
								},
							},
							ReturnTypes: []*compareGoVariableMeta{
								{
									Expression:           `*T1`,
									TypeExpression:       "*T1",
									TypeUnderlyingString: "pointer",
									TypeUnderlyingEnum:   UNDERLYING_TYPE_POINTER,
								},
								{
									Expression:           `*T2`,
									TypeExpression:       "*T2",
									TypeUnderlyingString: "pointer",
									TypeUnderlyingEnum:   UNDERLYING_TYPE_POINTER,
								},
							},
						},
					},
					"TypeConstraintsTemplateFunc": {
						compareGoFunctionDeclMeta: compareGoFunctionDeclMeta{
							FunctionName: "TypeConstraintsTemplateFunc",
							TypeParams: []*compareGoVariableMeta{
								{
									Expression:           `T TypeConstraints`,
									Name:                 "T",
									TypeExpression:       "TypeConstraints",
									TypeUnderlyingString: "TypeConstraints",
									TypeUnderlyingEnum:   UNDERLYING_TYPE_IDENT,
								},
							},
							Params: []*compareGoVariableMeta{
								{
									Expression:           `tv T`,
									Name:                 "tv",
									TypeExpression:       "T",
									TypeUnderlyingString: "T",
									TypeUnderlyingEnum:   UNDERLYING_TYPE_IDENT,
								},
							},
							ReturnTypes: []*compareGoVariableMeta{
								{
									Expression:           `*T`,
									TypeExpression:       "*T",
									TypeUnderlyingString: "pointer",
									TypeUnderlyingEnum:   UNDERLYING_TYPE_POINTER,
								},
							},
						},
					},
					"CannotInferTypeFunc1": {
						compareGoFunctionDeclMeta: compareGoFunctionDeclMeta{
							FunctionName: "CannotInferTypeFunc1",
							TypeParams: []*compareGoVariableMeta{
								{
									Expression:           `T any`,
									Name:                 "T",
									TypeExpression:       "any",
									TypeUnderlyingString: "any",
									TypeUnderlyingEnum:   UNDERLYING_TYPE_IDENT,
								},
							},
						},
					},
					"CannotInferTypeFunc2": {
						compareGoFunctionDeclMeta: compareGoFunctionDeclMeta{
							FunctionName: "CannotInferTypeFunc2",
							TypeParams: []*compareGoVariableMeta{
								{
									Expression:           `K comparable`,
									Name:                 "K",
									TypeExpression:       "comparable",
									TypeUnderlyingString: "comparable",
									TypeUnderlyingEnum:   UNDERLYING_TYPE_IDENT,
								},
								{
									Expression:           `V any`,
									Name:                 "V",
									TypeExpression:       "any",
									TypeUnderlyingString: "any",
									TypeUnderlyingEnum:   UNDERLYING_TYPE_IDENT,
								},
							},
							ReturnTypes: []*compareGoVariableMeta{
								{
									Expression:           `*K`,
									TypeExpression:       "*K",
									TypeUnderlyingString: "pointer",
									TypeUnderlyingEnum:   UNDERLYING_TYPE_POINTER,
								},
								{
									Expression:           `*V`,
									TypeExpression:       "*V",
									TypeUnderlyingString: "pointer",
									TypeUnderlyingEnum:   UNDERLYING_TYPE_POINTER,
								},
							},
						},
					},
				},
				pkgStructMeta: map[string]*compareGoStructMeta{
					"TemplateStruct": {
						StructName: "TemplateStruct",
						TypeParams: []*compareGoVariableMeta{
							{
								Expression:           `T any`,
								Name:                 "T",
								TypeExpression:       "any",
								TypeUnderlyingString: "any",
								TypeUnderlyingEnum:   UNDERLYING_TYPE_IDENT,
							},
						},
						StructMemberMeta: map[string]*compareGoVariableMeta{
							"v": {
								Expression:           `v T`,
								Name:                 "v",
								TypeExpression:       `T`,
								TypeUnderlyingString: "T",
								TypeUnderlyingEnum:   UNDERLYING_TYPE_IDENT,
							},
						},
						StructMethodMeta: map[string]*compareGoMethodMeta{
							"V": {
								compareGoFunctionMeta: &compareGoFunctionMeta{
									compareGoFunctionDeclMeta: compareGoFunctionDeclMeta{
										FunctionName: "V",
										ReturnTypes: []*compareGoVariableMeta{
											{
												Expression:           `T`,
												TypeExpression:       `T`,
												TypeUnderlyingString: "T",
												TypeUnderlyingEnum:   UNDERLYING_TYPE_IDENT,
											},
										},
									},
								},
								RecvStruct:      `TemplateStruct`,
								PointerReceiver: true,
							},
						},
					},
					"TwoTypeTemplateStruct": {
						StructName: "TwoTypeTemplateStruct",
						TypeParams: []*compareGoVariableMeta{
							{
								Expression:           `K TypeConstraints`,
								Name:                 "K",
								TypeExpression:       "TypeConstraints",
								TypeUnderlyingString: "TypeConstraints",
								TypeUnderlyingEnum:   UNDERLYING_TYPE_IDENT,
							},
							{
								Expression:           `V any`,
								Name:                 "V",
								TypeExpression:       "any",
								TypeUnderlyingString: "any",
								TypeUnderlyingEnum:   UNDERLYING_TYPE_IDENT,
							},
						},
						StructMemberMeta: map[string]*compareGoVariableMeta{
							"v": {
								Expression:           `v map[K]V`,
								Name:                 "v",
								TypeExpression:       `map[K]V`,
								TypeUnderlyingString: "map",
								TypeUnderlyingEnum:   UNDERLYING_TYPE_MAP,
							},
						},
						StructMethodMeta: map[string]*compareGoMethodMeta{
							"KVSlice": {
								compareGoFunctionMeta: &compareGoFunctionMeta{
									compareGoFunctionDeclMeta: compareGoFunctionDeclMeta{
										FunctionName: "KVSlice",
										Params: []*compareGoVariableMeta{
											{
												Expression:           `k K`,
												Name:                 "k",
												TypeExpression:       `K`,
												TypeUnderlyingString: "K",
												TypeUnderlyingEnum:   UNDERLYING_TYPE_IDENT,
											},
											{
												Expression:           `v V`,
												Name:                 "v",
												TypeExpression:       `V`,
												TypeUnderlyingString: "V",
												TypeUnderlyingEnum:   UNDERLYING_TYPE_IDENT,
											},
										},
										ReturnTypes: []*compareGoVariableMeta{
											{
												Expression:           `[]K`,
												TypeExpression:       `[]K`,
												TypeUnderlyingString: "array",
												TypeUnderlyingEnum:   UNDERLYING_TYPE_ARRAY,
											},
											{
												Expression:           `[]V`,
												TypeExpression:       `[]V`,
												TypeUnderlyingString: "array",
												TypeUnderlyingEnum:   UNDERLYING_TYPE_ARRAY,
											},
										},
									},
								},
								RecvStruct:      `TwoTypeTemplateStruct`,
								PointerReceiver: true,
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

	for _pkgName, _gpm := range standardProjectMeta.PackageMap {
		gpm := goProjectMeta.SearchPackageMeta(_pkgName)
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

			// member
			for _memberName, _gvm := range _gsm.StructMemberMeta {
				gvm := gsm.SearchMemberMeta(_memberName)
				if gvm == nil {
					Panic(gvm, _gvm)
				}
				checkVariableMeta(gvm, _gvm)
			}

			// method
			for _methodName, _gmm := range _gsm.StructMethodMeta {
				gmm := gpm.SearchMethodMeta(_structName, _methodName)
				if gmm == nil {
					Panic(gmm, _gmm)
				}
				checkMethodMeta(gmm, _gmm)

				// unit test
				var unittestFuncName string
				var unittestByte []byte
				if l := len(gmm.TypeParams()); l == 0 {
					unittestFuncName, unittestByte = gmm.MakeUnitTest(nil)
				} else {
					testTypeArgs := []string{"string", "[]string", "map[string]string"}
					typeArgs := make([]string, 0, l)
					for i := 0; i < l; i++ {
						typeArgs = append(typeArgs, testTypeArgs[i%len(testTypeArgs)])
					}
					unittestFuncName, unittestByte = gmm.MakeUnitTest(typeArgs)
				}
				fmt.Printf("unit test func %v:\n%v\n", unittestFuncName, string(unittestByte))

				// benchmark
				var benchmarkFuncName string
				var benchmarkByte []byte
				if l := len(gmm.TypeParams()); l == 0 {
					benchmarkFuncName, benchmarkByte = gmm.MakeBenchmark(nil)
				} else {
					testTypeArgs := []string{"string", "[]string", "map[string]string"}
					typeArgs := make([]string, 0, l)
					for i := 0; i < l; i++ {
						typeArgs = append(typeArgs, testTypeArgs[i%len(testTypeArgs)])
					}
					benchmarkFuncName, benchmarkByte = gmm.MakeBenchmark(typeArgs)
				}
				fmt.Printf("benchmark func %v:\n%v\n", benchmarkFuncName, string(benchmarkByte))
			}
		}

		for _interfaceName, _gim := range _gpm.pkgInterfaceMeta {
			gim := gpm.SearchInterfaceMeta(_interfaceName)
			if gim == nil {
				Panic(gim, _gim)
			}
			checkInterfaceMeta(gim, _gim)

			for _interfaceMethodName, _gimm := range _gim.InterfaceMethodMeta {
				gimm := gim.SearchMethodDecl(_interfaceMethodName)
				if gimm == nil {
					Panic(gimm, _gimm)
				}
				checkInterfaceMethodMeta(gimm, _gimm)

				// unit test
				var unittestFuncName string
				var unittestByte []byte
				if l := len(gimm.TypeParams()); l == 0 {
					unittestFuncName, unittestByte = gimm.MakeUnitTest(nil)
				} else {
					testTypeArgs := []string{"string", "[]string", "map[string]string"}
					typeArgs := make([]string, 0, l)
					for i := 0; i < l; i++ {
						typeArgs = append(typeArgs, testTypeArgs[i%len(testTypeArgs)])
					}
					unittestFuncName, unittestByte = gimm.MakeUnitTest(typeArgs)
				}
				fmt.Printf("unit test func %v:\n%v\n", unittestFuncName, string(unittestByte))

				// benchmark
				var benchmarkFuncName string
				var benchmarkByte []byte
				if l := len(gimm.TypeParams()); l == 0 {
					benchmarkFuncName, benchmarkByte = gimm.MakeBenchmark(nil)
				} else {
					testTypeArgs := []string{"string", "[]string", "map[string]string"}
					typeArgs := make([]string, 0, l)
					for i := 0; i < l; i++ {
						typeArgs = append(typeArgs, testTypeArgs[i%len(testTypeArgs)])
					}
					benchmarkFuncName, benchmarkByte = gimm.MakeBenchmark(typeArgs)
				}
				fmt.Printf("benchmark func %v:\n%v\n", benchmarkFuncName, string(benchmarkByte))

				// implement
				var implementFuncName string
				var implementFuncMeta *GoFunctionMeta
				receiverIdent, receiverType := strings.ToLower(string(gimm.interfaceMeta.InterfaceName()[0])), gimm.interfaceMeta.InterfaceName()+"Implement"
				if l := len(gimm.TypeParams()); l == 0 {
					implementFuncName, implementFuncMeta = gimm.MakeImplementMethodMeta(receiverIdent, receiverType)
				} else {
					// TODO: support
					// testTypeArgs := []string{"string", "[]string", "map[string]string"}
					// typeArgs := make([]string, 0, l)
					// for i := 0; i < l; i++ {
					// 	typeArgs = append(typeArgs, testTypeArgs[i%len(testTypeArgs)])
					// }
					// benchmarkFuncName, benchmarkByte = gimm.MakeBenchmark(typeArgs)
				}
				if implementFuncMeta != nil {
					fmt.Printf("implement func %v:\n%v\n", implementFuncName, implementFuncMeta.Format())
				}
			}
		}

		for _funcName, _gfm := range _gpm.pkgFunctionMeta {
			gfm := gpm.SearchFunctionMeta(_funcName)
			if gfm == nil {
				Panic(gfm, _gfm)
			}
			checkFunctionMeta(gfm, _gfm)

			// unit test
			var unittestFuncName string
			var unittestByte []byte
			if l := len(gfm.TypeParams()); l == 0 {
				unittestFuncName, unittestByte = gfm.MakeUnitTest(nil)
			} else {
				testTypeArgs := []string{"string", "[]string", "map[string]string"}
				typeArgs := make([]string, 0, l)
				for i := 0; i < l; i++ {
					typeArgs = append(typeArgs, testTypeArgs[i%len(testTypeArgs)])
				}
				unittestFuncName, unittestByte = gfm.MakeUnitTest(typeArgs)
			}
			fmt.Printf("unit test func %v:\n%v\n", unittestFuncName, string(unittestByte))

			// benchmark
			var benchmarkFuncName string
			var benchmarkByte []byte
			if l := len(gfm.TypeParams()); l == 0 {
				benchmarkFuncName, benchmarkByte = gfm.MakeBenchmark(nil)
			} else {
				testTypeArgs := []string{"string", "[]string", "map[string]string"}
				typeArgs := make([]string, 0, l)
				for i := 0; i < l; i++ {
					typeArgs = append(typeArgs, testTypeArgs[i%len(testTypeArgs)])
				}
				benchmarkFuncName, benchmarkByte = gfm.MakeBenchmark(typeArgs)
			}
			fmt.Printf("benchmark func %v:\n%v\n", benchmarkFuncName, string(benchmarkByte))

			testFileByte := MakeTestFile(fmt.Sprintf("%v_test.go", strings.Trim(gfm.path, ".go")), nil)
			fmt.Printf("unit test file:\n%v\n", string(testFileByte))
		}
	}

	// // arg type
	// for pkgImportPath, gpm := range goProjectMeta.packageMap {
	// 	// function
	// 	for funcName, gfm := range gpm.pkgFunctionDecl {
	// 		for call, gcms := range gfm.callMeta {
	// 			for _, gcm := range gcms {
	// 				for _, arg := range gcm.Args() {
	// 					if pkgImportPath == "main" && funcName == "main" && call == "pkg.ExampleFunc" {
	// 						fmt.Printf("in pkg %v, func %v, call %v\n", pkgImportPath, funcName, call)
	// 						argType := goProjectMeta.SearchArgType(arg)
	// 						fmt.Printf("arg %v type %v\n", arg.Expression(), argType)
	// 						fmt.Println()
	// 					}
	// 				}
	// 			}
	// 		}
	// 	}

	// 	// method
	// }
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
	_packages := stp.Key(_gpm.PackageMap)
	sort.Strings(_packages)
	if !stp.Compare(packages, _packages) {
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
	_fileNames := stp.Key(_gpm.pkgFileMap)
	sort.Strings(_fileNames)
	if !stp.Compare(fileNames, _fileNames) {
		Panic(fileNames, _fileNames)
	}

	// struct
	structNames := gpm.StructNames()
	sort.Strings(structNames)
	_structNames := stp.Key(_gpm.pkgStructMeta)
	sort.Strings(_structNames)
	if !stp.Compare(structNames, _structNames) {
		Panic(structNames, _structNames)
	}

	// interface
	interfaceNames := gpm.InterfaceNames()
	sort.Strings(interfaceNames)
	_interfaceNames := stp.Key(_gpm.pkgInterfaceMeta)
	sort.Strings(_interfaceNames)
	if !stp.Compare(interfaceNames, _interfaceNames) {
		Panic(interfaceNames, _interfaceNames)
	}

	// function
	functionNames := gpm.FunctionNames()
	sort.Strings(functionNames)
	_functionNames := stp.Key(_gpm.pkgFunctionMeta)
	sort.Strings(_functionNames)
	if !stp.Compare(functionNames, _functionNames) {
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
	if !stp.Compare(gsm.Doc(), _gsm.Doc) {
		Panic(gsm.Doc(), _gsm.Doc)
	}

	// type params
	if len(gsm.TypeParams()) != len(_gsm.TypeParams) {
		Panic(len(gsm.TypeParams()), len(_gsm.TypeParams))
	}
	for i := range _gsm.TypeParams {
		checkVariableMeta(gsm.TypeParams()[i], _gsm.TypeParams[i])
	}

	// member
	memberNames := gsm.Members()
	sort.Strings(memberNames)
	_memberNames := stp.Key(_gsm.StructMemberMeta)
	sort.Strings(_memberNames)
	if !stp.Compare(memberNames, _memberNames) {
		Panic(memberNames, _memberNames)
	}
}

func checkInterfaceMeta(gim *GoInterfaceMeta, _gim *compareGoInterfaceMeta) {
	// basic
	if gim.InterfaceName() != _gim.InterfaceName {
		Panic(gim.InterfaceName(), _gim.InterfaceName)
	}
	if !stp.Compare(gim.Doc(), _gim.Doc) {
		Panic(gim.Doc(), _gim.Doc)
	}

	// type params
	if len(gim.TypeParams()) != len(_gim.TypeParams) {
		Panic(len(gim.TypeParams()), len(_gim.TypeParams))
	}
	for i := range _gim.TypeParams {
		checkVariableMeta(gim.TypeParams()[i], _gim.TypeParams[i])
	}
}

func checkFunctionMeta(gfm *GoFunctionMeta, _gfm *compareGoFunctionMeta) {
	// basic
	if gfm.FunctionName() != _gfm.FunctionName {
		Panic(gfm.FunctionName(), _gfm.FunctionName)
	}
	if !stp.Compare(gfm.Doc(), _gfm.Doc) {
		Panic(gfm.Doc(), _gfm.Doc)
	}

	if len(gfm.TypeParams()) != len(_gfm.TypeParams) {
		Panic(len(gfm.TypeParams()), len(_gfm.TypeParams))
	}
	for i := range _gfm.TypeParams {
		checkVariableMeta(gfm.TypeParams()[i], _gfm.TypeParams[i])
	}

	if len(gfm.Params()) != len(_gfm.Params) {
		Panic(len(gfm.Params()), len(_gfm.Params))
	}
	for i := range _gfm.Params {
		checkVariableMeta(gfm.Params()[i], _gfm.Params[i])
	}

	if len(gfm.ReturnTypes()) != len(_gfm.ReturnTypes) {
		Panic(len(gfm.ReturnTypes()), len(_gfm.ReturnTypes))
	}
	for i := range _gfm.ReturnTypes {
		checkVariableMeta(gfm.ReturnTypes()[i], _gfm.ReturnTypes[i])
	}

	// // call
	// calls := stp.Key(gfm.Calls())
	// sort.Strings(calls)
	// _calls := stp.Key(_gfm.CallMeta)
	// sort.Strings(_calls)
	// if !stp.Compare(calls, _calls) {
	// 	gfm.Calls()
	// 	Panic(calls, _calls)
	// }

	// for _call, _gcmSlice := range _gfm.CallMeta {
	// 	gcmSlice := gfm.SearchCallMeta(_call)
	// 	if len(gcmSlice) != len(_gcmSlice) {
	// 		Panic(gcmSlice, _gcmSlice)
	// 	}
	// 	for index, _gcm := range _gcmSlice {
	// 		gcm := gcmSlice[index]
	// 		checkCallMeta(gcm, _gcm)
	// 	}
	// }
}

// func checkMemberMeta(gmm *GoVariableMeta, _gmm *compareGoVariableMeta) {
// 	// basic
// 	if gmm.MemberName() != _gmm.MemberName {
// 		Panic(gmm.MemberName(), _gmm.MemberName)
// 	}
// 	if gmm.Tag() != _gmm.Tag {
// 		Panic(gmm.Tag(), _gmm.Tag)
// 	}
// 	if gmm.Comment() != _gmm.Comment {
// 		Panic(gmm.Comment(), _gmm.Comment)
// 	}
// 	stp.Compare(gmm.Doc(), _gmm.Doc)
// }

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

	// unit test
	// b := MakeUnitTest(gmm)
	// fmt.Printf("unit test func:\n%v\n", string(b))
}

func checkInterfaceMethodMeta(gimm *GoInterfaceMethodMeta, _gimm *compareGoInterfaceMethodMeta) {
	if gimm.FunctionName() != _gimm.FunctionName {
		Panic(gimm.FunctionName(), _gimm.FunctionName)
	}
	if !stp.Compare(gimm.Doc(), _gimm.Doc) {
		Panic(gimm.Doc(), _gimm.Doc)
	}

	if len(gimm.TypeParams()) != len(_gimm.TypeParams) {
		Panic(len(gimm.TypeParams()), len(_gimm.TypeParams))
	}
	for i := range _gimm.TypeParams {
		checkVariableMeta(gimm.TypeParams()[i], _gimm.TypeParams[i])
	}

	if len(gimm.Params()) != len(_gimm.Params) {
		Panic(len(gimm.Params()), len(_gimm.Params))
	}
	for i := range _gimm.Params {
		checkVariableMeta(gimm.Params()[i], _gimm.Params[i])
	}

	if len(gimm.ReturnTypes()) != len(_gimm.ReturnTypes) {
		Panic(len(gimm.ReturnTypes()), len(_gimm.ReturnTypes))
	}
	for i := range _gimm.ReturnTypes {
		checkVariableMeta(gimm.ReturnTypes()[i], _gimm.ReturnTypes[i])
	}

	// unit test
	// b := MakeUnitTest(gmm)
	// fmt.Printf("unit test func:\n%v\n", string(b))
}

// func checkCallMeta(gcm *GoCallMeta, _gcm *compareGoCallMeta) {
// 	// basic
// 	if (gcm != nil) != (_gcm != nil) {
// 		Panic(gcm, _gcm)
// 	}
// 	if gcm == nil {
// 		return
// 	}
// 	if gcm.From() != _gcm.From {
// 		Panic(gcm.From(), _gcm.From)
// 	}
// 	if gcm.Call() != _gcm.Call {
// 		Panic(gcm.Call(), _gcm.Call)
// 	}

// 	// args
// 	if len(gcm.Args()) != len(_gcm.Args) {
// 		Panic(gcm.Args(), _gcm.Args)
// 	}
// 	args := gcm.Args()
// 	for index, _arg := range _gcm.Args {
// 		arg := args[index]
// 		if arg.Expression() != strings.TrimSpace(_arg.Expression) {
// 			Panic(arg.Expression(), _arg.Expression)
// 		}
// 		checkVariableMeta(arg.Head(), _arg.Head)
// 		fmt.Printf("Expression %v\n", _arg.Expression)
// 		for i, s := range arg.Slice() {
// 			fmt.Printf("slice index %v, s %+v, type %v\n", i, s, s.typeMeta.Expression())
// 		}
// 		fmt.Println()
// 		// if _arg.Expression == "module.NewExampleStruct(10)" {
// 		// 	arg.Slice()
// 		// }
// 	}
// }

func checkVariableMeta(gvm *GoVariableMeta, _gvm *compareGoVariableMeta) {
	if _gvm.Expression != "<IGNORE>" && gvm.Expression() != strings.TrimSpace(_gvm.Expression) {
		Panic(gvm.Expression(), _gvm.Expression)
	}
	if gvm.Name() != _gvm.Name {
		Panic(gvm.Name(), _gvm.Name)
	}
	te, tus, tue := gvm.Type()
	if te != _gvm.TypeExpression {
		Panic(te, _gvm.TypeExpression)
	}
	if tus != _gvm.TypeUnderlyingString {
		Panic(tus, _gvm.TypeUnderlyingString)
	}
	if tue != _gvm.TypeUnderlyingEnum {
		Panic(tue, _gvm.TypeUnderlyingEnum)
	}
	if gvm.Tag() != _gvm.Tag {
		Panic(gvm.Tag(), _gvm.Tag)
	}
	if gvm.Comment() != _gvm.Comment {
		Panic(gvm.Comment(), _gvm.Comment)
	}
	if !stp.Compare(gvm.Doc(), _gvm.Doc) {
		Panic(gvm.Doc(), _gvm.Doc)
	}
}

func checkImportMeta(gim *GoImportMeta, _gim *compareGoImportMeta) {
	if gim.Alias() != _gim.Alias {
		Panic(gim.Alias(), _gim.Alias)
	}
	if gim.Alias() != _gim.Alias {
		Panic(gim.Alias(), _gim.Alias)
	}
	if gim.Alias() != _gim.Alias {
		Panic(gim.Alias(), _gim.Alias)
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

// func TestReplaceGoProjectMeta(t *testing.T) {
// 	goProjectMeta, err := ExtractGoProjectMeta(standardProjectRelPath, standardProjectIgnorePathMap)
// 	if err != nil {
// 		panic(err)
// 	}

// 	for pkgName, replaceFunctionDoc := range replaceDoc {
// 		gpm, has := goProjectMeta.PackageMap[pkgName]
// 		if gpm == nil || !has {
// 			panic(pkgName)
// 		}
// 		for funcName, _replace := range replaceFunctionDoc {
// 			gpm.SearchFunctionMeta(funcName)
// 			gfm, has := gpm.pkgFunctionDecl[funcName]
// 			if gfm == nil || !has {
// 				panic(funcName)
// 			}
// 			checkDoc(gfm.Doc(), _replace.originDoc)
// 			originContent, replaceContent, err := gfm.ReplaceFunctionDoc(_replace.replaceDoc)
// 			if err != nil {
// 				panic(err)
// 			}
// 			if originContent != _replace.originContent {
// 				panic(originContent)
// 			}
// 			if replaceContent != _replace.replaceContent {
// 				panic(replaceContent)
// 			}
// 		}
// 	}
// }

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

// func TestParseGoCallMeta(t *testing.T) {
// 	for _, _gcm := range compareGoCallMetaSlice {
// 		gcm := ParseGoCallMeta(_gcm.Expression)
// 		gcm.PrintAST()

// 		if gcm.Expression() != _gcm.Expression {
// 			panic(gcm.Expression())
// 		}
// 		if gcm.Call() != _gcm.Call {
// 			panic(gcm.Call())
// 		}

// 		if len(gcm.Args()) != len(_gcm.Args) {
// 			panic(len(gcm.Args()))
// 		}
// 		for _, _arg := range _gcm.Args {
// 			for _, arg := range gcm.Args() {
// 				if reflect.DeepEqual(arg, _arg) {
// 					goto NEXT_PARAM
// 				}
// 			}
// 			panic(_arg)
// 		NEXT_PARAM:
// 		}
// 	}
// }
