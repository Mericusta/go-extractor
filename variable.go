package extractor

import (
	"fmt"
	"go/ast"
)

// 包别名
// 全局变量
// 局部变量
// 参数表
// 返回值
// 常量
type GoVariableMeta struct {
	*meta
	name     string
	typeMeta Meta
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
	var typeMeta Meta
	ast.Inspect(m.node, func(n ast.Node) bool {
		if IsVarNode(n) {
			decl := n.(*ast.ValueSpec)
			if decl.Names[0].String() == variableName {
				variableDecl = decl
				// TODO: typeMeta
			}
		}
		return variableDecl == nil
	})
	if variableDecl == nil {
		return nil
	}
	return &GoVariableMeta{
		meta:     m.newMeta(variableDecl),
		name:     variableName,
		typeMeta: typeMeta,
	}
}

// TODO:
func IsVarNode(n ast.Node) bool {
	valueSpec, ok := n.(*ast.ValueSpec)
	if !ok {
		return false
	}
	return valueSpec.Names[0].Obj.Kind == ast.Var
}

func (gvm *GoVariableMeta) Name() string {
	return gvm.name
}

func (gvm *GoVariableMeta) Type() string {
	return gvm.typeMeta.Expression()
}
