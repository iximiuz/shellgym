package bus

import (
	"testing"
	"time"
)

func TestPubSub(t *testing.T) {
	b := New()
	ch1, unsub1 := b.Subscribe()
	ch2, unsub2 := b.Subscribe()
	defer unsub2()

	b.Publish(Event{Type: "x", Data: 1})
	for _, ch := range []<-chan Event{ch1, ch2} {
		select {
		case ev := <-ch:
			if ev.Type != "x" {
				t.Errorf("ev = %+v", ev)
			}
		case <-time.After(time.Second):
			t.Fatal("no event")
		}
	}

	unsub1()
	if _, ok := <-ch1; ok {
		t.Error("channel not closed after unsubscribe")
	}
	unsub1() // double-unsubscribe must not panic
	b.Publish(Event{Type: "y"})
	select {
	case ev := <-ch2:
		if ev.Type != "y" {
			t.Errorf("ev = %+v", ev)
		}
	case <-time.After(time.Second):
		t.Fatal("no event after other unsubscribed")
	}
}

func TestSlowSubscriberDoesNotBlock(t *testing.T) {
	b := New()
	_, unsub := b.Subscribe() // never read
	defer unsub()
	done := make(chan struct{})
	go func() {
		for i := 0; i < 1000; i++ {
			b.Publish(Event{Type: "spam"})
		}
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("publisher blocked by slow subscriber")
	}
}
