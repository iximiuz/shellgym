package engine

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
)

// execWait runs one /exec/wait request against a synthetic watcher.
func execWaitOnce(t *testing.T, api *checkAPI, req ExecWaitRequest) bool {
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
	return out.Matched
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
	if execWaitOnce(t, api, req) {
		t.Fatal("argc=2 matched a 3-element argv")
	}
	// Without the argc constraint the same event matches.
	if !execWaitOnce(t, api, ExecWaitRequest{Regex: `^date \+%A %d$`, TimeoutSec: 0.05}) {
		t.Fatal("regex alone should match the unquoted form")
	}

	// Quoted form: two argv elements, same joined string.
	w.publish(ExecEvent{PID: 2, UID: 1000, TTYNr: 3, Argv: []string{"date", "+%A %d"}})
	if !execWaitOnce(t, api, req) {
		t.Fatal("argc=2 did not match the quoted form")
	}
}
