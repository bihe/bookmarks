package security

import (
	"testing"
	"time"
)

// copied and derived from https://github.com/goenning/go-cache-demo
// see https://github.com/goenning/go-cache-demo/blob/master/LICENSE
// initially created by: https://github.com/goenning

func TestMemCacheGetEmpty(t *testing.T) {
	cache := newMemCache(parse("5s"))
	user := cache.get("MY_KEY")
	if user != nil {
		t.Fatalf("exptected empty/nil user!")
	}
}

func TestMemCacheGetValue(t *testing.T) {
	cache := newMemCache(parse("5s"))
	u := User{
		DisplayName: "a",
		Email:       "a.b@c.de",
		UserID:      "1",
		Username:    "u",
		Roles:       []string{"role"},
	}
	cache.set("MY_KEY", &u)
	user := cache.get("MY_KEY")
	if user == nil {
		t.Fatalf("exptected cached used got nil!")
	}
	if user != &u {
		t.Fatalf("the returned object/address is not the same!")
	}
}

func TestMemCacheGetExpiredValue(t *testing.T) {
	cache := newMemCache(parse("1s"))
	u := User{
		DisplayName: "a",
		Email:       "a.b@c.de",
		UserID:      "1",
		Username:    "u",
		Roles:       []string{"role"},
	}
	cache.set("MY_KEY", &u)

	time.Sleep(parse("1s200ms"))

	user := cache.get("MY_KEY")
	if user != nil {
		t.Fatalf("exptected expiry of user object!")
	}
}

func parse(s string) time.Duration {
	d, _ := time.ParseDuration(s)
	return d
}
