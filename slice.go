package godelin

import "slices"

type Slice[T any] []T

func (s Slice[T]) All(predicate func(T) bool) bool {
	return slices.IndexFunc(s, func(e T) bool { return !predicate(e) }) == -1
}

func (s Slice[T]) Any(predicate func(T) bool) bool {
	return slices.IndexFunc(s, predicate) != -1
}

func Associate[T any, K comparable, V any](s Slice[T], transform func(T) (K, V)) map[K]V {
	m := make(map[K]V, len(s))
	for _, e := range s {
		k, v := transform(e)
		m[k] = v
	}
	return m
}
