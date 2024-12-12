package extractor

import (
	"fmt"
	"go/ast"
	"testing"

	stp "github.com/Mericusta/go-stp"
)

func TNilMetaPanic[T any](search string, meta *T) {
	if meta == nil {
		panic(fmt.Sprintf("search %v meta is nil", search))
	}
}

func TNotEqualPanic[T comparable](compared, value T) {
	if compared != value {
		panic(fmt.Sprintf("%+v != %+v", compared, value))
	}
}

func TSliceNotEqualPanic[CT, VT any](cs []CT, vs []VT, tCompare func(CT, VT)) {
	if len(cs) != len(vs) {
		panic(fmt.Sprintf("compare len %v not equal value len %v", len(cs), len(vs)))
	}
	for i, c := range cs {
		tCompare(c, vs[i])
	}
}

func TMapKeyNotExistPanic[CT any, VT any](cm map[string]*CT, vm map[string]*VT) {
	if len(cm) != len(vm) {
		panic(fmt.Sprintf("compare len %v not equal value len %v", len(cm), len(vm)))
	}
	for key := range cm {
		if value, has := vm[key]; value == nil || !has {
			panic(fmt.Sprintf("key %v value not exists", key))
		}
	}
}

func TMapNotEqualPanic[T comparable](cm, vm map[string]T) {
	if len(cm) != len(vm) {
		panic(fmt.Sprintf("compare len %v not equal value len %v", len(cm), len(vm)))
	}
	for key := range cm {
		TNotEqualPanic(cm[key], vm[key])
	}
}

type compareGoProjectMeta struct {
	absolutePath string
	moduleName   string
	packageMap   map[string]*compareGoPackageMeta
}

func (cgpm *compareGoProjectMeta) compare(gpm *GoProjectMeta) {
	TNotEqualPanic(cgpm.absolutePath, gpm.AbsolutePath())
	TNotEqualPanic(cgpm.moduleName, gpm.ModuleName())
	TMapKeyNotExistPanic(cgpm.packageMap, gpm.PackageMap())
}

type compareGoPackageMeta struct {
	ident                  string
	absolutePath           string
	importPath             string
	fileMetaMap            map[string]*compareGoFileMeta
	varMetaMap             map[string]*compareGoVarMeta[*ast.ValueSpec]
	funcMetaMap            map[string]*compareGoFuncMeta
	structMetaMap          map[string]*compareGoStructMeta
	interfaceMetaMap       map[string]*compareGoInterfaceMeta
	typeConstraintsMetaMap map[string]*compareGoTypeConstraintsMeta
}

func (cgpm *compareGoPackageMeta) compare(gpm *GoPackageMeta) {
	TNotEqualPanic(cgpm.ident, gpm.Ident())
	TNotEqualPanic(cgpm.absolutePath, gpm.AbsolutePath())
	TNotEqualPanic(cgpm.importPath, gpm.ImportPath())
	TMapKeyNotExistPanic(cgpm.fileMetaMap, gpm.FileMetaMap())
	TMapKeyNotExistPanic(cgpm.varMetaMap, gpm.VariableMetaMap())
	TMapKeyNotExistPanic(cgpm.funcMetaMap, gpm.FuncMetaMap())
	TMapKeyNotExistPanic(cgpm.structMetaMap, gpm.StructMetaMap())
	TMapKeyNotExistPanic(cgpm.interfaceMetaMap, gpm.InterfaceMetaMap())
	TMapKeyNotExistPanic(cgpm.typeConstraintsMetaMap, gpm.TypeConstraintsMetaMap()) // TODO:
}

type compareGoFileMeta struct {
	ident       string
	packageName string
}

func (cgfm *compareGoFileMeta) compare(gfm *GoFileMeta[*ast.File]) {
	TNotEqualPanic(cgfm.ident, gfm.Ident())
	TNotEqualPanic(cgfm.packageName, gfm.PackageName())
}

type compareGoVarMetaTypeConstraints interface {
	*ast.ValueSpec | *ast.Field
	ast.Node
}

type compareGoVarMeta[T compareGoVarMetaTypeConstraints] struct {
	expression     string
	ident          string
	typeIdent      string
	typeExpression string
	// TODO:
	TypeUnderlyingString string
	TypeUnderlyingEnum   UnderlyingType
	Tag                  string
	Comment              string
	Doc                  []string
}

func (cgvm *compareGoVarMeta[T]) compare(gvm *GoVarMeta[T]) {
	TNotEqualPanic(cgvm.expression, gvm.Expression())
	TNotEqualPanic(cgvm.ident, gvm.Ident())
	TNotEqualPanic(cgvm.typeIdent, gvm.TypeIdent())
	TNotEqualPanic(cgvm.typeExpression, gvm.TypeExpression())
}

type compareGoFuncDeclMeta struct {
	ident   string
	params  []*compareGoVarMeta[*ast.Field]
	returns []*compareGoVarMeta[*ast.Field]
	// TODO:
	Doc        []string
	TypeParams []*compareGoVarMeta[*ast.Field]
}

type iGoFuncDeclMeta interface {
	Ident() string
	Params() []*GoVarMeta[*ast.Field]
	Returns() []*GoVarMeta[*ast.Field]
}

func (cgfdm *compareGoFuncDeclMeta) compare(igfdm iGoFuncDeclMeta) {
	TNotEqualPanic(cgfdm.ident, igfdm.Ident())
	TSliceNotEqualPanic(cgfdm.params, igfdm.Params(), func(c *compareGoVarMeta[*ast.Field], v *GoVarMeta[*ast.Field]) { c.compare(v) })
	TSliceNotEqualPanic(cgfdm.returns, igfdm.Returns(), func(c *compareGoVarMeta[*ast.Field], v *GoVarMeta[*ast.Field]) { c.compare(v) })
}

