# go-cache-LRU

Usage

```go
import lru "github.com/too-rusty/go-cache-LRU"

func main() {
    cache := lru.NewCache[int]().WithCapacity(10)
    cache.Add(42)
}
```
