package godelin

type Pair[F, S any] struct {
	First  F
	Second S
}

func All[T any](slice []T, predicate func(T) bool) bool {
	for _, element := range slice {
		if !predicate(element) {
			return false
		}
	}
	return true
}

func Any[T any](slice []T, predicate func(T) bool) bool {
	for _, element := range slice {
		if predicate(element) {
			return true
		}
	}
	return false
}

func GetOrPut[M ~map[K]V, K comparable, V any](m M, key K, defaultValue func(K) V) V {
	if value, exists := m[key]; exists {
		return value
	}
	newValue := defaultValue(key)
	m[key] = newValue
	return newValue
}

func GroupBy[T any, K comparable, V any](slice []T, transform func(T) (K, V)) map[K][]V {
	result := make(map[K][]V, len(slice))
	for _, element := range slice {
		key, value := transform(element)
		result[key] = append(result[key], value)
	}
	return result
}

func ChunkedBy[T any](slice []T, groupingFn func(T, T) bool) [][]T {
	if len(slice) == 0 {
		return [][]T{} // return an empty slice, not nil
	}
	estimated := len(slice) / 2
	result := make([][]T, 0, estimated)
	currentChunk := make([]T, 0, len(slice))
	currentChunk = append(currentChunk, slice[0])
	for i := 1; i < len(slice); i++ {
		prev := currentChunk[len(currentChunk)-1]
		curr := slice[i]
		if groupingFn(prev, curr) {
			currentChunk = append(currentChunk, curr)
		} else {
			result = append(result, currentChunk)
			currentChunk = []T{curr}
		}
	}
	result = append(result, currentChunk)
	return result
}

func Distinct[T comparable](slice []T) []T {
	if len(slice) == 0 {
		return []T{}
	}
	seen := make(map[T]struct{}, len(slice))
	result := make([]T, 0, len(slice))
	for _, element := range slice {
		if _, exists := seen[element]; !exists {
			seen[element] = struct{}{}
			result = append(result, element)
		}
	}
	return result
}

func DistinctBy[T any, K comparable](slice []T, keySelector func(T) K) []T {
	if len(slice) == 0 {
		return []T{}
	}
	seen := make(map[K]struct{}, len(slice))
	result := make([]T, 0, len(slice))
	for _, element := range slice {
		key := keySelector(element)
		if _, exists := seen[key]; !exists {
			seen[key] = struct{}{}
			result = append(result, element)
		}
	}
	return result
}

func DropLastWhile[T any](slice []T, predicate func(T) bool) []T {
	for i := len(slice) - 1; i >= 0; i-- {
		if !predicate(slice[i]) {
			return slice[:i+1]
		}
	}
	return []T{}
}

func DropWhile[T any](slice []T, predicate func(T) bool) []T {
	for i, element := range slice {
		if !predicate(element) {
			return slice[i:]
		}
	}
	return []T{}
}

func Filter[T any](slice []T, predicate func(T) bool) []T {
	result := make([]T, 0, len(slice))
	for _, element := range slice {
		if predicate(element) {
			result = append(result, element)
		}
	}
	return result
}

func FilterIndexed[T any](slice []T, predicate func(int, T) bool) []T {
	result := make([]T, 0, len(slice))
	for i, element := range slice {
		if predicate(i, element) {
			result = append(result, element)
		}
	}
	return result
}

func FlatMap[T1, T2 any](slice []T1, transform func(T1) []T2) []T2 {
	result := make([]T2, 0, len(slice))
	for _, element := range slice {
		result = append(result, transform(element)...)
	}
	return result
}

func FlatMapIndexed[T1, T2 any](slice []T1, transform func(int, T1) []T2) []T2 {
	result := make([]T2, 0, len(slice))
	for i, element := range slice {
		result = append(result, transform(i, element)...)
	}
	return result
}

func Fold[T, R any](slice []T, initial R, combine func(R, T) R) R {
	acc := initial
	for _, element := range slice {
		acc = combine(acc, element)
	}
	return acc
}

func FoldIndexed[T, R any](slice []T, initial R, combine func(int, R, T) R) R {
	acc := initial
	for i, element := range slice {
		acc = combine(i, acc, element)
	}
	return acc
}

