package cache_test

import (
	"cache"
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestJanoitor(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// cache.Items = make(map[interface{}]cache.Item)
	janitor := cache.NewJanitor(ctx, 1*time.Microsecond)
	checkDone := make(chan bool)
	janitor.Stop = checkDone

	calledClean := int64(0)
	janitor.Run(func() { atomic.AddInt64(&calledClean, 1) })

	// waiting for cleanup
	time.Sleep(10 * time.Millisecond)
	cancel()

	select {
	case <-checkDone:
	case <-time.After(time.Second):
		t.Fatalf("failed to call done channel")
	}

	got := atomic.LoadInt64(&calledClean)
	if got <= 1 {
		t.Fatalf("failed to call clean callback in janitor: %d", got)
	}
}
