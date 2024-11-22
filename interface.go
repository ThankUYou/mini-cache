package cache

type Interface[K comparable, V any] interface {
	Get(key K) (value V, ok bool)
	Set(key K, value V)
	Delete(key K)
	Keys() []K
	Len() int
}
