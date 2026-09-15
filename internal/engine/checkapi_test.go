package engine

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/iximiuz/labs-content/tools/shellgym/internal/bus"
	"github.com/iximiuz/labs-content/tools/shellgym/internal/content"
	"github.com/iximiuz/labs-content/tools/shellgym/internal/state"
)

// execWait runs one /exec/wait request against a synthetic watcher.
func execWaitOnce(t *testing.T, api *checkAPI, req ExecWaitRequest) ExecWaitResponse {
	t.Helper()
	body, _ := json.Marshal(req)
	r := httptest.NewRequest("POST", "http://gym/exec/wait", bytes.NewReader(body))
	w := httptest.NewRecorder()
	api.handleExecWait(w, r)
	if w.Code != 200 {
		t.Fatalf("exec/wait returned %d: %s", w.Code, w.Body.String())
	}
	var out ExecWaitResponse
	if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
		t.Fatal(err)
	}
	return out
}

// The joined-argv regex cannot tell `date '+%A %d'` (one quoted argument)
// from `date +%A %d` (two arguments) - both join to the same string. The
// argc filter can, and that is what quoting reps rely on.
func TestExecWaitArgcDistinguishesQuoting(t *testing.T) {
	w := NewExecWatcher()
	api := &checkAPI{watcher: w, shellUID: 1000}

	// Unquoted form: three argv elements.
	w.publish(ExecEvent{PID: 1, UID: 1000, TTYNr: 3, Argv: []string{"date", "+%A", "%d"}})

	req := ExecWaitRequest{Regex: `^date \+%A %d$`, Argc: 2, TimeoutSec: 0.05}
	if execWaitOnce(t, api, req).Matched {
		t.Fatal("argc=2 matched a 3-element argv")
	}
	// Without the argc constraint the same event matches.
	if !execWaitOnce(t, api, ExecWaitRequest{Regex: `^date \+%A %d$`, TimeoutSec: 0.05}).Matched {
		t.Fatal("regex alone should match the unquoted form")
	}

	// Quoted form: two argv elements, same joined string.
	w.publish(ExecEvent{PID: 2, UID: 1000, TTYNr: 3, Argv: []string{"date", "+%A %d"}})
	if !execWaitOnce(t, api, req).Matched {
		t.Fatal("argc=2 did not match the quoted form")
	}
}

// A check that branches on WHICH command matched needs two things: the
// response must carry the matched argv, and --latest must prefer the
// newest buffered answer - otherwise the first (wrong) answer since
// activation would keep winning after every hint_exit restart.
func TestExecWaitLatestPrefersNewestMatch(t *testing.T) {
	w := NewExecWatcher()
	api := &checkAPI{watcher: w, shellUID: 1000}
	w.publish(ExecEvent{PID: 1, UID: 1000, TTYNr: 3, Argv: []string{"whoami"}})
	w.publish(ExecEvent{PID: 2, UID: 1000, TTYNr: 3, Argv: []string{"hostname"}})

	re := `(^|/)(hostname|whoami)$`
	oldest := execWaitOnce(t, api, ExecWaitRequest{Regex: re, TimeoutSec: 0.05})
	if !oldest.Matched || oldest.Event.Argv[0] != "whoami" {
		t.Fatalf("default should match the oldest event, got %+v", oldest.Event)
	}
	newest := execWaitOnce(t, api, ExecWaitRequest{Regex: re, Latest: true, TimeoutSec: 0.05})
	if !newest.Matched || newest.Event.Argv[0] != "hostname" {
		t.Fatalf("--latest should match the newest event, got %+v", newest.Event)
	}
}

// readSSEEvent scans lines from an SSE stream up to the next blank line and
// returns the "event:" and "data:" fields.
func readSSEEvent(t *testing.T, scanner *bufio.Scanner) (event, data string) {
	t.Helper()
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "event: "):
			event = strings.TrimPrefix(line, "event: ")
		case strings.HasPrefix(line, "data: "):
			data = strings.TrimPrefix(line, "data: ")
		case line == "" && event != "":
			return event, data
		}
	}
	t.Fatal("SSE stream ended before a full event was received")
	return "", ""
}

