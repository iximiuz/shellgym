package engine

import (
	"context"
	"os/exec"
	"strings"
	"testing"
	"time"
)

// startWatcher starts a real netlink watcher, skipping the test when the
// proc connector is unavailable (no CAP_NET_ADMIN / not root). Run the
// suite as root on the playground to exercise these.
func startWatcher(t *testing.T) *ExecWatcher {
	t.Helper()
	w := NewExecWatcher()
	w.Start()
	if w.Source != "netlink" {
		t.Skip("proc connector unavailable (needs CAP_NET_ADMIN); run as root")
	}
	t.Cleanup(w.Close)
	return w
}

func TestWatcherSeesExec(t *testing.T) {
	w := startWatcher(t)
	after := w.Seq()

	marker := "shellgym-test-marker-4711"
	cmd := exec.Command("bash", "-c", "exec -a "+marker+" sleep 2")
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = cmd.Process.Kill(); _, _ = cmd.Process.Wait() }()

	ev, ok := w.WaitMatch(context.Background(), after, time.Now().Add(5*time.Second), func(ev ExecEvent) bool {
		return strings.Contains(strings.Join(ev.Argv, " "), marker)
	})
	if !ok {
		t.Fatal("exec event not observed by watcher")
	}
	if ev.PID <= 0 {
		t.Errorf("bad event: %+v", ev)
	}
}

func TestWaitMatchTimeout(t *testing.T) {
	w := NewExecWatcher()
	start := time.Now()
	_, ok := w.WaitMatch(context.Background(), w.Seq(), time.Now().Add(600*time.Millisecond), func(ExecEvent) bool {
		return false
	})
	if ok {
		t.Fatal("unexpected match")
	}
	if time.Since(start) > 3*time.Second {
		t.Fatal("deadline not honored")
	}
}

func TestWaitMatchSeesBufferedEvents(t *testing.T) {
	w := NewExecWatcher()
	w.publish(ExecEvent{PID: 1, Argv: []string{"earlier"}})
	// A match request with after=0 must see the already-buffered event.
	ev, ok := w.WaitMatch(context.Background(), 0, time.Now().Add(time.Second), func(ev ExecEvent) bool {
		return ev.Argv[0] == "earlier"
	})
	if !ok || ev.PID != 1 {
		t.Fatalf("buffered event not found: %+v ok=%v", ev, ok)
	}
}

// A killed check script drops its API connection; the waiter must unblock on
// that cancellation instead of lingering until the deadline (a leak that,
// accumulated, made every broadcast crawl and delayed live matches).
func TestWaitMatchUnblocksOnContextCancel(t *testing.T) {
	w := NewExecWatcher()
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan bool, 1)
	go func() {
		_, ok := w.WaitMatch(ctx, w.Seq(), time.Now().Add(time.Hour), func(ExecEvent) bool { return false })
		done <- ok
	}()
	time.Sleep(50 * time.Millisecond)
	cancel()
	select {
	case ok := <-done:
		if ok {
			t.Fatal("unexpected match")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("WaitMatch did not unblock on context cancellation")
	}
}
