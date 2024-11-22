# mini-cache

mini-cache is an in-memory key:value store/cache.Any object can be stored, for a given duration or forever, and the cache can be safely used by multiple goroutines.

- a thread-safe
- TTL supported (with expiration times)
- Simple cache is like `map[string]interface{}`
- Many cache evication policies