// A supervisor watching /units/watch must never miss a unit completion, even
// across a reconnect - the snapshot on connect and one event per subsequent
// change are what make that possible.
func TestUnitsWatchStreamsSnapshotThenUpdates(t *testing.T) {
	path := &content.Path{ID: "p", Modules: []*content.Module{{ID: "m", Units: []*content.Unit{
		{ID: "m/a", ModuleID: "m"},
		{ID: "m/b", ModuleID: "m", Unsupported: true},
	}}}}
	st, err := state.Open(t.TempDir(), "p")
	if err != nil {
		t.Fatal(err)
	}
	b := bus.New()
	eng := New(path, st, b, NewExecWatcher(), Options{})
	api := &checkAPI{eng: eng}

	srv := httptest.NewServer(http.HandlerFunc(api.handleUnitsWatch))
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", srv.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)

	event, data := readSSEEvent(t, scanner)
	if event != "snapshot" {
		t.Fatalf("first event = %q, want %q", event, "snapshot")
	}
	var snap unitsSnapshot
	if err := json.Unmarshal([]byte(data), &snap); err != nil {
		t.Fatalf("decode snapshot: %v", err)
	}
	if snap.Path != "p" {
		t.Fatalf("snapshot path = %q, want %q", snap.Path, "p")
	}
	wantUnits := map[string]string{"m/a": "pending", "m/b": "unsupported"}
	if !reflect.DeepEqual(snap.Units, wantUnits) {
		t.Fatalf("snapshot units = %+v, want %+v", snap.Units, wantUnits)
	}

	// A non-unit event must not be forwarded; the unit event that follows
	// must be the very next thing the client sees.
	b.Publish(bus.Event{Type: "task", Data: TaskEvent{Unit: "m/a", Task: "x", Status: "completed"}})
	b.Publish(bus.Event{Type: "unit", Data: UnitEvent{Unit: "m/a", Status: "completed"}})

	event, data = readSSEEvent(t, scanner)
	if event != "unit" {
		t.Fatalf("second event = %q, want %q", event, "unit")
	}
	var us unitStatus
	if err := json.Unmarshal([]byte(data), &us); err != nil {
		t.Fatalf("decode unit status: %v", err)
	}
	if us.ID != "m/a" || us.Status != "completed" {
		t.Fatalf("unit status = %+v, want {m/a completed}", us)
	}

	cancel()
	for scanner.Scan() {
		// drain until the cancelled context closes the body
	}
}

func lineWaitOnce(t *testing.T, api *checkAPI, req LineWaitRequest) (LineWaitResponse, int) {
	t.Helper()
	body, _ := json.Marshal(req)
	r := httptest.NewRequest("POST", "http://gym/line/wait", bytes.NewReader(body))
	w := httptest.NewRecorder()
	api.handleLineWait(w, r)
	var out LineWaitResponse
	if w.Code == 200 {
		if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
			t.Fatal(err)
		}
	}
	return out, w.Code
}

// wait_line judges the typed line itself, so `sleep 2; hostname` and
// `sleep 2 && hostname` - identical as exec events - are told apart; the
// usual student-activity scoping (uid, tty, horizon) applies.
func TestLineWaitMatchesTypedOperator(t *testing.T) {
	w := NewLineWatcher()
	w.Source = "uprobe"
	api := &checkAPI{lines: w, shellUID: 1000}

	w.publish(LineEvent{PID: 1, UID: 1000, TTYNr: 3, Line: "sleep 2; hostname"})
	req := LineWaitRequest{Regex: `^sleep 2 && hostname$`, TimeoutSec: 0.05}
	if out, _ := lineWaitOnce(t, api, req); out.Matched {
		t.Fatal("`;` line matched the && regex")
	}
	w.publish(LineEvent{PID: 1, UID: 1000, TTYNr: 3, Line: "  sleep 2 && hostname  "})
	out, code := lineWaitOnce(t, api, req)
	if code != 200 || !out.Matched || out.Event.Line != "  sleep 2 && hostname  " {
		t.Fatalf("&& line not matched: code=%d %+v", code, out)
	}

	// Foreign uid and tty-less shells never count; the horizon hides
	// lines typed before the unit was activated.
	w.publish(LineEvent{PID: 2, UID: 0, TTYNr: 3, Line: "whoami"})
	w.publish(LineEvent{PID: 3, UID: 1000, TTYNr: 0, Line: "whoami"})
	if out, _ := lineWaitOnce(t, api, LineWaitRequest{Regex: `^whoami$`, TimeoutSec: 0.05}); out.Matched {
		t.Fatal("foreign-uid or tty-less line matched")
	}
	horizon := w.Seq()
	if out, _ := lineWaitOnce(t, api, LineWaitRequest{After: horizon, Regex: `hostname`, TimeoutSec: 0.05}); out.Matched {
		t.Fatal("line before the horizon matched")
	}

	// --latest prefers the newest match (right/wrong branching).
	w.publish(LineEvent{PID: 1, UID: 1000, TTYNr: 3, Line: "date --bogus; whoami"})
	w.publish(LineEvent{PID: 1, UID: 1000, TTYNr: 3, Line: "date --bogus || whoami"})
	out, _ = lineWaitOnce(t, api, LineWaitRequest{After: horizon, Regex: `^date --bogus ?(;|\|\|) whoami$`, TimeoutSec: 0.05})
	if !out.Matched || out.Event.Line != "date --bogus; whoami" {
		t.Fatalf("oldest-first: %+v", out)
	}
	out, _ = lineWaitOnce(t, api, LineWaitRequest{After: horizon, Regex: `^date --bogus ?(;|\|\|) whoami$`, Latest: true, TimeoutSec: 0.05})
	if !out.Matched || out.Event.Line != "date --bogus || whoami" {
		t.Fatalf("--latest: %+v", out)
	}
}

