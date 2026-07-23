// Package engine implements the validation engine: it activates units, runs
// their init and verification tasks, tracks rich task statuses, and
// publishes events. It never touches markdown or rendering.
package engine

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/iximiuz/labs-content/tools/shellgym/internal/bus"
	"github.com/iximiuz/labs-content/tools/shellgym/internal/content"
	"github.com/iximiuz/labs-content/tools/shellgym/internal/state"
)

// TaskStatus values (in-memory + persisted as strings).
const (
	StatusPending     = "pending"     // waiting for deps
	StatusRunning     = "running"     // supervisor active, condition not met yet
	StatusSatisfied   = "satisfied"   // level: currently true
	StatusUnsatisfied = "unsatisfied" // level: currently false
	StatusCompleted   = "completed"   // edge: terminal / level: frozen at unit completion
)

// Options tunes the engine.
type Options struct {
	ChecksDir     string
	SockPath      string
	EdgeTimeout   time.Duration // per-attempt timeout for edge tasks
	LevelInterval time.Duration // poll interval for level tasks
	LevelTimeout  time.Duration // per-poll timeout for level tasks
	HintInterval  time.Duration // min gap between dynamic hint refreshes
	RestartDelay  time.Duration // delay before restarting a failed edge attempt
}

func (o *Options) defaults() {
	if o.EdgeTimeout == 0 {
		o.EdgeTimeout = 30 * time.Second
	}
	if o.LevelInterval == 0 {
		o.LevelInterval = 1 * time.Second
	}
	if o.LevelTimeout == 0 {
		o.LevelTimeout = 10 * time.Second
	}
	if o.HintInterval == 0 {
		o.HintInterval = 10 * time.Second
	}
	if o.RestartDelay == 0 {
		o.RestartDelay = 500 * time.Millisecond
	}
}

// TaskEvent is published on every task status change.
type TaskEvent struct {
	Unit   string `json:"unit"`
	Task   string `json:"task"`
	Status string `json:"status"`
}

// HintEvent carries a dynamic hint update.
type HintEvent struct {
	Unit string `json:"unit"`
	Task string `json:"task"`
	Hint string `json:"hint"`
}

// UnitEvent is published when a unit's status changes.
type UnitEvent struct {
	Unit   string `json:"unit"`
	Status string `json:"status"`
}

// InitEvent reports init task progress/failures (author-facing).
type InitEvent struct {
	Unit  string `json:"unit"`
	Task  string `json:"task"`
	OK    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
}

// RunEvent streams task run records for the debug panel.
type RunEvent struct {
	Unit string        `json:"unit"`
	Task string        `json:"task"`
	Run  state.TaskRun `json:"run"`
}

// Engine drives one learning path.
type Engine struct {
	Path    *content.Path
	Store   *state.Store
	Bus     *bus.Bus
	Watcher *ExecWatcher
	Opts    Options

	runner *scriptRunner

	mu         sync.Mutex
	activeUnit string
	cancel     context.CancelFunc
	statuses   map[string]string // task name -> status (active unit only)
	sinceSeq   uint64
	wg         sync.WaitGroup
}

func New(p *content.Path, st *state.Store, b *bus.Bus, w *ExecWatcher, opts Options) *Engine {
	opts.defaults()
	return &Engine{
		Path:    p,
		Store:   st,
		Bus:     b,
		Watcher: w,
		Opts:    opts,
		runner:  &scriptRunner{checksDir: opts.ChecksDir, sockPath: opts.SockPath},
	}
}

// Resume re-activates the persisted current unit (after daemon restart).
func (e *Engine) Resume() {
	var current string
	e.Store.View(func(d *state.Data) { current = d.CurrentUnit })
	if current == "" {
		return
	}
	if u := e.Path.Unit(current); u != nil {
		var completed bool
		e.Store.View(func(d *state.Data) {
			completed = d.Unit(current).Status == state.UnitCompleted
		})
		if !completed {
			if err := e.ActivateUnit(current); err != nil {
				log.Printf("engine: resume %s: %v", current, err)
			}
		}
	}
}

// ErrUnitLocked is returned when a locked unit is activated: a unit whose
// `needs:` dependencies are not all completed cannot start, because its
// scene builds on state those units leave behind.
var ErrUnitLocked = errors.New("unit is locked")

