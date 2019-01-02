package cache_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/bihe/bookmarks/cache"
)

// copied from https://github.com/goenning/go-cache-demo
// see https://github.com/goenning/go-cache-demo/blob/master/LICENSE
// initially created by: https://github.com/goenning

func parse(s string) time.Duration {
	d, _ := time.ParseDuration(s)
	return d
}

func TestGetEmpty(t *testing.T) {
	cache := cache.NewCache()
	content := cache.Get("MY_KEY")

	assertContentEquals(t, content, []byte(""))
}

func TestGetValue(t *testing.T) {
	cache := cache.NewCache()
	cache.Set("MY_KEY", []byte("123456"), parse("5s"))
	content := cache.Get("MY_KEY")

	assertContentEquals(t, content, []byte("123456"))
}

func TestGetExpiredValue(t *testing.T) {
	cache := cache.NewCache()
	cache.Set("MY_KEY", []byte("123456"), parse("1s"))
	time.Sleep(parse("1s200ms"))
	content := cache.Get("MY_KEY")

	assertContentEquals(t, content, []byte(""))
}

func assertContentEquals(t *testing.T, content, expected []byte) {
	if !bytes.Equal(content, expected) {
		t.Errorf("content should '%s', but was '%s'", expected, content)
	}
}
