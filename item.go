package cache

import "time"

type Item[K comparable, V any] struct {
	Key        K
	Value      V
	Expiration time.Time // if Expiration is zero, the Item never expires; otherwise it expires at the Expiration time
}

func newItem[K comparable, V any](key K, value V, opts ...ItemOption) *Item[K, V] {
	o := new(itemOption)
	for _, optFunc := range opts {
		optFunc(o)
	}
	return &Item[K, V]{
		Key:        key,
		Value:      value,
		Expiration: o.expiration,
	}
}

// return true if the item has an expiration
func (item *Item[K, V]) hasExpiration() bool {
	return !item.Expiration.IsZero()
}

// return true if the item has expired
func (item *Item[K, V]) Exipred() bool {
	if !item.hasExpiration() {
		return false
	}
	return time.Now().After(item.Expiration)
}
