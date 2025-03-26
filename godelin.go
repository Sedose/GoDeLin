package godelin

func All[T any](inputSlice []T, predicate func(T) bool) bool {
	for _, element := range inputSlice {
		if !predicate(element) {
			return false
		}
	}
	return true
}

func Any[T any](inputSlice []T, predicate func(T) bool) bool {
	for _, element := range inputSlice {
		if predicate(element) {
			return true
		}
	}
	return false
}

func GetOrPut[M ~map[K]V, K comparable, V any](theMap M, key K, defaultValue func(K) V) V {
	if value, exists := theMap[key]; exists {
		return value
	}
	newValue := defaultValue(key)
	theMap[key] = newValue
	return newValue
}

func GroupBy[T any, K comparable, V any](inputSlice []T, transform func(T) (K, V)) map[K][]V {
	resultMap := make(map[K][]V, len(inputSlice))
	for _, element := range inputSlice {
		key, value := transform(element)
		resultMap[key] = append(resultMap[key], value)
	}
	return resultMap
}

func ChunkedBy[T any](inputSlice []T, groupingFn func(T, T) bool) [][]T {
	if len(inputSlice) == 0 {
		return nil
	}
	estimatedCapacity := len(inputSlice) / 2
	resultChunks := make([][]T, 0, estimatedCapacity)
	currentChunk := make([]T, 0, len(inputSlice))
	currentChunk = append(currentChunk, inputSlice[0])
	for _, currentElement := range inputSlice[1:] {
		lastElement := currentChunk[len(currentChunk)-1]
		if groupingFn(lastElement, currentElement) {
			currentChunk = append(currentChunk, currentElement)
		} else {
			resultChunks = append(resultChunks, currentChunk)
			currentChunk = []T{currentElement}
		}
	}
	resultChunks = append(resultChunks, currentChunk)
	return resultChunks
}

func Distinct[T comparable](inputSlice []T) []T {
	if len(inputSlice) < 1 {
		return nil
	}
	seen := make(map[T]struct{}, len(inputSlice))
	result := make([]T, 0, len(inputSlice))

	for _, element := range inputSlice {
		if _, exists := seen[element]; exists {
			continue
		}
		seen[element] = struct{}{}
		result = append(result, element)
	}
	return result
}

func DistinctBy[T any, K comparable](inputSlice []T, keySelector func(T) K) []T {
	if len(inputSlice) == 0 {
		return nil
	}

	seenKeys := make(map[K]struct{}, len(inputSlice))
	result := make([]T, 0, len(inputSlice))

	for _, element := range inputSlice {
		key := keySelector(element)
		if _, exists := seenKeys[key]; !exists {
			seenKeys[key] = struct{}{}
			result = append(result, element)
		}
	}

	return result
}

func DropLastWhile[T any](inputSlice []T, predicate func(T) bool) []T {
	for i := len(inputSlice) - 1; i >= 0; i-- {
		if !predicate(inputSlice[i]) {
			return inputSlice[:i+1]
		}
	}
	return []T{}
}

func DropWhile[T any](inputSlice []T, predicate func(T) bool) []T {
	for index, element := range inputSlice {
		if !predicate(element) {
			return inputSlice[index:]
		}
	}
	return []T{}
}

func Filter[T any](inputSlice []T, predicate func(T) bool) []T {
	result := make([]T, 0, len(inputSlice))
	for _, element := range inputSlice {
		if predicate(element) {
			result = append(result, element)
		}
	}
	return result
}

func FilterIndexed[T any](inputSlice []T, predicate func(int, T) bool) []T {
	result := make([]T, 0, len(inputSlice))
	for index, element := range inputSlice {
		if predicate(index, element) {
			result = append(result, element)
		}
	}
	return result
}

func FlatMap[T1, T2 any](inputSlice []T1, transform func(T1) []T2) []T2 {
	result := make([]T2, 0)
	for _, element := range inputSlice {
		result = append(result, transform(element)...)
	}
	return result
}

func FlatMapIndexed[T1, T2 any](inputSlice []T1, transform func(int, T1) []T2) []T2 {
	result := make([]T2, 0, len(inputSlice))
	for index, element := range inputSlice {
		result = append(result, transform(index, element)...)
	}
	return result
}

func Fold[T, R any](inputSlice []T, initial R, operation func(R, T) R) R {
	accumulator := initial
	for _, element := range inputSlice {
		accumulator = operation(accumulator, element)
	}
	return accumulator
}

func FoldIndexed[T, R any](inputSlice []T, initial R, operation func(int, R, T) R) R {
	accumulator := initial
	for index, element := range inputSlice {
		accumulator = operation(index, accumulator, element)
	}
	return accumulator
}

