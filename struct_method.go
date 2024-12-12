package extractor

import (
	"go/ast"
)

type GoMethodMetaTypeConstraints interface {
	*ast.FuncDecl

	ast.Node
}

type GoMethodMeta[T GoMethodMetaTypeConstraints] struct {
	*GoFuncMeta[T]

	// receiver 的 meta 数据
	receiver *GoVarMeta[*ast.Field]
}

func NewGoMethodMeta[T GoMethodMetaTypeConstraints](m *meta[T], ident string, stopExtract ...bool) *GoMethodMeta[T] {
	gmm := &GoMethodMeta[T]{GoFuncMeta: NewGoFunctionMeta(m, ident)}
	if len(stopExtract) == 0 {
		gmm.ExtractAll()
	}
	return gmm
}

// -------------------------------- extractor --------------------------------

func (gmm *GoMethodMeta[T]) ExtractAll() {
	// 提取 func
	gmm.GoFuncMeta.ExtractAll()

	// 提取 receiver
	gmm.extractReceiver()
}

func (gmm *GoMethodMeta[T]) extractReceiver() {
	var (
		funcDecl     *ast.FuncDecl = gmm.node
		receiverNode               = funcDecl.Recv.List[0]
		receiverName string
	)
	if len(receiverNode.Names) == 1 {
		receiverName = receiverNode.Names[0].String()
	}
	gmm.receiver = NewGoVarMeta(newMeta(receiverNode, gmm.path), receiverName)
}

func extractMethodRecvStruct(methodDecl *ast.FuncDecl) (string, bool) {
	if len(methodDecl.Recv.List) < 1 {
		return "", false
	}
	recvTypeIdentNode, pointerReceiver := traitReceiverStruct(methodDecl.Recv.List[0].Type)
	if recvTypeIdentNode == nil {
		return "", false
	}
	return recvTypeIdentNode.String(), pointerReceiver
}

func traitReceiverStruct(node ast.Node) (*ast.Ident, bool) {
	var pointerReceiver bool
	var identNode *ast.Ident
	ast.Inspect(node, func(n ast.Node) bool {
		switch _n := n.(type) {
		case *ast.StarExpr:
			pointerReceiver = true
			identNode, _ = traitReceiverStruct(_n.X)
		case *ast.IndexExpr:
			identNode, _ = traitReceiverStruct(_n.X)
		case *ast.IndexListExpr:
			identNode, _ = traitReceiverStruct(_n.X)
		case *ast.Ident:
			identNode = _n
		}
		return identNode == nil
	})
	return identNode, pointerReceiver
}

// -------------------------------- extractor --------------------------------

// -------------------------------- unit test --------------------------------

func (gmm *GoMethodMeta[T]) Receiver() *GoVarMeta[*ast.Field] { return gmm.receiver }

// -------------------------------- unit test --------------------------------

// func ExtractGoMethodMeta[T GoMethodMetaTypeConstraints](extractFilepath, structName, methodName string) (*GoMethodMeta[T], error) {
// 	gfm, err := ExtractGoFileMeta(extractFilepath)
// 	if err != nil {
// 		return nil, err
// 	}

// 	gmm := SearchGoMethodMeta(gfm, structName, methodName)
// 	if gmm == nil {
// 		return nil, fmt.Errorf("can not find method node")
// 	}

// 	return gmm, nil
// }

// func SearchGoMethodMeta(gfm *GoFileMeta, structName, methodName string) *GoMethodMeta {
// 	var methodDecl *ast.FuncDecl
// 	ast.Inspect(gfm.node, func(n ast.Node) bool {
// 		if IsMethodNode(n) {
// 			decl := n.(*ast.FuncDecl)
// 			if decl.Name.String() == methodName {
// 				recvStructName, _ := extractMethodRecvStruct(decl)
// 				if recvStructName == structName {
// 					methodDecl = decl
// 				}
// 			}
// 		}
// 		return methodDecl == nil
// 	})
// 	if methodDecl == nil {
// 		return nil
// 	}
// 	return &GoMethodMeta{
// 		GoFuncMeta: NewGoFunctionMeta(gfm.copyMeta(methodDecl), ""),
// 	}
// }

// func (gmm *GoMethodMeta) TypeParams() []*GoVarMeta {
// 	// check if method receiver or method function declaration has type parameters
// 	// - method receiver means struct has type parameters
// 	// - function declaration means function has type parameters (go 1.20 not supported yet)

// 	funcDecl := gmm.node.(*ast.FuncDecl)
// 	if funcDecl.Type == nil {
// 		return nil
// 	}

