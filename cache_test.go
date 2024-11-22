package cache_test

import (
	"cache"
	"context"
	"testing"
	"time"
)

func TestCacheWithNew(t *testing.T) {
	c := cache.New[int, string]()
	// cache.Len()
	// println(cache.Len())
	c.Set(1, "a")
	// println(cache.Len())
	c.Set(2, "b")
	value1, ok1 := c.Get(1)

	if !ok1 || value1 != "a" {
		t.Errorf("cache.Get(1) = %v, %v; want %v, %v", value1, ok1, "a", true)
	}

	value2, ok2 := c.Get(2)
	if !ok2 || value2 != "b" {
		t.Errorf("cache.Get(2) = %v, %v; want %v, %v", value2, ok2, "b", true)
	}

	length := c.Len()
	if length != 2 {
		t.Errorf("cache.Len() = %v; want %v", length, 2)
	}
	// delete
	c.Delete(1)
	value1, ok1 = c.Get(1)
	if ok1 {
		t.Errorf("cache.Get(1) = %v, %v; want %v, %v", value1, ok1, "", false)
	}
}

func TestCacheWithExpiration(t *testing.T) {
	c := cache.New[int, string]()
	c.Set(1, "a", cache.WithExpiration(3*time.Second))
	c.Set(2, "b", cache.WithExpiration(5*time.Second))
	time.Sleep(4 * time.Second)
	_, ok1 := c.Get(1)
	value2, ok2 := c.Get(2)
	if ok1 {
		t.Errorf("cache.Get(1) shoule be expired")
	}
	if !ok2 || value2 != "b" {
		t.Errorf("cache.Get(2) = %v, %v; want %v, %v", value2, ok2, "b", true)
	}
}

func TestNewContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// use simple cache algorithm without options.
	// an internal janitor will be stopped if specified the context is cancelled.
	c := cache.NewContext(ctx, cache.WithJanitorInterval[string, int](50*time.Second))
	c.Set("a", 1, cache.WithExpiration(3*time.Second))
	gota, aok := c.Get("a")
	gotb, bok := c.Get("b")

	if !aok || gota != 1 {
		t.Errorf("cache.Get(a) = %v, %v; want %v, %v", gota, aok, 1, true)
	}
	if bok {
		t.Errorf("cache.Get(b) = %v, %v; want %v, %v", gotb, bok, 0, false)
	}
}
