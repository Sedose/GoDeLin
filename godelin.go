package godelin

import "fmt"

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

func GetOrPut[K comparable, V any](table map[K]V, key K, computeValue func() V) V {
	if value, exists := table[key]; exists {
		return value
	}
	table[key] = computeValue()
	return table[key]
}

func GroupBy[T any, K comparable, V any](inputSlice []T, transform func(T) (K, V)) map[K][]V {
	resultMap := make(map[K][]V, len(inputSlice))
	for _, element := range inputSlice {
		key, value := transform(element)
		resultMap[key] = append(resultMap[key], value)
	}
	return resultMap
}

func Chunked[T any](inputSlice []T, chunkSize int) [][]T {
	if chunkSize <= 0 {
		panic("chunkSize must be positive")
	}
	totalItems := len(inputSlice)
	result := make([][]T, 0, (totalItems+chunkSize-1)/chunkSize)
	for start := 0; start < totalItems; start += chunkSize {
		end := start + chunkSize
		if end > totalItems {
			end = totalItems
		}
		result = append(result, inputSlice[start:end])
	}
	return result
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

func Drop[T any](inputSlice []T, numToDrop int) []T {
	if numToDrop <= 0 || len(inputSlice) == 0 {
		return inputSlice
	}
	if numToDrop >= len(inputSlice) {
		return nil
	}
	return inputSlice[numToDrop:]
}

func DropLast[T any](inputSlice []T, numToDrop int) []T {
	if numToDrop <= 0 || len(inputSlice) == 0 {
		return inputSlice
	}
	if numToDrop >= len(inputSlice) {
		return nil
	}
	return inputSlice[:len(inputSlice)-numToDrop]
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
	for i, element := range inputSlice {
		if !predicate(element) {
			return inputSlice[i:]
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

func FlatMapIndexed[T1, T2 any](s []T1, fn func(int, T1) []T2) []T2 {
	var ret []T2
	for i, e := range s {
		ret = append(ret, fn(i, e)...)
	}
	return ret
}

func Fold[T, R any](s []T, initial R, fn func(R, T) R) R {
	acc := initial
	for _, e := range s {
		acc = fn(acc, e)
	}
	return acc
}

func FoldIndexed[T, R any](s []T, initial R, fn func(int, R, T) R) R {
	acc := initial
	for i, e := range s {
		acc = fn(i, acc, e)
	}
	return acc
}

func FoldItems[M ~map[K]V, K comparable, V, R any](
	m M,
	initial R,
	fn func(R, K, V) R,
) R {
	acc := initial
	for k, v := range m {
		acc = fn(acc, k, v)
	}
	return acc
}

func GetOrInsert[M ~map[K]V, K comparable, V any](m M, k K, fn func(K) V) V {
	v, ok := m[k]
	if ok {
		// present, return existing value
		return v
	}
	// not present; get value, insert in map and return the new value
	v = fn(k)
	m[k] = v
	return v
}

func Items[M ~map[K]V, K comparable, V any](m M) []*Pair[K, V] {
	ret := make([]*Pair[K, V], 0, len(m))
	for k, v := range m {
		ret = append(ret, &Pair[K, V]{k, v})
	}
	return ret
}

func Map[T1, T2 any](s []T1, fn func(T1) T2) []T2 {
	ret := make([]T2, 0, len(s))
	for _, e := range s {
		ret = append(ret, fn(e))
	}
	return ret
}

func MapIndexed[T1, T2 any](s []T1, fn func(int, T1) T2) []T2 {
	ret := make([]T2, 0, len(s))
	for i, e := range s {
		ret = append(ret, fn(i, e))
	}
	return ret
}

func Partition[T any](s []T, predicate func(T) bool) ([]T, []T) {
	first := make([]T, 0)
	second := make([]T, 0)
	for _, e := range s {
		if predicate(e) {
			first = append(first, e)
		} else {
			second = append(second, e)
		}
	}
	return first, second
}

func Reduce[T any](s []T, fn func(T, T) T) T {
	if len(s) == 1 {
		return s[0]
	}
	return Fold(s[1:], s[0], fn)
}

func ReduceIndexed[T any](s []T, fn func(int, T, T) T) T {
	if len(s) == 1 {
		return s[0]
	}
	acc := s[0]
	for i, e := range s[1:] {
		acc = fn(i+1, acc, e)
	}
	return acc
}

func Reverse[T any](s []T) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func Reversed[T any](s []T) []T {
	ret := make([]T, 0, len(s))
	for i := len(s) - 1; i >= 0; i-- {
		ret = append(ret, s[i])
	}
	return ret
}

func Take[T any](s []T, n int) []T {
	if len(s) <= n {
		return s
	}
	return s[:n]
}

func TakeLast[T any](s []T, n int) []T {
	if len(s) <= n {
		return s
	}
	return s[len(s)-n:]
}

func TakeLastWhile[T any](s []T, fn func(T) bool) []T {
	if len(s) == 0 {
		return s
	}
	i := len(s) - 1
	for ; i >= 0; i-- {
		if !fn(s[i]) {
			break
		}
	}
	return s[i+1:]
}

func TakeWhile[T any](s []T, fn func(T) bool) []T {
	if len(s) == 0 {
		return s
	}
	i := 0
	for ; i < len(s); i++ {
		if !fn(s[i]) {
			break
		}
	}
	return s[:i]
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
		s1 = append(s1, p.Fst)
		s2 = append(s2, p.Snd)
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
			Fst: s1[i],
			Snd: s2[i],
		})
	}
	return ret
}

type Pair[T1, T2 any] struct {
	Fst T1
	Snd T2
}

func (p Pair[T1, T2]) String() string {
	return fmt.Sprintf("(%v, %v)", p.Fst, p.Snd)
}
