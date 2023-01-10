package pkg

import (
	"fmt"
	"standardProject/pkg/module"
)

// ExampleFunc this is example function
func ExampleFunc(s *module.ExampleStruct) {
	fmt.Println("pkg.ExampleFunc, Hello go-extractor,", s.V()) // `ast:init,default=bool ast:assign`
}
