package extractor

import (
	"fmt"
	"go/ast"
)

// GoFuncMeta go func 的 meta 数据
type GoFuncMeta struct {
	// 组合基本 meta 数据
	// ast 节点，要求为 *ast.FuncDecl
	// 以 ast 节点 为单位执行 AST/PrintAST/Expression/Format
	*meta

	// func 标识
	ident string

	// func 参数
	params []*GoVarMeta

	// func 返回值
	returns []*GoVarMeta

	// func 模板参数
	typeParams []*GoVarMeta

	// func 生成时的 block stmt
	makeUpBlockStmt string

	// callMeta map[string][]*GoCallMeta
	// nonSelectorCallMeta map[string][]*GoCallMeta
	// selectorCallMeta    map[string]map[string][]*GoCallMeta
}

// newGoFuncMeta 通过 ast 构造 func 的 meta 数据
func newGoFuncMeta(m *meta, ident string, stopExtract ...bool) *GoFuncMeta {
	gfm := &GoFuncMeta{meta: m, ident: ident}
	if len(stopExtract) == 0 {
		gfm.ExtractAll()
	}
	return gfm
}

func (gfm *GoFuncMeta) funcDecl() *ast.FuncDecl {
	return gfm.node.(*ast.FuncDecl)
}

// -------------------------------- unit test --------------------------------

func (gfm *GoFuncMeta) Ident() string         { return gfm.ident }
func (gfm *GoFuncMeta) Params() []*GoVarMeta  { return gfm.params }
func (gfm *GoFuncMeta) Returns() []*GoVarMeta { return gfm.returns }

// -------------------------------- unit test --------------------------------

// -------------------------------- extractor --------------------------------

// ExtractGoFuncMeta 通过文件的绝对路径和 func 的 标识 提取文件中 func 的 meta 数据
func ExtractGoFuncMeta(extractFilepath, funcIdent string) (*GoFuncMeta, error) {
	// 提取 package
	gpm, err := ExtractGoPackageMeta(extractFilepath, nil)
	if err != nil {
		return nil, err
	}

	// 提取 func
	gpm.extractFunc()

	// 搜索 func
	gfm := gpm.SearchFuncMeta(funcIdent)
	if gfm == nil {
		return nil, fmt.Errorf("can not find func node")
	}

	return gfm, nil
}

// ExtractAll 提取 func 内所有 params，returns，typeParams 的 meta 数据
func (gfm *GoFuncMeta) ExtractAll() {
	// 提取 params
	gfm.extractParams()

	// 提取 returns
	gfm.extractReturns()

	// 提取 typeParams
	gfm.extractTypeParams()
}

func (gfm *GoFuncMeta) extractParams() {
	funcDecl := gfm.funcDecl()
	if funcDecl.Type == nil || funcDecl.Type.Params == nil || len(funcDecl.Type.Params.List) == 0 {
		return
	}

	pLen := len(funcDecl.Type.Params.List)
	gfm.params = make([]*GoVarMeta, 0, pLen)
	for _, field := range funcDecl.Type.Params.List {
		for _, name := range field.Names {
			gfm.params = append(gfm.params, newGoVarMeta(newMeta(field, gfm.path), name.String()))
		}
	}
}

func (gfm *GoFuncMeta) extractReturns() {
	funcDecl := gfm.funcDecl()
	if funcDecl.Type == nil || funcDecl.Type.Results == nil || len(funcDecl.Type.Results.List) == 0 {
		return
	}

	rLen := len(funcDecl.Type.Results.List)
	gfm.returns = make([]*GoVarMeta, 0, rLen)
	for _, field := range funcDecl.Type.Results.List {
		if len(field.Names) > 0 {
			for _, name := range field.Names {
				gfm.returns = append(gfm.returns, newGoVarMeta(newMeta(field, gfm.path), name.String()))
			}
		} else {
			gfm.returns = append(gfm.returns, newGoVarMeta(newMeta(field, gfm.path), ""))
		}
	}
}

func (gfm *GoFuncMeta) extractTypeParams() {
}

// -------------------------------- extractor --------------------------------

// -------------------------------- maker --------------------------------

// MakeUpFuncMeta 通过参数生成 func 的 meta 数据
func MakeUpFuncMeta(ident string, params []*GoVarMeta, returns []*GoVarMeta) *GoFuncMeta {
	astFuncDecl := &ast.FuncDecl{
		Name: ast.NewIdent(ident),
		Type: &ast.FuncType{
			// TypeParams: makeFieldList(typeParams),
			Params:  makeFieldList(params),
			Results: makeFieldList(returns),
		},
	}

	gfm := newGoFuncMeta(newMeta(astFuncDecl, ""), ident, true)

	return gfm
}

// SetBlockStmt 设置生成时的 stmt
func (gfm *GoFuncMeta) SetBlockStmt(bs string) {
	gfm.makeUpBlockStmt = bs
}

// Make 生成 func
func (gfm *GoFuncMeta) Make() string {
	gfmFormat := gfm.Format()
	var funcDecl *ast.FuncDecl = gfm.node.(*ast.FuncDecl)
	if funcDecl.Body == nil {
		return gfmFormat + " " + gfm.makeUpBlockStmt
	}
	return gfmFormat
}

// -------------------------------- maker --------------------------------

// func (gfm *GoFuncMeta) Doc() []string {
// 	if gfm.node.(*ast.FuncDecl) == nil || gfm.node.(*ast.FuncDecl).Doc == nil || len(gfm.node.(*ast.FuncDecl).Doc.List) == 0 {
// 		return nil
// 	}
// 	commentSlice := make([]string, 0, len(gfm.node.(*ast.FuncDecl).Doc.List))
// 	for _, comment := range gfm.node.(*ast.FuncDecl).Doc.List {
// 		commentSlice = append(commentSlice, comment.Text)
// 	}
// 	return commentSlice
// }

