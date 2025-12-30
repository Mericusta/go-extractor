package extractor

import "go/ast"

// import (
// 	"fmt"
// 	"go/ast"
// 	"go/token"
// )

// func makeField(names []string, typeName string, from string, pointer bool) *ast.Field {
// 	var typeExpr ast.Expr = ast.NewIdent(typeName)
// 	if len(from) != 0 {
// 		typeExpr = &ast.SelectorExpr{
// 			X:   ast.NewIdent(from),
// 			Sel: ast.NewIdent(typeName),
// 		}
// 	}
// 	if pointer {
// 		typeExpr = &ast.StarExpr{
// 			X: typeExpr,
// 		}
// 	}
// 	namesIdent := make([]*ast.Ident, 0, len(names))
// 	for _, name := range names {
// 		namesIdent = append(namesIdent, ast.NewIdent(name))
// 	}
// 	return &ast.Field{
// 		Names: namesIdent,
// 		Type:  typeExpr,
// 	}
// }

type fieldMaker interface {
	*field | *GoVarMeta
	make() *ast.Field
}

func makeFieldList[T fieldMaker](vs []T) *ast.FieldList {
	l := len(vs)
	if l == 0 {
		return nil
	}
	fieldList := &ast.FieldList{List: make([]*ast.Field, l)}
	for i, v := range vs {
		fieldList.List[i] = v.make()
	}
	return fieldList
}

type field struct {
	names    []string
	typeName string
	from     string
	pointer  bool
}

// func newField(names []string, typeName, from string, pointer bool) *field {
// 	return &field{
// 		names:    names,
// 		typeName: typeName,
// 		from:     from,
// 		pointer:  pointer,
// 	}
// }

// func (f *field) make() *ast.Field {
// 	return makeField(f.names, f.typeName, f.from, f.pointer)
// }

// func (f *field) makeFieldList() *ast.FieldList {
// 	return makeFieldList([]*field{f})
// }

// // func (gvm *GoVariableMeta) make() *ast.Field {
// // 	typeExpr, ok := gvm.typeMeta.(*meta).node.(ast.Expr)
// // 	_, typeUnderlyingString, _ := gvm.Type()
// // 	if !ok {
// // 		// interface method
// // 		typeExpr = &ast.Ident{Name: gvm.Name()}
// // 	}
// // 	return &ast.Field{
// // 		Names: []*ast.Ident{ast.NewIdent(strings.ToLower(string(typeUnderlyingString[0])))},
// // 		Type:  typeExpr,
// // 	}
// // }

// func (gvm *GoVarMeta) make() *ast.Field {
// 	return &ast.Field{
// 		Names: []*ast.Ident{ast.NewIdent(gvm.Ident())},
// 		Type:  gvm.typeMeta.(*meta).node.(ast.Expr),
// 	}
// }

// func (gvm *GoVarMeta) makeFieldList() *ast.FieldList {
// 	return makeFieldList([]*GoVarMeta{gvm})
// }

// // ----------------------------------------------------------------

// type funcDeclMaker interface {
// 	fieldMaker
// 	makeFieldList() *ast.FieldList
// }

// func makeFuncDecl[T funcDeclMaker](name string, recv T, typeParams, params, returns []T) *ast.FuncDecl {
// 	if len(name) == 0 {
// 		return nil
// 	}

// 	var recvFieldList *ast.FieldList
// 	if recv != nil {
// 		recvFieldList = recv.makeFieldList()
// 	}

// 	paramsFieldList := makeFieldList(params)
// 	if paramsFieldList == nil {
// 		paramsFieldList = &ast.FieldList{}
// 	}

// 	return &ast.FuncDecl{
// 		Name: ast.NewIdent(name),
// 		Recv: recvFieldList,
// 		Type: &ast.FuncType{
// 			TypeParams: makeFieldList(typeParams),
// 			Params:     paramsFieldList,
// 			Results:    makeFieldList(returns),
// 		},
// 	}
// }

// type funcDecl struct {
// 	name       string
// 	doc        []string
// 	recv       *field
// 	typeParams []*field
// 	params     []*field
// 	returns    []*field
// }

// func newFuncDecl(name string, recv *field, typeParams, params, returns []*field) *funcDecl {
// 	return &funcDecl{
// 		name:       name,
// 		recv:       recv,
// 		typeParams: typeParams,
// 		params:     params,
// 		returns:    returns,
// 	}
// }

// func (fd *funcDecl) make() *ast.FuncDecl {
// 	return makeFuncDecl(fd.name, fd.recv, fd.typeParams, fd.params, fd.returns)
// }

