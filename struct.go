package extractor

import (
	"fmt"
	"go/ast"
)

type GoStructMeta struct {
	*meta        // *ast.TypeSpec
	commentGroup *ast.CommentGroup
	memberDecl   map[string]*GoVariableMeta
	methodDecl   map[string]*GoMethodMeta
}

func ExtractGoStructMeta(extractFilepath string, structName string) (*GoStructMeta, error) {
	gfm, err := ExtractGoFileMeta(extractFilepath)
	if err != nil {
		return nil, err
	}

	gsm := SearchGoStructMeta(gfm, structName)
	if gsm == nil {
		return nil, fmt.Errorf("can not find struct node")
	}

	return gsm, nil
}

func SearchGoStructMeta(gfm *GoFileMeta, structName string) *GoStructMeta {
	var structDecl *ast.TypeSpec
	var commentDecl *ast.CommentGroup
	ast.Inspect(gfm.node, func(n ast.Node) bool {
		if genDecl, ok := n.(*ast.GenDecl); ok {
			ast.Inspect(genDecl, func(n ast.Node) bool {
				if IsStructNode(n) {
					typeSpec := n.(*ast.TypeSpec)
					if typeSpec.Name.String() == structName {
						structDecl = typeSpec
						commentDecl = genDecl.Doc
						return false
					}
				}
				return true
			})
			return false // genDecl traverse done
		}
		return structDecl == nil // already found
	})
	if structDecl == nil {
		return nil
	}
	return &GoStructMeta{
		meta:         gfm.newMeta(structDecl),
		commentGroup: commentDecl,
		methodDecl:   make(map[string]*GoMethodMeta),
		memberDecl:   make(map[string]*GoVariableMeta),
	}
}

func IsStructNode(n ast.Node) bool {
	typeSpec, ok := n.(*ast.TypeSpec)
	if !ok {
		return false
	}
	if typeSpec.Type == nil {
		return false
	}
	_, ok = typeSpec.Type.(*ast.StructType)
	return ok
}

func (gsm *GoStructMeta) StructName() string {
	return gsm.node.(*ast.TypeSpec).Name.String()
}

func (gsm *GoStructMeta) Doc() []string {
	if gsm.node == nil || gsm.commentGroup == nil || len(gsm.commentGroup.List) == 0 {
		return nil
	}
	commentSlice := make([]string, 0, len(gsm.commentGroup.List))
	for _, comment := range gsm.commentGroup.List {
		commentSlice = append(commentSlice, comment.Text)
	}
	return commentSlice
}

func (gsm *GoStructMeta) Members() []string {
	if gsm.node.(*ast.TypeSpec) == nil || gsm.node.(*ast.TypeSpec).Type == nil {
		return nil
	}
	structType, ok := gsm.node.(*ast.TypeSpec).Type.(*ast.StructType)
	if structType == nil || !ok || structType.Fields == nil || len(structType.Fields.List) == 0 {
		return nil
	}
	members := make([]string, 0, len(structType.Fields.List))
	for _, field := range structType.Fields.List {
		if memberName := searchMemberName(field); len(memberName) != 0 {
			members = append(members, memberName)
		}
	}
	return members
}

func (gsm *GoStructMeta) SearchMemberMeta(member string) *GoVariableMeta {
	if gmm, has := gsm.memberDecl[member]; gmm != nil && has {
		return gmm
	}

	structType := gsm.node.(*ast.TypeSpec).Type.(*ast.StructType)
	gmm := SearchGoMemberMeta(gsm.meta, structType, member)
	if gmm == nil {
		return nil
	}
	gsm.memberDecl[member] = gmm

	return gsm.memberDecl[member]
}
