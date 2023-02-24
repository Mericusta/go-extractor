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

type testStmtMaker func(string, *GoVariableMeta, []*GoVariableMeta, []*GoVariableMeta, int, int) ast.Stmt

type testMaker struct {
	prefix           string
	funcParam        *field
	preStmtMaker     testStmtMaker
	runningStmtMaker testStmtMaker
	postStmtMaker    testStmtMaker
}

var (
	unittestMaker = &testMaker{
		prefix:           "Test",
		funcParam:        newField([]string{"t"}, "T", "testing", true),
		runningStmtMaker: makeUnitTestRunningStmt,
	}
	benchmarkMaker = &testMaker{
		prefix:           "Benchmark",
		funcParam:        newField([]string{"b"}, "B", "testing", true),
		preStmtMaker:     makeBenchmarkPreStmt,
		runningStmtMaker: makeBenchmarkRunningStmt,
	}
)

// makeTest generate unit test func
// @param1  *GoFunctionMeta/*GoMethodMeta
// @param2  specify unittest func name
// @param3  specify type args for type params if needs
// @return1 unit test func name
// @return2 func declaration content
func makeTest(tm *testMaker, gutm GoUnitTestMaker, testFuncName string, typeArgs []string) (string, []byte) {
	if tm == nil {
		return "", nil
	}

	funcName := gutm.FunctionName()
	typeParams := gutm.TypeParams()
	if len(typeParams) > len(typeArgs) {
		return "", nil
	}
	if len(testFuncName) == 0 {
		testFuncName = fmt.Sprintf("%v_%v", tm.prefix, gutm.testFuncName(typeArgs))
	}
	recv := gutm.Recv()
	params := gutm.Params()
	paramsLen := len(params)
	returnTypes := gutm.ReturnTypes()
	returnTypeLen := len(returnTypes)

	// func decl
	funcDecl := makeFuncDecl(testFuncName, nil, nil, []*field{tm.funcParam}, nil)
	funcDecl.Body = &ast.BlockStmt{}

	// arg struct stmt
	if len(params) > 0 {
		funcDecl.Body.List = append(funcDecl.Body.List, makeTestArgsStructStmt(params, typeParams, typeArgs))
	}
	// test cases stmt
	testCasesStmt := makeTestCasesAssignStmt(recv, params, typeParams, returnTypes, returnTypeLen, typeArgs)
	if testCasesStmt != nil {
		funcDecl.Body.List = append(funcDecl.Body.List, testCasesStmt)
	}

	if tm.preStmtMaker != nil {
		// pre stmt
		preStmt := tm.preStmtMaker(funcName, recv, params, returnTypes, paramsLen, returnTypeLen)
		if preStmt != nil {
			funcDecl.Body.List = append(funcDecl.Body.List, preStmt)
		}
	}

	// test cases for-range stmt
	forRangeStmt := makeTestCasesForRangeStmt(funcName, recv, params, returnTypes, paramsLen, returnTypeLen)

	// test case running stmt
	if tm.runningStmtMaker != nil {
		forRangeStmt.Body.List = append(forRangeStmt.Body.List, tm.runningStmtMaker(funcName, recv, params, returnTypes, paramsLen, returnTypeLen))
	}

	if forRangeStmt != nil {
		funcDecl.Body.List = append(funcDecl.Body.List, forRangeStmt)
	}

	if tm.postStmtMaker != nil {
		// post stmt
		postStmt := tm.postStmtMaker(funcName, recv, params, returnTypes, paramsLen, returnTypeLen)
		if postStmt != nil {
			funcDecl.Body.List = append(funcDecl.Body.List, postStmt)
		}
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
func makeTestArgsStructStmt(params, typeParams []*GoVariableMeta, typeArgs []string) ast.Stmt {
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
	return &ast.DeclStmt{Decl: argStructDecl}
}

// test cases stmt
func makeTestCasesAssignStmt(recv *GoVariableMeta, params, typeParams, returnTypes []*GoVariableMeta, returnTypeLen int, typeArgs []string) ast.Stmt {
	return &ast.AssignStmt{
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
	}
}

// benchmark pre stmt
func makeBenchmarkPreStmt(funcName string, recv *GoVariableMeta, params, returnTypes []*GoVariableMeta, paramsLen, returnTypeLen int) ast.Stmt {
	return &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("b"),
				Sel: ast.NewIdent("ResetTimer"),
			},
		},
	}
}

// test cases for-range stmt
func makeTestCasesForRangeStmt(funcName string, recv *GoVariableMeta, params, returnTypes []*GoVariableMeta, paramsLen, returnTypeLen int) *ast.RangeStmt {
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

// unittest test case running stmt
func makeUnitTestRunningStmt(funcName string, recv *GoVariableMeta, params, returnTypes []*GoVariableMeta, paramsLen, returnTypeLen int) ast.Stmt {
	return &ast.ExprStmt{
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
							if paramsLen > 0 {
								args := make([]ast.Expr, 0, paramsLen)
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
	}
}

// benchmark test case running stmt
func makeBenchmarkRunningStmt(funcName string, recv *GoVariableMeta, params, returnTypes []*GoVariableMeta, paramsLen, returnTypeLen int) ast.Stmt {
	iteratorIdent := ast.NewIdent("i")
	forStmt := &ast.ForStmt{
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
				list := make([]ast.Stmt, 0, 1+returnTypeLen)
				placeHolderIdent := ast.NewIdent("_")

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
				if paramsLen > 0 {
					args := make([]ast.Expr, 0, paramsLen)
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
							for i := 0; i < returnTypeLen; i++ {
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
	}
	return forStmt
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

type GoBenchmarkMaker interface {
	GoUnitTestMaker
	BenchmarkFuncName([]string) string
}

// func makeBenchmark(gbm GoBenchmarkMaker, benchmarkFuncName string, typeArgs []string) (string, []byte) {
// 	funcName := gbm.FunctionName()
// 	typeParams := gbm.TypeParams()
// 	if len(typeParams) > len(typeArgs) {
// 		return "", nil
// 	}
// 	if len(benchmarkFuncName) == 0 {
// 		benchmarkFuncName = gbm.BenchmarkFuncName(typeArgs)
// 	}
// 	recv := gbm.Recv()
// 	params := gbm.Params()
// 	returnTypes := gbm.ReturnTypes()
// 	returnTypeLen := len(returnTypes)

// 	funcBodyBlockStmt := &ast.BlockStmt{}

// 	return "", nil
// }

// func (gfm *GoFunctionMeta) BenchmarkFuncName(typeArgs []string) string {
// 	return gfm.wrapTypeArgs(gfm.wrapPrefix("Benchmark_", fmt.Sprintf("%v", gfm.FunctionName())), typeArgs)
// }

// func (gmm *GoMethodMeta) BenchmarkFuncName(typeArgs []string) string {
// 	recvStruct, _ := gmm.RecvStruct()
// 	return gmm.wrapTypeArgs(gmm.wrapPrefix("Benchmark_", fmt.Sprintf("%v_%v", recvStruct, gmm.FunctionName())), typeArgs)
// }

// ----

func MakeFile() {

}