// func (gfm *GoFuncMeta) makeFuncDecl() *ast.FuncDecl {
// 	return gfm.node.(*ast.FuncDecl)
// }

// // ----------------------------------------------------------------

// type genDeclMaker interface {
// 	make() ast.Spec
// }

// func makeDecl[T genDeclMaker](tok token.Token, vs []T) *ast.GenDecl {
// 	l := len(vs)
// 	if l == 0 {
// 		return nil
// 	}
// 	decl := &ast.GenDecl{Tok: tok, Specs: make([]ast.Spec, l)}
// 	for i, v := range vs {
// 		decl.Specs[i] = v.make()
// 	}
// 	return decl
// }

// // ----------------------------------------------------------------

// func makeTypeSpec[T fieldMaker](name string, typeParams, members []T) *ast.TypeSpec {
// 	structMemberFields := makeFieldList(members)
// 	if structMemberFields == nil {
// 		structMemberFields = &ast.FieldList{}
// 	}
// 	return &ast.TypeSpec{
// 		Name:       ast.NewIdent(name),
// 		TypeParams: makeFieldList(typeParams),
// 		Type: &ast.StructType{
// 			Fields: structMemberFields,
// 		},
// 	}
// }

// type structDeclMaker[T fieldMaker] interface {
// 	*typeSpec[T] | *GoStructMeta
// 	genDeclMaker
// }

// // makeStructDecl 约束传递 token.Type 的类型
// func makeStructDecl[T1 fieldMaker, T2 structDeclMaker[T1]](vs []T2) *ast.GenDecl {
// 	return makeDecl(token.TYPE, vs)
// }

// type typeSpec[T fieldMaker] struct {
// 	name       string
// 	typeParams []T
// 	members    []T
// }

// func newTypeSpec[T fieldMaker](name string, typeParams, members []T) *typeSpec[T] {
// 	return &typeSpec[T]{
// 		name:       name,
// 		typeParams: typeParams,
// 		members:    members,
// 	}
// }

// func (ss *typeSpec[T]) make() ast.Spec {
// 	return makeTypeSpec(ss.name, ss.typeParams, ss.members)
// }

// func (ss *typeSpec[T]) makeDecl() *ast.GenDecl {
// 	return makeStructDecl[T]([]*typeSpec[T]{ss})
// }

// func (gsm *GoStructMeta) make() ast.Spec {
// 	return gsm.node.(*ast.TypeSpec)
// }

// func (gsm *GoStructMeta) makeDecl() *ast.GenDecl {
// 	return makeStructDecl[*GoVarMeta]([]*GoStructMeta{gsm}) // ???
// 	// return makeStructDecl[*field]([]*GoStructMeta{gsm}) // ???
// }

// // ----------------------------------------------------------------

// func makeImportSpec(alias, importPath string) *ast.ImportSpec {
// 	importSpec := &ast.ImportSpec{
// 		Path: &ast.BasicLit{
// 			Kind:  token.STRING,
// 			Value: fmt.Sprintf("\"%v\"", importPath),
// 		},
// 	}
// 	if len(alias) > 0 {
// 		importSpec.Name = ast.NewIdent(alias)
// 	}
// 	return importSpec
// }

// type importDeclMaker interface {
// 	*importSpec | *GoImportMeta
// 	genDeclMaker
// }

// // makeStructDecl 约束传递 token.IMPORT 的类型
// func makeImportDecl[T importDeclMaker](vs []T) *ast.GenDecl {
// 	return makeDecl(token.IMPORT, vs)
// }

// type importSpec struct {
// 	alias string
// 	path  string
// }

// func newImportSpec(alias, path string) *importSpec {
// 	return &importSpec{alias: alias, path: path}
// }

// func (is *importSpec) make() ast.Spec {
// 	return makeImportSpec(is.alias, is.path)
// }

// func (is *importSpec) makeDecl() *ast.GenDecl {
// 	return makeImportDecl([]*importSpec{is})
// }

// func (gim *GoImportMeta) make() ast.Spec {
// 	return gim.node.(*ast.ImportSpec)
// }

// func (gim *GoImportMeta) makeDecl() *ast.GenDecl {
// 	return makeImportDecl([]*GoImportMeta{gim})
// }

// // ----------------------------------------------------------------

// // func makeSelectorExpr(v string, sel *selectorExpr) *ast.SelectorExpr {
// // 	// if sel == nil {
// // 	// 	return ast.NewIdent(v)
// // 	// }
// // 	// return

// // 	// selectorExpr := &ast.SelectorExpr{}

// // }

// // type selectorExpr struct {
// // 	v string
// // 	sel  *selectorExpr
// // }

// func makeCallExpr() {

// }

// type callExpr struct {
// }
