package pkg

import (
	"fmt"
	"standardProject/pkg/module"
)

func ExampleFunc(s *module.ExampleStruct) {
	fmt.Println("Hello go-extractor,", s.V)
}
