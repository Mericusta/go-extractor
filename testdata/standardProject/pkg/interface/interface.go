package pkgInterface

type ExampleInterface interface {
	// This is ExampleFunc Doc
	ExampleFunc(int)
	// This is AnotherExampleFunc Doc
	AnotherExampleFunc(int, []int) (int, []int)
}

type ExampleTemplateInterface[T any] interface {
	// This is ExampleFunc Doc
	ExampleFunc(T)
	// This is AnotherExampleFunc Doc
	AnotherExampleFunc(T, []T) (T, []T)
}
