package security

import (
	"sync"
	"time"
)

// copied and derived from https://github.com/goenning/go-cache-demo
// see https://github.com/goenning/go-cache-demo/blob/master/LICENSE
// initially created by: https://github.com/goenning

func newMemCache(duration time.Duration) *memoryCache {
	return &memoryCache{
		items:         make(map[string]cacheItem),
		cacheDuration: duration,
	}
}

// --------------------------------------------------------------------------
// memoryCache to store User objects
// --------------------------------------------------------------------------

type memoryCache struct {
	sync.Mutex
	items         map[string]cacheItem
	cacheDuration time.Duration
}

func (s *memoryCache) get(key string) *User {
	s.Lock()
	defer s.Unlock()

	item := s.items[key]
	if item.expired() {
		delete(s.items, key)
		return nil
	}
	return item.user
}

func (s *memoryCache) set(key string, user *User) {
	s.Lock()
	defer s.Unlock()

	s.items[key] = cacheItem{
		user:       user,
		expiration: time.Now().Add(s.cacheDuration).UnixNano(),
	}
}

// --------------------------------------------------------------------------
// cacheItem holding the data
// --------------------------------------------------------------------------

type cacheItem struct {
	user       *User
	expiration int64
}

func (item cacheItem) expired() bool {
	if item.expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > item.expiration
}
