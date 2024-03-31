package utils

import (
	"cmp"
	"slices"

	"golang.org/x/exp/maps"
)

func SortedKeysByValue[M ~map[K]V, K comparable, V cmp.Ordered](data M) []K {
	keys := maps.Keys(data)

	slices.SortFunc(keys, func(first, second K) int {
		return -1 * cmp.Compare(data[first], data[second])
	})

	return keys
}
