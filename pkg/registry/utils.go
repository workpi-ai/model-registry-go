package registry

func CopyPtr[T any](ptr *T) *T {
	if ptr == nil {
		return nil
	}
	v := *ptr
	return &v
}

func SetIfNotZero[T comparable](target *T, source T) {
	var zero T
	if source != zero {
		*target = source
	}
}
