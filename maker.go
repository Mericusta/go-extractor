package extractor

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
)

type GoUnitTestMaker interface {
	UnitTestFuncName() string
	FunctionName() string
	TypeParams() []*GoVariableMeta
	Params() []*GoVariableMeta
	ReturnTypes() []*GoVariableMeta
	Recv() *GoVariableMeta
}

// makeUnitTest generate unit test func
// @param1 *GoFunctionMeta/*GoMethodMeta
// @param2 type args for type params if needs
// @return func declaration content
func makeUnitTest(gutm GoUnitTestMaker, typeArgs []string) []byte {
	unitTestFuncName := gutm.UnitTestFuncName()
	funcName := gutm.FunctionName()
	typeParams := gutm.TypeParams()
	if len(typeParams) > len(typeArgs) {
		return nil
	}
	if len(typeArgs) > 0 {
		suffix := ""
		for _, typeArg := range typeArgs {
			if len(typeArg) > 0 {
				suffix = fmt.Sprintf("%v_%v", suffix, typeArg)
			}
		}
		if len(suffix) > 0 {
			hashStr := md5.Sum([]byte(suffix))
			unitTestFuncName = fmt.Sprintf("%v_%v", unitTestFuncName, hex.EncodeToString(hashStr[:]))
		}
	}
	recv := gutm.Recv()
	params := gutm.Params()
	returnTypes := gutm.ReturnTypes()
	returnTypeLen := len(returnTypes)

	funcBodyBlockStmt := &ast.BlockStmt{}
	// arg struct decl
	if len(params) > 0 {
		argStructDecl := newTypeSpec("args", nil, params).makeDecl()
		for paramIndex, param := range params {
			for typeParamIndex, typeParam := range typeParams {
				isTypeParam := false
				ast.Inspect(param.typeNode(), func(n ast.Node) bool {
					ident, ok := n.(*ast.Ident)
					if ident != nil && ok && ident.String() == typeParam.Name() {
						isTypeParam = true
						return false
					}
					return true
				})
				if isTypeParam {
					typeArg := typeArgs[typeParamIndex]
					ast.Inspect(argStructDecl.Specs[0].(*ast.TypeSpec).Type.(*ast.StructType).Fields.List[paramIndex].Type, func(n ast.Node) bool {
						ident, ok := n.(*ast.Ident)
						if ident != nil && ok && ident.String() == typeParam.Name() {
							ident.Name = typeArg
							return false
						}
						return true
					})
				}
			}
		}
		funcBodyBlockStmt.List = append(funcBodyBlockStmt.List, &ast.DeclStmt{Decl: argStructDecl})
	}
	funcBodyBlockStmt.List = append(funcBodyBlockStmt.List,
		// test cases
		&ast.AssignStmt{
			Lhs: []ast.Expr{ast.NewIdent("tests")},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CompositeLit{
					Type: &ast.ArrayType{
						Elt: &ast.StructType{
							Fields: &ast.FieldList{
								List: func() []*ast.Field {
									list := make([]*ast.Field, 0, 2+returnTypeLen)
									nameField := field{names: []string{"name"}, typeName: "string"}
									list = append(list, nameField.make())
									if recv != nil {
										list = append(list, recv.make())
									}
									if len(params) > 0 {
										list = append(list, newField([]string{"args"}, "args", "", false).make())
									}
									for i, rt := range returnTypes {
										list = append(list, &ast.Field{
											Names: []*ast.Ident{ast.NewIdent(fmt.Sprintf("want%v", i))},
											Type: func() ast.Expr {
												// TODO: tmp, compare and search if field type is in type params, replace by index
												for typeParamIndex, typeParam := range typeParams {
													isTypeParam := false
													ast.Inspect(rt.typeNode(), func(n ast.Node) bool {
														ident, ok := n.(*ast.Ident)
														if ident != nil && ok && ident.String() == typeParam.Name() {
															isTypeParam = true
															return false
														}
														return true
													})
													if isTypeParam {
														typeArg := typeArgs[typeParamIndex]
														rtTypeNode := rt.typeNode()
														ast.Inspect(rtTypeNode, func(n ast.Node) bool {
															ident, ok := n.(*ast.Ident)
															if ident != nil && ok && ident.String() == typeParam.Name() {
																ident.Name = typeArg
																return false
															}
															return true
														})
														return rtTypeNode
													}
												}
												return rt.typeNode()
											}(),
										})
									}
									return list
								}(),
							},
						},
					},
				},
			},
		},
		// for range
		func() *ast.RangeStmt {
			keyIdent := ast.NewIdent("_")
			valueIdent := ast.NewIdent("tt")
			rangeIdent := ast.NewIdent("tests")
			assignStmt := &ast.AssignStmt{
				Lhs: []ast.Expr{
					keyIdent,
					valueIdent,
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.UnaryExpr{
						Op: token.RANGE,
						X:  rangeIdent,
					},
				},
			}
			keyIdent.Obj = &ast.Object{
				Kind: ast.Var,
				Name: "_",
				Decl: assignStmt,
			}
			valueIdent.Obj = &ast.Object{
				Kind: ast.Var,
				Name: "tt",
				Decl: assignStmt,
			}
			rangeStmt := &ast.RangeStmt{
				Key:   keyIdent,
				Value: valueIdent,
				Tok:   token.DEFINE,
				X:     rangeIdent,
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   ast.NewIdent("t"),
									Sel: ast.NewIdent("Run"),
								},
								Args: []ast.Expr{
									&ast.SelectorExpr{
										X:   ast.NewIdent("tt"),
										Sel: ast.NewIdent("name"),
									},
									&ast.FuncLit{
										Type: &ast.FuncType{
											Params: &ast.FieldList{
												List: []*ast.Field{
													func() *ast.Field {
														tField := field{
															names:    []string{"t"},
															typeName: "T",
															from:     "testing",
															pointer:  true,
														}
														return tField.make()
													}(),
												},
											},
										},
										Body: &ast.BlockStmt{
											List: func() []ast.Stmt {
												list := make([]ast.Stmt, 0, 1+returnTypeLen)

												// call
												callExpr := &ast.CallExpr{
													Fun: ast.NewIdent(funcName),
												}
												if recv != nil {
													callExpr.Fun = &ast.SelectorExpr{
														X: &ast.SelectorExpr{
															X:   ast.NewIdent("tt"),
															Sel: ast.NewIdent(recv.Name()),
														},
														Sel: ast.NewIdent(funcName),
													}
												}

												// args
												if len(params) > 0 {
													args := make([]ast.Expr, 0, len(params))
													for _, param := range params {
														args = append(args, &ast.SelectorExpr{
															X: &ast.SelectorExpr{
																X:   ast.NewIdent("tt"),
																Sel: ast.NewIdent("args"),
															},
															Sel: ast.NewIdent(param.Name()),
														})
													}
													callExpr.Args = args
												}

												// returns
												if returnTypeLen > 0 {
													list = append(list, &ast.AssignStmt{
														Lhs: func() []ast.Expr {
															lhs := make([]ast.Expr, 0, returnTypeLen)
															for i := range returnTypes {
																lhs = append(lhs, ast.NewIdent(fmt.Sprintf("got%v", i)))
															}
															return lhs
														}(),
														Tok: token.DEFINE,
														Rhs: []ast.Expr{
															callExpr,
														},
													})

													// compare
													for i := range returnTypes {
														got := fmt.Sprintf("got%v", i)
														want := fmt.Sprintf("want%v", i)
														list = append(list, &ast.IfStmt{
															Cond: &ast.UnaryExpr{
																Op: token.NOT,
																X: &ast.CallExpr{
																	Fun: &ast.SelectorExpr{
																		X:   ast.NewIdent("reflect"),
																		Sel: ast.NewIdent("DeepEqual"),
																	},
																	Args: []ast.Expr{
																		ast.NewIdent(got),
																		&ast.SelectorExpr{
																			X:   ast.NewIdent("tt"),
																			Sel: ast.NewIdent(want),
																		},
																	},
																},
															},
															Body: &ast.BlockStmt{
																List: []ast.Stmt{
																	&ast.ExprStmt{
																		X: &ast.CallExpr{
																			Fun: &ast.SelectorExpr{
																				X:   ast.NewIdent("t"),
																				Sel: ast.NewIdent("Errorf"),
																			},
																			Args: []ast.Expr{
																				&ast.BasicLit{
																					Kind:  token.STRING,
																					Value: fmt.Sprintf("\"%v() %v = %%v, %v %%v\"", funcName, got, want),
																				},
																				ast.NewIdent(got),
																				&ast.SelectorExpr{
																					X:   ast.NewIdent("tt"),
																					Sel: ast.NewIdent(want),
																				},
																			},
																		},
																	},
																},
															},
														})
													}
												} else {
													list = append(list, &ast.ExprStmt{
														X: callExpr,
													})
												}
												return list
											}(),
										},
									},
								},
							},
						},
					},
				},
			}
			return rangeStmt
		}(),
	)

	funcDecl := makeFuncDecl(unitTestFuncName, nil, nil, []*field{newField([]string{"t"}, "T", "testing", true)}, nil)
	funcDecl.Body = funcBodyBlockStmt

	buffer := &bytes.Buffer{}
	err := format.Node(buffer, token.NewFileSet(), funcDecl)
	if err != nil {
		panic(err)
	}
	return buffer.Bytes()
}

