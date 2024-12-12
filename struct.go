package extractor

import (
	"fmt"
	"go/ast"
)

type GoStructMetaTypeConstraints interface {
	*ast.TypeSpec

	ast.Node
}

// GoStructMeta go struct 的 meta 数据
type GoStructMeta[T GoStructMetaTypeConstraints] struct {
	// 组合基本 meta 数据
	// ast 节点，要求为 满足 IsStructNode 的 *ast.TypeSpec
	// 以 ast 节点 为单位执行 AST/PrintAST/Expression/Format
	*meta[T]

	// struct 标识
	ident string

	// struct 内所有 member 的 meta 数据
	// - key: member 标识
	memberMetaMap map[string]*GoVarMeta[*ast.Field]

	// struct 内所有 method 的 meta 数据
	// - key: method 标识
	methodMetaMap map[string]*GoMethodMeta[*ast.FuncDecl]

	commentGroup *ast.CommentGroup
}

// newGoStructMeta 通过 ast 构造 struct 的 meta 数据
func newGoStructMeta[T GoStructMetaTypeConstraints](m *meta[T], ident string, stopExtract ...bool) *GoStructMeta[T] {
	gsm := &GoStructMeta[T]{
		meta:          m,
		ident:         ident,
		memberMetaMap: make(map[string]*GoVarMeta[*ast.Field]),
		methodMetaMap: make(map[string]*GoMethodMeta[*ast.FuncDecl]),
	}
	if len(stopExtract) == 0 {
		gsm.ExtractAll()
	}
	return gsm
}

// -------------------------------- extractor --------------------------------

// ExtractGoStructMeta 通过文件的绝对路径和 struct 的 标识 提取文件中的 struct 的 meta 数据
func ExtractGoStructMeta[T GoStructMetaTypeConstraints](extractFilepath, structIdent string) (*GoStructMeta[*ast.TypeSpec], error) {
	// 提取 package
	gpm, err := ExtractGoPackageMeta[T](extractFilepath, nil)
	if err != nil {
		return nil, err
	}

	// 提取 struct
	gpm.extractStruct()

	// 搜索 struct
	gsm := gpm.SearchStructMeta(structIdent)
	if gsm == nil {
		return nil, fmt.Errorf("can not find struct node")
	}

	return gsm, nil
}

// ExtractAll 提取 struct 内所有 member，method 的 meta 数据
func (gsm *GoStructMeta[T]) ExtractAll() {
	var typeSpec *ast.TypeSpec = gsm.node
	structType := typeSpec.Type.(*ast.StructType)
	if structType == nil || structType.Fields == nil {
		return
	}

	for _, member := range structType.Fields.List {
		if len(member.Names) > 0 {
			// 非匿名成员
			for _, name := range member.Names {
				memberIdent := name.String()
				gsm.memberMetaMap[memberIdent] = newGoVarMeta(newMeta(member, gsm.path), memberIdent)
			}
		} else {
			// 匿名成员
			// TODO: 使用 GoVariableMeta
			gvm := newGoVarMeta(newMeta(member, gsm.path), "")
			gvm.ident = gvm.typeIdent
			gsm.memberMetaMap[gvm.ident] = gvm
			// var (
			// 	typeIdent           string
			// 	nodeHandler         func(n ast.Node, post ...func(ast.Node) bool) bool
			// 	starExprHandler     func(n ast.Node) bool
			// 	selectorExprHandler func(n ast.Node) bool
			// 	indexExprHandler    func(n ast.Node) bool
			// )
			// nodeHandler = func(n ast.Node, posts ...func(ast.Node) bool) bool {
			// 	ident, ok := n.(*ast.Ident)
			// 	if ident != nil && ok {
			// 		typeIdent = ident.String()
			// 		return false
			// 	} else {
			// 		for _, post := range posts {
			// 			if post != nil && !post(n) {
			// 				return false
			// 			}
			// 		}
			// 	}
			// 	return true
			// }
			// starExprHandler = func(n ast.Node) bool {
			// 	// 遇到 StarExpr 取 X
			// 	starExpr, ok := n.(*ast.StarExpr)
			// 	if starExpr == nil || !ok {
			// 		return true
			// 	}
			// 	return nodeHandler(starExpr.X, starExprHandler, selectorExprHandler, indexExprHandler)
			// }
			// selectorExprHandler = func(n ast.Node) bool {
			// 	// 遇到 SelectorExpr 取 Sel
			// 	selectorExpr, ok := n.(*ast.SelectorExpr)
			// 	if selectorExpr == nil || !ok {
			// 		return true
			// 	}
			// 	return nodeHandler(selectorExpr.Sel, starExprHandler, selectorExprHandler, indexExprHandler)
			// }
			// indexExprHandler = func(n ast.Node) bool {
			// 	// 遇到 IndexExp 取 X
			// 	indexExpr, ok := n.(*ast.IndexExpr)
			// 	if indexExpr == nil || !ok {
			// 		return true
			// 	}
			// 	return nodeHandler(indexExpr.X, starExprHandler, selectorExprHandler, indexExprHandler)
			// }
			// ast.Inspect(member.Type, func(n ast.Node) bool {
			// 	return n != nil && nodeHandler(n, starExprHandler, selectorExprHandler, indexExprHandler)
			// })
			// gsm.memberMetaMap[typeIdent] = NewGoVarMeta(newMeta(member, gsm.path), typeIdent)
		}
	}
}