type compareGoFuncMeta struct {
	compareGoFuncDeclMeta
	// TODO:
	callMeta map[string][]*compareGoCallMeta
	// VarMeta map[string]
}

func (cgfm *compareGoFuncMeta) compare(gfm *GoFuncMeta[*ast.FuncDecl]) {
	cgfm.compareGoFuncDeclMeta.compare(gfm)
}

type compareGoStructMeta struct {
	ident         string
	memberMetaMap map[string]*compareGoVarMeta[*ast.Field]
	methodMetaMap map[string]*compareGoMethodMeta
	// TODO:
	expression string
	doc        []string
	// typeParams       []*compareGoVarMeta
}

func (cgsm compareGoStructMeta) compare(gsm *GoStructMeta[*ast.TypeSpec]) {
	TNotEqualPanic(cgsm.ident, gsm.Ident())
	TMapKeyNotExistPanic(cgsm.memberMetaMap, gsm.MemberMetaMap())
	TMapKeyNotExistPanic(cgsm.methodMetaMap, gsm.MethodMetaMap())
}

type compareGoMethodMeta struct {
	*compareGoFuncMeta
	// TODO:
	RecvStruct      string
	PointerReceiver bool
}

func (cgmm *compareGoMethodMeta) compare(gmm *GoMethodMeta[*ast.FuncDecl]) {
	cgmm.compareGoFuncDeclMeta.compare(gmm)
}

type compareGoInterfaceMeta struct {
	ident         string
	methodMetaMap map[string]*compareGoInterfaceMethodMeta
	// TODO:
	Expression string
	Doc        []string
	// TypeParams []*compareGoVarMeta
}

func (cgim compareGoInterfaceMeta) compare(gim *GoInterfaceMeta[*ast.TypeSpec]) {
	TNotEqualPanic(cgim.ident, gim.Ident())
	TMapKeyNotExistPanic(cgim.methodMetaMap, gim.MethodMetaMap())
}

type compareGoInterfaceMethodMeta struct {
	compareGoFuncDeclMeta
}

func (cgimm *compareGoInterfaceMethodMeta) compare(gimm *GoInterfaceMethodMeta[*ast.Field, *ast.TypeSpec]) {
	cgimm.compareGoFuncDeclMeta.compare(gimm)
}

// type compareGoMemberMeta struct {
// 	Expression string
// 	MemberName string
// 	Tag        string
// 	Comment    string
// 	Doc        []string
// }

