package godelin

type Slice[T any] []T

func (s Slice[T]) All(predicate func(T) bool) bool {
	for _, e := range s {
		if !predicate(e) {
			return false
		}
	}
	return true
}
