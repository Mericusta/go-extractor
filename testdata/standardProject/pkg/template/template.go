package template

import stpmap "github.com/Mericusta/go-stp/map"

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

func CannotInferTypeFunc1[T any]() {

}

func CannotInferTypeFunc2[K comparable, V any]() (*K, *V) {
	return nil, nil
}

type TemplateStruct[T any] struct {
	v T
}

func (t *TemplateStruct[T]) V() T {
	return t.v
}

type TwoTypeTemplateStruct[K TypeConstraints, V any] struct {
	v map[K]V
}

func (t *TwoTypeTemplateStruct[K, V]) KVSlice(k K, v V) ([]K, []V) {
	vs := make([]V, 0, len(t.v))
	for _, v := range t.v {
		vs = append(vs, v)
	}
	return stpmap.Key(t.v), vs
}
