package extractor

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

type GoStructMeta struct {
	typeSpec    *ast.TypeSpec
	commentDecl *ast.CommentGroup
	methodDecl  map[string]*GoMethodMeta
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
	var commentDecl *ast.CommentGroup
	ast.Inspect(fileAST, func(n ast.Node) bool {
		if n == fileAST {
			return true
		}
		if n == nil || structDecl != nil {
			return false
		}
		genDecl, ok := n.(*ast.GenDecl)
		if !ok {
			return false
		}
		ast.Inspect(genDecl, func(n ast.Node) bool {
			if n == genDecl {
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
				commentDecl = genDecl.Doc
				return false
			}
			return true
		})
		return false
	})
	return &GoStructMeta{
		typeSpec:    structDecl,
		commentDecl: commentDecl,
		methodDecl:  make(map[string]*GoMethodMeta),
	}
}

func (gsm *GoStructMeta) PrintAST() {
	ast.Print(token.NewFileSet(), gsm.typeSpec)
}

func (gsm *GoStructMeta) StructName() string {
	return gsm.typeSpec.Name.String()
}

func (gsm *GoStructMeta) Comments() []string {
	if gsm.typeSpec == nil || gsm.commentDecl == nil || len(gsm.commentDecl.List) == 0 {
		return nil
	}
	commentSlice := make([]string, 0, len(gsm.commentDecl.List))
	for _, comment := range gsm.commentDecl.List {
		commentSlice = append(commentSlice, comment.Text)
	}
	return commentSlice
}
