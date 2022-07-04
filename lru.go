package LRU

import (
	"container/list"
	"sync"
	"time"
)

type elementWrapper struct {
	element   *list.Element
	timestamp time.Time
}

type lruCache[Value comparable] struct {
	list       *list.List
	lastAccess map[Value]elementWrapper
	capacity   uint
	mu         sync.Mutex
}

type LruElement[Value comparable] struct {
	V Value
	T time.Time
}

// initiate a new cache
func NewCache[Value comparable]() *lruCache[Value] {
	return &lruCache[Value]{
		list:       list.New(),
		lastAccess: make(map[Value]elementWrapper),
	}
}

// Set capacity of the cache, if 0 then it's infinite
func (cache *lruCache[Value]) WithCapacity(capacity uint) *lruCache[Value] {
	cache.capacity = capacity
	return cache
}

// Capacity of the cache
func (cache *lruCache[Value]) Capacity() uint {
	return cache.capacity
}

// Current Length of the cache
func (cache *lruCache[Value]) Len() uint {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	return uint(cache.list.Len())
}

// add a value with default timestamp as time.Now
func (cache *lruCache[Value]) Add(value Value) {
	cache.AddLruElement(LruElement[Value]{V: value, T: time.Now()})
}

// add manually with timestamp if time.Now is not threadsafe
func (cache *lruCache[Value]) AddLruElement(data LruElement[Value]) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	if data.T.Equal(time.Time{}) {
		data.T = time.Now()
	}

	if ptr, found := cache.lastAccess[data.V]; found {
		cache.list.Remove(ptr.element)
	}

	newElement := cache.list.PushBack(data.V)
	cache.lastAccess[data.V] = elementWrapper{
		element:   newElement,
		timestamp: data.T,
	}

	if int(cache.capacity) > 0 && cache.list.Len() > int(cache.capacity) {
		// remove any extra element from the cache
		head := cache.list.Front()
		key, _ := head.Value.(Value)
		cache.list.Remove(head)
		delete(cache.lastAccess, key)

	}

}

// Remove all elements from the cache before a certain timestamp
func (cache *lruCache[Value]) RemoveBefore(moment time.Time) (ret []LruElement[Value]) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	for cache.list.Len() > 0 {
		key := cache.list.Front().Value.(Value)
		ptr := cache.lastAccess[key]
		if ptr.timestamp.After(moment) {
			break
		}
		ret = append(ret, LruElement[Value]{key, ptr.timestamp})
		cache.list.Remove(ptr.element)
		delete(cache.lastAccess, key)
	}
	return
}

// Remove first n elements from the cache
func (cache *lruCache[Value]) RemoveFirstN(n int) (ret []LruElement[Value]) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	for cache.list.Len() > 0 && n > 0 {
		head := cache.list.Front()
		key := head.Value.(Value)
		tstamp := cache.lastAccess[key].timestamp
		ret = append(ret, LruElement[Value]{key, tstamp})
		cache.list.Remove(head)
		delete(cache.lastAccess, key)
		n--
	}

	return
}

// Remove first n elements from the cache
func (cache *lruCache[Value]) ClearCache() []LruElement[Value] {
	return cache.RemoveFirstN(int(cache.Len()))
}
