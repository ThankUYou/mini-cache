package cache

import (
	"cache/evication/simple"
	"time"
)

// ---------------------------------------------------------------//
// --------------------option for item----------------------------//
// ---------------------------------------------------------------//

// ItemOption is a function that configures an item
type ItemOption func(*itemOption)

type itemOption struct {
	expiration time.Time // default none 
}

// WithExpiration sets the expiration time for the item
func WithExpiration(ttl time.Duration) ItemOption {
	return func(io *itemOption) {
		io.expiration = time.Now().Add(ttl)
	}
}

// ---------------------------------------------------------------//
// --------------------option for cache---------------------------//
// ---------------------------------------------------------------//

// / Option is a function that configures an option
type Option[K comparable, V any] func(*option[K, V])

type option[K comparable, V any] struct {
	cache           Interface[K, *Item[K, V]] // why use *Item[K, V] instead of V? because we need to set the expiration time for the item
	janitorInterval time.Duration             // interval for janitor
}

func newOption[K comparable, V any]() *option[K, V] {
	return &option[K, V]{
		cache:           simple.NewCache[K, *Item[K, V]](),
		janitorInterval: time.Minute,
	}
}

func WithJanitorInterval[K comparable, V any](ttl time.Duration) Option[K, V] {
	return func(o *option[K, V]) {
		o.janitorInterval = ttl
	}
}

// new evication policy
