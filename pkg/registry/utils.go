package registry

func SetIfNotZero[T comparable](target *T, source T) {
	var zero T
	if source != zero {
		*target = source
	}
}

func CopySlice[T any](src []T) []T {
	if len(src) == 0 {
		return nil
	}
	dst := make([]T, len(src))
	copy(dst, src)
	return dst
}
