package cmd

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestParallelRespectsLimit(t *testing.T) {
	const limit = 3
	var running, peak, done atomic.Int32
	items := make([]int, 20)
	parallel(t.Context(), limit, items, func(context.Context, int) {
		n := running.Add(1)
		for {
			p := peak.Load()
			if n <= p || peak.CompareAndSwap(p, n) {
				break
			}
		}
		time.Sleep(5 * time.Millisecond)
		running.Add(-1)
		done.Add(1)
	})
	if got := done.Load(); got != int32(len(items)) {
		t.Fatalf("ran %d items, want %d", got, len(items))
	}
	if got := peak.Load(); got != limit {
		t.Fatalf("peak concurrency %d, want %d", got, limit)
	}
}

func TestParallelStopsOnCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	var started atomic.Int32
	parallel(ctx, 1, make([]int, 10), func(context.Context, int) {
		started.Add(1)
		cancel()
	})
	if got := started.Load(); got != 1 {
		t.Fatalf("started %d items after cancel, want 1", got)
	}
}