// UnitLocked reports whether the unit's `needs:` dependencies are not all
// completed yet. Completed units are never locked, so solved reps stay
// browsable.
func (e *Engine) UnitLocked(id string) bool {
	var locked bool
	e.Store.View(func(d *state.Data) {
		locked = UnitLockedIn(e.Path, d, id)
	})
	return locked
}

// UnitLockedIn is UnitLocked against an already-locked view of the state
// (for callers inside their own Store.View/Update).
func UnitLockedIn(p *content.Path, d *state.Data, id string) bool {
	u := p.Unit(id)
	if u == nil {
		return false
	}
	if us, ok := d.Units[id]; ok && us.Status == state.UnitCompleted {
		return false
	}
	for _, need := range u.Front.Needs {
		if ds, ok := d.Units[u.ModuleID+"/"+need]; !ok || ds.Status != state.UnitCompleted {
			return true
		}
	}
	return false
}

// ActivateUnit makes the unit current: resolves vars (once), runs init tasks
// (once), and starts task supervisors. Completed units only become "current"
// for viewing; their tasks are not restarted. Locked units (see UnitLocked)
// are rejected with ErrUnitLocked so a unit never arms on a half-built scene.
func (e *Engine) ActivateUnit(id string) error {
	u := e.Path.Unit(id)
	if u == nil {
		return fmt.Errorf("unknown unit %q", id)
	}
	if e.UnitLocked(id) {
		return fmt.Errorf("unit %q: %w - it builds on units that are not solved yet", id, ErrUnitLocked)
	}

	e.mu.Lock()
	if e.activeUnit == id {
		e.mu.Unlock()
		return nil
	}
	e.stopActiveLocked()
	e.activeUnit = id
	e.statuses = map[string]string{}
	ctx, cancel := context.WithCancel(context.Background())
	e.cancel = cancel
	e.mu.Unlock()

	vars, varsErr := e.EnsureVars(id)
	if varsErr != nil {
		log.Printf("engine: %s: vars: %v", id, varsErr)
		vars = map[string]string{}
	}
	var completed bool
	var initDone bool
	taskStatuses := map[string]string{}
	err := e.Store.Update(func(d *state.Data) {
		d.CurrentUnit = id
		us := d.Unit(id)
		completed = us.Status == state.UnitCompleted
		if completed {
			return
		}
		if us.Status == state.UnitPending {
			us.Status = state.UnitActive
			us.ActivatedAt = time.Now()
		}
		initDone = us.InitDone
		for name, ts := range us.Tasks {
			taskStatuses[name] = ts.Status
		}
	})
	if err != nil {
		return err
	}
	if completed {
		e.Bus.Publish(bus.Event{Type: "unit", Data: UnitEvent{Unit: id, Status: "completed"}})
		return nil
	}

	e.mu.Lock()
	e.sinceSeq = 0
	if e.Watcher != nil {
		e.sinceSeq = e.Watcher.Seq()
	}
	for _, t := range u.Tasks {
		if taskStatuses[t.Name] == StatusCompleted {
			e.statuses[t.Name] = StatusCompleted
		} else {
			e.statuses[t.Name] = StatusPending
		}
	}
	e.mu.Unlock()

	e.wg.Add(1)
	go func() {
		defer e.wg.Done()
		if !initDone {
			if !e.runInit(ctx, u, vars) {
				// Init failed: tasks are not started (they'd misfire on a
				// half-built environment). Clear the active marker so a
				// later activate retries the init from scratch.
				e.mu.Lock()
				if e.activeUnit == u.ID {
					e.activeUnit = ""
				}
				e.mu.Unlock()
				return
			}
		}
		for _, t := range u.Tasks {
			if e.status(t.Name) == StatusCompleted {
				e.publishTask(id, t.Name, StatusCompleted)
				continue
			}
			// Announce the fresh slate: after a reset, a task box that was
			// green in some connected UI gets no event until its supervisor
			// runs - which for needs-gated tasks may be never.
			e.publishTask(id, t.Name, StatusPending)
			e.wg.Add(1)
			go func(t *content.Task) {
				defer e.wg.Done()
				e.superviseTask(ctx, u, t, vars)
			}(t)
		}
		e.maybeCompleteUnit(u)
	}()
	e.Bus.Publish(bus.Event{Type: "unit", Data: UnitEvent{Unit: id, Status: "active"}})
	return nil
}

