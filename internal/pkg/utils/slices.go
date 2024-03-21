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

func Map[T any, U any](input []T, mapFunc func(*T) *U) []U {
	result := make([]U, 0, len(input))
	for _, item := range input {
		result = append(result, *mapFunc(&item))
	}
	return result
}

func Zip[T any, U any, V any](first []T, second []U, mapFunc func(T, U) V) []V {
	size := min(len(first), len(second))
	result := make([]V, 0, size)

	for i := range size {
		result = append(result, mapFunc(first[i], second[i]))
	}
	return result
}
