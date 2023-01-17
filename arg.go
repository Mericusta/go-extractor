package extractor

const (
	ARG_TYPE_VARIABLE = 1 // 变量
	ARG_TYPE_BUILTIN  = 2 // 内建类型
	ARG_TYPE_CALL     = 3 // 函数调用
)

type GoArgMeta struct {
	*meta
	// node     ast.Node    // ast 节点
	// from     string      // SEL 信息
	// argType  int         // 变量类型
	// arg      string      // 变量
	// value    interface{} // 内建类型变量的值
	// callMeta *GoCallMeta // 函数调用
}

// func (gam *GoArgMeta) Arg() string {
// 	return gam.arg
// }

// func (gam *GoArgMeta) ArgType() int {
// 	return gam.argType
// }

// func (gam *GoArgMeta) From() string {
// 	return gam.from
// }

// func (gam *GoArgMeta) Value() interface{} {
// 	return gam.value
// }

// func (gam *GoArgMeta) CallMeta() *GoCallMeta {
// 	return gam.callMeta
// }