// EnsureVars returns the unit's resolved vars, resolving and persisting
// them on first access (so a unit renders identically before and after
// activation). `from:` references recursively resolve the referenced
// unit's vars first, so dependent units share randomized state.
func (e *Engine) EnsureVars(id string) (map[string]string, error) {
	u := e.Path.Unit(id)
	if u == nil {
		return nil, fmt.Errorf("unknown unit %q", id)
	}
	var existing map[string]string
	e.Store.View(func(d *state.Data) {
		if us, ok := d.Units[id]; ok {
			existing = us.Vars
		}
	})
	if existing != nil {
		return existing, nil
	}
	lookup := func(unitName, varName string) (string, error) {
		depVars, err := e.EnsureVars(u.ModuleID + "/" + unitName)
		if err != nil {
			return "", err
		}
		v, ok := depVars[varName]
		if !ok {
			return "", fmt.Errorf("unit %q has no var %q", unitName, varName)
		}
		return v, nil
	}
	resolved, err := content.ResolveVars(u.Front.Vars, lookup)
	if err != nil {
		return nil, err
	}
	var vars map[string]string
	err = e.Store.Update(func(d *state.Data) {
		us := d.Unit(id)
		if us.Vars == nil {
			us.Vars = resolved
		}
		vars = us.Vars
	})
	return vars, err
}

// ResetUnit forgets unit progress and re-activates it.
func (e *Engine) ResetUnit(id string) error {
	u := e.Path.Unit(id)
	if u == nil {
		return fmt.Errorf("unknown unit %q", id)
	}
	e.mu.Lock()
	if e.activeUnit == id {
		e.stopActiveLocked()
	}
	e.mu.Unlock()
	if err := e.Store.DropUnit(id); err != nil {
		return err
	}
	return e.ActivateUnit(id)
}

// MarkModuleSeen records a module intro scene as viewed.
func (e *Engine) MarkModuleSeen(id string) error {
	return e.Store.Update(func(d *state.Data) {
		if d.SeenModules == nil {
			d.SeenModules = map[string]bool{}
		}
		d.SeenModules[id] = true
	})
}

func (e *Engine) Shutdown() {
	e.mu.Lock()
	e.stopActiveLocked()
	e.mu.Unlock()
	e.wg.Wait()
}

func (e *Engine) stopActiveLocked() {
	if e.cancel != nil {
		e.cancel()
		e.cancel = nil
	}
	e.activeUnit = ""
}

func (e *Engine) status(task string) string {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.statuses[task]
}

func (e *Engine) setStatus(unit, task, status string) {
	e.mu.Lock()
	if e.activeUnit != unit {
		e.mu.Unlock()
		return
	}
	e.statuses[task] = status
	e.mu.Unlock()
	_ = e.Store.Update(func(d *state.Data) {
		d.Unit(unit).Task(task).Status = status
	})
	e.publishTask(unit, task, status)
}

func (e *Engine) publishTask(unit, task, status string) {
	e.Bus.Publish(bus.Event{Type: "task", Data: TaskEvent{Unit: unit, Task: task, Status: status}})
}

// taskEnv builds the environment for task scripts.
func (e *Engine) taskEnv(u *content.Unit, taskName string, vars map[string]string) map[string]string {
	env := map[string]string{
		"GYM_UNIT":      u.ID,
		"GYM_TASK":      taskName,
		"GYM_USER":      e.Path.ShellUser,
		"GYM_SINCE_SEQ": fmt.Sprintf("%d", e.currentSinceSeq()),
	}
	for k, v := range vars {
		env[k] = v
	}
	return env
}

func (e *Engine) currentSinceSeq() uint64 {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.sinceSeq
}

