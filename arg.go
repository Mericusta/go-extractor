package extractor

import "go/ast"

type GoArgMeta struct {
	expression string
	node       ast.Node
	arg        string
	from       string
	value      interface{}
	callMeta   *GoCallMeta
}

func (gam *GoArgMeta) Expression() string {
	return gam.expression
}

func (gam *GoArgMeta) Arg() string {
	return gam.arg
}

func (gam *GoArgMeta) From() string {
	return gam.from
}

func (gam *GoArgMeta) Value() interface{} {
	return gam.value
}

func (gam *GoArgMeta) CallMeta() *GoCallMeta {
	return gam.callMeta
}
