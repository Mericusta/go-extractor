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
		if searchMemberName(field) == memberName {
			memberDecl = field
			break
		}
	}
	return &GoMemberMeta{
		field: memberDecl,
	}
}

func (gmm *GoMemberMeta) MemberName() string {
	return searchMemberName(gmm.field)
}

func searchMemberName(field *ast.Field) string {
	if field.Names == nil {
		return field.Type.(*ast.StarExpr).X.(*ast.Ident).Name
	}
	return field.Names[0].Name
}

func (gmm *GoMemberMeta) Tag() string {
	if gmm.field.Tag == nil {
		return ""
	}
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
