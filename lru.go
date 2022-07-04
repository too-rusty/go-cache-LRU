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

type lruCache[T comparable] struct {
	list       *list.List
	lastAccess map[T]elementWrapper
	capacity   uint
	mu         sync.Mutex
}

// initiate a new cache
func NewCache[T comparable]() *lruCache[T] {
	return &lruCache[T]{
		list:       list.New(),
		lastAccess: make(map[T]elementWrapper),
	}
}

// Set capacity of the cache, if 0 then it's infinite
func (cache *lruCache[T]) WithCapacity(capacity uint) *lruCache[T] {
	cache.capacity = capacity
	return cache
}

// Capacity of the cache
func (cache *lruCache[T]) Capacity() uint {
	return cache.capacity
}

// Current Length of the cache
func (cache *lruCache[T]) Len() uint {
	return uint(cache.list.Len())
}

// add a value with default timestamp as time.Now
func (cache *lruCache[T]) Add(value T) {
	cache.AddWithTimeStamp(value, time.Now())
}

// add manually with timestamp if time.Now is not threadsafe
func (cache *lruCache[T]) AddWithTimeStamp(value T, moment time.Time) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	if moment.Equal(time.Time{}) {
		moment = time.Now()
	}

	if ptr, found := cache.lastAccess[value]; found {
		cache.list.Remove(ptr.element)
	}

	newElement := cache.list.PushBack(value)
	cache.lastAccess[value] = elementWrapper{
		element:   newElement,
		timestamp: moment,
	}

	if int(cache.capacity) > 0 && cache.list.Len() > int(cache.capacity) {
		// remove any extra element from the cache
		head := cache.list.Front()
		key, _ := head.Value.(T)
		cache.list.Remove(head)
		delete(cache.lastAccess, key)

	}

}

// Remove all elements from the cache before a certain timestamp
func (cache *lruCache[T]) RemoveBefore(moment time.Time) (ret []T) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	for cache.list.Len() > 0 {
		key := cache.list.Front().Value.(T)
		ptr := cache.lastAccess[key]
		if ptr.timestamp.After(moment) {
			break
		}
		ret = append(ret, key)
		cache.list.Remove(ptr.element)
		delete(cache.lastAccess, key)
	}
	return
}

// Remove first n elements from the cache
func (cache *lruCache[T]) RemoveFirstN(n int) (ret []T) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	for cache.list.Len() > 0 && n > 0 {
		head := cache.list.Front()
		key := head.Value.(T)
		ret = append(ret, key)
		cache.list.Remove(head)
		delete(cache.lastAccess, key)
		n--
	}

	return
}
