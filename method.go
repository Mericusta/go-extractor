package extractor

import (
	"fmt"
	"go/ast"
)

type GoMethodMeta struct {
	*GoFunctionMeta
}

func ExtractGoMethodMeta(extractFilepath string, structName, methodName string) (*GoMethodMeta, error) {
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

	var pointerReceiver bool
	var recvTypeIdentNode ast.Node
	switch methodDecl.Recv.List[0].Type.(type) {
	case *ast.Ident:
		pointerReceiver = false
		recvTypeIdentNode = methodDecl.Recv.List[0].Type
	case *ast.StarExpr:
		pointerReceiver = true
		recvTypeIdentNode = methodDecl.Recv.List[0].Type.(*ast.StarExpr).X
	}

	recvTypeIdent, ok := recvTypeIdentNode.(*ast.Ident)
	if !ok {
		return "", false
	}
	return recvTypeIdent.String(), pointerReceiver
}
