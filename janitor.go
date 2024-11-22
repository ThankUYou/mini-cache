package cache

import (
	"context"
	"time"
)

type janitor struct {
	ctx      context.Context // context for the janitor
	Stop     chan bool       // signal to stop the janitor
	interval time.Duration   // how often to clean the cache
}

func NewJanitor(ctx context.Context, interval time.Duration) *janitor {
	return &janitor{
		ctx:      ctx,
		Stop:     make(chan bool),
		interval: interval,
	}
}

// stop the janitor
func (j *janitor) kill() {
	close(j.Stop)
}

// run the janitor
func (j *janitor) Run(clean func()) {
	go func() {
		ticker := time.NewTicker(j.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				clean()
			case <-j.ctx.Done():
				j.kill()
			case <-j.Stop:
				clean()
				return
			}
		}
	}()
}