func (gfm *GoFunctionMeta) UnitTestFuncName() string {
	return fmt.Sprintf("Test_%v", gfm.FunctionName())
}

func (gfm *GoFunctionMeta) Recv() *GoVariableMeta {
	return nil
}

func (gmm *GoMethodMeta) UnitTestFuncName() string {
	recvStruct, _ := gmm.RecvStruct()
	return fmt.Sprintf("Test_%v_%v", recvStruct, gmm.FunctionName())
}

func MakeUnitTestFile(pkg string, importMetas []*GoImportMeta) []byte {
	// importDecl := makeImportDecl([]*importSpec{
	// 	newImportSpec("", "reflect"),
	// 	newImportSpec("", "reflect"),
	// })

	fileDecl := &ast.File{
		Name: ast.NewIdent(pkg),
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok: token.IMPORT,
				Specs: func() []ast.Spec {
					spec := make([]ast.Spec, 0, 2+len(importMetas))
					spec = append(spec,
						&ast.ImportSpec{
							Path: &ast.BasicLit{
								Kind:  token.STRING,
								Value: "\"reflect\"",
							},
						},
						&ast.ImportSpec{
							Path: &ast.BasicLit{
								Kind:  token.STRING,
								Value: "\"testing\"",
							},
						},
					)
					for _, importMeta := range importMetas {
						if importMeta.Alias() == "reflect" || importMeta.ImportPath() == "reflect" || importMeta.Alias() == "testing" || importMeta.ImportPath() == "testing" {
							continue
						}
						spec = append(spec, importMeta.node.(*ast.ImportSpec))
					}
					return spec
				}(),
			},
		},
	}
	buffer := &bytes.Buffer{}
	err := format.Node(buffer, token.NewFileSet(), fileDecl)
	if err != nil {
		panic(err)
	}
	return buffer.Bytes()
}

type GoBenchmarkMaker interface {
	BenchmarkFuncName() string
}

func MakeBenchmark(gbm GoBenchmarkMaker) []byte {
	// benchmarkFuncName := gbm.BenchmarkFuncName()

	// funcDecl := &ast.FuncDecl{}

	return nil
}

// ----

func MakeFile() {

}
