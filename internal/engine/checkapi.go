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

	"github.com/iximiuz/labs-content/tools/shellgym/internal/checkclient"
)

// The daemon exposes an internal HTTP API over a unix socket for the
// built-in check commands (shellgym check ...) that run inside task
// scripts. This keeps checks tiny and centralizes the expensive machinery
// (exec watching, shell discovery) in one process.

type checkAPI struct {
	watcher   *ExecWatcher
	lines     *LineWatcher // nil or inactive: /line/* unavailable
	shellUser string
	shellUID  int
	hintSink  HintSink
	varSink   VarSink
	eng       *Engine
}

// HintSink receives hints posted by the hint_exit built-in from inside
// task scripts.
type HintSink func(unit, task, message string) error

// VarSink receives task vars posted by the set_var built-in from inside
// task scripts.
type VarSink func(unit, name, value string) error

// ServeCheckAPI starts the unix-socket listener at sockPath.
func ServeCheckAPI(sockPath, shellUser string, watcher *ExecWatcher, eng *Engine) error {
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
	api := &checkAPI{
		watcher:   watcher,
		lines:     eng.Lines,
		shellUser: shellUser,
		shellUID:  uid,
		hintSink:  eng.PublishHint,
		varSink:   eng.SetVar,
		eng:       eng,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/shells", api.handleShells)
	mux.HandleFunc("/hint", api.handleHint)
	mux.HandleFunc("/vars", api.handleSetVar)
	mux.HandleFunc("/exec/seq", api.handleSeq)
	mux.HandleFunc("/events/seq", api.handleEventsSeq)
	mux.HandleFunc("/exec/wait", api.handleExecWait)
	mux.HandleFunc("/exec/snapshot", api.handleSnapshot)
	mux.HandleFunc("/line/seq", api.handleLineSeq)
	mux.HandleFunc("/line/wait", api.handleLineWait)
	mux.HandleFunc("/line/snapshot", api.handleLineSnapshot)
	mux.HandleFunc("/units/watch", api.handleUnitsWatch)
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

// handleEventsSeq serves the event_seq built-in: the current value of the
// clock every event ring stamps from. Any event observed from now on -
// exec or typed line - gets a higher number, so the value is a mark a
// later check can pass as --after.
func (a *checkAPI) handleEventsSeq(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, map[string]uint64{"seq": eventClock.Load()})
}

// ExecWaitRequest asks the daemon to block until a matching exec event.
type ExecWaitRequest struct {
	After      uint64  `json:"after"`      // only events with Seq > After
	Regex      string  `json:"regex"`      // matched against the joined argv
	Argc       int     `json:"argc"`       // >0: argv must have exactly this many elements
	Latest     bool    `json:"latest"`     // prefer the newest buffered match over the oldest
	Cwd        string  `json:"cwd"`        // if set: the command must have run from this directory (path or regex, as wait_cwd)
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
		// Argv-count matching: the joined-argv regex cannot tell a quoted
		// single argument from the same text split into several arguments
		// (both join to the same string) - the count can.
		if req.Argc > 0 && len(ev.Argv) != req.Argc {
			return false
		}
		// Working-directory scoping: the event horizon is per unit, so a
		// task gated on "the shell stands in X" would otherwise be
		// satisfied by a matching command run anywhere before the
		// student moved. An unknown cwd (the process was gone before the
		// link could be read) never satisfies the filter - the point of
		// the filter is proof of location.
		if req.Cwd != "" && (ev.Cwd == "" || !checkclient.PathMatch(req.Cwd, ev.Cwd)) {
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
	wait := a.watcher.WaitMatch
	if req.Latest {
		wait = a.watcher.WaitMatchLatest
	}
	ev, ok := wait(r.Context(), req.After, time.Now().Add(timeout), match)
	writeJSON(w, ExecWaitResponse{Matched: ok, Event: ev})
}

// LineWaitRequest asks the daemon to block until an interactive shell of
// the observed user reads a command line matching the regex.
type LineWaitRequest struct {
	After      uint64  `json:"after"`      // only events with Seq > After
	Regex      string  `json:"regex"`      // matched against the typed line (surrounding whitespace trimmed)
	Latest     bool    `json:"latest"`     // prefer the newest buffered match over the oldest
	TimeoutSec float64 `json:"timeoutSec"` // <=0: practically forever
}

type LineWaitResponse struct {
	Matched bool      `json:"matched"`
	Event   LineEvent `json:"event"`
}

// linesAvailable reports whether the readline watcher is active; when it
// is not, the /line/* endpoints answer 501 so a wait_line in a unit that
// forgot `requires: [readline]` fails loudly instead of hanging.
func (a *checkAPI) linesAvailable(w http.ResponseWriter) bool {
	if a.lines == nil || a.lines.Source == "" {
		http.Error(w, "command line watching is unavailable on this host (the unit should declare requires: [readline])", 501)
		return false
	}
	return true
}

func (a *checkAPI) handleLineSeq(w http.ResponseWriter, r *http.Request) {
	if !a.linesAvailable(w) {
		return
	}
	writeJSON(w, map[string]uint64{"seq": a.lines.Seq()})
}

func (a *checkAPI) handleLineSnapshot(w http.ResponseWriter, r *http.Request) {
	if !a.linesAvailable(w) {
		return
	}
	writeJSON(w, a.lines.Snapshot(0, 200))
}

func (a *checkAPI) handleLineWait(w http.ResponseWriter, r *http.Request) {
	if !a.linesAvailable(w) {
		return
	}
	var req LineWaitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	timeout := time.Duration(req.TimeoutSec * float64(time.Second))
	if timeout <= 0 || timeout > time.Hour {
		timeout = time.Hour
	}
	re, err := regexp.Compile(req.Regex)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	match := func(ev LineEvent) bool {
		// Same student-activity scoping as exec events: only a CONFIRMED
		// foreign uid or missing tty disqualifies. The daemon's own task
		// scripts never call readline (non-interactive bash), so there is
		// no self-match to guard against here anyway.
		if ev.UID != -1 && ev.UID != a.shellUID {
			return false
		}
		if ev.TTYNr == 0 {
			return false
		}
		return re.MatchString(strings.TrimSpace(ev.Line))
	}
	wait := a.lines.WaitMatch
	if req.Latest {
		wait = a.lines.WaitMatchLatest
	}
	ev, ok := wait(r.Context(), req.After, time.Now().Add(timeout), match)
	writeJSON(w, LineWaitResponse{Matched: ok, Event: ev})
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

// unitsSnapshot is the "snapshot" SSE event payload.
type unitsSnapshot struct {
	Path  string            `json:"path"`
	Units map[string]string `json:"units"`
}

// unitStatus is the "unit" SSE event payload.
type unitStatus struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

// handleUnitsWatch streams unit status changes as Server-Sent Events, for
// supervisors (the iximiuz Labs examiner) that need to learn when
// learning-path units complete without polling.
//
// On connect it sends a `snapshot` event with every unit's current status,
// then one `unit` event per subsequent change. A client that reconnects
// resyncs from the fresh snapshot, so no event is ever permanently missed
// even across a dropped connection.
func (a *checkAPI) handleUnitsWatch(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", 500)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Subscribe before taking the snapshot so no unit change in between is
	// lost.
	events, unsubscribe := a.eng.Bus.Subscribe()
	defer unsubscribe()

	snapshot, err := json.Marshal(unitsSnapshot{Path: a.eng.Path.ID, Units: a.eng.UnitStatuses()})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if _, err := fmt.Fprintf(w, "event: snapshot\ndata: %s\n\n", snapshot); err != nil {
		return
	}
	flusher.Flush()

	keepalive := time.NewTicker(15 * time.Second)
	defer keepalive.Stop()
	for {
		select {
		case ev, ok := <-events:
			if !ok {
				return
			}
			if ev.Type != "unit" {
				continue
			}
			ue, ok := ev.Data.(UnitEvent)
			if !ok {
				continue
			}
			data, err := json.Marshal(unitStatus{ID: ue.Unit, Status: ue.Status})
			if err != nil {
				continue
			}
			if _, err := fmt.Fprintf(w, "event: unit\ndata: %s\n\n", data); err != nil {
				return
			}
			flusher.Flush()
		case <-keepalive.C:
			if _, err := fmt.Fprint(w, ": keepalive\n\n"); err != nil {
				return
			}
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
