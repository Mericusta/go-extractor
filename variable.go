package extractor

import (
	"fmt"
	"go/ast"
)

type GoVariableMeta struct {
	*meta
}

func ExtractGoVariableMeta(extractFilepath string, variableName string) (*GoVariableMeta, error) {
	goFileMeta, err := ExtractGoFileMeta(extractFilepath)
	if err != nil {
		return nil, err
	}

	gam := SearchGoVariableMeta(goFileMeta.meta, variableName)
	if gam.node == nil {
		return nil, fmt.Errorf("can not find function node")
	}

	return gam, nil
}

func SearchGoVariableMeta(m *meta, variableName string) *GoVariableMeta {
	var variableDecl *ast.ValueSpec
	ast.Inspect(m.node, func(n ast.Node) bool {
		if IsVarNode(n) {
			decl := n.(*ast.ValueSpec)
			if decl.Names[0].String() == variableName {
				variableDecl = decl
			}
		}
		return variableDecl == nil
	})
	if variableDecl == nil {
		return nil
	}
	return &GoVariableMeta{
		meta: m.newMeta(variableDecl),
	}
}

func IsVarNode(n ast.Node) bool {
	valueSpec, ok := n.(*ast.ValueSpec)
	if !ok {
		return false
	}
	return valueSpec.Names[0].Obj.Kind == ast.Var
}