// 	tParams := make([]*GoVarMeta, 0)

// 	// method receiver
// 	if funcDecl.Recv != nil && len(funcDecl.Recv.List) != 0 && funcDecl.Recv.List[0].Type != nil {
// 		var typeParamExpr ast.Expr
// 		ast.Inspect(funcDecl.Recv.List[0].Type, func(n ast.Node) bool {
// 			switch _n := n.(type) {
// 			case *ast.IndexExpr, *ast.IndexListExpr:
// 				typeParamExpr = _n.(ast.Expr)
// 			}
// 			return typeParamExpr == nil
// 		})

// 		switch _tpe := typeParamExpr.(type) {
// 		case *ast.IndexExpr: // 因为这里直接使用了 struct 的
// 			tParams = append(tParams, &GoVarMeta{
// 				meta:     gmm.copyMeta(_tpe.Index),
// 				ident:    _tpe.Index.(*ast.Ident).String(),
// 				typeMeta: gmm.copyMeta(_tpe.Index),
// 			})
// 		case *ast.IndexListExpr:
// 			for _, _i := range _tpe.Indices {
// 				tParams = append(tParams, &GoVarMeta{
// 					meta:     gmm.copyMeta(_i),
// 					ident:    _i.(*ast.Ident).String(),
// 					typeMeta: gmm.copyMeta(_i),
// 				})
// 			}
// 		}
// 	}

// 	// receiverType, ok := funcDecl.Recv.List[0].Type.(ast.Expr)
// 	// function declaration
// 	if funcDecl.Type.TypeParams != nil && len(funcDecl.Type.TypeParams.List) != 0 {
// 		for _, field := range gmm.node.(*ast.FuncDecl).Type.TypeParams.List {
// 			for _, name := range field.Names {
// 				tParams = append(tParams, &GoVarMeta{
// 					meta:     gmm.copyMeta(field),
// 					ident:    name.String(),
// 					typeMeta: gmm.copyMeta(field.Type),
// 				})
// 			}
// 		}
// 	}

// 	return tParams
// }

// func (gmm *GoMethodMeta) RecvStruct() (string, bool) {
// 	return extractMethodRecvStruct(gmm.node.(*ast.FuncDecl))
// }

// func (gmm *GoMethodMeta) Recv() *GoVarMeta {
// 	recv := gmm.node.(*ast.FuncDecl).Recv.List[0]
// 	return &GoVarMeta{
// 		meta:     gmm.copyMeta(recv),
// 		ident:    recv.Names[0].String(),
// 		typeMeta: gmm.copyMeta(recv.Type),
// 	}
// }

// func extractMethodRecvStruct(methodDecl *ast.FuncDecl) (string, bool) {
// 	if len(methodDecl.Recv.List) < 1 {
// 		return "", false
// 	}
// 	recvTypeIdentNode, pointerReceiver := traitReceiverStruct(methodDecl.Recv.List[0].Type)
// 	if recvTypeIdentNode == nil {
// 		return "", false
// 	}
// 	return recvTypeIdentNode.String(), pointerReceiver
// }

// func traitReceiverStruct(node ast.Node) (*ast.Ident, bool) {
// 	var pointerReceiver bool
// 	var identNode *ast.Ident
// 	ast.Inspect(node, func(n ast.Node) bool {
// 		switch _n := n.(type) {
// 		case *ast.StarExpr:
// 			pointerReceiver = true
// 			identNode, _ = traitReceiverStruct(_n.X)
// 		case *ast.IndexExpr:
// 			identNode, _ = traitReceiverStruct(_n.X)
// 		case *ast.IndexListExpr:
// 			identNode, _ = traitReceiverStruct(_n.X)
// 		case *ast.Ident:
// 			identNode = _n
// 		}
// 		return identNode == nil
// 	})
// 	return identNode, pointerReceiver
// }

// // func trait

// func (gmm *GoMethodMeta) MakeUnitTest(typeArgs []string) (string, []byte) {
// 	return makeTest(unittestMaker, gmm, "", typeArgs)
// }

// func (gmm *GoMethodMeta) UnittestFuncName(typeArgs []string) string {
// 	return wrapTestType(UNITTEST, gmm.testFuncName(typeArgs))
// }

// func (gmm *GoMethodMeta) MakeBenchmark(typeArgs []string) (string, []byte) {
// 	return makeTest(benchmarkMaker, gmm, "", typeArgs)
// }

// func (gmm *GoMethodMeta) BenchmarkFuncName(typeArgs []string) string {
// 	return wrapTestType(BENCHMARK, gmm.testFuncName(typeArgs))
// }
