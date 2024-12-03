# mini-cache

mini-cache is an in-memory key:value store/cache.Any object can be stored, for a given duration or forever, and the cache can be safely used by multiple goroutines.

- thread-safe
- TTL supported (with expiration times)
- Many cache evication policies(FIFO)

# TODO
- [x] LRU、LFU
- [x] Add benchmark
- [x] Try to Improve efficiency(shares) 
- [ ] MRU、CLOCK

# Result
```
goos: darwin
goarch: arm64
pkg: cache/benchmark
BenchmarkRWMutexCacheSet-8              58110638                20.65 ns/op            0 B/op          0 allocs/op
BenchmarkRWMutexMapSetConcurrent-8       9827103               127.9 ns/op             0 B/op          0 allocs/op
BenchmarkMiniCacheSet-8                 12181406                96.97 ns/op          104 B/op          3 allocs/op
BenchmarkRWMutexCacheGet-8              59744590                19.48 ns/op            0 B/op          0 allocs/op
BenchmarkRWMutexMapGetConcurrent-8      16569025                74.31 ns/op            0 B/op          0 allocs/op
BenchmarkMiniCacheGet-8                 78443338                14.74 ns/op            0 B/op          0 allocs/op
BenchmarkMiniCacheDelete-8              100000000               11.72 ns/op            0 B/op          0 allocs/op
BenchmarkMiniCacheDeleteExpired-8       100000000               11.00 ns/op            0 B/op          0 allocs/op
```
