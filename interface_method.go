package extractor

import (
	"fmt"
	"go/ast"
)

// GoInterfaceMethodMeta go interface 的 method 的 meta 数据
type GoInterfaceMethodMeta struct {
	// 组合基本 meta 数据
	// ast 节点，要求为 满足 IsInterfaceMethodNode 的 *ast.Field
	// 以 ast 节点 为单位执行 AST/PrintAST/Expression/Format
	*meta

	// method 所属的 interface 的 meta 数据
	interfaceMeta *GoInterfaceMeta

	// interface 的 method 的标识
	ident string

	// func 参数
	params []*GoVarMeta

	// func 返回值
	returns []*GoVarMeta

	// func 模板参数
	typeParams []*GoVarMeta
}

// newGoInterfaceMethodMeta 通过 ast 构造 interface 的 method 的 meta 数据
func newGoInterfaceMethodMeta(m *meta, ident string, gim *GoInterfaceMeta, stopExtract ...bool) *GoInterfaceMethodMeta {
	gimm := &GoInterfaceMethodMeta{meta: m, ident: ident, interfaceMeta: gim}
	if len(stopExtract) == 0 {
		gimm.ExtractAll()
	}
	return gimm
}

func (gimm *GoInterfaceMethodMeta) funcType() *ast.FuncType {
	var field *ast.Field = gimm.node.(*ast.Field)
	if field.Type == nil {
		return nil
	}
	return field.Type.(*ast.FuncType)
}

// -------------------------------- extractor --------------------------------

// ExtractGoInterfaceMethodMeta 通过文件的绝对路径和 interface 以及 method 的 标识 提取文件中的 interface 的 method 的 meta 数据
func ExtractGoInterfaceMethodMeta(extractFilepath, interfaceIdent, methodIdent string) (*GoInterfaceMethodMeta, error) {
	// 提取 package
	gpm, err := ExtractGoPackageMeta(extractFilepath, nil)
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

	// 搜索 method
	gimm := gim.SearchMethodMeta(interfaceIdent)
	if gimm == nil {
		return nil, fmt.Errorf("can not find interface method node")
	}

	return gimm, nil
}

// ExtractAll 提取 func 内所有 params，returns，typeParams 的 meta 数据
func (gimm *GoInterfaceMethodMeta) ExtractAll() {
	// 提取 params
	gimm.extractParams()

	// 提取 returns
	gimm.extractReturns()

	// 提取 typeParams
	gimm.extractTypeParams()
}

func (gimm *GoInterfaceMethodMeta) extractParams() {
	funcType := gimm.funcType()
	if funcType == nil || funcType.Params == nil || len(funcType.Params.List) == 0 {
		return
	}

	pLen := len(funcType.Params.List)
	gimm.params = make([]*GoVarMeta, 0, pLen)
	for _, field := range funcType.Params.List {
		if len(field.Names) > 0 {
			// 定义参数名称的方法
			for _, name := range field.Names {
				gimm.params = append(gimm.params, newGoVarMeta(newMeta(field, gimm.path), name.String()))
			}
		} else {
			// 未定义参数名称的方法
			gimm.params = append(gimm.params, newGoVarMeta(newMeta(field, gimm.path), ""))
		}
	}
}

func (gimm *GoInterfaceMethodMeta) extractReturns() {
	funcType := gimm.funcType()
	if funcType == nil || funcType.Results == nil || len(funcType.Results.List) == 0 {
		return
	}

	rLen := len(funcType.Results.List)
	gimm.returns = make([]*GoVarMeta, 0, rLen)
	for _, field := range funcType.Results.List {
		if len(field.Names) > 0 {
			// 定义返回值名称的方法
			for _, name := range field.Names {
				gimm.returns = append(gimm.returns, newGoVarMeta(newMeta(field, gimm.path), name.String()))
			}
		} else {
			// 未定义返回值名称的方法
			gimm.returns = append(gimm.returns, newGoVarMeta(newMeta(field, gimm.path), ""))
		}
	}
}

func (gimm *GoInterfaceMethodMeta) extractTypeParams() {
}

// -------------------------------- extractor --------------------------------

// -------------------------------- unit test --------------------------------

func (gimm *GoInterfaceMethodMeta) Ident() string         { return gimm.ident }
func (gimm *GoInterfaceMethodMeta) Params() []*GoVarMeta  { return gimm.params }
func (gimm *GoInterfaceMethodMeta) Returns() []*GoVarMeta { return gimm.returns }

// -------------------------------- unit test --------------------------------

// func (gimm *GoInterfaceMethodMeta) FunctionName() string {
// 	return gimm.node.(*ast.Field).Names[0].String()
// }

// func (gimm *GoInterfaceMethodMeta) Doc() []string {
// 	if gimm.node.(*ast.Field) == nil || gimm.node.(*ast.Field).Doc == nil || len(gimm.node.(*ast.Field).Doc.List) == 0 {
// 		return nil
// 	}
// 	commentSlice := make([]string, 0, len(gimm.node.(*ast.Field).Doc.List))
// 	for _, comment := range gimm.node.(*ast.Field).Doc.List {
// 		commentSlice = append(commentSlice, comment.Text)
// 	}
// 	return commentSlice
// }

// func (gimm *GoInterfaceMethodMeta) TypeParams() []*GoVarMeta {
// 	return gimm.interfaceMeta.TypeParams()
// }

// func (gimm *GoInterfaceMethodMeta) Params() []*GoVarMeta {
// 	if gimm.node.(*ast.Field).Type == nil || gimm.node.(*ast.Field).Type.(*ast.FuncType).Params == nil || len(gimm.node.(*ast.Field).Type.(*ast.FuncType).Params.List) == 0 {
// 		return nil
// 	}

