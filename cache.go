package cache

import (
	"context"
	"sync"
)

type Cache[K comparable, V any] struct {
	cache         Interface[K, *Item[K, V]] // cache
	mu            sync.Mutex                // mutex, lock for cache
	janitor       *janitor                  // janitor for cache expiration
	expireManager *expireManager[K]         // expire manager
}

func New[K comparable, V any](opts ...Option[K, V]) *Cache[K, V] {
	return NewContext[K, V](context.Background(), opts...)
}

func NewContext[K comparable, V any](ctx context.Context, opts ...Option[K, V]) *Cache[K, V] {
	option := newOption[K, V]()
	for _, optFunc := range opts {
		optFunc(option)
	}
	cache := &Cache[K, V]{
		cache:         option.cache,
		janitor:       NewJanitor(ctx, option.janitorInterval),
		expireManager: newExpirationManager[K](),
	}
	cache.janitor.Run(cache.DeleteExpired)
	return cache
}

func (c *Cache[K, V]) Get(key K) (zero V, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	item, ok := c.cache.Get(key)
	if !ok {
		return
	}
	if item.Exipred() {
		return zero, false
	}
	return item.Value, true
}

// GetOrSet atomically gets a key's value from the cache, or if the
// key is not present, sets the given value.
// The loaded result is true if the value was loaded, false if stored.
func (c *Cache[K, V]) GetOrSet(key K, value V, opts ...ItemOption) (actual V, loaded bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	item, ok := c.cache.Get(key)
	if ok && !item.Exipred() {
		item := newItem(key, value, opts...)
		c.cache.Set(key, item)
		return value, false
	}
	return item.Value, true
}

// Set a value to the cache
func (c *Cache[K, V]) Set(key K, value V, opts ...ItemOption) {
	c.mu.Lock()
	defer c.mu.Unlock()
	item := newItem(key, value, opts...)
	if item.hasExpiration() {
		c.expireManager.update(key, item.Expiration)
	}
	c.cache.Set(key, item)

}

// Delete an item from the cache according to the key
func (c *Cache[K, V]) Delete(key K) {
	c.mu.Lock()
	c.cache.Delete(key)
	c.expireManager.remove(key)
	c.mu.Unlock()
}

func (c *Cache[K, V]) DeleteExpired() {
	c.mu.Lock()
	l := c.expireManager.len()
	c.mu.Unlock()
	evict := func() bool {
		key := c.expireManager.pop()
		item, ok := c.cache.Get(key)
		if ok {
			if item.Exipred() {
				c.cache.Delete(key)
				return false
			}
			c.expireManager.update(key, item.Expiration)
		}
		return true
	}
	for i := 0; i < l; i++ {
		c.mu.Lock()
		shouldBreak := evict()
		c.mu.Unlock()
		if shouldBreak {
			break
		}
	}
}

// return sorted keys in the cache
func (c *Cache[K, V]) Keys() []K {
	c.mu.Lock()
	count := c.cache.Keys()
	c.mu.Unlock()
	return count
}

// return the number of items in the cache
func (c *Cache[K, V]) Len() int {
	c.mu.Lock()
	len := c.cache.Len()
	c.mu.Unlock()
	return len
}

// check if the cache contains the key
func (c *Cache[K, V]) Contains(key K) bool {
	c.mu.Lock()
	_, ok := c.cache.Get(key)
	c.mu.Unlock()
	return ok
}
