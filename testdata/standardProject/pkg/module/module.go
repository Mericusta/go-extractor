package module

import "fmt"

type ExampleStruct struct {
	v int
}

func NewExampleStruct(v int) *ExampleStruct {
	return &ExampleStruct{v: v}
}

func (es ExampleStruct) ExampleFunc(v int) {
	fmt.Println("module.ExampleStruct.ExampleFunc Hello go-extractor,", v)
}

func (es ExampleStruct) V() int {
	return es.v
}

func ExampleFunc(s *ExampleStruct) {
	s.ExampleFunc(s.v)
}
