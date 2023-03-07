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

type GoTestMaker interface {
	testFuncName([]string) string
	FunctionName() string
	TypeParams() []*GoVariableMeta
	Params() []*GoVariableMeta
	ReturnTypes() []*GoVariableMeta
	Recv() *GoVariableMeta
}

type testType int

const (
	UNITTEST = iota + 1
	BENCHMARK
)

type testStmtMaker func(string, GoTestMaker, []*typeArgInfer) []ast.Stmt

type testMaker struct {
	prefix             string
	funcParam          *field
	testCasesStmtMaker testStmtMaker
	preStmtMaker       testStmtMaker
	runningStmtMaker   testStmtMaker
	postStmtMaker      testStmtMaker
}

var (
	unittestMaker = &testMaker{
		prefix:             "Test",
		funcParam:          newField([]string{"t"}, "T", "testing", true),
		testCasesStmtMaker: makeUnittestTestCasesAssignStmt,
		runningStmtMaker:   makeUnitTestTestCaseRunningStmt,
	}
	benchmarkMaker = &testMaker{
		prefix:             "Benchmark",
		funcParam:          newField([]string{"b"}, "B", "testing", true),
		testCasesStmtMaker: makeBenchmarkTestCasesAssignStmt,
		preStmtMaker:       makeBenchmarkTestCasePreStmt,
		runningStmtMaker:   makeBenchmarkTestCaseRunningStmt,
		postStmtMaker:      makeBenchmarkTestCasePostStmt,
	}
)

type typeArgInfer struct {
	typeArg  string
	canInfer bool
}

// makeTest generate unit test func
// @param1  *GoFunctionMeta/*GoMethodMeta
// @param2  specify unittest func name
// @param3  specify type args for type params if needs
// @return1 unit test func name
// @return2 func declaration content
func makeTest(tm *testMaker, gtm GoTestMaker, testFuncName string, typeArgs []string) (string, []byte) {
	if tm == nil {
		return "", nil
	}

	funcName := gtm.FunctionName()
	if len(gtm.TypeParams()) > len(typeArgs) {
		return "", nil
	}
	if len(testFuncName) == 0 {
		testFuncName = fmt.Sprintf("%v_%v", tm.prefix, gtm.testFuncName(typeArgs))
	}
	typeArgInfers := make([]*typeArgInfer, 0, len(typeArgs))
	for _, typeArg := range typeArgs {
		typeArgInfers = append(typeArgInfers, &typeArgInfer{
			typeArg:  typeArg,
			canInfer: false,
		})
	}

	// func decl
	funcDecl := makeFuncDecl(testFuncName, nil, nil, []*field{tm.funcParam}, nil)
	funcDecl.Body = &ast.BlockStmt{}

	// arg struct stmt
	if len(gtm.Params()) > 0 {
		testArgsStructStmt := makeTestArgsStructStmt(funcName, gtm, typeArgInfers)
		funcDecl.Body.List = append(funcDecl.Body.List, testArgsStructStmt)
	}

	// test cases stmt
	testCasesStmts := tm.testCasesStmtMaker(funcName, gtm, typeArgInfers)
	funcDecl.Body.List = append(funcDecl.Body.List, testCasesStmts...)

	// test cases for-range stmt
	forRangeStmt := makeTestCasesForRangeStmt(funcName, gtm, typeArgs)
	funcDecl.Body.List = append(funcDecl.Body.List, forRangeStmt)

	// test case pre stmt
	if tm.preStmtMaker != nil {
		preStmts := tm.preStmtMaker(funcName, gtm, typeArgInfers)
		forRangeStmt.Body.List = append(forRangeStmt.Body.List, preStmts...)
	}

	// test case running stmt
	if tm.runningStmtMaker != nil {
		runningStmts := tm.runningStmtMaker(funcName, gtm, typeArgInfers)
		forRangeStmt.Body.List = append(forRangeStmt.Body.List, runningStmts...)
	}

	// test case post stmt
	if tm.postStmtMaker != nil {
		postStmts := tm.postStmtMaker(funcName, gtm, typeArgInfers)
		forRangeStmt.Body.List = append(forRangeStmt.Body.List, postStmts...)
	}

	// output
	buffer := &bytes.Buffer{}
	err := format.Node(buffer, token.NewFileSet(), funcDecl)
	if err != nil {
		panic(err)
	}
	return testFuncName, buffer.Bytes()
}

