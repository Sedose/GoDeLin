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

// Chunked splits the slice into chunks of the given size.
func (s Slice[T]) Chunked(size int) []Slice[T] {
	if size <= 0 {
		panic("size must be greater than zero")
	}

	if len(s) == 0 {
		return []Slice[T]{}
	}

	var chunks []Slice[T]
	for i := 0; i < len(s); i += size {
		end := i + size
		if end > len(s) {
			end = len(s)
		}
		chunks = append(chunks, s[i:end])
	}
	return chunks
}
