// Package bus is a minimal pub/sub used to decouple the validation engine
// from UIs: the engine publishes events, any number of UI frontends
// subscribe.
package bus

import "sync"

// Event is a loosely-typed engine event. Type examples: "task", "unit",
// "init", "hint", "shells".
type Event struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

type Bus struct {
	mu   sync.Mutex
	subs map[int]chan Event
	next int
}

func New() *Bus {
	return &Bus{subs: map[int]chan Event{}}
}

// Subscribe returns a buffered event channel and an unsubscribe func.
// Slow subscribers drop events instead of blocking the engine.
func (b *Bus) Subscribe() (<-chan Event, func()) {
	b.mu.Lock()
	defer b.mu.Unlock()
	id := b.next
	b.next++
	ch := make(chan Event, 256)
	b.subs[id] = ch
	return ch, func() {
		b.mu.Lock()
		defer b.mu.Unlock()
		if c, ok := b.subs[id]; ok {
			delete(b.subs, id)
			close(c)
		}
	}
}

func (b *Bus) Publish(ev Event) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, ch := range b.subs {
		select {
		case ch <- ev:
		default: // drop for slow subscriber
		}
	}
}