// arg struct stmt
func makeTestArgsStructStmt(funcName string, gtm GoTestMaker, typeArgInfers []*typeArgInfer) ast.Stmt {
	argStructDecl := newTypeSpec("args", nil, gtm.Params()).makeDecl()
	for paramIndex, param := range gtm.Params() {
		for typeParamIndex, typeParam := range gtm.TypeParams() {
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
				typeArgInfer := typeArgInfers[typeParamIndex]
				ast.Inspect(argStructDecl.Specs[0].(*ast.TypeSpec).Type.(*ast.StructType).Fields.List[paramIndex].Type, func(n ast.Node) bool {
					ident, ok := n.(*ast.Ident)
					if ident != nil && ok && ident.String() == typeParam.Name() {
						ident.Name = typeArgInfer.typeArg // Note: 修改了原 node 的值
						typeArgInfer.canInfer = true
						return false
					}
					return true
				})
			}
		}
	}
	return &ast.DeclStmt{Decl: argStructDecl}
}

// unittest test cases stmt
func makeUnittestTestCasesAssignStmt(funcName string, gtm GoTestMaker, typeArgInfers []*typeArgInfer) []ast.Stmt {
	return []ast.Stmt{
		&ast.AssignStmt{
			Lhs: []ast.Expr{ast.NewIdent("tests")},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CompositeLit{
					Type: &ast.ArrayType{
						Elt: &ast.StructType{
							Fields: &ast.FieldList{
								List: func() []*ast.Field {
									list := make([]*ast.Field, 0, 2+len(gtm.ReturnTypes()))
									nameField := field{names: []string{"name"}, typeName: "string"}
									list = append(list, nameField.make())
									if gtm.Recv() != nil {
										list = append(list, gtm.Recv().make())
									}
									if len(gtm.Params()) > 0 {
										list = append(list, newField([]string{"args"}, "args", "", false).make())
									}
									for i, rt := range gtm.ReturnTypes() {
										list = append(list, &ast.Field{
											Names: []*ast.Ident{ast.NewIdent(fmt.Sprintf("want%v", i))},
											Type: func() ast.Expr {
												// TODO: tmp, compare and search if field type is in type params, replace by index
												for typeParamIndex, typeParam := range gtm.TypeParams() {
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
														typeArgInfer := typeArgInfers[typeParamIndex]
														rtTypeNode := rt.typeNode()
														ast.Inspect(rtTypeNode, func(n ast.Node) bool {
															ident, ok := n.(*ast.Ident)
															if ident != nil && ok && ident.String() == typeParam.Name() {
																ident.Name = typeArgInfer.typeArg
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
	}
}

// benchmark test cases stmt
func makeBenchmarkTestCasesAssignStmt(funcName string, gtm GoTestMaker, typeArgInfers []*typeArgInfer) []ast.Stmt {
	return []ast.Stmt{
		&ast.AssignStmt{
			Lhs: []ast.Expr{ast.NewIdent("tests")},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CompositeLit{
					Type: &ast.ArrayType{
						Elt: &ast.StructType{
							Fields: &ast.FieldList{
								List: func() []*ast.Field {
									list := make([]*ast.Field, 0, 2+len(gtm.ReturnTypes()))
									list = append(list, newField([]string{"name"}, "string", "", false).make())
									list = append(list, newField([]string{"limit"}, "Duration", "time", false).make())
									if gtm.Recv() != nil {
										list = append(list, gtm.Recv().make())
									}
									if len(gtm.Params()) > 0 {
										list = append(list, newField([]string{"args"}, "args", "", false).make())
									}
									return list
								}(),
							},
						},
					},
				},
			},
		},
	}
}

// test cases for-range stmt
func makeTestCasesForRangeStmt(funcName string, gtm GoTestMaker, typeArgs []string) *ast.RangeStmt {
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
		Body:  &ast.BlockStmt{},
	}
	return rangeStmt
}

// benchmark test case pre stmt
func makeBenchmarkTestCasePreStmt(funcName string, gtm GoTestMaker, typeArgInfers []*typeArgInfer) []ast.Stmt {
	return []ast.Stmt{
		&ast.ExprStmt{
			X: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("b"),
					Sel: ast.NewIdent("ResetTimer"),
				},
			},
		},
	}
}

// unittest test case running stmt
func makeUnitTestTestCaseRunningStmt(funcName string, gtm GoTestMaker, typeArgInfers []*typeArgInfer) []ast.Stmt {
	return []ast.Stmt{
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
								list := make([]ast.Stmt, 0, 1+len(gtm.ReturnTypes()))

								// call
								callExpr := &ast.CallExpr{
									Fun: ast.NewIdent(funcName),
								}
								if gtm.Recv() != nil {
									callExpr.Fun = &ast.SelectorExpr{
										X: &ast.SelectorExpr{
											X:   ast.NewIdent("tt"),
											Sel: ast.NewIdent(gtm.Recv().Name()),
										},
										Sel: ast.NewIdent(funcName),
									}
								} else {
									// TODO: support method
									for _, typeArgInfer := range typeArgInfers {
										if !typeArgInfer.canInfer {
											inferTypeArgsExplicitly(callExpr, typeArgInfers)
											break
										}
									}
								}

								// args
								if l := len(gtm.Params()); l > 0 {
									args := make([]ast.Expr, 0, l)
									for _, param := range gtm.Params() {
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
								if l := len(gtm.ReturnTypes()); l > 0 {
									list = append(list, &ast.AssignStmt{
										Lhs: func() []ast.Expr {
											lhs := make([]ast.Expr, 0, l)
											for i := range gtm.ReturnTypes() {
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
									for i := range gtm.ReturnTypes() {
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
	}
}

// benchmark test case running stmt
func makeBenchmarkTestCaseRunningStmt(funcName string, gtm GoTestMaker, typeArgInfers []*typeArgInfer) []ast.Stmt {
	iteratorIdent := ast.NewIdent("i")
	return []ast.Stmt{
		&ast.ForStmt{
			Init: &ast.AssignStmt{
				Lhs: []ast.Expr{iteratorIdent},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{&ast.BasicLit{Kind: token.INT, Value: "0"}},
			},
			Cond: &ast.BinaryExpr{
				X:  iteratorIdent,
				Op: token.LSS,
				Y: &ast.SelectorExpr{
					X:   ast.NewIdent("b"),
					Sel: ast.NewIdent("N"),
				},
			},
			Post: &ast.IncDecStmt{
				X:   iteratorIdent,
				Tok: token.INC,
			},
			Body: &ast.BlockStmt{
				List: func() []ast.Stmt {
					list := make([]ast.Stmt, 0, 1+len(gtm.ReturnTypes()))
					placeHolderIdent := ast.NewIdent("_")

					// call
					callExpr := &ast.CallExpr{
						Fun: ast.NewIdent(funcName),
					}
					if gtm.Recv() != nil {
						callExpr.Fun = &ast.SelectorExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("tt"),
								Sel: ast.NewIdent(gtm.Recv().Name()),
							},
							Sel: ast.NewIdent(funcName),
						}
					} else {
						// TODO: support method
						for _, typeArgInfer := range typeArgInfers {
							if !typeArgInfer.canInfer {
								inferTypeArgsExplicitly(callExpr, typeArgInfers)
								break
							}
						}
					}

					// args
					if len(gtm.Params()) > 0 {
						args := make([]ast.Expr, 0, len(gtm.Params()))
						for _, param := range gtm.Params() {
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
					if len(gtm.ReturnTypes()) > 0 {
						list = append(list, &ast.AssignStmt{
							Lhs: func() []ast.Expr {
								lhs := make([]ast.Expr, 0, len(gtm.ReturnTypes()))
								for i := 0; i < len(gtm.ReturnTypes()); i++ {
									lhs = append(lhs, placeHolderIdent)
								}
								return lhs
							}(),
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								callExpr,
							},
						})
					} else {
						list = append(list, &ast.ExprStmt{
							X: callExpr,
						})
					}
					return list
				}(),
			},
		},
	}
}

// benchmark test case post stmt
func makeBenchmarkTestCasePostStmt(funcName string, gtm GoTestMaker, typeArgInfers []*typeArgInfer) []ast.Stmt {
	return []ast.Stmt{
		&ast.ExprStmt{
			X: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("b"),
					Sel: ast.NewIdent("StopTimer"),
				},
			},
		},
		&ast.IfStmt{
			Cond: &ast.BinaryExpr{
				X: &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("b"),
						Sel: ast.NewIdent("Elapsed"),
					},
				},
				Op: token.GTR,
				Y: &ast.BinaryExpr{
					X: &ast.SelectorExpr{
						X:   ast.NewIdent("tt"),
						Sel: ast.NewIdent("limit"),
					},
					Op: token.MUL,
					Y: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("time"),
							Sel: ast.NewIdent("Duration"),
						},
						Args: []ast.Expr{
							&ast.SelectorExpr{
								X:   ast.NewIdent("b"),
								Sel: ast.NewIdent("N"),
							},
						},
					},
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ExprStmt{
						X: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("b"),
								Sel: ast.NewIdent("Fatalf"),
							},
							Args: []ast.Expr{
								&ast.BasicLit{
									Kind:  token.STRING,
									Value: "\"overtime limit %v, got %.2f\\n\"",
								},
								&ast.SelectorExpr{
									X:   ast.NewIdent("tt"),
									Sel: ast.NewIdent("limit"),
								},
								&ast.BinaryExpr{
									X: &ast.CallExpr{
										Fun: ast.NewIdent("float64"),
										Args: []ast.Expr{
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X:   ast.NewIdent("b"),
													Sel: ast.NewIdent("Elapsed"),
												},
											},
										},
									},
									Op: token.QUO,
									Y: &ast.CallExpr{
										Fun: ast.NewIdent("float64"),
										Args: []ast.Expr{
											&ast.SelectorExpr{
												X:   ast.NewIdent("b"),
												Sel: ast.NewIdent("N"),
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
}

func wrapTestType(t testType, s string) string {
	switch t {
	case UNITTEST:
		s = fmt.Sprintf("Test_%v", s)
	case BENCHMARK:
		s = fmt.Sprintf("Benchmark_%v", s)
	}
	return s
}

func wrapTypeArgs(unitTestFuncName string, typeArgs []string) string {
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
	return unitTestFuncName
}

func inferTypeArgsExplicitly(callExpr *ast.CallExpr, typeArgInfers []*typeArgInfer) {
	l := len(typeArgInfers)
	if l == 1 {
		callExpr.Fun = &ast.IndexExpr{
			X:     callExpr.Fun,
			Index: ast.NewIdent(typeArgInfers[0].typeArg),
		}
	} else {
		indexListExpr := &ast.IndexListExpr{
			X:       callExpr.Fun,
			Indices: make([]ast.Expr, 0, l),
		}
		for _, typeArgInfer := range typeArgInfers {
			indexListExpr.Indices = append(indexListExpr.Indices, ast.NewIdent(typeArgInfer.typeArg))
		}
		callExpr.Fun = indexListExpr
	}
}

func (gfm *GoFunctionMeta) testFuncName(typeArgs []string) string {
	return wrapTypeArgs(fmt.Sprintf("%v", gfm.FunctionName()), typeArgs)
}

func (gfm *GoFunctionMeta) Recv() *GoVariableMeta {
	return nil
}

func (gmm *GoMethodMeta) testFuncName(typeArgs []string) string {
	recvStruct, _ := gmm.RecvStruct()
	return wrapTypeArgs(fmt.Sprintf("%v_%v", recvStruct, gmm.FunctionName()), typeArgs)
}

func MakeTestFile(pkg string, importMetas []*GoImportMeta) []byte {
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
