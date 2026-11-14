package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestExpandVars(t *testing.T) {
	vars := map[string]string{"PAUSE": "2", "PAUSE_MAX": "9", "BOGUS": "carrots", "LO": "1", "HI": "3", "REF": "$PAUSE"}
	cases := map[string]string{
		"sleep $PAUSE && hostname":         "sleep 2 && hostname",
		"sleep ${PAUSE}s; echo $PAUSE_MAX": "sleep 2s; echo 9",
		"date --$BOGUS || whoami":          "date --carrots || whoami",
		"echo $PAUSEX $? $HOME":            "echo $PAUSEX $? $HOME",
		"#!type sleep $PAUSE":              "#!type sleep 2",
		// Adjacent references resolve independently, in either spelling.
		"ls sales-q[$LO$HI].csv":     "ls sales-q[13].csv",
		"ls sales-q[${LO}${HI}].csv": "ls sales-q[13].csv",
		"echo $HI$LO$HI":             "echo 313",
		// A substituted value is typed verbatim, never expanded again.
		"echo $REF":         "echo $PAUSE",
		"echo $((PAUSE+1))": "echo $((PAUSE+1))",
		"echo ${PAUSE":      "echo ${PAUSE",
	}
	for in, want := range cases {
		if got := expandVars(in, vars); got != want {
			t.Errorf("%q: got %q, want %q", in, got, want)
		}
	}
}

func TestWaitInit(t *testing.T) {
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// init finishes on the third poll
		fmt.Fprintf(w, `{"status":"active","initDone":%t}`, calls.Add(1) >= 3)
	}))
	defer srv.Close()
	c := &apiClient{base: srv.URL}

	no := false
	if err := waitInit(c, "m/u", &no, time.Now().Add(5*time.Second)); err != nil {
		t.Fatalf("waitInit: %v", err)
	}
	if n := calls.Load(); n != 3 {
		t.Errorf("polled %d times, want 3", n)
	}

	// An older daemon does not report initDone: no waiting at all.
	calls.Store(0)
	if err := waitInit(c, "m/u", nil, time.Now().Add(5*time.Second)); err != nil || calls.Load() != 0 {
		t.Errorf("nil initDone: err=%v polls=%d, want no polls", err, calls.Load())
	}

	// An init that never finishes runs into the deadline.
	stuck := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"status":"active","initDone":false}`)
	}))
	defer stuck.Close()
	if err := waitInit(&apiClient{base: stuck.URL}, "m/u", &no, time.Now().Add(500*time.Millisecond)); err == nil {
		t.Error("stuck init: no error, want a deadline error")
	}
}
