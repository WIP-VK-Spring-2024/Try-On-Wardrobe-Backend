package utils

import (
	"atomicgo.dev/constraints"
)

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

func Count[T any, S ~[]T](slice S, predicate func(T) bool) int {
	count := 0
	for _, item := range slice {
		if predicate(item) {
			count += 1
		}
	}
	return count
}

func Filter[T any, S ~[]T](slice S, predicate func(T) bool) S {
	result := make(S, 0, len(slice))
	for _, elem := range slice {
		if predicate(elem) {
			result = append(result, elem)
		}
	}
	return result
}

func Every[T any, S ~[]T](slice S, predicate func(T) bool) bool {
	for _, elem := range slice {
		if predicate(elem) {
			return false
		}
	}
	return true
}
