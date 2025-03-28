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

func MapEntries[M ~map[K]V, K comparable, V any](
	inputMap M,
	transform func(K, V) (K, V, bool),
) M {
	result := make(map[K]V)
	for key, value := range inputMap {
		newKey, newValue, include := transform(key, value)
		if include {
			result[newKey] = newValue
		}
	}
	return result
}

func Unzip[T1 any, T2 any](pairList []*Pair[T1, T2]) ([]T1, []T2) {
	firstElements := make([]T1, 0, len(pairList))
	secondElements := make([]T2, 0, len(pairList))
	for _, pair := range pairList {
		firstElements = append(firstElements, pair.First)
		secondElements = append(secondElements, pair.Second)
	}
	return firstElements, secondElements
}

func Zip[T1 any, T2 any](firstSlice []T1, secondSlice []T2) []*Pair[T1, T2] {
	minLen := min(len(firstSlice), len(secondSlice))
	if minLen == 0 {
		return nil
	}
	result := make([]*Pair[T1, T2], 0, minLen)
	for index := 0; index < minLen; index++ {
		result = append(result, &Pair[T1, T2]{
			First:  firstSlice[index],
			Second: secondSlice[index],
		})
	}
	return result
}

func Windowed[T any](input []T, size, step int) [][]T {
	if len(input) == 0 {
		return nil
	}
	if size <= 0 || step <= 0 {
		panic("Windowed: size and step must be positive")
	}
	result := make([][]T, 0, (len(input)+step-1)/step)
	for index := 0; index < len(input); index += step {
		end := index + size
		if end > len(input) {
			end = len(input)
		}
		window := make([]T, end-index)
		copy(window, input[index:end])
		result = append(result, window)
	}
	return result
}

type Pair[F, S any] struct {
	First  F
	Second S
}
