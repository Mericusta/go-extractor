package template

func OneTemplateFunc[T any](tv *T) *T {
	return nil
}

func DoubleSameTemplateFunc[T1, T2 any](tv1 T1, tv2 T2) (*T1, *T2) {
	return nil, nil
}

func DoubleDifferenceTemplateFunc[T1 any, T2 comparable](tv1 T1, tv2 T2) (*T1, *T2) {
	return nil, nil
}

type TypeConstraints interface {
	int8 | int16 | uint8 | uint16
}

func TypeConstraintsTemplateFunc[T TypeConstraints](tv T) *T {
	return nil
}
