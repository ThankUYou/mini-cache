package cache

type share[K ~string, V any] struct {
	seed   uint32
	mask   uint32
	bucket []Cache[K, V]
	// janitor *janitor
}

func (s *share[K, V]) getBucket(key K) *Cache[K, V] {
	hash := djb33(s.seed, string(key))
	return &s.bucket[hash&s.mask]
}

func (s *share[K, V]) Get(key K) (V, bool) {
	return s.getBucket(key).Get(key)
}

func (s *share[K, V]) Set(key K, value V) {
	s.getBucket(key).Set(key, value)
}

func (s *share[K, V]) Delete(key K) {
	s.getBucket(key).Delete(key)
}

func (s *share[K, V]) DeleteExpired() {
	for _, b := range s.bucket {
		b.DeleteExpired()
	}
}



// djb2 with better shuffling. 5x faster than FNV with the hash.Hash overhead.
func djb33(seed uint32, k string) uint32 {
	var (
		l = uint32(len(k))
		d = 5381 + seed + l
		i = uint32(0)
	)
	// Why is all this 5x faster than a for loop?
	if l >= 4 {
		for i < l-4 {
			d = (d * 33) ^ uint32(k[i])
			d = (d * 33) ^ uint32(k[i+1])
			d = (d * 33) ^ uint32(k[i+2])
			d = (d * 33) ^ uint32(k[i+3])
			i += 4
		}
	}
	switch l - i {
	case 1:
	case 2:
		d = (d * 33) ^ uint32(k[i])
	case 3:
		d = (d * 33) ^ uint32(k[i])
		d = (d * 33) ^ uint32(k[i+1])
	case 4:
		d = (d * 33) ^ uint32(k[i])
		d = (d * 33) ^ uint32(k[i+1])
		d = (d * 33) ^ uint32(k[i+2])
	}
	return d ^ (d >> 16)
}