// runInit executes init tasks sequentially. Returns false if any failed.
func (e *Engine) runInit(ctx context.Context, u *content.Unit, vars map[string]string) bool {
	allOK := true
	for _, it := range u.Front.Init {
		name := it.Name
		if name == "" {
			name = "init"
		}
		res := e.runner.Run(ctx, it.Run, e.taskEnv(u, name, vars), 60*time.Second)
		_ = e.Store.AddRun(u.ID, "init:"+name, toTaskRun(res, "init"))
		ev := InitEvent{Unit: u.ID, Task: name, OK: res.ExitCode == 0}
		if res.ExitCode != 0 {
			allOK = false
			ev.Error = strings.TrimSpace(res.Stderr + "\n" + res.Stdout)
			log.Printf("engine: %s: init %s failed (exit %d): %s", u.ID, name, res.ExitCode, ev.Error)
		}
		e.Bus.Publish(bus.Event{Type: "init", Data: ev})
		e.Bus.Publish(bus.Event{Type: "run", Data: RunEvent{Unit: u.ID, Task: "init:" + name, Run: toTaskRun(res, "init")}})
		if ctx.Err() != nil {
			return false
		}
	}
	if allOK {
		_ = e.Store.Update(func(d *state.Data) {
			d.Unit(u.ID).InitDone = true
		})
	}
	return allOK
}

// depsReady reports whether all of t's needs are in an OK state for t to run.
func (e *Engine) depsReady(t *content.Task) bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	for _, need := range t.Needs {
		st := e.statuses[need]
		if st != StatusCompleted && st != StatusSatisfied {
			return false
		}
	}
	return true
}

func (e *Engine) superviseTask(ctx context.Context, u *content.Unit, t *content.Task, vars map[string]string) {
	// Wait for dependencies.
	for !e.depsReady(t) {
		select {
		case <-ctx.Done():
			return
		case <-time.After(300 * time.Millisecond):
		}
	}
	if t.Mode == content.ModeEdge {
		e.superviseEdge(ctx, u, t, vars)
	} else {
		e.superviseLevel(ctx, u, t, vars)
	}
}

func (e *Engine) superviseEdge(ctx context.Context, u *content.Unit, t *content.Task, vars map[string]string) {
	e.setStatus(u.ID, t.Name, StatusRunning)
	timeout := e.Opts.EdgeTimeout
	if t.Timeout > 0 {
		timeout = time.Duration(t.Timeout) * time.Second
	}
	var lastHint time.Time
	for {
		if ctx.Err() != nil {
			return
		}
		res := e.runner.Run(ctx, t.Check, e.taskEnv(u, t.Name, vars), timeout)
		e.recordRun(u.ID, t.Name, res, "run")
		if ctx.Err() != nil {
			return
		}
		if res.ExitCode == 0 {
			e.setStatus(u.ID, t.Name, StatusCompleted)
			e.maybeCompleteUnit(u)
			return
		}
		// Attempt failed (usually: wait_* timed out). Maybe refresh the hint.
		if t.Hint != "" && time.Since(lastHint) >= e.Opts.HintInterval {
			lastHint = time.Now()
			e.runHint(ctx, u, t, vars, res)
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(e.Opts.RestartDelay):
		}
	}
}

func (e *Engine) superviseLevel(ctx context.Context, u *content.Unit, t *content.Task, vars map[string]string) {
	timeout := e.Opts.LevelTimeout
	if t.Timeout > 0 {
		timeout = time.Duration(t.Timeout) * time.Second
	}
	last := ""
	var unsatisfiedSince time.Time
	var lastHint time.Time
	for {
		if ctx.Err() != nil {
			return
		}
		if !e.depsReady(t) {
			select {
			case <-ctx.Done():
				return
			case <-time.After(e.Opts.LevelInterval):
			}
			continue
		}
		res := e.runner.Run(ctx, t.Check, e.taskEnv(u, t.Name, vars), timeout)
		if ctx.Err() != nil {
			return
		}
		status := StatusUnsatisfied
		if res.ExitCode == 0 {
			status = StatusSatisfied
		}
		if status != last {
			e.recordRun(u.ID, t.Name, res, "run")
			e.setStatus(u.ID, t.Name, status)
			last = status
			if status == StatusSatisfied {
				unsatisfiedSince = time.Time{}
				e.maybeCompleteUnit(u)
				if e.unitDone(u) {
					return
				}
			} else {
				unsatisfiedSince = time.Now()
			}
		}
		if status == StatusUnsatisfied && t.Hint != "" &&
			!unsatisfiedSince.IsZero() && time.Since(unsatisfiedSince) > 5*time.Second &&
			time.Since(lastHint) >= e.Opts.HintInterval {
			lastHint = time.Now()
			e.runHint(ctx, u, t, vars, res)
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(e.Opts.LevelInterval):
		}
	}
}

