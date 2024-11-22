package fifo

import "container/list"

type Cache[K comparable, V any] struct {
	items    map[K]*list.Element
	queue    *list.List
	capacity int
}

type entry[K comparable, V any] struct {
	key   K
	value V
}

type Option func(*option)

type option struct {
	capacity int
}

func newOption() *option {
	return &option{
		capacity: 128,
	}
}

func WithCapacity(capacity int) Option {
	return func(o *option) {
		o.capacity = capacity
	}
}

func NewCache[K comparable, V any](opts ...Option) *Cache[K, V] {
	o := newOption()
	for _, optFunc := range opts {
		optFunc(o)
	}
	return &Cache[K, V]{
		items:    make(map[K]*list.Element, o.capacity),
		queue:    list.New(),
		capacity: o.capacity,
	}
}

func (c *Cache[K, V]) Set(key K, value V) {
	if c.queue.Len() == c.capacity {
		e := c.dequeue()
		delete(c.items, e.Value.(*entry[K, V]).key)
	}
	c.Delete(key) // delete old key if already exists specified key.
	entry := &entry[K, V]{
		key:   key,
		value: value,
	}
	e := c.queue.PushBack(entry)
	c.items[key] = e
}

func (c *Cache[K, V]) dequeue() *list.Element {
	e := c.queue.Front()
	c.queue.Remove(e)
	return e
}

func (c *Cache[K, V]) Keys() []K {
	keys := make([]K, 0, len(c.items))
	for e := c.queue.Front(); e != nil; e = e.Next() {
		keys = append(keys, e.Value.(*entry[K, V]).key)
	}
	return keys
}

func (c *Cache[K, V]) Get(key K) (value V, ok bool) {
	got, ok := c.items[key]
	if !ok {
		return
	}
	return got.Value.(*entry[K, V]).value, true
}

func (c *Cache[K, V]) Delete(key K) {
	if e, ok := c.items[key]; ok {
		c.queue.Remove(e)
		delete(c.items, key)
	}
}

// Len returns the number of items in the cache.
func (c *Cache[K, V]) Len() int {
	return c.queue.Len()
}
