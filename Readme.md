# GoDeLin 🚀

### Bringing Kotlin's Expressiveness to Go

Welcome to GoDelin, the Go library inspired by Kotlin. Tired of writing endless loops and `if err != nil`? GoDelin helps you write expressive, declarative code while maintaining Go's speed and reliability.

📖 **Overview**
GoDelin transforms Go into a more expressive language with Kotlin-inspired utility functions for slices, maps, and collections. With GoDelin, you can:

- Write clear, expressive code using functions like `All`, `Any`, `GroupBy`, `Chunked`, `Distinct`, `Filter`, `Map`, `FlatMap`, and more.
- Replace repetitive loops with declarative operations.

🎯 **Features**
**Kotlin-Style Functions:**
- `All` / `Any`: Check if all or any elements satisfy a predicate.
- `GroupBy`: Turn a slice into a map based on a transformation.
- `Chunked` & `ChunkedBy`: Break your slice into fixed-size or conditionally grouped chunks.
- `Distinct` / `DistinctBy`: Filter out duplicates while keeping the order.
- `Filter`, `Map`, `FlatMap`: Transform and filter collections.
- Additional helpers like `Drop`, `Take`, `Zip`, and `Unzip`.

**Declarative Syntax:**
Chain operations in a functional style.

**Test-Driven Quality:**
Comprehensive test suite ensures reliability (see `slice_test.go`).

**Optimized for Go:**
Functions are crafted for performance using efficient algorithms.

📦 **Installation**
To add GoDelin to your project:

```bash
go get github.com/Sedose/godelin
```

Then import it in your Go files:

```go
import "github.com/Sedose/godelin"
```

🔧 **Usage Examples**

**Grouping Strings by First Letter:**

```go
fruits := []string{"apple", "apricot", "banana", "avocado"}
groups := GroupBy(fruits, func(fruit string) (string, string) {
    return fruit[:1], fruit
})
// groups → map[string][]string{"a": {"apple", "apricot", "avocado"}, "b": {"banana"}}
```

**Chunking a Slice with Custom Logic:**

```go
nums := []int{1, 2, 3, 7, 8, 10, 11, 12, 20}
groups := ChunkedBy(nums, func(prev, curr int) bool {
    return curr == prev+1
})
// groups → [][]int{ {1, 2, 3}, {7, 8}, {10, 11, 12}, {20} }
```

**Filtering Unique Elements:**

```go
input := []int{3, 1, 2, 3, 2, 1}
unique := Distinct(input)
// unique → []int{3, 1, 2}
```

**Transforming and Filtering Collections:**

```go
nums := []int{1, 2, 3, 4, 5, 6}
result := Map(
    Filter(nums, func(n int) bool { return n%2 == 0 }),
    func(n int) int { return n * n },
)
// result → []int{4, 16, 36}
```

🧪 **Testing**
Run tests with:

```bash
go test ./...
```

🚀 **Conclusion**
GoDelin brings Kotlin’s expressiveness to Go. Write cleaner, more maintainable code.

Happy coding! 🚀
