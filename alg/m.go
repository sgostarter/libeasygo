package alg

import "golang.org/x/exp/constraints"

type MaxCheck[T constraints.Ordered] struct {
}

func (MaxCheck[T]) NewIsTop(oldV, newV T) bool {
	return newV >= oldV
}

type MinCheck[T constraints.Ordered] struct {
}

func (MinCheck[T]) NewIsTop(oldV, newV T) bool {
	return newV <= oldV
}
