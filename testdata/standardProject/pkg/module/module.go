package module

import "fmt"

type ExampleStruct struct {
	v int
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
