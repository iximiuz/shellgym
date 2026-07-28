package engine

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
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
