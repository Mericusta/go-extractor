package extractor

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

type GoStructMeta struct {
	typeSpec     *ast.TypeSpec
	commentGroup *ast.CommentGroup
	memberDecl   map[string]*GoMemberMeta
	methodDecl   map[string]*GoMethodMeta
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
		typeSpec:     structDecl,
		commentGroup: commentDecl,
		methodDecl:   make(map[string]*GoMethodMeta),
		memberDecl:   make(map[string]*GoMemberMeta),
	}
}

func (gsm *GoStructMeta) PrintAST() {
	ast.Print(token.NewFileSet(), gsm.typeSpec)
}

func (gsm *GoStructMeta) StructName() string {
	return gsm.typeSpec.Name.String()
}

func (gsm *GoStructMeta) Comments() []string {
	if gsm.typeSpec == nil || gsm.commentGroup == nil || len(gsm.commentGroup.List) == 0 {
		return nil
	}
	commentSlice := make([]string, 0, len(gsm.commentGroup.List))
	for _, comment := range gsm.commentGroup.List {
		commentSlice = append(commentSlice, comment.Text)
	}
	return commentSlice
}

func (gsm *GoStructMeta) Members() []string {
	if gsm.typeSpec == nil || gsm.typeSpec.Type == nil {
		return nil
	}
	structType, ok := gsm.typeSpec.Type.(*ast.StructType)
	if structType == nil || !ok || structType.Fields == nil || len(structType.Fields.List) == 0 {
		return nil
	}
	members := make([]string, 0, len(structType.Fields.List))
	for _, field := range structType.Fields.List {
		members = append(members, field.Names[0].Name)
	}
	return members
}

func (gsm *GoStructMeta) SearchMemberMeta(member string) *GoMemberMeta {
	if gmm, has := gsm.memberDecl[member]; gmm != nil && has {
		return gmm
	}

	structType := gsm.typeSpec.Type.(*ast.StructType)
	gmm := SearchGoMemberMeta(structType, member)
	if gmm == nil {
		return nil
	}
	gsm.memberDecl[member] = gmm

	return gsm.memberDecl[member]
}
