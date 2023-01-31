package extractor

import (
	"go/ast"
)

type GoArgMeta struct {
	*meta
}

func (gam *GoArgMeta) Head() *GoVariableMeta {
	n := gam.node
SWITCH:
	switch headExpr := n.(type) {
	case *ast.SelectorExpr:
		n = headExpr.X
		goto SWITCH
	case *ast.CallExpr:
		n = headExpr.Fun
		goto SWITCH
	case *ast.Ident:
		// 没有 Obj：包名
		if headExpr.Obj == nil {
			gim, err := ExtractGoImportMeta(gam.path, headExpr.String())
			if err != nil {
				// fmt.Printf("extract go import meta from path %v by alias %v occurs error: %v", gam.path, headExpr.String(), err)
				return nil
			}
			return &GoVariableMeta{
				meta:     gam.newMeta(headExpr),
				name:     headExpr.String(),
				typeMeta: gim,
			}
		}

		// Obj 的类型：
		// - *ast.Field：父级函数/方法的接收器或参数表
		//   - Names 是签名
		//   - Type 是类型名
		// - *ast.AssignStmt：函数内局部变量
		//   - Lhs 是左值
		//   - Rhs 是右值
		// - *ast.ValueSpec：包级全局变量
		// - *ast.FuncDecl：包级函数
		// gvm := &GoVariableMeta{meta: gam.newMeta(gam.node)}
		// fmt.Printf("headExpr = %v\n", headExpr.String())
		// if headExpr.Obj != nil {
		// 	// 	gvm.node = headExpr.Obj
		// 	fmt.Printf("headExpr.Obj = %+v\n", headExpr.Obj)
		// }
		// return gvm

		switch obj := headExpr.Obj.Decl.(type) {
		case *ast.ValueSpec:
			return &GoVariableMeta{
				meta:     gam.newMeta(obj),
				name:     headExpr.String(),
				typeMeta: gam.newMeta(obj.Type),
			}
		case *ast.Field:
			return &GoVariableMeta{
				meta:     gam.newMeta(obj),
				name:     headExpr.String(),
				typeMeta: gam.newMeta(obj.Type),
			}
		case *ast.AssignStmt:
			ll, rl := len(obj.Lhs), len(obj.Rhs)
			if ll > rl && rl == 1 { // := 在右值为多返回值函数时有且仅有一个
				return &GoVariableMeta{
					meta:     gam.newMeta(obj),
					name:     headExpr.String(),
					typeMeta: gam.newMeta(obj.Rhs[0]),
				}
			} else {
				for index, lh := range obj.Lhs {
					if lh.(*ast.Ident).String() == headExpr.String() {
						return &GoVariableMeta{
							meta:     gam.newMeta(obj),
							name:     headExpr.String(),
							typeMeta: gam.newMeta(obj.Rhs[index]),
						}
					}
				}
			}
		case *ast.FuncDecl:
			return &GoVariableMeta{
				meta:     gam.newMeta(obj.Name),
				name:     headExpr.String(),
				typeMeta: gam.newMeta(obj.Type.Results.List[0].Type),
			}
		}
		return nil
	case *ast.BasicLit:
		constantMeta := gam.newMeta(headExpr)
		return &GoVariableMeta{
			meta:     constantMeta,
			name:     headExpr.Value,
			typeMeta: constantMeta,
		}
	default:
		return nil
	}
}

func (gam *GoArgMeta) Next(last string) string {
	n := gam.node
FOR:
	for {
		switch headExpr := n.(type) {
		case *ast.SelectorExpr:
			n = headExpr.X
			goto FOR
		case *ast.CallExpr:
			n = headExpr.Fun
			goto FOR
		case *ast.Ident:
			if headExpr.String() == last {

			}
		}
	}

	return ""
}
