package module

import "fmt"

type ParentStruct struct {
	P int
}

// ExampleStruct this is an example struct
// this is struct comment
// this is another struct comment
type ExampleStruct struct {
	*ParentStruct
	// v this is member doc line1
	// v this is member doc line2
	v   int `ast:init,default=1` // this is member single comment line
	sub *ExampleStruct
}

var globalExampleStruct *ExampleStruct

// NewExampleStruct this is new example struct
// @param           value
// @return          pointer to ExampleStruct
func NewExampleStruct(v int) *ExampleStruct {
	globalExampleStruct = &ExampleStruct{v: v + 1}
	return &ExampleStruct{v: v, sub: globalExampleStruct}
}

func (es ExampleStruct) ExampleFunc(v int) {
	nes := NewExampleStruct(v)
	fmt.Println("module.ExampleStruct.ExampleFunc Hello go-extractor", es, es.v, es.V(), nes, nes.v, nes.V(), nes.sub.v, es.sub.V(), globalExampleStruct)
}

func (es *ExampleStruct) ExampleFuncWithPointerReceiver(v int) {
	fmt.Println("module.ExampleStruct.ExampleFuncWithPointerReceiver Hello go-extractor")
}

func (es ExampleStruct) V() int {
	return es.v
}

func ExampleFunc(s *ExampleStruct) {
	s.ExampleFunc(s.v)
}
