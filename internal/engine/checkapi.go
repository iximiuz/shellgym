package engine

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

// The daemon exposes an internal HTTP API over a unix socket for the
// built-in check commands (shellgym check ...) that run inside task
// scripts. This keeps checks tiny and centralizes the expensive machinery
// (exec watching, shell discovery) in one process.

type checkAPI struct {
	watcher   *ExecWatcher
	shellUser string
	shellUID  int
	hintSink  HintSink
	varSink   VarSink
}

// HintSink receives hints posted by the hint_exit built-in from inside
// task scripts.
type HintSink func(unit, task, message string) error

// VarSink receives task vars posted by the set_var built-in from inside
// task scripts.
type VarSink func(unit, name, value string) error

// ServeCheckAPI starts the unix-socket listener at sockPath.
func ServeCheckAPI(sockPath, shellUser string, watcher *ExecWatcher, hints HintSink, vars VarSink) error {
	uid, err := lookupUID(shellUser)
	if err != nil {
		return fmt.Errorf("shell user %q: %w", shellUser, err)
	}
	_ = os.Remove(sockPath)
	ln, err := net.Listen("unix", sockPath)
	if err != nil {
		return err
	}
	if err := os.Chmod(sockPath, 0o600); err != nil {
		return err
	}
	api := &checkAPI{watcher: watcher, shellUser: shellUser, shellUID: uid, hintSink: hints, varSink: vars}
	mux := http.NewServeMux()
	mux.HandleFunc("/shells", api.handleShells)
	mux.HandleFunc("/hint", api.handleHint)
	mux.HandleFunc("/vars", api.handleSetVar)
	mux.HandleFunc("/exec/seq", api.handleSeq)
	mux.HandleFunc("/exec/wait", api.handleExecWait)
	mux.HandleFunc("/exec/snapshot", api.handleSnapshot)
	go func() { _ = http.Serve(ln, mux) }()
	return nil
}

func (a *checkAPI) handleShells(w http.ResponseWriter, r *http.Request) {
	shells, err := FindShells(a.shellUser)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	writeJSON(w, shells)
}

func (a *checkAPI) handleSeq(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, map[string]uint64{"seq": a.watcher.Seq()})
}

// ExecWaitRequest asks the daemon to block until a matching exec event.
type ExecWaitRequest struct {
	After      uint64  `json:"after"`      // only events with Seq > After
	Regex      string  `json:"regex"`      // matched against the joined argv
	EnvName    string  `json:"envName"`    // if set: match process env instead
	EnvRegex   string  `json:"envRegex"`   //
	TimeoutSec float64 `json:"timeoutSec"` // <=0: practically forever
}

type ExecWaitResponse struct {
	Matched bool      `json:"matched"`
	Event   ExecEvent `json:"event"`
}

func (a *checkAPI) handleExecWait(w http.ResponseWriter, r *http.Request) {
	var req ExecWaitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	// Cap the wait: attempts are killed by the engine's task timeout long
	// before an hour, and an uncapped wait would pin a handler goroutine.
	timeout := time.Duration(req.TimeoutSec * float64(time.Second))
	if timeout <= 0 || timeout > time.Hour {
		timeout = time.Hour
	}
	var argvRe, envRe *regexp.Regexp
	var err error
	if req.Regex != "" {
		if argvRe, err = regexp.Compile(req.Regex); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}
	if req.EnvRegex != "" {
		if envRe, err = regexp.Compile(req.EnvRegex); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}
	match := func(ev ExecEvent) bool {
		// -1 = unknown (the process was gone before it could be read);
		// only a CONFIRMED foreign uid / missing tty disqualifies - fast
		// interactive commands routinely die before inspection.
		if ev.UID != -1 && ev.UID != a.shellUID {
			return false
		}
		if ev.TTYNr == 0 {
			return false
		}
		if argvRe != nil && !argvRe.MatchString(strings.Join(ev.Argv, " ")) {
			return false
		}
		if req.EnvName != "" {
			v, ok := envOf(ev, req.EnvName)
			if !ok {
				return false
			}
			if envRe != nil && !envRe.MatchString(v) {
				return false
			}
		}
		return true
	}
	ev, ok := a.watcher.WaitMatch(r.Context(), req.After, time.Now().Add(timeout), match)
	writeJSON(w, ExecWaitResponse{Matched: ok, Event: ev})
}

// HintRequest is posted by the hint_exit built-in.
type HintRequest struct {
	Unit    string `json:"unit"`
	Task    string `json:"task"`
	Message string `json:"message"`
}

func (a *checkAPI) handleHint(w http.ResponseWriter, r *http.Request) {
	var req HintRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if req.Unit == "" || req.Task == "" || req.Message == "" {
		http.Error(w, "unit, task, and message are required", 400)
		return
	}
	if a.hintSink == nil {
		http.Error(w, "hints not supported", 501)
		return
	}
	if err := a.hintSink(req.Unit, req.Task, req.Message); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	writeJSON(w, map[string]bool{"ok": true})
}

// SetVarRequest is posted by the set_var built-in.
type SetVarRequest struct {
	Unit  string `json:"unit"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (a *checkAPI) handleSetVar(w http.ResponseWriter, r *http.Request) {
	var req SetVarRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if req.Unit == "" || req.Name == "" {
		http.Error(w, "unit and name are required", 400)
		return
	}
	if a.varSink == nil {
		http.Error(w, "task vars not supported", 501)
		return
	}
	if err := a.varSink(req.Unit, req.Name, req.Value); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	writeJSON(w, map[string]bool{"ok": true})
}

// envOf reads an env var from the event's eagerly-captured environment,
// falling back to a live /proc read (the process may already be gone).
func envOf(ev ExecEvent, name string) (string, bool) {
	prefix := name + "="
	for _, kv := range ev.Env {
		if strings.HasPrefix(kv, prefix) {
			return kv[len(prefix):], true
		}
	}
	if len(ev.Env) == 0 {
		return ShellEnvOf(ev.PID, name)
	}
	return "", false
}

func (a *checkAPI) handleSnapshot(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, a.watcher.Snapshot(0, 200))
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
