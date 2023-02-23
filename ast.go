package extractor

import (
	"go/ast"
	"go/token"
)

func makeField(names []string, typeName string, from string, pointer bool) *ast.Field {
	var typeExpr ast.Expr = ast.NewIdent(typeName)
	if len(from) != 0 {
		typeExpr = &ast.SelectorExpr{
			X:   ast.NewIdent(from),
			Sel: ast.NewIdent(typeName),
		}
	}
	if pointer {
		typeExpr = &ast.StarExpr{
			X: typeExpr,
		}
	}
	namesIdent := make([]*ast.Ident, 0, len(names))
	for _, name := range names {
		namesIdent = append(namesIdent, ast.NewIdent(name))
	}
	return &ast.Field{
		Names: namesIdent,
		Type:  typeExpr,
	}
}

type fieldMaker interface {
	*field | *GoVariableMeta
	makeField() *ast.Field
	makeFieldList() *ast.FieldList
}

func makeFieldList[T fieldMaker](vs []T) *ast.FieldList {
	l := len(vs)
	if l == 0 {
		return nil
	}
	fieldList := &ast.FieldList{List: make([]*ast.Field, l)}
	for i, v := range vs {
		fieldList.List[i] = v.makeField()
	}
	return fieldList
}

type field struct {
	names    []string
	typeName string
	from     string
	pointer  bool
}

func (f *field) makeField() *ast.Field {
	return makeField(f.names, f.typeName, f.from, f.pointer)
}

func (f *field) makeFieldList() *ast.FieldList {
	return &ast.FieldList{
		List: []*ast.Field{
			makeField(f.names, f.typeName, f.from, f.pointer),
		},
	}
}

func (gvm *GoVariableMeta) makeField() *ast.Field {
	return &ast.Field{
		Names: []*ast.Ident{ast.NewIdent(gvm.Name())},
		Type:  gvm.typeMeta.(*meta).node.(ast.Expr),
	}
}

func (gvm *GoVariableMeta) makeFieldList() *ast.FieldList {
	return &ast.FieldList{
		List: []*ast.Field{gvm.makeField()},
	}
}

// ----------------------------------------------------------------

type funcDeclMaker interface {
	fieldMaker
}

func makeFuncDecl[T funcDeclMaker](name string, recv T, typeParams, params, returns []T) *ast.FuncDecl {
	if len(name) == 0 {
		return nil
	}

	var recvFieldList *ast.FieldList
	if recv != nil {
		recvFieldList = recv.makeFieldList()
	}

	paramsFieldList := makeFieldList(params)
	if paramsFieldList == nil {
		paramsFieldList = &ast.FieldList{}
	}

	return &ast.FuncDecl{
		Name: ast.NewIdent(name),
		Recv: recvFieldList,
		Type: &ast.FuncType{
			TypeParams: makeFieldList(typeParams),
			Params:     paramsFieldList,
			Results:    makeFieldList(returns),
		},
	}
}

type funcDecl struct {
	name       string
	doc        []string
	recv       *field
	typeParams []*field
	params     []*field
	returns    []*field
}

func (fd *funcDecl) make() *ast.FuncDecl {
	return makeFuncDecl(fd.name, fd.recv, fd.typeParams, fd.params, fd.returns)
}

func (gfm *GoFunctionMeta) makeFuncDecl() *ast.FuncDecl {
	return gfm.node.(*ast.FuncDecl)
}

// ----------------------------------------------------------------

type structDeclMaker interface {
	fieldMaker
}

// structDecl GenDecl - TypeSpec - StructType
func makeStructDecl[T structDeclMaker](name string, typeParams, members []T) *ast.GenDecl {
	structMemberFields := makeFieldList(members)
	if structMemberFields == nil {
		structMemberFields = &ast.FieldList{}
	}
	return &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name:       ast.NewIdent(name),
				TypeParams: makeFieldList(typeParams),
				Type: &ast.StructType{
					Fields: structMemberFields,
				},
			},
		},
	}
}