func (e *Engine) runHint(ctx context.Context, u *content.Unit, t *content.Task, vars map[string]string, lastRes RunResult) {
	env := e.taskEnv(u, t.Name, vars)
	env["GYM_TASK_EXIT"] = fmt.Sprintf("%d", lastRes.ExitCode)
	env["GYM_TASK_STDOUT"] = clip(lastRes.Stdout, 1024)
	env["GYM_TASK_STDERR"] = clip(lastRes.Stderr, 1024)
	res := e.runner.Run(ctx, t.Hint, env, 10*time.Second)
	e.recordRun(u.ID, t.Name, res, "hint")
	hint := strings.TrimSpace(res.Stdout)
	if res.ExitCode != 0 || hint == "" {
		return
	}
	_ = e.PublishHint(u.ID, t.Name, hint)
}

// PublishHint stores and broadcasts a hint for a task (used both by hint:
// block runs and by the hint_exit built-in via the check API).
func (e *Engine) PublishHint(unit, task, hint string) error {
	u := e.Path.Unit(unit)
	if u == nil {
		return fmt.Errorf("unknown unit %q", unit)
	}
	if _, ok := u.Front.Tasks[task]; !ok {
		return fmt.Errorf("unit %s has no task %q", unit, task)
	}
	_ = e.Store.Update(func(d *state.Data) {
		d.Unit(unit).Task(task).Hint = hint
	})
	e.Bus.Publish(bus.Event{Type: "hint", Data: HintEvent{Unit: unit, Task: task, Hint: hint}})
	return nil
}

func clip(s string, n int) string {
	if len(s) > n {
		return s[:n]
	}
	return s
}

func (e *Engine) recordRun(unit, task string, res RunResult, kind string) {
	run := toTaskRun(res, kind)
	_ = e.Store.AddRun(unit, task, run)
	e.Bus.Publish(bus.Event{Type: "run", Data: RunEvent{Unit: unit, Task: task, Run: run}})
}

func (e *Engine) unitDone(u *content.Unit) bool {
	var done bool
	e.Store.View(func(d *state.Data) {
		done = d.Unit(u.ID).Status == state.UnitCompleted
	})
	return done
}

// maybeCompleteUnit completes the unit when every edge task is completed and
// every level task is currently satisfied.
func (e *Engine) maybeCompleteUnit(u *content.Unit) {
	e.mu.Lock()
	if e.activeUnit != u.ID {
		e.mu.Unlock()
		return
	}
	if len(u.Tasks) == 0 {
		e.mu.Unlock()
		return
	}
	for _, t := range u.Tasks {
		st := e.statuses[t.Name]
		if t.Mode == content.ModeEdge && st != StatusCompleted {
			e.mu.Unlock()
			return
		}
		if t.Mode == content.ModeLevel && st != StatusSatisfied && st != StatusCompleted {
			e.mu.Unlock()
			return
		}
	}
	for _, t := range u.Tasks {
		e.statuses[t.Name] = StatusCompleted
	}
	cancel := e.cancel
	e.cancel = nil
	e.mu.Unlock()
	if cancel != nil {
		cancel() // stop supervisors; completion is terminal
	}

	_ = e.Store.Update(func(d *state.Data) {
		us := d.Unit(u.ID)
		us.Status = state.UnitCompleted
		us.CompletedAt = time.Now()
		for _, t := range u.Tasks {
			us.Task(t.Name).Status = StatusCompleted
		}
	})
	for _, t := range u.Tasks {
		e.publishTask(u.ID, t.Name, StatusCompleted)
	}
	e.Bus.Publish(bus.Event{Type: "unit", Data: UnitEvent{Unit: u.ID, Status: "completed"}})
	log.Printf("engine: unit %s completed", u.ID)
}