func FoldMapEntries[M ~map[K]V, K comparable, V, R any](
	sourceMap M,
	initialValue R,
	combine func(R, K, V) R,
) R {
	accumulator := initialValue
	for key, value := range sourceMap {
		accumulator = combine(accumulator, key, value)
	}
	return accumulator
}

func Items[M ~map[K]V, K comparable, V any](inputMap M) []*Pair[K, V] {
	if len(inputMap) == 0 {
		return nil
	}
	pairs := make([]*Pair[K, V], 0, len(inputMap))
	for key, value := range inputMap {
		pairs = append(pairs, &Pair[K, V]{First: key, Second: value})
	}
	return pairs
}

func Map[T, R any](inputSlice []T, transform func(T) R) []R {
	if len(inputSlice) == 0 {
		return nil
	}
	mapped := make([]R, 0, len(inputSlice))
	for _, element := range inputSlice {
		mapped = append(mapped, transform(element))
	}
	return mapped
}

func MapIndexed[T, R any](inputSlice []T, transform func(int, T) R) []R {
	if len(inputSlice) == 0 {
		return nil
	}
	mapped := make([]R, 0, len(inputSlice))
	for index, element := range inputSlice {
		mapped = append(mapped, transform(index, element))
	}
	return mapped
}

func Partition[T any](s []T, predicate func(T) bool) ([]T, []T) {
	first := make([]T, 0)
	second := make([]T, 0)
	for _, element := range s {
		if predicate(element) {
			first = append(first, element)
		} else {
			second = append(second, element)
		}
	}
	return first, second
}

func Reduce[T any](inputSlice []T, combine func(T, T) T) T {
	if len(inputSlice) == 0 {
		panic("Reduce called on empty slice")
	}
	if len(inputSlice) == 1 {
		return inputSlice[0]
	}
	return Fold(inputSlice[1:], inputSlice[0], combine)
}

func ReduceIndexed[T any](inputSlice []T, combine func(int, T, T) T) T {
	if len(inputSlice) == 0 {
		panic("ReduceIndexed called on empty slice")
	}
	if len(inputSlice) == 1 {
		return inputSlice[0]
	}
	return FoldIndexed(inputSlice[1:], inputSlice[0],
		func(index int, acc T, element T) T {
			return combine(index+1, acc, element)
		})
}

func TakeLastWhile[T any](inputSlice []T, predicate func(T) bool) []T {
	if len(inputSlice) == 0 {
		return nil
	}
	index := len(inputSlice) - 1
	for ; index >= 0; index-- {
		if !predicate(inputSlice[index]) {
			break
		}
	}
	return inputSlice[index+1:]
}

func TakeWhile[T any](inputSlice []T, predicate func(T) bool) []T {
	if len(inputSlice) == 0 {
		return nil
	}
	var index int
	for ; index < len(inputSlice); index++ {
		if !predicate(inputSlice[index]) {
			break
		}
	}
	return inputSlice[:index]
}

func TransformMap[M ~map[K]V, K comparable, V any](
	m M,
	fn func(k K, v V) (K, V, bool),
) M {
	ret := make(map[K]V)
	for k, v := range m {
		newK, newV, include := fn(k, v)
		if include {
			ret[newK] = newV
		}
	}
	return ret
}

func Unzip[T1 any, T2 any](ps []*Pair[T1, T2]) ([]T1, []T2) {
	l := len(ps)
	s1 := make([]T1, 0, l)
	s2 := make([]T2, 0, l)
	for _, p := range ps {
		s1 = append(s1, p.First)
		s2 = append(s2, p.Second)
	}
	return s1, s2
}

func Windowed[T any](s []T, size, step int) [][]T {
	ret := make([][]T, 0)
	sz := len(s)
	if sz == 0 {
		return ret
	}
	start := 0
	end := 0
	updateEnd := func() {
		e := start + size
		if e >= sz {
			e = sz
		}
		end = e
	}
	updateStart := func() {
		s := start + step
		if s >= sz {
			s = sz
		}
		start = s
	}
	updateEnd()

	for {
		sub := make([]T, 0, end)
		for i := start; i < end; i++ {
			sub = append(sub, s[i])
		}
		ret = append(ret, sub)
		updateStart()
		updateEnd()
		if start == end {
			break
		}
	}
	return ret
}

func Zip[T1 any, T2 any](s1 []T1, s2 []T2) []*Pair[T1, T2] {
	minLen := len(s1)
	if minLen > len(s2) {
		minLen = len(s2)
	}

	ret := make([]*Pair[T1, T2], 0, minLen)

	for i := 0; i < minLen; i++ {
		ret = append(ret, &Pair[T1, T2]{
			First:  s1[i],
			Second: s2[i],
		})
	}
	return ret
}

type Pair[F, S any] struct {
	First  F
	Second S
}
