package module

import (
	"fmt"
	"math/rand"
)

type ParentStruct struct {
	P int // parent value
}

// ExampleStruct this is an example struct
// this is struct comment
// this is another struct comment
type ExampleStruct struct {
	*ParentStruct // parent struct
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
	es := &ExampleStruct{
		ParentStruct: &ParentStruct{P: v * 10},
		v:            v,
	}
	globalExampleStruct = &ExampleStruct{
		ParentStruct: &ParentStruct{P: rand.Intn(v)},
		v:            v / 10,
	}
	es.sub = globalExampleStruct
	globalExampleStruct.ParentStruct = es.ParentStruct
	return es
}

func (es ExampleStruct) ExampleFunc(v int) {
	nes := NewExampleStruct(v)
	fmt.Println("module.ExampleStruct.ExampleFunc Hello go-extractor",
		es, es.v, es.V(),
		nes, nes.v, nes.V(),
		nes.sub.v, es.sub.V(),
		globalExampleStruct,
		NewExampleStruct(nes.sub.V()),
	)
}

func (es *ExampleStruct) ExampleFuncWithPointerReceiver(v int) {
	fmt.Println("module.ExampleStruct.ExampleFuncWithPointerReceiver Hello go-extractor")
}

func (es ExampleStruct) V() int {
	return es.v
}

func (es *ExampleStruct) Sub() *ExampleStruct {
	return es.sub
}

func ExampleFunc(s *ExampleStruct) {
	s.ExampleFunc(s.v)
}
