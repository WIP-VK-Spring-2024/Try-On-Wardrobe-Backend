package utils

func Map[T any, U any](input []T, mapFunc func(*T) *U) []U {
	result := make([]U, 0, len(input))
	for _, item := range input {
		result = append(result, *mapFunc(&item))
	}
	return result
}