// Without the readline capability the endpoint refuses instead of hanging,
// so a wait_line in a unit that forgot `requires: [readline]` fails fast.
func TestLineWaitUnavailable(t *testing.T) {
	for _, api := range []*checkAPI{{shellUID: 1000}, {lines: NewLineWatcher(), shellUID: 1000}} {
		if _, code := lineWaitOnce(t, api, LineWaitRequest{Regex: `x`, TimeoutSec: 0.05}); code != 501 {
			t.Fatalf("want 501 without an active line watcher, got %d", code)
		}
	}
}

// The event horizon is per unit, so a task gated on "the shell stands in
// X" (needs: [at_x]) would otherwise be satisfied by a matching command the
// student ran anywhere before moving. The --cwd filter scopes the match to
// commands executed from the directory; an unknown cwd never satisfies it.
func TestExecWaitCwdScopesToDirectory(t *testing.T) {
	w := NewExecWatcher()
	api := &checkAPI{watcher: w, shellUID: 1000}
	w.publish(ExecEvent{PID: 1, UID: 1000, TTYNr: 3, Argv: []string{"ls", "-l"}, Cwd: "/home/u"})
	w.publish(ExecEvent{PID: 2, UID: 1000, TTYNr: 3, Argv: []string{"ls", "-l"}, Cwd: ""})

	re := `(^|/)ls( +\S+)*$`
	if !execWaitOnce(t, api, ExecWaitRequest{Regex: re, TimeoutSec: 0.05}).Matched {
		t.Fatal("without --cwd the event should match")
	}
	req := ExecWaitRequest{Regex: re, Cwd: "/home/u/projects/data", TimeoutSec: 0.05}
	if execWaitOnce(t, api, req).Matched {
		t.Fatal("a command run elsewhere (or with an unknown cwd) matched the --cwd filter")
	}

	w.publish(ExecEvent{PID: 3, UID: 1000, TTYNr: 3, Argv: []string{"ls", "-l"}, Cwd: "/home/u/projects/data"})
	out := execWaitOnce(t, api, req)
	if !out.Matched || out.Event.PID != 3 {
		t.Fatalf("the command run from the directory should match, got %+v", out.Event)
	}
	// The pattern follows wait_cwd's rules: a regex when it contains
	// metacharacters, anchored to the whole path.
	if !execWaitOnce(t, api, ExecWaitRequest{Regex: re, Cwd: "/home/u/projects/(data|notes)", TimeoutSec: 0.05}).Matched {
		t.Fatal("regex --cwd did not match")
	}
	if execWaitOnce(t, api, ExecWaitRequest{Regex: re, Cwd: "/home/u/projects", TimeoutSec: 0.05}).Matched {
		t.Fatal("--cwd must match the whole path, not a prefix")
	}
}

// Exec and line events share one sequence clock, and event_seq reads it:
// a mark taken after one event is exceeded by every later event of either
// kind, so "the `echo $?` line typed after the sleep" is `wait_line
// --after <mark taken once the sleep exec was seen>`.
func TestEventSeqMarkOrdersExecAndLineEvents(t *testing.T) {
	w := NewExecWatcher()
	l := NewLineWatcher()
	l.Source = "uprobe"
	api := &checkAPI{watcher: w, lines: l, shellUID: 1000}

	l.publish(LineEvent{PID: 1, UID: 1000, TTYNr: 3, Line: "echo $?"}) // typed before the sleep
	w.publish(ExecEvent{PID: 2, UID: 1000, TTYNr: 3, Argv: []string{"sleep", "350"}})
	if !execWaitOnce(t, api, ExecWaitRequest{Regex: `(^|/)sleep 350$`, TimeoutSec: 0.05}).Matched {
		t.Fatal("sleep exec not seen")
	}
	r := httptest.NewRequest("GET", "http://gym/events/seq", nil)
	rec := httptest.NewRecorder()
	api.handleEventsSeq(rec, r)
	var mark struct {
		Seq uint64 `json:"seq"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &mark); err != nil || mark.Seq == 0 {
		t.Fatalf("bad mark: %v %s", err, rec.Body.String())
	}
	if out, _ := lineWaitOnce(t, api, LineWaitRequest{After: mark.Seq, Regex: `^echo \$\?$`, TimeoutSec: 0.05}); out.Matched {
		t.Fatal("a line typed before the mark must not satisfy --after")
	}
	l.publish(LineEvent{PID: 1, UID: 1000, TTYNr: 3, Line: "echo $?"}) // typed after
	out, _ := lineWaitOnce(t, api, LineWaitRequest{After: mark.Seq, Regex: `^echo \$\?$`, TimeoutSec: 0.05})
	if !out.Matched || out.Event.Seq <= mark.Seq {
		t.Fatalf("the later line should match above the mark: %+v (mark %d)", out, mark.Seq)
	}
}
