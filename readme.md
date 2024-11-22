# mini-cache

mini-cache is an in-memory key:value store/cache.Any object can be stored, for a given duration or forever, and the cache can be safely used by multiple goroutines.

- thread-safe
- TTL supported (with expiration times)
- Many cache evication policies(FIFO)

# TODO
- [x] LRU、LFU、MRU、CLOCK
- [x] Add benchmark
- [x] Try to Improve efficiency(shares or some other ways) 