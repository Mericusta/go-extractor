package extractor

import (
	"fmt"
	"go/ast"
)

type GoMethodMeta struct {
	*GoFunctionMeta
}

func ExtractGoMethodMeta(extractFilepath, structName, methodName string) (*GoMethodMeta, error) {
	gfm, err := ExtractGoFileMeta(extractFilepath)
	if err != nil {
		return nil, err
	}

	gmm := SearchGoMethodMeta(gfm, structName, methodName)
	if gmm == nil {
		return nil, fmt.Errorf("can not find method node")
	}

	return gmm, nil
}

func SearchGoMethodMeta(gfm *GoFileMeta, structName, methodName string) *GoMethodMeta {
	var methodDecl *ast.FuncDecl
	ast.Inspect(gfm.node, func(n ast.Node) bool {
		if IsMethodNode(n) {
			decl := n.(*ast.FuncDecl)
			if decl.Name.String() == methodName {
				recvStructName, _ := extractMethodRecvStruct(decl)
				if recvStructName == structName {
					methodDecl = decl
				}
			}
		}
		return methodDecl == nil
	})
	if methodDecl == nil {
		return nil
	}
	return &GoMethodMeta{
		GoFunctionMeta: NewGoFunctionMeta(gfm.newMeta(methodDecl)),
	}
}

func IsMethodNode(n ast.Node) bool {
	decl, ok := n.(*ast.FuncDecl)
	return ok && decl.Recv != nil && len(decl.Recv.List) > 0
}

func (gmm *GoMethodMeta) TypeParams() []*GoVariableMeta {
	// check if method receiver or method function declaration has type parameters
	// - method receiver means struct has type parameters
	// - function declaration means function has type parameters (go 1.20 not supported yet)

	funcDecl := gmm.node.(*ast.FuncDecl)
	if funcDecl.Type == nil {
		return nil
	}

	tParams := make([]*GoVariableMeta, 0)

	// method receiver
	if funcDecl.Recv != nil && len(funcDecl.Recv.List) != 0 && funcDecl.Recv.List[0].Type != nil {
		var typeParamExpr ast.Expr
		ast.Inspect(funcDecl.Recv.List[0].Type, func(n ast.Node) bool {
			switch _n := n.(type) {
			case *ast.IndexExpr, *ast.IndexListExpr:
				typeParamExpr = _n.(ast.Expr)
			}
			return typeParamExpr == nil
		})

		switch _tpe := typeParamExpr.(type) {
		case *ast.IndexExpr: // 因为这里直接使用了 struct 的
			tParams = append(tParams, &GoVariableMeta{
				meta:     gmm.newMeta(_tpe.Index),
				name:     _tpe.Index.(*ast.Ident).String(),
				typeMeta: gmm.newMeta(_tpe.Index),
			})
		case *ast.IndexListExpr:
			for _, _i := range _tpe.Indices {
				tParams = append(tParams, &GoVariableMeta{
					meta:     gmm.newMeta(_i),
					name:     _i.(*ast.Ident).String(),
					typeMeta: gmm.newMeta(_i),
				})
			}
		}
	}

	// receiverType, ok := funcDecl.Recv.List[0].Type.(ast.Expr)
	// function declaration
	if funcDecl.Type.TypeParams != nil && len(funcDecl.Type.TypeParams.List) != 0 {
		for _, field := range gmm.node.(*ast.FuncDecl).Type.TypeParams.List {
			for _, name := range field.Names {
				tParams = append(tParams, &GoVariableMeta{
					meta:     gmm.newMeta(field),
					name:     name.String(),
					typeMeta: gmm.newMeta(field.Type),
				})
			}
		}
	}

	return tParams
}

func (gmm *GoMethodMeta) RecvStruct() (string, bool) {
	return extractMethodRecvStruct(gmm.node.(*ast.FuncDecl))
}

func (gmm *GoMethodMeta) Recv() *GoVariableMeta {
	recv := gmm.node.(*ast.FuncDecl).Recv.List[0]
	return &GoVariableMeta{
		meta:     gmm.newMeta(recv),
		name:     recv.Names[0].String(),
		typeMeta: gmm.newMeta(recv.Type),
	}
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

// func trait

func (gmm *GoMethodMeta) MakeUnitTest(typeArgs []string) (string, []byte) {
	return makeTest(unittestMaker, gmm, "", typeArgs)
}

func (gmm *GoMethodMeta) UnittestFuncName(typeArgs []string) string {
	return wrapTestType(UNITTEST, gmm.testFuncName(typeArgs))
}

func (gmm *GoMethodMeta) MakeBenchmark(typeArgs []string) (string, []byte) {
	return makeTest(benchmarkMaker, gmm, "", typeArgs)
}

func (gmm *GoMethodMeta) BenchmarkFuncName(typeArgs []string) string {
	return wrapTestType(BENCHMARK, gmm.testFuncName(typeArgs))
}
