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

type field struct {
	name     string
	typeName string
	from     string
	pointer  bool
}

func (f *field) make() *ast.Field {
	var typeExpr ast.Expr = ast.NewIdent(f.typeName)
	if len(f.from) != 0 {
		typeExpr = &ast.SelectorExpr{
			X:   ast.NewIdent(f.from),
			Sel: ast.NewIdent(f.typeName),
		}
	}
	if f.pointer {
		typeExpr = &ast.StarExpr{
			X: typeExpr,
		}
	}
	return &ast.Field{
		Names: []*ast.Ident{ast.NewIdent(f.name)},
		Type:  typeExpr,
	}
}

func (gvm *GoVariableMeta) makeField() *ast.Field {
	var typeExpr ast.Expr = gvm.node.(*ast.Field).Type
	return &ast.Field{
		Names: []*ast.Ident{ast.NewIdent(gvm.Name())},
		Type:  typeExpr,
	}
}

func makeFieldList(fieldList []*GoVariableMeta) []*ast.Field {
	field := make([]*ast.Field, len(fieldList))
	for i, param := range fieldList {
		field[i] = param.makeField()
	}
	return field
}

type GoUnitTestMaker interface {
	UnitTestFuncName() string
	FunctionName() string
	TypeParams() []*GoVariableMeta
	Params() []*GoVariableMeta
	ReturnTypes() []*GoVariableMeta
	recvField() *ast.Field // TODO: remove
}

func MakeUnitTest(gutm GoUnitTestMaker) []byte {
	unitTestFuncName := gutm.UnitTestFuncName()
	funcName := gutm.FunctionName()
	typeParams := gutm.TypeParams()
	params := gutm.Params()
	returnTypes := gutm.ReturnTypes()
	returnTypeLen := len(returnTypes)
	recvField := gutm.recvField()
	funcDecl := &ast.FuncDecl{
		Name: ast.NewIdent(unitTestFuncName),
		Type: &ast.FuncType{
			TypeParams: func() *ast.FieldList {
				if len(typeParams) == 0 {
					return nil
				}
				return &ast.FieldList{
					List: makeFieldList(typeParams),
				}
			}(),
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("t")},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("testing"),
								Sel: ast.NewIdent("T"),
							},
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				// arg struct decl
				&ast.DeclStmt{
					Decl: &ast.GenDecl{
						Tok: token.TYPE,
						Specs: []ast.Spec{
							&ast.TypeSpec{
								Name: ast.NewIdent("args"),
								Type: &ast.StructType{
									Fields: &ast.FieldList{
										List: makeFieldList(params),
									},
								},
							},
						},
					},
				},
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
											nameField := field{name: "name", typeName: "string"}
											list = append(list, nameField.make())
											if recvField != nil {
												list = append(list, recvField)
											}
											argsField := field{name: "args", typeName: "args"}
											list = append(list, argsField.make())
											for i, rt := range returnTypes {
												list = append(list, &ast.Field{
													Names: []*ast.Ident{ast.NewIdent(fmt.Sprintf("want%v", i))},
													Type:  rt.typeNode(),
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
																	name:     "t",
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
														if recvField != nil {
															callExpr.Fun = &ast.SelectorExpr{
																X: &ast.SelectorExpr{
																	X:   ast.NewIdent("tt"),
																	Sel: recvField.Names[0],
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
			},
		},
	}

	buffer := &bytes.Buffer{}
	err := format.Node(buffer, token.NewFileSet(), funcDecl)
	if err != nil {
		panic(err)
	}
	return buffer.Bytes()
}

func MakeUnitTestWithTypeArgs(gutm GoUnitTestMaker, typeArgs []string) []byte {
	unitTestFuncName := gutm.UnitTestFuncName()
	funcName := gutm.FunctionName()
	typeParams := gutm.TypeParams()
	if len(typeParams) > len(typeArgs) {
		return nil
	}
	suffix := unitTestFuncName
	for _, typeArg := range typeArgs {
		suffix = fmt.Sprintf("%v_%v", suffix, typeArg)
	}
	hashStr := md5.Sum([]byte(suffix))
	unitTestFuncName = fmt.Sprintf("%v_%v", unitTestFuncName, hex.EncodeToString(hashStr[:]))
	params := gutm.Params()
	returnTypes := gutm.ReturnTypes()
	returnTypeLen := len(returnTypes)
	recvField := gutm.recvField()
	funcDecl := &ast.FuncDecl{
		Name: ast.NewIdent(unitTestFuncName),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("t")},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("testing"),
								Sel: ast.NewIdent("T"),
							},
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				// arg struct decl
				&ast.DeclStmt{
					Decl: &ast.GenDecl{
						Tok: token.TYPE,
						Specs: []ast.Spec{
							&ast.TypeSpec{
								Name: ast.NewIdent("args"),
								Type: &ast.StructType{
									Fields: &ast.FieldList{
										List: func() []*ast.Field {
											fieldList := makeFieldList(params)
											// TODO: tmp, compare and search if field type is in type params, replace by index
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
														ast.Inspect(fieldList[paramIndex].Type, func(n ast.Node) bool {
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
											return fieldList
										}(),
									},
								},
							},
						},
					},
				},
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
											nameField := field{name: "name", typeName: "string"}
											list = append(list, nameField.make())
											if recvField != nil {
												list = append(list, recvField)
											}
											argsField := field{name: "args", typeName: "args"}
											list = append(list, argsField.make())

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
																	name:     "t",
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
														if recvField != nil {
															callExpr.Fun = &ast.SelectorExpr{
																X: &ast.SelectorExpr{
																	X:   ast.NewIdent("tt"),
																	Sel: recvField.Names[0],
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
			},
		},
	}

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

func (gfm *GoFunctionMeta) recvField() *ast.Field {
	return nil
}

func (gmm *GoMethodMeta) UnitTestFuncName() string {
	recvStruct, _ := gmm.RecvStruct()
	return fmt.Sprintf("Test_%v_%v", recvStruct, gmm.FunctionName())
}

func MakeUnitTestFile(pkg string, importMetas []*GoImportMeta) []byte {
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
						if importMeta.alias == "reflect" || importMeta.name == "reflect" || importMeta.alias == "testing" || importMeta.name == "testing" {
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

// func MakeUpFuncDecl(funcName string, params []*MakeField, returns []*MakeField) {
// 	funcDecl := &ast.FuncDecl{
// 		Name: ast.NewIdent(funcName),
// 		Type: &ast.FuncType{Params: &ast.FieldList{
// 			List: makeFieldList(params),
// 		}},
// 	}
// }

// func MakeUpStructDecl(structName string, members []*MakeField) {
// 	structDecl := &ast.GenDecl{
// 		Tok: token.TYPE,
// 		Specs: []ast.Spec{
// 			&ast.TypeSpec{
// 				Name: ast.NewIdent(structName),
// 				Type: &ast.StructType{
// 					Fields: &ast.FieldList{
// 						List: makeFieldList(members),
// 					},
// 				},
// 			},
// 		},
// 	}
// }

// func MakeUpAssign(left []string, right []string) {
// 	assignStmt := &ast.AssignStmt{
// 		Lhs: func() []ast.Expr {
// 			expr := make([]ast.Expr, len(left))
// 			for i, e := range left {
// 				expr[i] = ast.NewIdent(e)
// 			}
// 			return expr
// 		}(),
// 		Tok: token.DEFINE,
// 		Rhs: func() []ast.Expr {

// 			return nil
// 		}(),
// 	}
// }
