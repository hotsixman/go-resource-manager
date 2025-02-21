package util

func CloneSlice[T any](original []T) []T {
	clone := make([]T, len(original))
	copy(clone, original)
	return clone
}
