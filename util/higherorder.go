package util

// Filter 함수
func Filter[T any](data []T, predicate func(T) bool) []T {
	result := make([]T, 0)
	for _, v := range data {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

// Map 함수
func Transition[T any, R any](data []T, mapper func(T) R) []R {
	result := make([]R, len(data))
	for i, v := range data {
		result[i] = mapper(v)
	}
	return result
}
