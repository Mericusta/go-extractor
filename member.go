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
