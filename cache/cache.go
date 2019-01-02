package cache

import (
	"sync"
	"time"
)

// MemoryCache defines the interface to store keys in the cache
// copied from https://github.com/goenning/go-cache-demo
// see https://github.com/goenning/go-cache-demo/blob/master/LICENSE
// initially created by: https://github.com/goenning
type MemoryCache struct {
	items map[string]Item
	mu    *sync.RWMutex
}

// NewCache creates a new in memory storage
func NewCache() *MemoryCache {
	return &MemoryCache{
		items: make(map[string]Item),
		mu:    &sync.RWMutex{},
	}
}

// Get a cached content by key
func (s *MemoryCache) Get(key string) []byte {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item := s.items[key]
	if item.Expired() {
		delete(s.items, key)
		return nil
	}
	return item.Content
}

// Set a cached content by key
func (s *MemoryCache) Set(key string, content []byte, duration time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.items[key] = Item{
		Content:    content,
		Expiration: time.Now().Add(duration).UnixNano(),
	}
}

// Item is a cached reference
type Item struct {
	Content    []byte
	Expiration int64
}

// Expired returns true if the item has expired.
func (item Item) Expired() bool {
	if item.Expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > item.Expiration
}
