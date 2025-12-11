package registry

func SetIfNotZero[T comparable](target *T, source T) {
	var zero T
	if source != zero {
		*target = source
	}
}
