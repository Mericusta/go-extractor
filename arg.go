package extractor

// import (
// 	"go/ast"
// )

// type GoArgMeta struct {
// 	*meta
// }

// func (gam *GoArgMeta) Head() *GoVariableMeta {
// 	n := gam.node
// SWITCH:
// 	switch headExpr := n.(type) {
// 	case *ast.SelectorExpr:
// 		n = headExpr.X
// 		goto SWITCH
// 	case *ast.CallExpr:
// 		n = headExpr.Fun
// 		goto SWITCH
// 	case *ast.Ident:
// 		// 没有 Obj：包名
// 		if headExpr.Obj == nil {
// 			gim, err := ExtractGoImportMeta(gam.path, headExpr.String())
// 			if err != nil {
// 				// fmt.Printf("extract go import meta from path %v by alias %v occurs error: %v", gam.path, headExpr.String(), err)
// 				return nil
// 			}
// 			return &GoVariableMeta{
// 				meta:     gam.newMeta(headExpr),
// 				name:     headExpr.String(),
// 				typeMeta: gim,
// 				typeEnum: TYPE_PKG_ALIAS,
// 			}
// 		}

// 		// Obj 的类型：
// 		// - *ast.Field：父级函数/方法的接收器或参数表
// 		//   - Names 是签名
// 		//   - Type 是类型名
// 		// - *ast.AssignStmt：函数内局部变量
// 		//   - Lhs 是左值
// 		//   - Rhs 是右值
// 		// - *ast.ValueSpec：包级全局变量
// 		// - *ast.FuncDecl：包级函数
// 		// gvm := &GoVariableMeta{meta: gam.newMeta(gam.node)}
// 		// fmt.Printf("headExpr = %v\n", headExpr.String())
// 		// if headExpr.Obj != nil {
// 		// 	// 	gvm.node = headExpr.Obj
// 		// 	fmt.Printf("headExpr.Obj = %+v\n", headExpr.Obj)
// 		// }
// 		// return gvm

// 		switch obj := headExpr.Obj.Decl.(type) {
// 		case *ast.ValueSpec:
// 			return &GoVariableMeta{
// 				meta:     gam.newMeta(obj),
// 				name:     headExpr.String(),
// 				typeMeta: gam.newMeta(obj.Type),
// 				typeEnum: TYPE_VAR_FIELD,
// 			}
// 		case *ast.Field:
// 			return &GoVariableMeta{
// 				meta:     gam.newMeta(obj),
// 				name:     headExpr.String(),
// 				typeMeta: gam.newMeta(obj.Type),
// 				typeEnum: TYPE_VAR_FIELD,
// 			}
// 		case *ast.AssignStmt:
// 			ll, rl := len(obj.Lhs), len(obj.Rhs)
// 			if ll > rl && rl == 1 { // := 在右值为多返回值函数时有且仅有一个
// 				return &GoVariableMeta{
// 					meta:     gam.newMeta(obj),
// 					name:     headExpr.String(),
// 					typeMeta: gam.newMeta(obj.Rhs[0]),
// 					typeEnum: TYPE_ASSIGNMENT,
// 				}
// 			} else {
// 				for index, lh := range obj.Lhs {
// 					if lh.(*ast.Ident).String() == headExpr.String() {
// 						return &GoVariableMeta{
// 							meta:     gam.newMeta(obj),
// 							name:     headExpr.String(),
// 							typeMeta: gam.newMeta(obj.Rhs[index]),
// 							typeEnum: TYPE_ASSIGNMENT,
// 						}
// 					}
// 				}
// 			}
// 		case *ast.FuncDecl:
// 			return &GoVariableMeta{
// 				meta:     gam.newMeta(obj),
// 				name:     headExpr.String(),
// 				typeMeta: gam.newMeta(obj.Type.Results.List[0].Type),
// 				typeEnum: TYPE_FUNC_CALL,
// 			}
// 		}
// 		return nil
// 	case *ast.BasicLit:
// 		constantMeta := gam.newMeta(headExpr)
// 		return &GoVariableMeta{
// 			meta:     constantMeta,
// 			name:     headExpr.Value,
// 			typeMeta: constantMeta,
// 			typeEnum: TYPE_CONSTANTS,
// 		}
// 	default:
// 		return nil
// 	}
// }