func FoldMapEntries[M ~map[K]V, K comparable, V, R any](
	m M,
	initial R,
	combine func(R, K, V) R,
) R {
	acc := initial
	for key, value := range m {
		acc = combine(acc, key, value)
	}
	return acc
}

func Items[M ~map[K]V, K comparable, V any](m M) []Pair[K, V] {
	if len(m) == 0 {
		return []Pair[K, V]{}
	}
	pairs := make([]Pair[K, V], 0, len(m))
	for key, value := range m {
		pairs = append(pairs, Pair[K, V]{First: key, Second: value})
	}
	return pairs
}

func Map[T, R any](slice []T, transform func(T) R) []R {
	if len(slice) == 0 {
		return []R{}
	}
	result := make([]R, 0, len(slice))
	for _, element := range slice {
		result = append(result, transform(element))
	}
	return result
}

func MapIndexed[T, R any](slice []T, transform func(int, T) R) []R {
	if len(slice) == 0 {
		return []R{}
	}
	result := make([]R, 0, len(slice))
	for i, element := range slice {
		result = append(result, transform(i, element))
	}
	return result
}

func Partition[T any](slice []T, predicate func(T) bool) ([]T, []T) {
	matching := make([]T, 0, len(slice))
	others := make([]T, 0, len(slice))
	for _, element := range slice {
		if predicate(element) {
			matching = append(matching, element)
		} else {
			others = append(others, element)
		}
	}
	return matching, others
}

func Reduce[T any](slice []T, combine func(T, T) T) T {
	if len(slice) == 0 {
		panic("Reduce called on empty slice")
	}
	if len(slice) == 1 {
		return slice[0]
	}
	return Fold(slice[1:], slice[0], combine)
}

func ReduceIndexed[T any](slice []T, combine func(int, T, T) T) T {
	if len(slice) == 0 {
		panic("ReduceIndexed called on empty slice")
	}
	if len(slice) == 1 {
		return slice[0]
	}
	return FoldIndexed(slice[1:], slice[0], func(i int, acc T, element T) T {
		return combine(i+1, acc, element)
	})
}

func TakeLastWhile[T any](slice []T, predicate func(T) bool) []T {
	if len(slice) == 0 {
		return []T{}
	}
	idx := len(slice) - 1
	for ; idx >= 0; idx-- {
		if !predicate(slice[idx]) {
			break
		}
	}
	return slice[idx+1:]
}

func TakeWhile[T any](slice []T, predicate func(T) bool) []T {
	if len(slice) == 0 {
		return []T{}
	}
	var i int
	for ; i < len(slice); i++ {
		if !predicate(slice[i]) {
			break
		}
	}
	return slice[:i]
}

func MapEntries[M ~map[K]V, K comparable, V any](
	m M,
	transform func(K, V) (K, V, bool),
) M {
	result := make(M)
	for key, value := range m {
		newKey, newValue, keep := transform(key, value)
		if keep {
			result[newKey] = newValue
		}
	}
	return result
}

func Unzip[T1, T2 any](pairs []Pair[T1, T2]) ([]T1, []T2) {
	firsts := make([]T1, 0, len(pairs))
	seconds := make([]T2, 0, len(pairs))
	for _, p := range pairs {
		firsts = append(firsts, p.First)
		seconds = append(seconds, p.Second)
	}
	return firsts, seconds
}

func Zip[T1, T2 any](first []T1, second []T2) []Pair[T1, T2] {
	minLen := min(len(first), len(second))
	if minLen == 0 {
		return []Pair[T1, T2]{}
	}
	result := make([]Pair[T1, T2], 0, minLen)
	for i := 0; i < minLen; i++ {
		result = append(result, Pair[T1, T2]{
			First:  first[i],
			Second: second[i],
		})
	}
	return result
}

func Windowed[T any](slice []T, size, step int) [][]T {
	if len(slice) == 0 {
		return [][]T{}
	}
	if size <= 0 || step <= 0 {
		panic("Windowed: size and step must be positive")
	}
	result := make([][]T, 0, (len(slice)+step-1)/step)
	for i := 0; i < len(slice); i += step {
		end := i + size
		if end > len(slice) {
			end = len(slice)
		}
		window := make([]T, 0, end-i)
		window = append(window, slice[i:end]...)
		result = append(result, window)
	}
	return result
}
