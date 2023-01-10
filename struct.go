package extractor

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

type GoStructMeta struct {
	typeSpec   *ast.TypeSpec
	methodDecl map[string]*GoMethodMeta
}

func ExtractGoStructMeta(extractFilepath string, structName string) (*GoStructMeta, error) {
	fileAST, err := parser.ParseFile(token.NewFileSet(), extractFilepath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	gsm := SearchGoStructMeta(fileAST, structName)
	if gsm.typeSpec == nil {
		return nil, fmt.Errorf("can not find struct decl")
	}

	return gsm, nil
}

func SearchGoStructMeta(fileAST *ast.File, structName string) *GoStructMeta {
	var structDecl *ast.TypeSpec
	ast.Inspect(fileAST, func(n ast.Node) bool {
		if n == fileAST {
			return true
		}
		if n == nil || structDecl != nil {
			return false
		}
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}
		if typeSpec.Type == nil {
			return false
		}
		_, ok = typeSpec.Type.(*ast.StructType)
		if !ok {
			return true
		}
		if typeSpec.Name.String() == structName {
			structDecl = typeSpec
			return false
		}
		return true
	})
	return &GoStructMeta{
		typeSpec: structDecl,
	}
}

func (gsm *GoStructMeta) PrintAST() {
	ast.Print(token.NewFileSet(), gsm.typeSpec)
}

func (gsm *GoStructMeta) StructName() string {
	return gsm.typeSpec.Name.String()
}
