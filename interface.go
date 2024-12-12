package extractor

import (
	"fmt"
	"go/ast"
)

type GoInterfaceMetaTypeConstraints interface {
	*ast.TypeSpec

	ast.Node
}

// GoInterfaceMeta go interface 的 meta 数据
type GoInterfaceMeta[T GoInterfaceMetaTypeConstraints] struct {
	// 组合基本 meta 数据
	// ast 节点，要求为 满足 IsInterfaceNode 的 *ast.TypeSpec
	// 以 ast 节点 为单位执行 AST/PrintAST/Expression/Format
	*meta[T]

	// interface 标识
	ident string

	// interface 内所有 method 的 meta 数据
	// - key: method 标识
	methodMetaMap map[string]*GoInterfaceMethodMeta[*ast.Field, T]

	commentGroup *ast.CommentGroup
}

// newGoInterfaceMeta 通过 ast 构造 interface 的 meta
func newGoInterfaceMeta[T GoInterfaceMetaTypeConstraints](m *meta[T], ident string, stopExtract ...bool) *GoInterfaceMeta[T] {
	gim := &GoInterfaceMeta[T]{
		meta:          m,
		ident:         ident,
		methodMetaMap: make(map[string]*GoInterfaceMethodMeta[*ast.Field, T]),
	}
	if len(stopExtract) == 0 {
		gim.ExtractAll()
	}
	return gim
}

// -------------------------------- extractor --------------------------------

// ExtractGoInterfaceMeta 通过文件的绝对路径和 interface 的 标识 提取文件中的 interface 的 meta 数据
func ExtractGoInterfaceMeta[T GoInterfaceMetaTypeConstraints](extractFilepath, interfaceIdent string) (*GoInterfaceMeta[*ast.TypeSpec], error) {
	// 提取 package
	gpm, err := ExtractGoPackageMeta[T](extractFilepath, nil)
	if err != nil {
		return nil, err
	}

	// 提取 interface
	gpm.extractInterface()

	// 搜索 interface
	gim := gpm.SearchInterfaceMeta(interfaceIdent)
	if gim == nil {
		return nil, fmt.Errorf("can not find interface node")
	}

	return gim, nil
}

// ExtractAll 提取 interface 内所有 method 的 meta 数据
func (gim *GoInterfaceMeta[T]) ExtractAll() {
	var typeSpec *ast.TypeSpec = gim.node
	interfaceType := typeSpec.Type.(*ast.InterfaceType)
	if interfaceType == nil || interfaceType.Methods == nil {
		return
	}

	for _, method := range interfaceType.Methods.List {
		if IsInterfaceMethodNode(method) {
			for _, name := range method.Names {
				methodIdent := name.String()
				gim.methodMetaMap[methodIdent] = newGoInterfaceMethodMeta[*ast.Field, T](newMeta(method, gim.path), methodIdent, gim)
			}
		}
	}
}

// // extractGoInterfaceMeta 通过 interface 的 标识 提取 文件 的 meta 数据中的 interface 的 meta 数据
// func extractGoInterfaceMeta(extractFilepath, interfaceIdent string) (*GoInterfaceMeta, error) {
// 	gpm, err := ExtractGoPackageMeta(extractFilepath, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = gpm.extractStruct()
// 	if err != nil {
// 		return nil, err
// 	}

// 	gim := gpm.SearchInterfaceMeta(interfaceIdent)
// 	if gim == nil {
// 		return nil, fmt.Errorf("can not find interface node")
// 	}

// 	return gim, nil
// }

// -------------------------------- extractor --------------------------------

func (gim *GoInterfaceMeta[T]) SearchMethodMeta(methodIdent string) *GoInterfaceMethodMeta[*ast.Field, T] {
	return gim.methodMetaMap[methodIdent]
}

