package simple

import (
	"sort"
	"time"
)

type Cache[K comparable, V any] struct {
	items map[K]*entry[V]
}

type entry[V any] struct {
	value     V
	createdAt time.Time
}

func NewCache[K comparable, V any]() *Cache[K, V] {
	return &Cache[K, V]{
		items: make(map[K]*entry[V], 0),
	}
}

func (c *Cache[K, V]) Set(key K, value V) {
	c.items[key] = &entry[V]{
		value:     value,
		createdAt: time.Now(),
	}
}

// Keys returns cache keys sorted by created time
func (c *Cache[K, _]) Keys() []K {
	ret := make([]K, 0, len(c.items))
	for key := range c.items {
		ret = append(ret, key)
	}
	sort.Slice(ret, func(i, j int) bool {
		i1 := c.items[ret[i]]
		i2 := c.items[ret[j]]
		return i1.createdAt.Before(i2.createdAt)
	})
	return ret
}

func (c *Cache[K, V]) Get(key K) (value V, ok bool) {
	entry, ok := c.items[key]
	if !ok {
		return
	}
	return entry.value, true
}

func (c *Cache[K, V]) Delete(key K) {
	delete(c.items, key)
}

func (c *Cache[K, V]) Len() int {
	return len(c.items)
}
