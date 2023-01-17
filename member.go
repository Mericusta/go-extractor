package extractor

import "go/ast"

type GoMemberMeta struct {
	*meta
}

func SearchGoMemberMeta(m *meta, structType *ast.StructType, memberName string) *GoMemberMeta {
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
		meta: m.newMeta(memberDecl),
	}
}

func (gmm *GoMemberMeta) MemberName() string {
	return searchMemberName(gmm.node.(*ast.Field))
}

func searchMemberName(field *ast.Field) string {
	if field.Names == nil { // anonymous struct member
		switch fieldType := field.Type.(type) {
		case *ast.Ident:
			return fieldType.Name
		case *ast.StarExpr:
			return fieldType.X.(*ast.Ident).Name
		}
	} else { // named struct member
		return field.Names[0].Name
	}
	return ""
}

func (gmm *GoMemberMeta) Tag() string {
	if gmm.node.(*ast.Field).Tag == nil {
		return ""
	}
	return gmm.node.(*ast.Field).Tag.Value
}

func (gmm *GoMemberMeta) Doc() []string {
	if gmm.node.(*ast.Field).Doc == nil {
		return nil
	}
	commentSlice := make([]string, 0, len(gmm.node.(*ast.Field).Doc.List))
	for _, comment := range gmm.node.(*ast.Field).Doc.List {
		commentSlice = append(commentSlice, comment.Text)
	}
	return commentSlice
}

func (gmm *GoMemberMeta) Comment() string {
	if gmm.node.(*ast.Field).Comment == nil || len(gmm.node.(*ast.Field).Comment.List) == 0 {
		return ""
	}
	return gmm.node.(*ast.Field).Comment.List[0].Text
}
