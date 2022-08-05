package main

import (
	"standardProject/pkg"
	"standardProject/pkg/module"
)

func main() {
	pkg.ExampleFunc(module.NewExampleStruct(10))
}
