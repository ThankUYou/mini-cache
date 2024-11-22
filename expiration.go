package cache

import (
	"container/heap"
	"time"
)

// expireManager is a manager for expiration keys
type expireManager[K comparable] struct {
	queue   expirationQueue[K]      // expiration queue
	mapping map[K]*expirationKey[K] // mapping from key to expiration key
}

type expirationQueue[K comparable] []*expirationKey[K]

type expirationKey[K comparable] struct {
	key        K
	expiration time.Time
	index      int
}

func newExpirationManager[K comparable]() *expireManager[K] {
	queue := make(expirationQueue[K], 0)
	heap.Init(&queue)
	return &expireManager[K]{
		queue:   queue,
		mapping: make(map[K]*expirationKey[K]),
	}
}

// if exist, update the expiration of the key; else add the key to the queue
func (m *expireManager[K]) update(key K, expiration time.Time) {
	if item, ok := m.mapping[key]; ok {
		item.expiration = expiration
		heap.Fix(&m.queue, item.index) // fix the heap
	} else {
		v := &expirationKey[K]{
			key:        key,
			expiration: expiration,
		}
		heap.Push(&m.queue, v)
		m.mapping[key] = v
	}
}

// returns the length of the queue
func (m *expireManager[K]) len() int {
	return m.queue.Len()
}

// pop the key with the earliest expiration
func (m *expireManager[K]) pop() K {
	v := heap.Pop(&m.queue)
	key := v.(*expirationKey[K]).key
	delete(m.mapping, key)
	return key
}

func (m *expireManager[K]) remove(key K) {
	if item, ok := m.mapping[key]; ok {
		heap.Remove(&m.queue, item.index)
		delete(m.mapping, key)
	}
}

var _ heap.Interface = (*expirationQueue[int])(nil)

// implements heap.Interface
func (pq expirationQueue[K]) Len() int { return len(pq) }

func (pq expirationQueue[K]) Less(i, j int) bool {
	return pq[i].expiration.Before(pq[j].expiration)
}

func (pq expirationQueue[K]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *expirationQueue[K]) Push(x interface{}) {
	n := len(*pq)
	item := x.(*expirationKey[K])
	item.index = n
	*pq = append(*pq, item)
}

func (pq *expirationQueue[K]) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}
