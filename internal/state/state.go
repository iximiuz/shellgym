// Package state persists learner progress and task run history.
//
// Layout: one directory per learning path -
//
//	<root>/<path-id>/progress.json             small, atomically-replaced doc
//	<root>/<path-id>/runs/<unit>/<task>.jsonl  bounded append-only run records
//
// progress.json stays tiny (statuses, vars, timestamps); potentially large
// task run records (stdout/stderr) live in per-task JSONL files that are
// appended to and compacted to the last maxRunsPerTask lines.
package state

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// UnitStatus is the persistent per-unit progress.
type UnitStatus string

const (
	UnitPending   UnitStatus = "pending"
	UnitActive    UnitStatus = "active"
	UnitCompleted UnitStatus = "completed"
)

// TaskRun records a single execution of a task script.
type TaskRun struct {
	StartedAt time.Time `json:"startedAt"`
	Duration  float64   `json:"durationSec"`
	ExitCode  int       `json:"exitCode"`
	Stdout    string    `json:"stdout"`
	Stderr    string    `json:"stderr"`
	Kind      string    `json:"kind"` // "run" | "hint" | "init"
	TimedOut  bool      `json:"timedOut,omitempty"`
}

// TaskState is the persistent per-task record (kept in progress.json;
// run history lives in the task's JSONL file).
type TaskState struct {
	Status string `json:"status"` // pending|running|satisfied|unsatisfied|completed
	Hint   string `json:"hint,omitempty"`
}

// UnitState is the persistent per-unit record.
type UnitState struct {
	Status      UnitStatus            `json:"status"`
	Vars        map[string]string     `json:"vars,omitempty"`
	InitDone    bool                  `json:"initDone"`
	Tasks       map[string]*TaskState `json:"tasks,omitempty"`
	ActivatedAt time.Time             `json:"activatedAt,omitempty"`
	CompletedAt time.Time             `json:"completedAt,omitempty"`
}

// Data is the progress.json document.
type Data struct {
	PathID      string                `json:"pathId"`
	CurrentUnit string                `json:"currentUnit"`
	SeenModules map[string]bool       `json:"seenModules,omitempty"`
	Units       map[string]*UnitState `json:"units"`
	UpdatedAt   time.Time             `json:"updatedAt"`
}

// Store owns one path's state directory.
type Store struct {
	mu   sync.Mutex
	dir  string // <root>/<path-id>
	data *Data
}

const maxRunsPerTask = 20
const maxStreamLen = 4096

func Open(root, pathID string) (*Store, error) {
	dir := filepath.Join(root, pathID)
	if err := os.MkdirAll(filepath.Join(dir, "runs"), 0o755); err != nil {
		return nil, err
	}
	s := &Store{dir: dir}
	progressPath := filepath.Join(dir, "progress.json")
	raw, err := os.ReadFile(progressPath)
	switch {
	case err == nil:
		var d Data
		if jErr := json.Unmarshal(raw, &d); jErr != nil {
			// Corrupt state: keep a backup, start fresh.
			_ = os.Rename(progressPath, progressPath+".corrupt")
			s.data = newData(pathID)
		} else {
			s.data = &d
		}
	case os.IsNotExist(err):
		s.data = newData(pathID)
	default:
		return nil, err
	}
	return s, nil
}

func newData(pathID string) *Data {
	return &Data{
		PathID:      pathID,
		Units:       map[string]*UnitState{},
		SeenModules: map[string]bool{},
	}
}

// Update mutates the document under lock and persists it atomically.
func (s *Store) Update(fn func(*Data)) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	fn(s.data)
	s.data.UpdatedAt = time.Now()
	return s.flushLocked()
}

// View reads the document under lock.
func (s *Store) View(fn func(*Data)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	fn(s.data)
}

// Unit returns the unit state, creating it if needed. Use inside
// Update/View callbacks only.
func (d *Data) Unit(id string) *UnitState {
	u, ok := d.Units[id]
	if !ok {
		u = &UnitState{Status: UnitPending, Tasks: map[string]*TaskState{}}
		d.Units[id] = u
	}
	if u.Tasks == nil {
		u.Tasks = map[string]*TaskState{}
	}
	return u
}

// Task returns the task state within a unit, creating it if needed.
func (u *UnitState) Task(name string) *TaskState {
	t, ok := u.Tasks[name]
	if !ok {
		t = &TaskState{Status: "pending"}
		u.Tasks[name] = t
	}
	return t
}

// --- run history (JSONL side files) ----------------------------------------

// slug flattens a unit id ("module/unit") into a directory name.
func slug(unitID string) string {
	return strings.ReplaceAll(unitID, "/", "__")
}

func (s *Store) runsFile(unitID, task string) string {
	return filepath.Join(s.dir, "runs", slug(unitID), task+".jsonl")
}

// AddRun appends a run record to the task's JSONL file, truncating streams
// and compacting the file when it grows past maxRunsPerTask records.
func (s *Store) AddRun(unitID, task string, r TaskRun) error {
	r.Stdout = truncate(r.Stdout)
	r.Stderr = truncate(r.Stderr)
	s.mu.Lock()
	defer s.mu.Unlock()
	path := s.runsFile(unitID, task)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	line, err := json.Marshal(r)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	_, werr := f.Write(append(line, '\n'))
	cerr := f.Close()
	if werr != nil {
		return werr
	}
	if cerr != nil {
		return cerr
	}
	return s.compactLocked(path)
}

// compactLocked keeps only the last maxRunsPerTask lines of a JSONL file.
func (s *Store) compactLocked(path string) error {
	lines, err := readLines(path)
	if err != nil || len(lines) <= maxRunsPerTask {
		return err
	}
	lines = lines[len(lines)-maxRunsPerTask:]
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, []byte(strings.Join(lines, "\n")+"\n"), 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

// Runs returns the recorded run history of a task, oldest first.
func (s *Store) Runs(unitID, task string) []TaskRun {
	s.mu.Lock()
	lines, _ := readLines(s.runsFile(unitID, task))
	s.mu.Unlock()
	out := make([]TaskRun, 0, len(lines))
	for _, l := range lines {
		var r TaskRun
		if json.Unmarshal([]byte(l), &r) == nil {
			out = append(out, r)
		}
	}
	return out
}

// RunTasks lists task names that have run history for a unit.
func (s *Store) RunTasks(unitID string) []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	entries, err := os.ReadDir(filepath.Join(s.dir, "runs", slug(unitID)))
	if err != nil {
		return nil
	}
	var out []string
	for _, e := range entries {
		if name, ok := strings.CutSuffix(e.Name(), ".jsonl"); ok {
			out = append(out, name)
		}
	}
	return out
}

// DropUnit forgets a unit entirely (progress + run history).
func (s *Store) DropUnit(unitID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data.Units, unitID)
	if err := os.RemoveAll(filepath.Join(s.dir, "runs", slug(unitID))); err != nil {
		return err
	}
	return s.flushLocked()
}

func readLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()
	var out []string
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for sc.Scan() {
		if line := strings.TrimSpace(sc.Text()); line != "" {
			out = append(out, line)
		}
	}
	return out, sc.Err()
}

func truncate(s string) string {
	if len(s) <= maxStreamLen {
		return s
	}
	return s[:maxStreamLen] + fmt.Sprintf("\n... (%d bytes truncated)", len(s)-maxStreamLen)
}

func (s *Store) flushLocked() error {
	raw, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}
	tmp := filepath.Join(s.dir, "progress.json.tmp")
	if err := os.WriteFile(tmp, raw, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, filepath.Join(s.dir, "progress.json"))
}
