package extractor

import "go/ast"

func SearchGoMemberMeta(m *meta, structType *ast.StructType, memberName string) *GoVariableMeta {
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
	return &GoVariableMeta{
		meta:     m.newMeta(memberDecl),
		name:     memberName,
		typeMeta: m.newMeta(memberDecl.Type),
	}
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
