package simple_test

import (
	"cache/evication/simple"
	"strings"
	"testing"
)

func TestSimple(t *testing.T) {
	c := simple.NewCache[string, int]()
	c.Set("a", 1)
	gota, aok := c.Get("a")
	_, bok := c.Get("b")
	if c.Len() != 1 {
		t.Errorf("Expected len to be 1, got %v", c.Len())
	}
	if !aok {
		t.Errorf("Expected aok to be true, got %v", aok)
	}
	if gota != 1 {
		t.Errorf("Expected gota to be 1, got %v", gota)
	}
	if bok {
		t.Errorf("Expected bok to be false, got %v", bok)
	}

	c.Delete("a")
	if c.Len() != 0 {
		t.Errorf("Expected len to be 0, got %v", c.Len())
	}

	_, aok = c.Get("a")
	if aok {
		t.Errorf("Expected aok to be false, got %v", aok)
	}
}

func TestKeys(t *testing.T) {
	cache := simple.NewCache[string, int]()
	cache.Set("foo", 1)
	cache.Set("bar", 2)
	cache.Set("baz", 3)
	cache.Set("bar", 4) // again
	cache.Set("foo", 5) // again

	got := strings.Join(cache.Keys(), ",")
	want := strings.Join([]string{
		"baz",
		"bar",
		"foo",
	}, ",")
	if got != want {
		t.Errorf("want %q, but got %q", want, got)
	}
	if len(cache.Keys()) != cache.Len() {
		t.Errorf("want number of keys %d, but got %d", len(cache.Keys()), cache.Len())
	}
}
