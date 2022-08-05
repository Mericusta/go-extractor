package module

type ExampleStruct struct {
	V int
}

func NewExampleStruct(v int) *ExampleStruct {
	return &ExampleStruct{V: v}
}
