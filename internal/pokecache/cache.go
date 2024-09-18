package pokecache

import (
		"sync"
		"time"
		)

//Cache 

type Cache struct {

	data	map[string]cacheEntry
	mut	    *sync.Mutex
}

//Cache Entry

type cacheEntry struct {

	createdAt 	time.Time
	val 	  	[]byte
}


// Methods 

func NewCache(interval time.Duration) Cache {
	cache := Cache{
		data: make(map[string]cacheEntry),
		mut: &sync.Mutex{},
	}
	go cache.reapLoop(interval)
	return cache
}

func (c Cache) Add(key string, val []byte) {
	c.mut.Lock()
	defer c.mut.Unlock()
	c.data[key] = cacheEntry {
		createdAt: time.Now().UTC(),
		val: val,
	}
}

func (c Cache) Get(key string) ([]byte, bool) {
	c.mut.Lock()
	defer c.mut.Unlock()
	output, exists := c.data[key]
	return output.val, exists
}

func (c Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		c.reap(time.Now().UTC(), interval)
	}
}

func (c Cache) reap(now time.Time, last time.Duration) {
	c.mut.Lock()
	defer c.mut.Unlock()
	for k, v := range c.data {
		if v.createdAt.Before(now.Add(-last)) {
			delete(c.data, k)
		}
	}
}