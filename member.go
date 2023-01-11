package extractor

import "go/ast"

type GoMemberMeta struct {
	field *ast.Field
}

func SearchGoMemberMeta(structType *ast.StructType, memberName string) *GoMemberMeta {
	if structType.Fields == nil || len(structType.Fields.List) == 0 {
		return nil
	}
	var memberDecl *ast.Field
	for _, field := range structType.Fields.List {
		if field.Names[0].Name == memberName {
			memberDecl = field
			break
		}
	}
	return &GoMemberMeta{
		field: memberDecl,
	}
}

func (gmm *GoMemberMeta) MemberName() string {
	return gmm.field.Names[0].Name
}

func (gmm *GoMemberMeta) Tag() string {
	return gmm.field.Tag.Value
}

func (gmm *GoMemberMeta) Doc() []string {
	if gmm.field.Doc == nil {
		return nil
	}
	commentSlice := make([]string, 0, len(gmm.field.Doc.List))
	for _, comment := range gmm.field.Doc.List {
		commentSlice = append(commentSlice, comment.Text)
	}
	return commentSlice
}

func (gmm *GoMemberMeta) Comment() string {
	if gmm.field.Comment == nil || len(gmm.field.Comment.List) == 0 {
		return ""
	}
	return gmm.field.Comment.List[0].Text
}
