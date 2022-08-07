package main

import (
	"standardProject/pkg"
	"standardProject/pkg/module"
)

func main() {
	pkg.ExampleFunc(module.NewExampleStruct(10))
	module.ExampleFunc(module.NewExampleStruct(11))
}