// func (gam *GoArgMeta) Slice() []*GoVariableMeta {
// 	s := make([]*GoVariableMeta, 0, 8)
// 	var joinFunc func(ast.Node) bool
// 	joinFunc = func(n ast.Node) bool {
// 		switch node := n.(type) {
// 		case *ast.Ident:
// 			if node.Obj == nil {
// 				if len(s) == 0 {
// 					gim, err := ExtractGoImportMeta(gam.path, node.String())
// 					if err != nil {
// 						return false
// 					}
// 					s = append(s, &GoVariableMeta{
// 						meta:     gam.newMeta(node),
// 						name:     node.String(),
// 						typeMeta: gim,
// 						typeEnum: TYPE_PKG_ALIAS,
// 					})
// 				} else {
// 					m := gam.newMeta(node)
// 					s = append(s, &GoVariableMeta{
// 						meta:     m,
// 						name:     node.String(),
// 						typeMeta: m,
// 						typeEnum: TYPE_VAR_FIELD,
// 					})
// 				}
// 			} else {
// 				switch obj := node.Obj.Decl.(type) {
// 				case *ast.ValueSpec:
// 					s = append(s, &GoVariableMeta{
// 						meta:     gam.newMeta(obj),
// 						name:     node.String(),
// 						typeMeta: gam.newMeta(obj.Type),
// 						typeEnum: TYPE_VAR_FIELD,
// 					})
// 				case *ast.Field:
// 					s = append(s, &GoVariableMeta{
// 						meta:     gam.newMeta(obj),
// 						name:     node.String(),
// 						typeMeta: gam.newMeta(obj.Type),
// 						typeEnum: TYPE_VAR_FIELD,
// 					})
// 				case *ast.AssignStmt:
// 					ll, rl := len(obj.Lhs), len(obj.Rhs)
// 					if ll > rl && rl == 1 { // := 在右值为多返回值函数时有且仅有一个
// 						s = append(s, &GoVariableMeta{
// 							meta:     gam.newMeta(obj),
// 							name:     node.String(),
// 							typeMeta: gam.newMeta(obj.Rhs[0]),
// 							typeEnum: TYPE_ASSIGNMENT,
// 						})
// 					} else {
// 						for index, lh := range obj.Lhs {
// 							if lh.(*ast.Ident).String() == node.String() {
// 								s = append(s, &GoVariableMeta{
// 									meta:     gam.newMeta(obj),
// 									name:     node.String(),
// 									typeMeta: gam.newMeta(obj.Rhs[index]),
// 									typeEnum: TYPE_ASSIGNMENT,
// 								})
// 							}
// 						}
// 					}
// 				case *ast.FuncDecl:
// 					s = append(s, &GoVariableMeta{
// 						meta:     gam.newMeta(obj),
// 						name:     node.String(),
// 						typeMeta: gam.newMeta(obj.Type.Results.List[0].Type),
// 						typeEnum: TYPE_FUNC_CALL,
// 					})
// 				}
// 			}
// 			// typeEnum = TYPE_VAR_FIELD
// 		case *ast.BasicLit:
// 			constantMeta := gam.newMeta(node)
// 			s = append(s, &GoVariableMeta{
// 				meta:     constantMeta,
// 				name:     node.Value,
// 				typeMeta: constantMeta,
// 				typeEnum: TYPE_CONSTANTS,
// 			})
// 		case *ast.CallExpr:
// 			// typeEnum = TYPE_FUNC_CALL
// 			ast.Inspect(node.Fun, joinFunc)
// 			if l := len(s); l > 1 {
// 				s[l-1].typeEnum = TYPE_FUNC_CALL
// 			}
// 			// typeEnum = TYPE_VAR_FIELD
// 			return false
// 		}
// 		return true
// 	}
// 	ast.Inspect(gam.node, joinFunc)
// 	return s
// }

// func (gam *GoArgMeta) next(last ast.Node) *GoVariableMeta {
// 	var nextExpr ast.Node
// 	var nextExprName string
// 	var nextTypeEnum VariableTypeEnum
// 	ast.Inspect(gam.node, func(n ast.Node) bool {
// 		if IsSelectorNode(n) {
// 			expr := n.(*ast.SelectorExpr)
// 			if expr.X == last {
// 				nextExpr = expr
// 				nextExprName = expr.Sel.String()
// 			}
// 		}
// 		return nextExpr == nil
// 	})

// 	// 再遍历一遍，把节点往上提升一层，为了处理 a.B() 的情况
// 	ast.Inspect(gam.node, func(n ast.Node) bool {
// 		if IsCallNode(n) {
// 			expr := n.(*ast.CallExpr)
// 			if expr.Fun == nextExpr {
// 				nextExpr = expr
// 				nextTypeEnum = TYPE_FUNC_CALL
// 				return false
// 			}
// 		}
// 		return true
// 	})

// 	return &GoVariableMeta{
// 		meta:     gam.newMeta(nextExpr),
// 		name:     nextExprName,
// 		typeEnum: nextTypeEnum,
// 	}
// }
