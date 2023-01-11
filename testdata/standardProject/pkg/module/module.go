package module

import "fmt"

// ExampleStruct this is an example struct
// this is struct comment
// this is another struct comment
type ExampleStruct struct {
	v int `ast:init,default=1`
}

func NewExampleStruct(v int) *ExampleStruct {
	return &ExampleStruct{v: v}
}

func (es ExampleStruct) ExampleFunc(v int) {
	nes := NewExampleStruct(v)
	fmt.Println("module.ExampleStruct.ExampleFunc Hello go-extractor,", nes.V())
}

func (es *ExampleStruct) AnotherExampleFunc(v int) {
	nes := NewExampleStruct(v)
	fmt.Println("module.ExampleStruct.ExampleFunc Hello go-extractor,", nes.V())
}

func (es ExampleStruct) V() int {
	return es.v
}

func ExampleFunc(s *ExampleStruct) {
	s.ExampleFunc(s.v)
}
