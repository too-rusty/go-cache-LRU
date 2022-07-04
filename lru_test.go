package LRU

import (
	"sort"
	"sync"
	"testing"
	"time"
)

type DummyMap struct {
	timestamps map[int64][]time.Time
	mu         sync.Mutex
}

func TestLRUOrder(t *testing.T) {
	capacity := uint(3)

	cache := NewCache[int64]().WithCapacity(capacity)
	m := DummyMap{timestamps: make(map[int64][]time.Time)}

	var wg sync.WaitGroup
	for i := int64(0); i < 10; i++ {
		wg.Add(1)

		func(value int64) {
			defer wg.Done()

			moment := time.Now()
			momentUnix := moment.Unix() + value*2

			value = value % int64(cache.Capacity())

			cache.AddWithTimeStamp(value, time.Unix(momentUnix, 0))

			m.mu.Lock()
			defer m.mu.Unlock()

			m.timestamps[value] = append(m.timestamps[value], time.Unix(momentUnix, 0))

		}(i)

	}
	wg.Wait()

	type TmpValue struct {
		value int64
		t     int64
	}

	var arr []TmpValue
	for k, timestamps := range m.timestamps {
		for _, v := range timestamps {
			arr = append(arr, TmpValue{value: k, t: v.Unix()})
		}
	}

	sort.Slice(arr, func(i, j int) bool {
		return arr[i].t < arr[j].t
	})

	// take all values

	tmpArr := arr[len(arr)-int(cache.Capacity()):]
	cacheValues := cache.RemoveFirstN(int(cache.Capacity()))

	if len(tmpArr) != len(cacheValues) {
		t.Errorf("got: %v, want: %v", len(cacheValues), len(tmpArr))
	}

	for i := 0; i < len(tmpArr); i++ {
		if cacheValues[i] != tmpArr[i].value {
			t.Error("value mismatch")
		}
	}

}

func TestLRUConcurrentFullCapacity(t *testing.T) {

	capacity := uint(100)
	cache := NewCache[int64]().WithCapacity(capacity)

	var wg sync.WaitGroup
	for i := int64(0); i < 100000; i++ {
		wg.Add(1)
		go func(value int64) {
			defer wg.Done()
			cache.Add(value)
		}(i)
	}
	wg.Wait()

	if cache.Len() != capacity {
		t.Errorf("got: %v, want: %v", cache.Len(), capacity)
	}

}

func TestLRUConcurrentLesserCapacity(t *testing.T) {

	capacity := uint(100)
	cache := NewCache[int64]().WithCapacity(capacity)

	var wg sync.WaitGroup
	var mxVal int64 = 10
	for i := int64(0); i < 100000; i++ {
		wg.Add(1)
		go func(value int64) {
			defer wg.Done()
			cache.Add(value)
		}(i % mxVal)
	}
	wg.Wait()

	if cache.Len() != uint(mxVal) {
		t.Errorf("got: %v, want: %v", cache.Len(), mxVal)
	}

}