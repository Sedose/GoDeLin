# GoDeLin 🚀

### Bringing Kotlin's std lib expressiveness to Go

- GoDeLin is a Go library inspired by Kotlin standard library.

- GoDeLin helps you write expressive, declarative Go code.

- The goal is to complement Go, not to provide alternative to things that are effectively covered by the Go standard library.

📦 **Installation**
To add GoDeLin to your project:

```bash
go get github.com/Sedose/GoDeLin
```

Then import it in your Go files:

```go
import "github.com/Sedose/GoDeLin"
```

🔧 **Usage Examples**
TODO

Happy coding! 🚀

📝 **These functions are not provided**
-   `Chunked`. Use `slices.Chunk` function.
-   `Drop`: Use standard Go slice syntax `slice[n:]`. Handle potential out-of-bounds access if needed (e.g., `slice[min(n, len(slice)):]`).
-   `DropLast`: Use standard Go slice syntax `slice[:len(slice)-n]`. Handle potential negative results if needed (e.g., `slice[:max(0, len(slice)-n)]`).
-   `Take`: Use standard Go slice syntax `slice[:n]`. Handle potential out-of-bounds access if needed (e.g., `slice[:min(n, len(slice))]`).
-   `TakeLast`: Use standard Go slice syntax `slice[len(slice)-n:]`. Handle potential negative results if needed (e.g., `slice[max(0, len(slice)-n):]`).