// 	pLen := len(gimm.node.(*ast.Field).Type.(*ast.FuncType).Params.List)
// 	params := make([]*GoVarMeta, 0, pLen)
// 	for index, field := range gimm.node.(*ast.Field).Type.(*ast.FuncType).Params.List {
// 		params = append(params, &GoVarMeta{
// 			meta:     gimm.copyMeta(field),
// 			ident:    fmt.Sprintf("p%v", index),
// 			typeMeta: gimm.copyMeta(field.Type),
// 		})

// 	}
// 	return params
// }

// func (gimm *GoInterfaceMethodMeta) ReturnTypes() []*GoVarMeta {
// 	if gimm.node.(*ast.Field).Type == nil || gimm.node.(*ast.Field).Type.(*ast.FuncType).Results == nil || len(gimm.node.(*ast.Field).Type.(*ast.FuncType).Results.List) == 0 {
// 		return nil
// 	}

// 	rLen := len(gimm.node.(*ast.Field).Type.(*ast.FuncType).Results.List)
// 	returns := make([]*GoVarMeta, 0, rLen)
// 	for _, field := range gimm.node.(*ast.Field).Type.(*ast.FuncType).Results.List {
// 		// TODO: not support named return value
// 		returns = append(returns, &GoVarMeta{
// 			meta:     gimm.copyMeta(field),
// 			ident:    "",
// 			typeMeta: gimm.copyMeta(field.Type),
// 		})
// 	}
// 	return returns
// }

// func (gimm *GoInterfaceMethodMeta) RecvInterface() (string, bool) {
// 	return gimm.interfaceMeta.Ident(), true
// }

// func (gimm *GoInterfaceMethodMeta) Recv() *GoVarMeta {
// 	if gimm.receiverMeta != nil {
// 		return gimm.receiverMeta
// 	}

// 	var receiverTypeExpr ast.Expr = ast.NewIdent(gimm.interfaceMeta.Ident())
// 	typeParams := gimm.TypeParams()
// 	if l := len(typeParams); l > 0 {
// 		typeParamsExpr := make([]ast.Expr, 0, l)
// 		for _, typeParam := range typeParams {
// 			typeParamsExpr = append(typeParamsExpr, ast.NewIdent(typeParam.Ident()))
// 		}
// 		if l == 1 {
// 			receiverTypeExpr = &ast.IndexExpr{
// 				X:     receiverTypeExpr,
// 				Index: typeParamsExpr[0],
// 			}
// 		} else {
// 			receiverTypeExpr = &ast.IndexListExpr{
// 				X:       receiverTypeExpr,
// 				Indices: typeParamsExpr,
// 			}
// 		}
// 	}
// 	gimm.receiverMeta = &GoVarMeta{
// 		meta:     gimm.copyMeta(gimm.interfaceMeta.node),
// 		ident:    "i",
// 		typeMeta: gimm.copyMeta(receiverTypeExpr),
// 	}

// 	return gimm.receiverMeta
// }

// func (gimm *GoInterfaceMethodMeta) MakeUnitTest(typeArgs []string) (string, []byte) {
// 	return makeTest(unittestMaker, gimm, "", typeArgs)
// }

// func (gimm *GoInterfaceMethodMeta) MakeBenchmark(typeArgs []string) (string, []byte) {
// 	return makeTest(benchmarkMaker, gimm, "", typeArgs)
// }

// func (gimm *GoInterfaceMethodMeta) MakeImplementMethodMeta(receiverIdent, receiverType string) (string, *GoFuncMeta) {
// 	if len(receiverIdent) == 0 || len(receiverType) == 0 {
// 		return "", nil
// 	}

// 	funcName := gimm.FunctionName()

// 	// func params
// 	params := gimm.Params()
// 	paramsFieldList := make([]*field, 0, len(params))
// 	for _, gvm := range params {
// 		typeNodeIdent, isPointer := traitReceiverStruct(gvm.typeMeta.(*meta).node)
// 		paramsFieldList = append(paramsFieldList, newField(
// 			[]string{gvm.Ident()}, typeNodeIdent.String(), "", isPointer,
// 		))
// 	}

// 	// func returns
// 	returnTypes := gimm.ReturnTypes()
// 	returnFieldSlice := make([]*field, 0, len(returnTypes))
// 	for _, rts := range returnTypes {
// 		typeNodeIdent, isPointer := traitReceiverStruct(rts.typeMeta.(*meta).node)
// 		returnFieldSlice = append(returnFieldSlice, newField(
// 			nil, typeNodeIdent.String(), "", isPointer,
// 		))
// 	}

// 	// func decl
// 	funcDecl := makeFuncDecl(
// 		funcName,
// 		newField([]string{receiverIdent}, receiverType, "", true),
// 		nil,
// 		paramsFieldList,
// 		returnFieldSlice,
// 	)
// 	funcDecl.Body = &ast.BlockStmt{}

// 	funcMeta := gimm.copyMeta(funcDecl)
// 	// funcMeta.PrintAST()

// 	// // output
// 	// buffer := &bytes.Buffer{}
// 	// err := format.Node(buffer, token.NewFileSet(), funcDecl)
// 	// if err != nil {
// 	// 	panic(err)
// 	// }

// 	// ex, e := parser.Parse(buffer.String())
// 	// fmt.Printf("ex = %v, e = %v\n", ex, e)
// 	// return funcName, buffer.Bytes()
// 	return funcName, NewGoFunctionMeta(funcMeta, "")
// }