// -------------------------------- extractor --------------------------------

func (gsm *GoStructMeta[T]) SearchMemberMeta(member string) *GoVarMeta[*ast.Field] {
	return gsm.memberMetaMap[member]
}

func (gsm *GoStructMeta[T]) SearchMethodMeta(method string) *GoMethodMeta[*ast.FuncDecl] {
	return gsm.methodMetaMap[method]
}

// func SearchGoStructMeta(gfm *GoFileMeta, structName string) *GoStructMeta {
// 	var structDecl *ast.TypeSpec
// 	var commentDecl *ast.CommentGroup
// 	ast.Inspect(gfm.node, func(n ast.Node) bool {
// 		if genDecl, ok := n.(*ast.GenDecl); ok {
// 			ast.Inspect(genDecl, func(n ast.Node) bool {
// 				if IsTypeNode(n) {
// 					typeSpec := n.(*ast.TypeSpec)
// 					if typeSpec.Name.String() == structName {
// 						structDecl = typeSpec
// 						commentDecl = genDecl.Doc
// 						return false
// 					}
// 				}
// 				return true
// 			})
// 			return false // genDecl traverse done
// 		}
// 		return structDecl == nil // already found
// 	})
// 	if structDecl == nil {
// 		return nil
// 	}
// 	return &GoStructMeta{
// 		meta:         gfm.copyMeta(structDecl),
// 		commentGroup: commentDecl,
// 		methodDecl:   make(map[string]*GoMethodMeta),
// 		memberMeta:   make(map[string]*GoVarMeta),
// 	}
// }

// -------------------------------- unit test --------------------------------

func (gsm *GoStructMeta[T]) Ident() string { return gsm.ident }
func (gsm *GoStructMeta[T]) MemberMetaMap() map[string]*GoVarMeta[*ast.Field] {
	return gsm.memberMetaMap
}
func (gsm *GoStructMeta[T]) MethodMetaMap() map[string]*GoMethodMeta[*ast.FuncDecl] {
	return gsm.methodMetaMap
}

// -------------------------------- unit test --------------------------------

func (gsm *GoStructMeta[T]) Doc() []string {
	if gsm.node == nil || gsm.commentGroup == nil || len(gsm.commentGroup.List) == 0 {
		return nil
	}
	commentSlice := make([]string, 0, len(gsm.commentGroup.List))
	for _, comment := range gsm.commentGroup.List {
		commentSlice = append(commentSlice, comment.Text)
	}
	return commentSlice
}

// func (gsm *GoStructMeta[T]) TypeParams() []*GoVarMeta {
// 	if gsm.node == nil || gsm.node.(*ast.TypeSpec).TypeParams == nil || len(gsm.node.(*ast.TypeSpec).TypeParams.List) == 0 {
// 		return nil
// 	}

// 	tParamLen := len(gsm.node.(*ast.TypeSpec).TypeParams.List)
// 	tParams := make([]*GoVarMeta, 0, tParamLen)
// 	for _, field := range gsm.node.(*ast.TypeSpec).TypeParams.List {
// 		for _, name := range field.Names {
// 			tParams = append(tParams, &GoVarMeta{
// 				meta:     gsm.copyMeta(field),
// 				ident:    name.String(),
// 				typeMeta: gsm.copyMeta(field.Type),
// 			})
// 		}
// 	}
// 	return tParams
// }

// func (gsm *GoStructMeta[T]) Members() []string {
// 	if gsm.node.(*ast.TypeSpec) == nil || gsm.node.(*ast.TypeSpec).Type == nil {
// 		return nil
// 	}
// 	structType, ok := gsm.node.(*ast.TypeSpec).Type.(*ast.StructType)
// 	if structType == nil || !ok || structType.Fields == nil || len(structType.Fields.List) == 0 {
// 		return nil
// 	}
// 	members := make([]string, 0, len(structType.Fields.List))
// 	for _, field := range structType.Fields.List {
// 		if memberName := searchMemberName(field); len(memberName) != 0 {
// 			members = append(members, memberName)
// 		}
// 	}
// 	return members
// }
