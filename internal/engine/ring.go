package engine

import (
	"context"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

// eventClock hands out sequence numbers to every ring in the process, so
// an exec event and a line event can be ordered against each other: a
// check that took a mark after one event (`event_seq`) can ask
// any event check for what came after it (`--after N`).
var eventClock atomic.Uint64

// eventRing is a bounded, sequence-numbered buffer of observed events with
// blocking matchers - the shared core of the exec watcher (proc connector)
// and the line watcher (readline uprobe). Events are stamped with a
// monotonically increasing Seq (from the process-wide eventClock) and the
// publish time; checks wait on the ring by Seq, so a unit's activation
// horizon is one number.
type eventRing[T any] struct {
	mu     sync.Mutex
	cond   *sync.Cond
	ring   []T
	seq    uint64
	closed bool

	size  int
	seqOf func(T) uint64
	stamp func(T, uint64, time.Time) T
}

func newEventRing[T any](size int, seqOf func(T) uint64, stamp func(T, uint64, time.Time) T) *eventRing[T] {
	r := &eventRing[T]{size: size, seqOf: seqOf, stamp: stamp}
	r.cond = sync.NewCond(&r.mu)
	return r
}

// Close wakes every waiter; subsequent waits return false immediately.
func (r *eventRing[T]) Close() {
	r.mu.Lock()
	r.closed = true
	r.cond.Broadcast()
	r.mu.Unlock()
}

func (r *eventRing[T]) isClosed() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.closed
}

// Seq returns the current sequence number; events published after a given
// point have Seq greater than this.
func (r *eventRing[T]) Seq() uint64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.seq
}

// WaitMatch blocks until an event with Seq > after matches fn, returning it.
// Among already-buffered events the OLDEST match wins. Returns false when
// the deadline passes, ctx is canceled (e.g. the waiting check script was
// killed and its API connection dropped - without this, every killed
// attempt would leak a waiter that keeps scanning the ring on each
// broadcast), or the ring closes.
func (r *eventRing[T]) WaitMatch(ctx context.Context, after uint64, deadline time.Time, fn func(T) bool) (T, bool) {
	return r.waitMatch(ctx, after, deadline, fn, false)
}

// WaitMatchLatest is WaitMatch, except that among already-buffered events
// the NEWEST match wins. A check that judges the student's most recent
// answer (branching right/wrong and restarting after a hint) needs this:
// with oldest-first matching, the first wrong answer since activation
// would keep winning every restart and a later correct one could never be
// seen.
func (r *eventRing[T]) WaitMatchLatest(ctx context.Context, after uint64, deadline time.Time, fn func(T) bool) (T, bool) {
	return r.waitMatch(ctx, after, deadline, fn, true)
}

func (r *eventRing[T]) waitMatch(ctx context.Context, after uint64, deadline time.Time, fn func(T) bool, latest bool) (T, bool) {
	var zero T
	timer := time.AfterFunc(time.Until(deadline), func() {
		r.mu.Lock()
		r.cond.Broadcast()
		r.mu.Unlock()
	})
	defer timer.Stop()
	stop := context.AfterFunc(ctx, func() {
		r.mu.Lock()
		r.cond.Broadcast()
		r.mu.Unlock()
	})
	defer stop()

	r.mu.Lock()
	defer r.mu.Unlock()
	scanned := after
	for {
		// The ring ascends by Seq - skip the already-scanned prefix instead
		// of re-running fn over the whole buffer on every wake-up.
		i := sort.Search(len(r.ring), func(i int) bool { return r.seqOf(r.ring[i]) > scanned })
		if latest {
			for j := len(r.ring) - 1; j >= i; j-- {
				if fn(r.ring[j]) {
					return r.ring[j], true
				}
			}
		} else {
			for _, ev := range r.ring[i:] {
				if fn(ev) {
					return ev, true
				}
			}
		}
		scanned = r.seq
		if r.closed || ctx.Err() != nil || time.Now().After(deadline) {
			return zero, false
		}
		r.cond.Wait()
	}
}

// Snapshot returns events with Seq > after (for debugging APIs).
func (r *eventRing[T]) Snapshot(after uint64, limit int) []T {
	r.mu.Lock()
	defer r.mu.Unlock()
	var out []T
	for _, ev := range r.ring {
		if r.seqOf(ev) > after {
			out = append(out, ev)
		}
	}
	if limit > 0 && len(out) > limit {
		out = out[len(out)-limit:]
	}
	return out
}

func (r *eventRing[T]) publish(ev T) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.seq = eventClock.Add(1)
	ev = r.stamp(ev, r.seq, time.Now())
	r.ring = append(r.ring, ev)
	if len(r.ring) > r.size {
		r.ring = r.ring[len(r.ring)-r.size:]
	}
	r.cond.Broadcast()
}