type compareGoTypeConstraintsMeta struct {
	ident string
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
	standardProjectAbsPath    = "/Users/dragonplus/Projects/github.com/Mericustar/go-extractor/testdata/standardProject"
	standardProjectModuleName = "standardProject"
	standardProjectMeta       = &compareGoProjectMeta{
		absolutePath: standardProjectAbsPath,
		moduleName:   standardProjectModuleName,
		packageMap: map[string]*compareGoPackageMeta{
			"main": {
				ident:        "main",
				absolutePath: stp.FormatFilePathWithOS(standardProjectAbsPath + "\\cmd"),
				fileMetaMap: map[string]*compareGoFileMeta{
					"main.go": {
						ident:       "main.go",
						packageName: "main",
					},
					"init.go": {
						ident:       "init.go",
						packageName: "main",
					},
				},
				varMetaMap: map[string]*compareGoVarMeta[*ast.ValueSpec]{
					"globalVariableInt": {
						expression:     `globalVariableInt        int    = 1`,
						ident:          "globalVariableInt",
						typeIdent:      "int",
						typeExpression: `int`,
					},
					"globalVariableString": {
						expression:     `globalVariableString     string = os.Getenv("ENV")`,
						ident:          "globalVariableString",
						typeIdent:      "string",
						typeExpression: `string`,
					},
					"globalVariableStruct": {
						expression:     `globalVariableStruct     *module.ExampleStruct`,
						ident:          "globalVariableStruct",
						typeIdent:      "ExampleStruct",
						typeExpression: `*module.ExampleStruct`,
					},
					"globalVariableTStruct": {
						expression:     `globalVariableTStruct    *template.TemplateStruct[int]`,
						ident:          "globalVariableTStruct",
						typeIdent:      "TemplateStruct",
						typeExpression: `*template.TemplateStruct[int]`,
					},
					"globalVariableInterface": {
						expression:     `globalVariableInterface  *pkgInterface.ExampleInterface`,
						ident:          "globalVariableInterface",
						typeIdent:      "ExampleInterface",
						typeExpression: `*pkgInterface.ExampleInterface`,
					},
					"globalVariableTInterface": {
						expression:     `globalVariableTInterface *pkgInterface.ExampleTemplateInterface[int]`,
						ident:          "globalVariableTInterface",
						typeIdent:      "ExampleTemplateInterface",
						typeExpression: `*pkgInterface.ExampleTemplateInterface[int]`,
					},
					"anotherGlobalVariableAny": {
						expression:     `anotherGlobalVariableAny interface{}`,
						ident:          "anotherGlobalVariableAny",
						typeIdent:      "interface{}",
						typeExpression: `interface{}`,
					},
				},
				funcMetaMap: map[string]*compareGoFuncMeta{
					"main": {
						compareGoFuncDeclMeta: compareGoFuncDeclMeta{
							ident: "main",
						},
						callMeta: map[string][]*compareGoCallMeta{
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
						compareGoFuncDeclMeta: compareGoFuncDeclMeta{
							ident: "Init",
						},
					},
				},
			},
			standardProjectModuleName + "/pkg": {
				ident:        "pkg",
				absolutePath: stp.FormatFilePathWithOS(standardProjectAbsPath + "\\pkg"),
				importPath:   standardProjectModuleName + "/pkg",
				fileMetaMap: map[string]*compareGoFileMeta{
					"pkg.go": {
						ident:       "pkg.go",
						packageName: "pkg",
					},
				},
				funcMetaMap: map[string]*compareGoFuncMeta{
					"ExampleFunc": {
						compareGoFuncDeclMeta: compareGoFuncDeclMeta{
							ident: "ExampleFunc",
							Doc: []string{
								"// ExampleFunc this is example function",
							},
							params: []*compareGoVarMeta[*ast.Field]{
								{
									expression:     `s *module.ExampleStruct`,
									ident:          "s",
									typeIdent:      "ExampleStruct",
									typeExpression: "*module.ExampleStruct",
								},
							},
						},
						callMeta: map[string][]*compareGoCallMeta{
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
						compareGoFuncDeclMeta: compareGoFuncDeclMeta{
							ident: "NoDocExampleFunc",
							params: []*compareGoVarMeta[*ast.Field]{
								{
									expression:     `s *module.ExampleStruct`,
									ident:          "s",
									typeIdent:      "ExampleStruct",
									typeExpression: "*module.ExampleStruct",
								},
							},
						},
						callMeta: map[string][]*compareGoCallMeta{
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
						compareGoFuncDeclMeta: compareGoFuncDeclMeta{
							ident: "OneLineDocExampleFunc",
							Doc: []string{
								"// OneLineDocExampleFunc this is one-line-doc example function",
							},
							params: []*compareGoVarMeta[*ast.Field]{
								{
									expression:     `s *module.ExampleStruct`,
									ident:          "s",
									typeIdent:      "ExampleStruct",
									typeExpression: "*module.ExampleStruct",
								},
							},
						},
						callMeta: map[string][]*compareGoCallMeta{
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
						compareGoFuncDeclMeta: compareGoFuncDeclMeta{
							ident: "ImportSelectorFunc",
							params: []*compareGoVarMeta[*ast.Field]{
								{
									expression:     `s *module.ExampleStruct`,
									ident:          "s",
									typeIdent:      "ExampleStruct",
									typeExpression: "*module.ExampleStruct",
								},
							},
						},
						callMeta: map[string][]*compareGoCallMeta{
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
				structMetaMap: map[string]*compareGoStructMeta{
					"ExampleTemplateStructWithTemplateParent": {
						ident: "ExampleTemplateStructWithTemplateParent",
						// typeParams: []*compareGoVarMeta[*ast.Field]{
						// 	{
						// 		expression:           `T any`,
						// 		ident:                "T",
						// 		TypeExpression:       "any",
						// 	},
						// },
						memberMetaMap: map[string]*compareGoVarMeta[*ast.Field]{
							"TemplateStruct": {
								expression:     `*template.TemplateStruct[map[string]*template.TemplateStruct[*T]]`,
								ident:          "TemplateStruct",
								typeIdent:      "TemplateStruct",
								typeExpression: `*template.TemplateStruct[map[string]*template.TemplateStruct[*T]]`,
							},
						},
					},
				},
				interfaceMetaMap: map[string]*compareGoInterfaceMeta{
					"ExampleTemplateInterfaceWithTypeConstraints": &compareGoInterfaceMeta{
						ident: "ExampleTemplateInterfaceWithTypeConstraints",
						methodMetaMap: map[string]*compareGoInterfaceMethodMeta{
							"Parse": {
								compareGoFuncDeclMeta{
									ident: "Parse",
									params: []*compareGoVarMeta[*ast.Field]{
										{
											expression:     `T`,
											ident:          "",
											typeIdent:      "T",
											typeExpression: `T`,
										},
									},
								},
							},
							"Format": &compareGoInterfaceMethodMeta{
								compareGoFuncDeclMeta{
									ident: "Format",
									returns: []*compareGoVarMeta[*ast.Field]{
										{
											expression:     `T`,
											ident:          "",
											typeIdent:      "T",
											typeExpression: `T`,
										},
									},
								},
							},
						},
					},
				},
			},
			standardProjectModuleName + "/pkg/pkgInterface": {
				ident:        "pkgInterface",
				absolutePath: stp.FormatFilePathWithOS(standardProjectAbsPath + "\\pkg\\interface"),
				importPath:   standardProjectModuleName + "/pkg/pkgInterface",
				fileMetaMap: map[string]*compareGoFileMeta{
					"interface.go": {
						ident:       "interface.go",
						packageName: "pkgInterface",
					},
				},
				interfaceMetaMap: map[string]*compareGoInterfaceMeta{
					"ExampleInterface": {
						ident: "ExampleInterface",
						methodMetaMap: map[string]*compareGoInterfaceMethodMeta{
							"ExampleFunc": {
								compareGoFuncDeclMeta{
									ident: "ExampleFunc",
									Doc:   []string{"// This is ExampleFunc Doc"},
									params: []*compareGoVarMeta[*ast.Field]{
										{
											expression:     `int`,
											ident:          "",
											typeIdent:      "int",
											typeExpression: `int`,
										},
									},
								},
							},
							"AnotherExampleFunc": {
								compareGoFuncDeclMeta{
									ident: "AnotherExampleFunc",
									Doc:   []string{"// This is AnotherExampleFunc Doc"},
									params: []*compareGoVarMeta[*ast.Field]{
										{
											expression:     `int`,
											ident:          "",
											typeIdent:      "int",
											typeExpression: `int`,
										},
										{
											expression:     `[]int`,
											ident:          "",
											typeIdent:      "[]int",
											typeExpression: `[]int`,
										},
									},
									returns: []*compareGoVarMeta[*ast.Field]{
										{
											expression:     `int`,
											ident:          "",
											typeIdent:      "int",
											typeExpression: `int`,
										},
										{
											expression:     `[]int`,
											ident:          "",
											typeIdent:      "[]int",
											typeExpression: `[]int`,
										},
									},
								},
							},
						},
					},
					"ExampleTemplateInterface": {
						ident: "ExampleTemplateInterface",
						// TypeParams: []*compareGoVarMeta[*ast.Field]{
						// 	{
						// 		expression:           `T any`,
						// 		ident:                "T",
						// 		TypeExpression:       "any",
						// 	},
						// },
						methodMetaMap: map[string]*compareGoInterfaceMethodMeta{
							"ExampleFunc": {
								compareGoFuncDeclMeta{
									ident: "ExampleFunc",
									Doc:   []string{"// This is ExampleFunc Doc"},
									TypeParams: []*compareGoVarMeta[*ast.Field]{
										{
											expression:     `T any`,
											ident:          "T",
											typeIdent:      "T",
											typeExpression: `any`,
										},
									},
									params: []*compareGoVarMeta[*ast.Field]{
										{
											expression:     `T`,
											ident:          "",
											typeIdent:      "T",
											typeExpression: `T`,
										},
									},
								},
							},
							"AnotherExampleFunc": {
								compareGoFuncDeclMeta{
									ident: "AnotherExampleFunc",
									Doc:   []string{"// This is AnotherExampleFunc Doc"},
									TypeParams: []*compareGoVarMeta[*ast.Field]{
										{
											expression:     `T any`,
											ident:          "T",
											typeIdent:      "any",
											typeExpression: `any`,
										},
									},
									params: []*compareGoVarMeta[*ast.Field]{
										{
											expression:     `T`,
											ident:          "",
											typeIdent:      "T",
											typeExpression: `T`,
										},
										{
											expression:     `[]T`,
											ident:          "",
											typeIdent:      "[]T",
											typeExpression: `[]T`,
										},
									},
									returns: []*compareGoVarMeta[*ast.Field]{
										{
											expression:     `T`,
											ident:          "",
											typeIdent:      "T",
											typeExpression: `T`,
										},
										{
											expression:     `[]T`,
											ident:          "",
											typeIdent:      "[]T",
											typeExpression: `[]T`,
										},
									},
								},
							},
						},
					},
				},
			},
			standardProjectModuleName + "/pkg/module": {
				ident:        "module",
				absolutePath: stp.FormatFilePathWithOS(standardProjectAbsPath + "\\pkg\\module"),
				importPath:   standardProjectModuleName + "/pkg/module",
				fileMetaMap: map[string]*compareGoFileMeta{
					"module.go": {
						ident:       "module.go",
						packageName: "module",
					},
				},
				varMetaMap: map[string]*compareGoVarMeta[*ast.ValueSpec]{
					"globalExampleStruct": {
						expression:     `globalExampleStruct *ExampleStruct`,
						ident:          "globalExampleStruct",
						typeIdent:      "ExampleStruct",
						typeExpression: "*ExampleStruct",
					},
				},
				structMetaMap: map[string]*compareGoStructMeta{
					"ParentStruct": {
						ident: "ParentStruct",
						memberMetaMap: map[string]*compareGoVarMeta[*ast.Field]{
							"p": {
								expression:     `p int`,
								ident:          "p",
								typeIdent:      "int",
								typeExpression: `int`,
								Comment:        "// parent value",
							},
						},
						methodMetaMap: map[string]*compareGoMethodMeta{
							"P": {
								compareGoFuncMeta: &compareGoFuncMeta{
									compareGoFuncDeclMeta: compareGoFuncDeclMeta{
										ident: "P",
										returns: []*compareGoVarMeta[*ast.Field]{
											&compareGoVarMeta[*ast.Field]{
												expression:     `int`,
												ident:          "",
												typeIdent:      "int",
												typeExpression: `int`,
											},
										},
									},
								},
								RecvStruct:      "ParentStruct",
								PointerReceiver: true,
							},
						},
					},
					"ExampleStruct": {
						ident: "ExampleStruct",
						doc: []string{
							"// ExampleStruct this is an example struct",
							"// this is struct comment",
							"// this is another struct comment",
						},
						memberMetaMap: map[string]*compareGoVarMeta[*ast.Field]{
							"ParentStruct": {
								expression:     `*ParentStruct`,
								ident:          "ParentStruct",
								typeIdent:      "ParentStruct",
								typeExpression: `*ParentStruct`,
								Comment:        "// parent struct",
							},
							"v": {
								expression:     "v   int `ast:init,default=1`",
								ident:          "v",
								typeIdent:      "int",
								typeExpression: `int`,
								Tag:            "`ast:init,default=1`",
								Doc: []string{
									"// v this is member doc line1",
									"// v this is member doc line2",
								},
								Comment: "// this is member single comment line",
							},
							"sub": {
								expression:     `sub *ExampleStruct`,
								ident:          "sub",
								typeIdent:      "ExampleStruct",
								typeExpression: `*ExampleStruct`,
							},
						},
						methodMetaMap: map[string]*compareGoMethodMeta{
							"ExampleFunc": {
								compareGoFuncMeta: &compareGoFuncMeta{
									compareGoFuncDeclMeta: compareGoFuncDeclMeta{
										ident: "ExampleFunc",
										params: []*compareGoVarMeta[*ast.Field]{
											{
												expression:     `v int`,
												ident:          "v",
												typeIdent:      "int",
												typeExpression: `int`,
											},
										},
									},
									callMeta: map[string][]*compareGoCallMeta{
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
								compareGoFuncMeta: &compareGoFuncMeta{
									compareGoFuncDeclMeta: compareGoFuncDeclMeta{
										ident: "ExampleFuncWithPointerReceiver",
										params: []*compareGoVarMeta[*ast.Field]{
											{
												expression:     `v int`,
												ident:          "v",
												typeIdent:      "int",
												typeExpression: `int`,
											},
										},
									},
									callMeta: map[string][]*compareGoCallMeta{
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
								compareGoFuncMeta: &compareGoFuncMeta{
									compareGoFuncDeclMeta: compareGoFuncDeclMeta{
										ident: "DoubleReturnFunc",
										returns: []*compareGoVarMeta[*ast.Field]{
											{
												expression:     `int`,
												ident:          "",
												typeIdent:      "int",
												typeExpression: `int`,
											},
											{
												expression:     `int`,
												ident:          "",
												typeIdent:      "int",
												typeExpression: `int`,
											},
										},
									},
									callMeta: map[string][]*compareGoCallMeta{
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
								compareGoFuncMeta: &compareGoFuncMeta{
									compareGoFuncDeclMeta: compareGoFuncDeclMeta{
										ident: "V",
										returns: []*compareGoVarMeta[*ast.Field]{
											{
												expression:     `int`,
												ident:          "",
												typeIdent:      "int",
												typeExpression: `int`,
											},
										},
									},
								},
								RecvStruct:      "ExampleStruct",
								PointerReceiver: false,
							},
							"Sub": {
								compareGoFuncMeta: &compareGoFuncMeta{
									compareGoFuncDeclMeta: compareGoFuncDeclMeta{
										ident: "Sub",
										returns: []*compareGoVarMeta[*ast.Field]{
											{
												expression:     `*ExampleStruct`,
												ident:          "",
												typeIdent:      "ExampleStruct",
												typeExpression: "*ExampleStruct",
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
				funcMetaMap: map[string]*compareGoFuncMeta{
					"NewExampleStruct": {
						compareGoFuncDeclMeta: compareGoFuncDeclMeta{
							ident: "NewExampleStruct",
							Doc: []string{
								"// NewExampleStruct this is new example struct",
								"// @param           value",
								"// @return          pointer to ExampleStruct",
							},
							params: []*compareGoVarMeta[*ast.Field]{
								{
									expression:     `v int`,
									ident:          "v",
									typeIdent:      "int",
									typeExpression: `int`,
								},
							},
							returns: []*compareGoVarMeta[*ast.Field]{
								{
									expression:     `*ExampleStruct`,
									ident:          "",
									typeIdent:      "ExampleStruct",
									typeExpression: "*ExampleStruct",
								},
							},
						},
						callMeta: map[string][]*compareGoCallMeta{
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
						compareGoFuncDeclMeta: compareGoFuncDeclMeta{
							ident: "ExampleFunc",
							params: []*compareGoVarMeta[*ast.Field]{
								{
									expression:     `s *ExampleStruct`,
									ident:          "s",
									typeIdent:      "ExampleStruct",
									typeExpression: "*ExampleStruct",
								},
							},
						},
						callMeta: map[string][]*compareGoCallMeta{
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
				ident:        "template",
				absolutePath: stp.FormatFilePathWithOS(standardProjectAbsPath + "\\pkg\\template"),
				importPath:   standardProjectModuleName + "/pkg/template",
				fileMetaMap: map[string]*compareGoFileMeta{
					"template.go": {
						ident:       "template.go",
						packageName: "template",
					},
				},
				funcMetaMap: map[string]*compareGoFuncMeta{
					"OneTemplateFunc": {
						compareGoFuncDeclMeta: compareGoFuncDeclMeta{
							ident: "OneTemplateFunc",
							TypeParams: []*compareGoVarMeta[*ast.Field]{
								{
									expression:     `T any`,
									ident:          "T",
									typeIdent:      "any",
									typeExpression: "any",
								},
							},
							params: []*compareGoVarMeta[*ast.Field]{
								{
									expression:     `tv *T`,
									ident:          "tv",
									typeIdent:      "T",
									typeExpression: "*T",
								},
							},
							returns: []*compareGoVarMeta[*ast.Field]{
								{
									expression:     `*T`,
									ident:          "",
									typeIdent:      "T",
									typeExpression: "*T",
								},
							},
						},
					},
					"DoubleSameTemplateFunc": {
						compareGoFuncDeclMeta: compareGoFuncDeclMeta{
							ident: "DoubleSameTemplateFunc",
							TypeParams: []*compareGoVarMeta[*ast.Field]{
								{
									expression:     `T1, T2 any`,
									ident:          "T1",
									typeIdent:      "any",
									typeExpression: "any",
								},
								{
									expression:     `T1, T2 any`,
									ident:          "T2",
									typeIdent:      "any",
									typeExpression: "any",
								},
							},
							params: []*compareGoVarMeta[*ast.Field]{
								{
									expression:     `tv1 T1`,
									ident:          "tv1",
									typeIdent:      "T1",
									typeExpression: "T1",
								},
								{
									expression:     `tv2 T2`,
									ident:          "tv2",
									typeIdent:      "T2",
									typeExpression: "T2",
								},
							},
							returns: []*compareGoVarMeta[*ast.Field]{
								{
									expression:     `*T1`,
									ident:          "",
									typeIdent:      "T1",
									typeExpression: "*T1",
								},
								{
									expression:     `*T2`,
									ident:          "",
									typeIdent:      "T2",
									typeExpression: "*T2",
								},
							},
						},
					},
					"DoubleDifferenceTemplateFunc": {
						compareGoFuncDeclMeta: compareGoFuncDeclMeta{
							ident: "DoubleDifferenceTemplateFunc",
							TypeParams: []*compareGoVarMeta[*ast.Field]{
								{
									expression:     `T1 any`,
									ident:          "T1",
									typeIdent:      "any",
									typeExpression: "any",
								},
								{
									expression:     `T2 comparable`,
									ident:          "T2",
									typeIdent:      "comparable",
									typeExpression: "comparable",
								},
							},
							params: []*compareGoVarMeta[*ast.Field]{
								{
									expression:     `tv1 T1`,
									ident:          "tv1",
									typeIdent:      "T1",
									typeExpression: "T1",
								},
								{
									expression:     `tv2 T2`,
									ident:          "tv2",
									typeIdent:      "T2",
									typeExpression: "T2",
								},
							},
							returns: []*compareGoVarMeta[*ast.Field]{
								{
									expression:     `*T1`,
									ident:          "",
									typeIdent:      "T1",
									typeExpression: "*T1",
								},
								{
									expression:     `*T2`,
									ident:          "",
									typeIdent:      "T2",
									typeExpression: "*T2",
								},
							},
						},
					},
					"TypeConstraintsTemplateFunc": {
						compareGoFuncDeclMeta: compareGoFuncDeclMeta{
							ident: "TypeConstraintsTemplateFunc",
							TypeParams: []*compareGoVarMeta[*ast.Field]{
								{
									expression:     `T TypeConstraints`,
									ident:          "T",
									typeIdent:      "TypeConstraints",
									typeExpression: `TypeConstraints`,
								},
							},
							params: []*compareGoVarMeta[*ast.Field]{
								{
									expression:     `tv T`,
									ident:          "tv",
									typeIdent:      "T",
									typeExpression: "T",
								},
							},
							returns: []*compareGoVarMeta[*ast.Field]{
								{
									expression:     `*T`,
									ident:          "",
									typeIdent:      "T",
									typeExpression: "*T",
								},
							},
						},
					},
					"CannotInferTypeFunc1": {
						compareGoFuncDeclMeta: compareGoFuncDeclMeta{
							ident: "CannotInferTypeFunc1",
							TypeParams: []*compareGoVarMeta[*ast.Field]{
								{
									expression:     `T any`,
									ident:          "T",
									typeIdent:      "any",
									typeExpression: "any",
								},
							},
						},
					},
					"CannotInferTypeFunc2": {
						compareGoFuncDeclMeta: compareGoFuncDeclMeta{
							ident: "CannotInferTypeFunc2",
							TypeParams: []*compareGoVarMeta[*ast.Field]{
								{
									expression:     `K comparable`,
									ident:          "K",
									typeIdent:      "comparable",
									typeExpression: "comparable",
								},
								{
									expression:     `V any`,
									ident:          "V",
									typeIdent:      "any",
									typeExpression: "any",
								},
							},
							returns: []*compareGoVarMeta[*ast.Field]{
								{
									expression:     `*K`,
									ident:          "",
									typeExpression: "*K",
									typeIdent:      "K",
								},
								{
									expression:     `*V`,
									ident:          "",
									typeExpression: "*V",
									typeIdent:      "V",
								},
							},
						},
					},
				},
				structMetaMap: map[string]*compareGoStructMeta{
					"TemplateStruct": {
						ident: "TemplateStruct",
						// typeParams: []*compareGoVarMeta[*ast.Field]{
						// 	{
						// 		expression:           `T any`,
						// 		ident:                "T",
						// 		TypeExpression:       "any",
						// 	},
						// },
						memberMetaMap: map[string]*compareGoVarMeta[*ast.Field]{
							"v": {
								expression:     `v T`,
								ident:          "v",
								typeIdent:      "T",
								typeExpression: `T`,
							},
						},
						methodMetaMap: map[string]*compareGoMethodMeta{
							"V": {
								compareGoFuncMeta: &compareGoFuncMeta{
									compareGoFuncDeclMeta: compareGoFuncDeclMeta{
										ident: "V",
										returns: []*compareGoVarMeta[*ast.Field]{
											{
												expression:     `T`,
												ident:          "",
												typeIdent:      "T",
												typeExpression: `T`,
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
						ident: "TwoTypeTemplateStruct",
						// typeParams: []*compareGoVarMeta[*ast.Field]{
						// 	{
						// 		expression:           `K TypeConstraints`,
						// 		ident:                "K",
						// 		TypeExpression:       "TypeConstraints",
						// 	},
						// 	{
						// 		expression:           `V any`,
						// 		ident:                "V",
						// 		TypeExpression:       "any",
						// 	},
						// },
						memberMetaMap: map[string]*compareGoVarMeta[*ast.Field]{
							"v": {
								expression:     `v map[K]V`,
								ident:          "v",
								typeIdent:      "map[K]V",
								typeExpression: `map[K]V`,
							},
						},
						methodMetaMap: map[string]*compareGoMethodMeta{
							"KVSlice": {
								compareGoFuncMeta: &compareGoFuncMeta{
									compareGoFuncDeclMeta: compareGoFuncDeclMeta{
										ident: "KVSlice",
										params: []*compareGoVarMeta[*ast.Field]{
											{
												expression:     `k K`,
												ident:          "k",
												typeIdent:      "K",
												typeExpression: `K`,
											},
											{
												expression:     `v V`,
												ident:          "v",
												typeIdent:      "V",
												typeExpression: `V`,
											},
										},
										returns: []*compareGoVarMeta[*ast.Field]{
											{
												expression:     `[]K`,
												ident:          "",
												typeIdent:      "[]K",
												typeExpression: `[]K`,
											},
											{
												expression:     `[]V`,
												ident:          "",
												typeIdent:      "[]V",
												typeExpression: `[]V`,
											},
										},
									},
								},
								RecvStruct:      `TwoTypeTemplateStruct`,
								PointerReceiver: true,
							},
						},
					},
					"TemplateStructWithParent": {
						ident: "TemplateStructWithParent",
						// typeParams: []*compareGoVarMeta[*ast.Field]{
						// 	{
						// 		expression:           `T any`,
						// 		ident:                "T",
						// 		TypeExpression:       "any",
						// 	},
						// },
						memberMetaMap: map[string]*compareGoVarMeta[*ast.Field]{
							"TemplateStruct": {
								expression:     `*TemplateStruct[T]`,
								ident:          "TemplateStruct",
								typeIdent:      "TemplateStruct",
								typeExpression: `*TemplateStruct[T]`,
							},
						},
					},
				},
			},
		},
	}
)

// 存在父级的 meta 数据，可以通过父级 meta 的子节点获得，也可以直接提取对应的文件的 meta 接口获得

func TestExtractGoProjectMeta(t *testing.T) {
	// 根据 标准项目 的相对路径 `提取` 标准项目
	goProjectMeta, err := ExtractGoProjectMeta(standardProjectRelPath, standardProjectIgnorePathMap)
	if err != nil {
		panic(err)
	}

	// 比较 标准项目 的 meta 数据
	standardProjectMeta.compare(goProjectMeta)

	// 逐个比较 package 的 meta 数据
	for comparePackageImportPath, cgpm := range standardProjectMeta.packageMap {
		// 在 项目 的 meta 数据中，根据 package 的 导入路径 `搜索` package 的 meta 数据
		gpm := goProjectMeta.SearchPackageMeta(comparePackageImportPath)
		TNilMetaPanic(comparePackageImportPath, gpm)

		// 比较 package 的 meta 数据
		cgpm.compare(gpm)

		// 输出所有文件的 AST
		for _, gfm := range gpm.fileMetaMap {
			gfm.OutputAST()
		}

		// 逐个比较 文件 的 meta 数据
		for compareFileName, cgfm := range cgpm.fileMetaMap {
			// 在 package 的 meta 数据中，根据 文件名 搜索 文件 的 meta 数据
			gfm := gpm.SearchFileMeta(compareFileName)
			TNilMetaPanic(compareFileName, gfm)
			cgfm.compare(gfm)
		}

		// 逐个比较 var 的 meta 数据
		for compareVarIdent, cgvm := range cgpm.varMetaMap {
			// 在 package 的 meta 数据中，根据 var 的 标识 搜索 var 的 meta 数据
			gvm := gpm.SearchVarMeta(compareVarIdent)
			TNilMetaPanic(compareVarIdent, gvm)
			cgvm.compare(gvm)
		}

		// 逐个比较 func 的 meta 数据
		for compareFuncName, cgfm := range cgpm.funcMetaMap {
			// 在 package 的 meta 数据中，根据 func 的 表示 搜索 func 的 meta 数据
			gfm := gpm.SearchFuncMeta(compareFuncName)
			TNilMetaPanic(compareFuncName, gfm)
			cgfm.compare(gfm)

			// make
			params := make([]*GoVarMeta[*ast.Field], 0, 8)
			returns := make([]*GoVarMeta[*ast.Field], 0, 8)
			for _, p := range gfm.Params() {
				ngvm := MakeUpVarMeta(p.ident, p.TypeExpression())
				if ngvm != nil {
					params = append(params, ngvm)
				}
			}
			for _, r := range gfm.Returns() {
				ngvm := MakeUpVarMeta(r.ident, r.TypeExpression())
				if ngvm != nil {
					returns = append(returns, ngvm)
				}
			}
			ngfm := MakeUpFuncMeta(compareFuncName, params, returns)
			ngfm.SetBlockStmt("{}")
			fmt.Printf("\n%v\n", ngfm.Make())

			// 	// unit test
			// 	var unittestFuncName string
			// 	var unittestByte []byte
			// 	if l := len(gfm.TypeParams()); l == 0 {
			// 		unittestFuncName, unittestByte = gfm.MakeUnitTest(nil)
			// 	} else {
			// 		testTypeArgs := []string{"string", "[]string", "map[string]string"}
			// 		typeArgs := make([]string, 0, l)
			// 		for i := 0; i < l; i++ {
			// 			typeArgs = append(typeArgs, testTypeArgs[i%len(testTypeArgs)])
			// 		}
			// 		unittestFuncName, unittestByte = gfm.MakeUnitTest(typeArgs)
			// 	}
			// 	fmt.Printf("unit test func %v:\n%v\n", unittestFuncName, string(unittestByte))

			// 	// benchmark
			// 	var benchmarkFuncName string
			// 	var benchmarkByte []byte
			// 	if l := len(gfm.TypeParams()); l == 0 {
			// 		benchmarkFuncName, benchmarkByte = gfm.MakeBenchmark(nil)
			// 	} else {
			// 		testTypeArgs := []string{"string", "[]string", "map[string]string"}
			// 		typeArgs := make([]string, 0, l)
			// 		for i := 0; i < l; i++ {
			// 			typeArgs = append(typeArgs, testTypeArgs[i%len(testTypeArgs)])
			// 		}
			// 		benchmarkFuncName, benchmarkByte = gfm.MakeBenchmark(typeArgs)
			// 	}
			// 	fmt.Printf("benchmark func %v:\n%v\n", benchmarkFuncName, string(benchmarkByte))

			// 	testFileByte := MakeTestFile(fmt.Sprintf("%v_test.go", strings.Trim(gfm.path, ".go")), nil)
			// 	fmt.Printf("unit test file:\n%v\n", string(testFileByte))
		}

		// 逐个比较 struct 的 meta 数据
		for compareStructIdent, cgsm := range cgpm.structMetaMap {
			// 在 package 的 meta 数据中，根据 struct 的 标识 搜索 struct 的 meta 数据
			gsm := gpm.SearchStructMeta(compareStructIdent)
			TNilMetaPanic(compareStructIdent, gsm)
			cgsm.compare(gsm)

			// 逐个比较 struct 的 member 数据
			for compareStructMemberIdent, cgsmm := range cgsm.memberMetaMap {
				// 在 struct 的 meta 数据中，根据 member 的 标识 搜索 member 的 meta 数据
				gsmm := gsm.SearchMemberMeta(compareStructMemberIdent)
				TNilMetaPanic(compareStructMemberIdent, gsmm)
				cgsmm.compare(gsmm)
			}

			// 逐个比较 struct 的 method 数据
			for compareStructMethodIdent, cgmm := range cgsm.methodMetaMap {
				// 在 struct 的 meta 数据中，根据 method 的 标识 搜索 method 的 meta 数据
				gmm := gsm.SearchMethodMeta(compareStructMethodIdent)
				TNilMetaPanic(compareStructMethodIdent, gmm)
				cgmm.compare(gmm)

				// // unit test
				// var unittestFuncName string
				// var unittestByte []byte
				// if l := len(gmm.TypeParams()); l == 0 {
				// 	unittestFuncName, unittestByte = gmm.MakeUnitTest(nil)
				// } else {
				// 	testTypeArgs := []string{"string", "[]string", "map[string]string"}
				// 	typeArgs := make([]string, 0, l)
				// 	for i := 0; i < l; i++ {
				// 		typeArgs = append(typeArgs, testTypeArgs[i%len(testTypeArgs)])
				// 	}
				// 	unittestFuncName, unittestByte = gmm.MakeUnitTest(typeArgs)
				// }
				// fmt.Printf("unit test func %v:\n%v\n", unittestFuncName, string(unittestByte))

				// // benchmark
				// var benchmarkFuncName string
				// var benchmarkByte []byte
				// if l := len(gmm.TypeParams()); l == 0 {
				// 	benchmarkFuncName, benchmarkByte = gmm.MakeBenchmark(nil)
				// } else {
				// 	testTypeArgs := []string{"string", "[]string", "map[string]string"}
				// 	typeArgs := make([]string, 0, l)
				// 	for i := 0; i < l; i++ {
				// 		typeArgs = append(typeArgs, testTypeArgs[i%len(testTypeArgs)])
				// 	}
				// 	benchmarkFuncName, benchmarkByte = gmm.MakeBenchmark(typeArgs)
				// }
				// fmt.Printf("benchmark func %v:\n%v\n", benchmarkFuncName, string(benchmarkByte))
			}
		}

		// 逐个比较 interface 的 meta 数据
		for compareInterfaceIdent, cgim := range cgpm.interfaceMetaMap {
			// 在 package 的 meta 数据中，根据 interface 的 标识 搜索 interface 的 meta 数据
			gim := gpm.SearchInterfaceMeta(compareInterfaceIdent)
			TNilMetaPanic(compareInterfaceIdent, gim)
			cgim.compare(gim)

			// 逐个比较 interface 的 method 数据
			for compareInterfaceMethodIdent, cgimm := range cgim.methodMetaMap {
				// 在 interface 的 meta 数据中，根据 method 的 标识 搜索 method 的 meta 数据
				gimm := gim.SearchMethodMeta(compareInterfaceMethodIdent)
				TNilMetaPanic(compareInterfaceMethodIdent, gimm)
				cgimm.compare(gimm)

				// // unit test
				// var unittestFuncName string
				// var unittestByte []byte
				// if l := len(gimm.TypeParams()); l == 0 {
				// 	unittestFuncName, unittestByte = gimm.MakeUnitTest(nil)
				// } else {
				// 	testTypeArgs := []string{"string", "[]string", "map[string]string"}
				// 	typeArgs := make([]string, 0, l)
				// 	for i := 0; i < l; i++ {
				// 		typeArgs = append(typeArgs, testTypeArgs[i%len(testTypeArgs)])
				// 	}
				// 	unittestFuncName, unittestByte = gimm.MakeUnitTest(typeArgs)
				// }
				// fmt.Printf("unit test func %v:\n%v\n", unittestFuncName, string(unittestByte))

				// // benchmark
				// var benchmarkFuncName string
				// var benchmarkByte []byte
				// if l := len(gimm.TypeParams()); l == 0 {
				// 	benchmarkFuncName, benchmarkByte = gimm.MakeBenchmark(nil)
				// } else {
				// 	testTypeArgs := []string{"string", "[]string", "map[string]string"}
				// 	typeArgs := make([]string, 0, l)
				// 	for i := 0; i < l; i++ {
				// 		typeArgs = append(typeArgs, testTypeArgs[i%len(testTypeArgs)])
				// 	}
				// 	benchmarkFuncName, benchmarkByte = gimm.MakeBenchmark(typeArgs)
				// }
				// fmt.Printf("benchmark func %v:\n%v\n", benchmarkFuncName, string(benchmarkByte))

				// // implement
				// var implementFuncName string
				// var implementFuncMeta *GoFunctionMeta
				// receiverIdent, receiverType := strings.ToLower(string(gimm.interfaceMeta.InterfaceName()[0])), gimm.interfaceMeta.InterfaceName()+"Implement"
				// if l := len(gimm.TypeParams()); l == 0 {
				// 	implementFuncName, implementFuncMeta = gimm.MakeImplementMethodMeta(receiverIdent, receiverType)
				// } else {
				// 	// TODO: support
				// 	// testTypeArgs := []string{"string", "[]string", "map[string]string"}
				// 	// typeArgs := make([]string, 0, l)
				// 	// for i := 0; i < l; i++ {
				// 	// 	typeArgs = append(typeArgs, testTypeArgs[i%len(testTypeArgs)])
				// 	// }
				// 	// benchmarkFuncName, benchmarkByte = gimm.MakeBenchmark(typeArgs)
				// }
				// if implementFuncMeta != nil {
				// 	fmt.Printf("implement func %v:\n%v\n", implementFuncName, implementFuncMeta.Format())
				// }
			}
		}
	}
}
