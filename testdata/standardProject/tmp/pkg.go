package pkg

import (
	"fmt"
	"standardProject/pkg/module"
)

// ExampleFunc this is example function
func ExampleFunc(s *module.ExampleStruct) {
	fmt.Println("pkg.ExampleFunc, Hello go-extractor", s.V())
}

func NoDocExampleFunc(s *module.ExampleStruct) {
	fmt.Println("pkg.NoDocExampleFunc, Hello go-extractor", s.V())
}

// OneLineDocExampleFunc this is one-line-doc example function
func OneLineDocExampleFunc(s *module.ExampleStruct) {
	fmt.Println("pkg.OneLineDocExampleFunc, Hello go-extractor", s.V())
}

func ImportSelectorFunc(s *module.ExampleStruct) {
	fmt.Println("pkg.ImportSelectorFunc, Hello go-extractor", module.NewExampleStruct(s.V()).Sub().ParentStruct.P)
}