// func SearchGoInterfaceMeta(gfm *GoFileMeta, interfaceName string) *GoInterfaceMeta {
// 	var interfaceDecl *ast.TypeSpec
// 	var commentDecl *ast.CommentGroup
// 	ast.Inspect(gfm.node, func(n ast.Node) bool {
// 		if genDecl, ok := n.(*ast.GenDecl); ok {
// 			ast.Inspect(genDecl, func(n ast.Node) bool {
// 				if IsTypeNode(n) {
// 					typeSpec := n.(*ast.TypeSpec)
// 					if typeSpec.Name.String() == interfaceName {
// 						interfaceDecl = typeSpec
// 						commentDecl = genDecl.Doc
// 						return false
// 					}
// 				}
// 				return true
// 			})
// 			return false // genDecl traverse done
// 		}
// 		return interfaceDecl == nil // already found
// 	})
// 	if interfaceDecl == nil {
// 		return nil
// 	}
// 	return &GoInterfaceMeta{
// 		meta:         gfm.copyMeta(interfaceDecl),
// 		commentGroup: commentDecl,
// 		methodMeta:   make(map[string]*GoInterfaceMethodMeta),
// 	}
// }

// -------------------------------- unit test --------------------------------

func (gim *GoInterfaceMeta[T]) Ident() string { return gim.ident }
func (gim *GoInterfaceMeta[T]) MethodMetaMap() map[string]*GoInterfaceMethodMeta[*ast.Field, T] {
	return gim.methodMetaMap
}

// -------------------------------- unit test --------------------------------

func (gim *GoInterfaceMeta[T]) Doc() []string {
	if gim.node == nil || gim.commentGroup == nil || len(gim.commentGroup.List) == 0 {
		return nil
	}
	commentSlice := make([]string, 0, len(gim.commentGroup.List))
	for _, comment := range gim.commentGroup.List {
		commentSlice = append(commentSlice, comment.Text)
	}
	return commentSlice
}

// // SearchMethodDecl search method decl from node.(*ast.InterfaceType)
// func (gim *GoInterfaceMeta[T]) SearchMethodDecl(methodName string) *GoInterfaceMethodMeta {
// 	gim.ForeachMethodDecl(func(f *ast.Field) bool {
// 		if f.Names[0].String() == methodName && IsInterfaceMethodNode(f) {
// 			gim.methodMeta[methodName] = NewGoInterfaceMethodMeta(
// 				gim.copyMeta(f), methodName, gim,
// 			)
// 			return false
// 		}
// 		return true
// 	})
// 	return gim.methodMeta[methodName]
// }

// func (gim *GoInterfaceMeta[T]) ForeachMethodDecl(f func(*ast.Field) bool) {
// 	interfaceType := gim.node.(*ast.TypeSpec).Type.(*ast.InterfaceType)
// 	if interfaceType.Methods == nil {
// 		return
// 	}
// 	for _, methodField := range interfaceType.Methods.List {
// 		_, ok := methodField.Type.(*ast.FuncType)
// 		if ok {
// 			if !f(methodField) {
// 				break
// 			}
// 		}
// 	}
// }

// func (gim *GoInterfaceMeta[T]) TypeParams() []*GoVarMeta {
// 	if gim.node == nil || gim.node.(*ast.TypeSpec).TypeParams == nil || len(gim.node.(*ast.TypeSpec).TypeParams.List) == 0 {
// 		return nil
// 	}

// 	tParamLen := len(gim.node.(*ast.TypeSpec).TypeParams.List)
// 	tParams := make([]*GoVarMeta, 0, tParamLen)
// 	for _, field := range gim.node.(*ast.TypeSpec).TypeParams.List {
// 		for _, name := range field.Names {
// 			tParams = append(tParams, &GoVarMeta{
// 				meta:     gim.copyMeta(field),
// 				ident:    name.String(),
// 				typeMeta: gim.copyMeta(field.Type),
// 			})
// 		}
// 	}
// 	return tParams
// }

// func (gim *GoInterfaceMeta[T]) AllMethodIdentSlice() []string {
// 	methodIdentSlice := make([]string, 0, 8)
// 	gim.ForeachMethodDecl(func(f *ast.Field) bool {
// 		methodIdentSlice = append(methodIdentSlice, f.Names[0].String())
// 		return true
// 	})
// 	return methodIdentSlice
// }
