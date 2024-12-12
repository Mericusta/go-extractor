package main

import (
	"os"

	pkgInterface "standardProject/pkg/interface"
	"standardProject/pkg/module"
	"standardProject/pkg/template"
)

var (
	globalVariableInt        int    = 1
	globalVariableString     string = os.Getenv("ENV")
	globalVariableStruct     *module.ExampleStruct
	globalVariableTStruct    *template.TemplateStruct[int]
	globalVariableInterface  *pkgInterface.ExampleInterface
	globalVariableTInterface *pkgInterface.ExampleTemplateInterface[int]
)

var anotherGlobalVariableAny interface{}

func Init() {
	_ = globalVariableInt
	_ = globalVariableString
	_ = globalVariableStruct
	_ = globalVariableTStruct
	_ = globalVariableInterface
	_ = globalVariableTInterface
	_ = anotherGlobalVariableAny
}