// func (gfm *GoFuncMeta) TypeParams() []*GoVarMeta {
// 	if gfm.node.(*ast.FuncDecl).Type == nil || gfm.node.(*ast.FuncDecl).Type.TypeParams == nil || len(gfm.node.(*ast.FuncDecl).Type.TypeParams.List) == 0 {
// 		return nil
// 	}

// 	tParamLen := len(gfm.node.(*ast.FuncDecl).Type.TypeParams.List)
// 	tParams := make([]*GoVarMeta, 0, tParamLen)
// 	for _, field := range gfm.node.(*ast.FuncDecl).Type.TypeParams.List {
// 		for _, name := range field.Names {
// 			tParams = append(tParams, &GoVarMeta{
// 				meta:     gfm.copyMeta(field),
// 				ident:    name.String(),
// 				typeMeta: gfm.copyMeta(field.Type),
// 			})
// 		}
// 	}
// 	return tParams
// }

// func (gfm *GoFuncMeta) ReplaceDecl(new *GoFuncMeta) {
// 	if new.node.(*ast.FuncDecl).Doc != nil {
// 		gfm.node.(*ast.FuncDecl).Doc = new.node.(*ast.FuncDecl).Doc
// 	}
// 	if new.node.(*ast.FuncDecl).Name != nil {
// 		gfm.node.(*ast.FuncDecl).Name = new.node.(*ast.FuncDecl).Name
// 	}
// 	if new.node.(*ast.FuncDecl).Recv != nil {
// 		gfm.node.(*ast.FuncDecl).Recv = new.node.(*ast.FuncDecl).Recv
// 	}
// 	if new.node.(*ast.FuncDecl).Type != nil {
// 		gfm.node.(*ast.FuncDecl).Type = new.node.(*ast.FuncDecl).Type
// 	}
// }

// type BlockMeta struct {
// 	*meta
// }

// func (gfm *GoFuncMeta) Body() *BlockMeta {
// 	return &BlockMeta{meta: gfm.copyMeta(gfm.node.(*ast.FuncDecl).Body)}
// }

// func (gfm *GoFuncMeta) ReplaceBody(new *BlockMeta) {
// 	gfm.node.(*ast.FuncDecl).Body = new.node.(*ast.BlockStmt)
// }

// func (gfm *GoFuncMeta) Expression() string {
// 	originPos := gfm.node.(*ast.FuncDecl).Pos()
// 	originEnd := gfm.node.(*ast.FuncDecl).End()
// 	if gfm.node.(*ast.FuncDecl).Doc != nil {
// 		originPos = gfm.node.(*ast.FuncDecl).Doc.Pos()
// 	}

// 	fileContent, err := os.ReadFile(gfm.path)
// 	if err != nil {
// 		return ""
// 	}
// 	fileContentLen := len(fileContent)
// 	if originPos > originEnd || int(originPos) >= fileContentLen || int(originEnd) > fileContentLen {
// 		return ""
// 	}
// 	return strings.TrimSpace(string(fileContent[originPos-1 : originEnd-1]))
// }

// func (gfm *GoFuncMeta) MakeUnitTest(typeArgs []string) (string, []byte) {
// 	return makeTest(unittestMaker, gfm, "", typeArgs)
// }

// func (gfm *GoFuncMeta) UnittestFuncName(typeArgs []string) string {
// 	return wrapTestType(UNITTEST, gfm.testFuncName(typeArgs))
// }

// func (gfm *GoFuncMeta) MakeBenchmark(typeArgs []string) (string, []byte) {
// 	return makeTest(benchmarkMaker, gfm, "", typeArgs)
// }

// func (gfm *GoFuncMeta) BenchmarkFuncName(typeArgs []string) string {
// 	return wrapTestType(BENCHMARK, gfm.testFuncName(typeArgs))
// }

// func (gfm *GoFunctionMeta) ReplaceFunctionDoc(newDoc []string) (string, string, error) {
// 	originContent := gfm.Expression()

// 	gfm.node.(*ast.FuncDecl).Doc = &ast.CommentGroup{
// 		List: make([]*ast.Comment, 0, len(newDoc)),
// 	}
// 	for _, comment := range newDoc {
// 		gfm.node.(*ast.FuncDecl).Doc.List = append(gfm.node.(*ast.FuncDecl).Doc.List, &ast.Comment{
// 			Text: comment,
// 		})
// 	}

// 	buffer := &bytes.Buffer{}
// 	err := format.Node(buffer, gfm.fileMeta.fileSet, gfm.node.(*ast.FuncDecl))
// 	if err != nil {
// 		panic(err)
// 	}

// 	return strings.ReplaceAll(originContent, "\r", ""), strings.ReplaceAll(buffer.String(), "\r", ""), nil
// }

// func (gfm *GoFunctionMeta) SearchCallMeta(call string) []*GoCallMeta {
// 	if gcm, has := gfm.callMeta[call]; gcm != nil && has {
// 		return gcm
// 	}

// 	if gfm.node.(*ast.FuncDecl) == nil {
// 		return nil
// 	}

// 	gfm.callMeta[call] = SearchGoCallMeta(gfm.meta, call)

// 	return gfm.callMeta[call]
// }

// func (gfm *GoFunctionMeta) Calls() map[string][]*GoCallMeta {
// 	return ExtractGoCallMeta(gfm.meta)
// }

// func (gfm *GoFunctionMeta) Search
