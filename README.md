# go-cache-LRU

Usage

```go
import lru "github.com/too-rusty/go-cache-LRU"

func main() {
    cache := lru.NewCache[int]().WithCapacity(10)
    cache.Add(42)
}
```


https://en.wikipedia.org/wiki/Cache_replacement_policies#Least_recently_used_(LRU)

Get and put operations should update the cache
currently only put operation is supported

## TODO

support get operation along with put