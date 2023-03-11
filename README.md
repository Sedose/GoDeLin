# fun
[![GoDoc](https://godoc.org/github.com/luraim/fun?status.svg)](https://godoc.org/github.com/luraim/fun)

### Simple generic utility functions to reduce golang boilerplate
- Inspired by Kotlin and Rust collection functions
- Supplement to the generic functions in golang.org/x/exp/slices and golang.org/x/exp/maps
- Note: The Go compiler does not currently inline generic callback functions. So please use your judgement while using functions from this library that involve callbacks. Use them when the expressiveness is worth any performance degration compared to handcoded *for loop* boilerplate.

## List of functions
- [fun](#fun)
    - [Simple generic utility functions to reduce golang boilerplate](#simple-generic-utility-functions-to-reduce-golang-boilerplate)
  - [List of functions](#list-of-functions)
    - [All](#all)
    - [Any](#any)
    - [AppendToGroup](#appendtogroup)
    - [Associate](#associate)
    - [Chunked](#chunked)
    - [ChunkedBy](#chunkedby)
    - [Distinct](#distinct)
    - [DistinctBy](#distinctby)
    - [Drop](#drop)
    - [DropLast](#droplast)
    - [DropWhile](#dropwhile)
    - [DropLastWhile](#droplastwhile)
    - [Filter](#filter)
    - [FilterIndexed](#filterindexed)
    - [FilterMap](#filtermap)
    - [Fold](#fold)
    - [FoldIndexed](#foldindexed)
    - [FoldItems](#folditems)
    - [GetOrInsert](#getorinsert)
    - [GroupBy](#groupby)
    - [Items](#items)
    - [Map](#map)
    - [MapIndexed](#mapindexed)
    - [Partition](#partition)
    - [Reduce](#reduce)
    - [ReduceIndexed](#reduceindexed)
    - [Reverse](#reverse)
    - [Reversed](#reversed)
    - [Take](#take)
    - [TakeLast](#takelast)
    - [TakeWhile](#takewhile)
    - [TakeLastWhile](#takelastwhile)
    - [TransformMap](#transformmap)
    - [Unzip](#unzip)
    - [Windowed](#windowed)
    - [Zip](#zip)

### All
- Returns true if all elements return true for given predicate
```go
All([]int{1, 2, 3, 4, 5}, func(i int)bool {return i < 7})
// true

All([]int{1, 2, 3, 4, 5}, func(i int)bool {return i % 2 == 0})
// false

```

### Any
- Returns true if at least one element returns true for given predicate
```go
Any([]int{1, 2, 3}, func(i int)bool {return i%2==0})
// true

Any([]int{1, 2, 3}, func(i int)bool {return i > 7})
// false
```

### AppendToGroup
- Adds the key, value to the given map where each key maps to a slice of values
```go
group := make(map[string][]int)

AppendToGroup(grp, "a", 1)
AppendToGroup(grp, "b", 2)
AppendToGroup(grp, "a", 10)
AppendToGroup(grp, "b", 20)
AppendToGroup(grp, "a", 100)
AppendToGroup(grp, "b", 200)

// {"a":[1, 10, 100], "b":[2, 20, 200]}
```

### Associate
- Returns a map containing key-value pairs returned by the given function applied to the elements of the given slice
```go
Associate([]int{1, 2, 3, 4}, func(i int) (string, int) {
    return fmt.Sprintf("M%d", i), i * 10
})
// {"M1": 10, "M2": 20, "M3": 30, "M4": 40}
```

### Chunked
- Splits the slice into a slice of slices, each not exceeding given chunk size
- The last slice might have fewer elements than the given chunk size
```go
Chunked([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}, 2)
// [[1, 2], [3, 4], [5, 6], [7, 8], [9]]
```

### ChunkedBy
- Splits the slice into a slice of slices, starting a new sub slice whenever the callback function returns false.
- The callback function is passed the previous and current element. 
```go
input := []int{10, 20, 30, 40, 31, 31, 33, 34, 21, 22, 23, 24, 11, 12, 13, 14}
ChunkedBy(input, func(prev, next int) bool { return prev < next})
// [[10, 20, 30, 40], [31], [31, 33, 34], [21, 22, 23, 24], [11, 12, 13, 14]]
```

### Distinct
- Returns a slice containing only distinct elements from the given slice
```go
Distinct([]int{1, 1, 2, 3, 3, 4, 4, 4, 4, 5, 5, 5})
// [1, 2, 3, 4, 5]
```

### DistinctBy
- Returns a slice containing only distinct elements from the given slice as distinguished by the given selector function
```go
DistinctBy([]string{"a", "A", "b", "B", "c", "C"},func(s string) string {
	return strings.ToLower(s)
})
// ["a", "b", "c"]
```

### Drop
- Returns a slice containing all elements except the first n.
```go
// letters = ['a'..'z']
Drop(letters, 23)
// ['x', 'y', 'z']
```

### DropLast
- Returns a slice containing all elements except the last n.
```go
// letters = ['a'..'z']
DropLast(letters, 23)
// ['a', 'b', 'c']
```

### DropWhile
- Returns a slice containing all elements except the first elements that satisfy the given predicate.
```go
// letters = ['a'..'z']
DropWhile(letters, func(r rune) bool { return r < 'x' })
// ['x', 'y', 'z']
```

### DropLastWhile
- Returns a slice containing all elements except the last elements that satisfy the given predicate.
```go
// letters = ['a'..'z']
DropLastWhile(letters, func(r rune) bool { return r > 'c' })
// ['a', 'b', 'c']
```

### Filter
- Returns the slice obtained after retaining only those elements in the given slice for which the given function returns true
```go
Filter([]int{1, 2, 3, 4, 5, 6, 7, 8}, func(i int)bool {return i%2==0})
// [2, 4, 6, 8]
```

### FilterIndexed
- Returns the slice obtained after retaining only those elements in the given slice for which the given function returns true
- Predicate function receives the value as well as its index in the slice.
```go
FilterIndexed([]int{0, 1, 2, 3, 4, 8, 6}, func(index int, v int) bool {
	return index == v
})
// [0, 1, 2, 3, 4, 6]
```

### FilterMap
- FilterMap returns the slice obtained after both filtering and mapping using the given function.
- The function should return two values - the result of the mapping operation and whether the element should be included or dropped.
- This is faster than doing separate filter and map operations, since it avoids extra allocations and slice traversals.
- Inspired by std::iter::filter_map in Rust
```go
FilterMap([]int{1, 2, 3, 4, 5},
    func(i int) (int, bool) {
        if i%2 != 0 {
            return i, false // drop odd numbers
        }
        return i * i, true // square even numbers
    })
// [4, 16]
```

### Fold
- Accumulates values starting with given initial value and applying given function to current accumulator and each element of the given slice.
```go
Fold([]int{1, 2, 3, 4, 5}, func(acc, v int) int { return acc + v })
// 15
```

### FoldIndexed
- Accumulates values starting with given initial value and applying given function to current accumulator and each element of the given slice.
- Function also receives index of current element.
```go
FoldIndexed([]int{1, 2, 3, 4, 5}, func(index, acc, v int) int {
	return acc + index*v
})
// 40
```

### FoldItems
- Accumulates values starting with given intial value and applying given function to current accumulator and each key, value of the given map.
- Accumulator can be of any reference type.
```go
m := map[int]int{1: 10, 2: 20, 3: 30}
FoldItems(m, func(acc map[string]string, k, v int) map[string]string {
    acc[fmt.Sprintf("entry_%d", k)] = fmt.Sprintf("%d->%d", k, v)
    return acc
})
// {"entry_1": "1->10", "entry_2": "2->20", "entry_3": "3->30"}
```

### GetOrInsert 
- checks if a value corresponding to the given key is present in the map. 
- If present it returns the existing value. 
- If not present, it invokes the given callback function to get a new value for the given key, inserts it in the map and returns the new value
```go
m := map[int]int{1:10, 2:20}
GetOrInsert(m, 3, func(i int) int {return i * 10})
// returns 30; m is updated to {1:10, 2:20, 3:30},
```  

### GroupBy
- Returns a map where each key maps to slices of elements all having the same key as returned by the given function
```go
GroupBy([]string{"a", "abc", "ab", "def", "abcd"}, func(s string) (int,string) {
	return len(s), s
})
// {1: ["a"], 2: ["ab"], 3: ["abc", "def"], 4: ["abcd"]},
```

### Items
- Returns the (key, value) pairs of the given map as a slice
```go
// m := map[string][]int{"a": {1, 2, 3, 4}, "b": {1, 2}, "c": {1, 2, 3}}
Items(m)
// []*Pair[string, []int]{
//      {"a", []int{1, 2, 3, 4}},
//      {"b", []int{1, 2}},
//      {"c", []int{1, 2, 3}},
// }
```

### Map
- Returns the slice obtained after applying the given function over every element in the given slice
```go
Map([]int{1, 2, 3, 4, 5}, func(i int) int { return i * i })
// [1, 4, 9, 16, 25]
```

### MapIndexed
- Returns the slice obtained after applying the given function over every element in the given slice
- The function also receives the index of each element in the slice.
```go
MapIndexed([]int{1, 2, 3, 4, 5}, func(index, i int) int { return index * i })
// [0, 2, 6, 12, 20]
```

### Partition
- Returns two slices where the first slice contains elements for which the predicate returned true and the second slice contains elements for which it returned false.
```go
type person struct {
    name string
    age  int
}

tom := &person{"Tom", 18}
andy := &person{"Andy", 32}
sarah := &person{"Sarah", 22}

Partition([]*person{tom, andy, sarah}, func(p *person) bool { return p.age < 30 })
// [tom, sarah], [andy]
```

### Reduce
- Accumulates the values starting with the first element and applying the operation from left to right to the current accumulator value and each element.
- The input slice must have at least one element.
```go
Reduce([]int{1, 2, 3, 4, 5}, func(acc, v int) int { return acc + v })
// 15
```

### ReduceIndexed
- Accumulates the values starting with the first element and applying the operation from left to right to the current accumulator value and each element.
- The input slice must have at least one element.
- The function also receives the index of each element.
```go
ReduceIndexed([]string{"a", "b", "c", "d"}, func(index int, acc, v string) string {
    return fmt.Sprintf("%s%s%d", acc, v, index)
})
// "ab1c2d3"
```

### Reverse
- Reverses the elements of the list in place.
```go
// s = [1, 2, 3, 4, 5, 6, 7]
Reverse(s)
// s = [7, 6, 5, 4, 3, 2, 1]
```

### Reversed
- Returns a new list with the elements in reverse order.
```go
// s = [1, 2, 3, 4, 5, 6, 7]
r := Reversed(s)
// r = [7, 6, 5, 4, 3, 2, 1]
// s = [1, 2, 3, 4, 5, 6, 7]
```

### Take
- Returns the slice obtained after taking the first n elements from the given slice.
```go
// letters = ['a'..'z']
Take(letters, 2)
// ['a', 'b']
```

### TakeLast
- Returns the slice obtained after taking the last n elements from the given slice.
```go
// letters = ['a'..'z']
TakeLast(letters, 2)
// ['y', 'z']
```

### TakeWhile
- Returns a slice containing the first elements satisfying the given predicate
```go
// letters = ['a'..'z']
TakeWhile(letters,  func(s rune) bool { return s < 'f' })
// ['a', 'b', 'c', 'd', 'e']
```

### TakeLastWhile
- Returns a slice containing the last elements satisfying the given predicate
```go
// letters = ['a'..'z']
TakeLastWhile(letters, func(s rune) bool { return s > 'w' })
// ['x', 'y', 'z']
```

### TransformMap
- Applies the given function to each key, value in the map, and returns a new map of the same type after transforming the keys and values depending on the callback functions return values. 
- If the last bool return value from the callback function is false, the entry is dropped
```go
// filtering a map
// m = {"a":[1, 2, 3, 4] "b":[1, 2] "c":[1, 2, 3]}
TransformMap(m, 
    func(k string, v []int) (string, []int, bool) {
        if len(v) < 3 {
            return k, v, false
        }
        return k, v, true
    })
// drops all values with length less than 3
// {"a":[1, 2, 3, 4]  "c":[1, 2, 3]}


// transforming keys and values
// m = {"a":[1, 2, 3, 4] "b":[5, 6]}
TransformMap(m, 
    func(k string, v []int) (string, []int, bool) {
        newK := strings.ToUpper(k)
        newV := Map(v, func(i int) int { return i * 10 })
        return newK, newV, true
    })
// {"A":[10, 20, 30, 40]  "B":[50, 60]}
```


### Unzip
- Returns two slices, where:
- the first slice is built from the first values of each pair from the input slice
- the second slice is built from the second values of each pair
```go
Unzip([]*Pair[string, int]{{"a", 1}, {"b", 2}, {"c", 3}})
// ["a", "b", "c"], [1, 2, 3]
```

### Windowed
- Returns a slice of sliding windows, each of the given size, and with the given step
- Several last slices may have fewer elements than the given size
```go
Windowed([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 5, 1)
// [
//     [1, 2, 3, 4, 5],
//     [2, 3, 4, 5, 6],
//     [3, 4, 5, 6, 7],
//     [4, 5, 6, 7, 8],
//     [5, 6, 7, 8, 9],
//     [6, 7, 8, 9, 10],
//     [7, 8, 9, 10],
//     [8, 9, 10],
//     [9, 10],
//     [10]
// ]

Windowed([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 5, 3)
// [
//     [1, 2, 3, 4, 5],
//     [4, 5, 6, 7, 8],
//     [7, 8, 9, 10],
//     [10]
// ]
```

### Zip
- Returns a slice of pairs from the elements of both slices with the same index
- The returned slice has the length of the shortest input slice
```go
Zip([]string{"a", "b", "c", "d"}, []int{1, 2, 3})
// []*Pair[string, int]{{"a", 1}, {"b", 2}, {"c", 3}}
```