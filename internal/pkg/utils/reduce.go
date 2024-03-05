package utils

import "atomicgo.dev/constraints"

func Accumulate[T constraints.Addable](slice []T) T {
	return Reduce(slice, func(first, second T) T {
		return first + second
	})
}

func Product[T constraints.Numeric](slice []T) T {
	return Reduce(slice, func(first, second T) T {
		return first * second
	})
}

func Reduce[T any](slice []T, reducer func(T, T) T) T {
	var result T
	for _, el := range slice {
		result = reducer(result, el)
	}
	return result
}
