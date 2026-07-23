package engine

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"testing"
	"time"

	"github.com/creack/pty"

	"github.com/iximiuz/labs-content/tools/shellgym/internal/bus"
	"github.com/iximiuz/labs-content/tools/shellgym/internal/content"
	"github.com/iximiuz/labs-content/tools/shellgym/internal/state"
)

// studentShell is a real interactive bash on a PTY, imitating the learner.
type studentShell struct {
	t   *testing.T
	pty *os.File
	cmd *exec.Cmd
}

func newStudentShell(t *testing.T) *studentShell {
	t.Helper()
	cmd := exec.Command("bash", "--norc", "-i")
	cmd.Env = append(os.Environ(), "PS1=$ ", "TERM=dumb")
	f, err := pty.Start(cmd)
	if err != nil {
		t.Skipf("cannot start pty shell: %v", err)
	}
	s := &studentShell{t: t, pty: f, cmd: cmd}
	t.Cleanup(func() {
		_ = cmd.Process.Kill()
		_, _ = cmd.Process.Wait()
		_ = f.Close()
	})
	go func() { // drain output so the shell never blocks on a full pty buffer
		buf := make([]byte, 4096)
		for {
			if _, err := f.Read(buf); err != nil {
				return
			}
		}
	}()
	time.Sleep(300 * time.Millisecond)
	return s
}

func (s *studentShell) Type(cmdline string) {
	s.t.Helper()
	if _, err := s.pty.WriteString(cmdline + "\n"); err != nil {
		s.t.Fatalf("type %q: %v", cmdline, err)
	}
	time.Sleep(250 * time.Millisecond)
}

// testEnv wires a full engine around a scratch content dir.
type testEnv struct {
	t      *testing.T
	eng    *Engine
	events <-chan bus.Event
	dir    string
}

func newTestEnv(t *testing.T, units map[string]string) *testEnv {
	t.Helper()
	me, err := user.Current()
	if err != nil {
		t.Skip("no current user")
	}

	dir := t.TempDir()
	mustWrite(t, filepath.Join(dir, "content", "path.yaml"),
		"id: itest\ntitle: itest\nshellUser: "+me.Username+"\n")
	for rel, data := range units {
		mustWrite(t, filepath.Join(dir, "content", rel), data)
	}
	path, err := content.Load(filepath.Join(dir, "content"), "ubuntu", []string{"debian"}, nil)
	if err != nil {
		t.Fatal(err)
	}

	st, err := state.Open(filepath.Join(dir, "state"), path.ID)
	if err != nil {
		t.Fatal(err)
	}

	selfExe := buildSelf(t)
	checksDir := filepath.Join(dir, "checks")
	if err := WriteCheckShims(checksDir, selfExe); err != nil {
		t.Fatal(err)
	}

	watcher := NewExecWatcher()
	watcher.Start() // netlink; Source stays "" without CAP_NET_ADMIN
	t.Cleanup(watcher.Close)

	sockPath := filepath.Join(dir, "gym.sock")

	b := bus.New()
	ch, unsub := b.Subscribe()
	t.Cleanup(unsub)

	eng := New(path, st, b, watcher, Options{
		ChecksDir:    checksDir,
		SockPath:     sockPath,
		EdgeTimeout:  5 * time.Second,
		RestartDelay: 200 * time.Millisecond,
		HintInterval: 1 * time.Second,
	})
	t.Cleanup(eng.Shutdown)
	if err := ServeCheckAPI(sockPath, path.ShellUser, watcher, eng.PublishHint); err != nil {
		t.Fatal(err)
	}
	return &testEnv{t: t, eng: eng, events: ch, dir: dir}
}

// buildSelf compiles the shellgym binary once per test run (needed by the
// check shims, which exec `shellgym check ...`).
var builtSelf string

func buildSelf(t *testing.T) string {
	t.Helper()
	if builtSelf != "" {
		return builtSelf
	}
	out := filepath.Join(os.TempDir(), fmt.Sprintf("shellgym-test-%d", os.Getpid()))
	cmd := exec.Command("go", "build", "-o", out, "../../cmd/shellgym")
	cmd.Env = os.Environ()
	if raw, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build self: %v\n%s", err, raw)
	}
	builtSelf = out
	return out
}

func mustWrite(t *testing.T, path, data string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
}

func (te *testEnv) waitTaskStatus(unit, task, wantStatus string, timeout time.Duration) bool {
	te.t.Helper()
	deadline := time.After(timeout)
	for {
		select {
		case ev := <-te.events:
			if d, ok := ev.Data.(TaskEvent); ok && ev.Type == "task" &&
				d.Unit == unit && d.Task == task && d.Status == wantStatus {
				return true
			}
		case <-deadline:
			return false
		}
	}
}

func (te *testEnv) waitEvent(kind, unit, wantStatus string, timeout time.Duration) bool {
	te.t.Helper()
	deadline := time.After(timeout)
	for {
		select {
		case ev := <-te.events:
			switch kind {
			case "unit":
				if d, ok := ev.Data.(UnitEvent); ok && ev.Type == "unit" &&
					d.Unit == unit && d.Status == wantStatus {
					return true
				}
			case "task":
				if d, ok := ev.Data.(TaskEvent); ok && ev.Type == "task" &&
					d.Unit == unit && d.Status == wantStatus {
					return true
				}
			case "hint":
				if _, ok := ev.Data.(HintEvent); ok && ev.Type == "hint" {
					return true
				}
			}
		case <-deadline:
			return false
		}
	}
}

// --- the tests --------------------------------------------------------------

const fileUnit = `---
title: Make a file
init:
  - name: prep
    run: |
      mkdir -p /tmp/shellgym-itest
tasks:
  made:
    check: |
      wait_file /tmp/shellgym-itest/made-it.txt
---
Create the file.

::task{name="made"}
Waiting...
::
`

func TestEngineCompletesUnitOnFileCreation(t *testing.T) {
	te := newTestEnv(t, map[string]string{"010.m/010.make-file/unit.md": fileUnit})
	defer os.RemoveAll("/tmp/shellgym-itest")

	if err := te.eng.ActivateUnit("m/make-file"); err != nil {
		t.Fatal(err)
	}
	sh := newStudentShell(t)
	sh.Type("touch /tmp/shellgym-itest/made-it.txt")

	if !te.waitEvent("unit", "m/make-file", "completed", 15*time.Second) {
		t.Fatal("unit did not complete")
	}
}

const cwdUnit = `---
title: Go there
init:
  - name: prep
    run: |
      mkdir -p /tmp/shellgym-cwd-test
tasks:
  there:
    check: |
      wait_cwd /tmp/shellgym-cwd-test
---
::task{name="there"}
Waiting...
::
`

func TestEngineSeesShellCwd(t *testing.T) {
	te := newTestEnv(t, map[string]string{"010.m/010.go-there/unit.md": cwdUnit})
	defer os.RemoveAll("/tmp/shellgym-cwd-test")

	if err := te.eng.ActivateUnit("m/go-there"); err != nil {
		t.Fatal(err)
	}
	sh := newStudentShell(t)
	sh.Type("cd /tmp/shellgym-cwd-test")

	if !te.waitEvent("unit", "m/go-there", "completed", 15*time.Second) {
		t.Fatal("cwd change not detected")
	}
}

const execUnit = `---
title: Run it
tasks:
  ran:
    check: |
      wait_exec 'shellgym-marker-cmd'
    hint: |
      echo "run the marker command"
---
::task{name="ran"}
Waiting...
::
`

func TestEngineSeesExecAndEmitsHints(t *testing.T) {
	te := newTestEnv(t, map[string]string{"010.m/010.run-it/unit.md": execUnit})
	if te.eng.Watcher.Source != "netlink" {
		t.Skip("exec watching needs the proc connector (CAP_NET_ADMIN); run as root")
	}
	if err := te.eng.ActivateUnit("m/run-it"); err != nil {
		t.Fatal(err)
	}
	// Hint should fire after the first failed attempt (5s edge timeout + 1s
	// hint interval).
	if !te.waitEvent("hint", "", "", 20*time.Second) {
		t.Fatal("no dynamic hint emitted")
	}
	sh := newStudentShell(t)
	// Long-lived so the watcher cannot miss it.
	sh.Type("bash -c 'exec -a shellgym-marker-cmd sleep 3' &")

	if !te.waitEvent("unit", "m/run-it", "completed", 20*time.Second) {
		t.Fatal("exec not detected")
	}
}

const depUnits = `---
title: Two steps
vars:
  STAMP: { value: itest-stamp }
tasks:
  first:
    check: |
      wait_file "/tmp/shellgym-dep/$STAMP-one"
  second:
    needs: [first]
    check: |
      wait_file "/tmp/shellgym-dep/$STAMP-two"
---
::task{name="first"}
one
::
::task{name="second"}
two
::
`

func TestTaskDependencyOrdering(t *testing.T) {
	te := newTestEnv(t, map[string]string{"010.m/010.two-steps/unit.md": depUnits})
	defer os.RemoveAll("/tmp/shellgym-dep")
	_ = os.MkdirAll("/tmp/shellgym-dep", 0o755)

	// Create the SECOND file first: task "second" must not complete until
	// "first" does, even though its condition is already true.
	mustWrite(t, "/tmp/shellgym-dep/itest-stamp-two", "x")
	if err := te.eng.ActivateUnit("m/two-steps"); err != nil {
		t.Fatal(err)
	}
	if te.waitEvent("unit", "m/two-steps", "completed", 3*time.Second) {
		t.Fatal("unit completed before dependency was satisfied")
	}
	mustWrite(t, "/tmp/shellgym-dep/itest-stamp-one", "x")
	if !te.waitEvent("unit", "m/two-steps", "completed", 15*time.Second) {
		t.Fatal("unit did not complete after both files exist")
	}
}

// Mirrors 080.final-lap/010.field-kit: edge tasks chained with needs.
// Solving ONLY the first task must complete it individually (its own task
// event, i.e. the box turns green right away) - not when the whole unit
// is solved.
const chainUnit = `---
title: Chain
tasks:
  first:
    check: |
      wait_file /tmp/shellgym-chain/one
  second:
    needs: [first]
    check: |
      wait_file /tmp/shellgym-chain/two
  third:
    needs: [second]
    check: |
      wait_file /tmp/shellgym-chain/three
---
::task{name="first"}
1
::
::task{name="second"}
2
::
::task{name="third"}
3
::
`

func TestFirstTaskCompletesIndividually(t *testing.T) {
	te := newTestEnv(t, map[string]string{"010.m/010.chain/unit.md": chainUnit})
	defer os.RemoveAll("/tmp/shellgym-chain")
	_ = os.MkdirAll("/tmp/shellgym-chain", 0o755)

	if err := te.eng.ActivateUnit("m/chain"); err != nil {
		t.Fatal(err)
	}
	// Solve only the first task.
	mustWrite(t, "/tmp/shellgym-chain/one", "x")
	if !te.waitTaskStatus("m/chain", "first", StatusCompleted, 5*time.Second) {
		t.Fatal("task 'first' did not complete individually")
	}
	// Second becomes running (deps met), third still pending.
	if !te.waitTaskStatus("m/chain", "second", StatusRunning, 5*time.Second) {
		t.Fatal("task 'second' did not start after 'first' completed")
	}
	if te.waitEvent("unit", "m/chain", "completed", 2*time.Second) {
		t.Fatal("unit completed with only one task solved")
	}
}

// A (re)activation must announce every task's fresh status: after a reset,
// a task box that was green in some connected UI gets no event until its
// supervisor runs - which for needs-gated tasks may be never.
func TestActivationAnnouncesFreshTaskStatuses(t *testing.T) {
	te := newTestEnv(t, map[string]string{"010.m/010.chain/unit.md": chainUnit})
	defer os.RemoveAll("/tmp/shellgym-chain")
	_ = os.MkdirAll("/tmp/shellgym-chain", 0o755)

	if err := te.eng.ActivateUnit("m/chain"); err != nil {
		t.Fatal(err)
	}
	mustWrite(t, "/tmp/shellgym-chain/one", "x")
	if !te.waitTaskStatus("m/chain", "first", StatusCompleted, 5*time.Second) {
		t.Fatal("task 'first' did not complete")
	}

	_ = os.Remove("/tmp/shellgym-chain/one")
	if err := te.eng.ResetUnit("m/chain"); err != nil {
		t.Fatal(err)
	}
	if !te.waitTaskStatus("m/chain", "third", StatusPending, 5*time.Second) {
		t.Fatal("re-activation did not announce needs-gated task 'third' as pending")
	}
}

func TestStatePersistsAcrossEngineRestart(t *testing.T) {
	units := map[string]string{"010.m/010.make-file/unit.md": fileUnit}
	te := newTestEnv(t, units)
	defer os.RemoveAll("/tmp/shellgym-itest")

	if err := te.eng.ActivateUnit("m/make-file"); err != nil {
		t.Fatal(err)
	}
	mustWrite(t, "/tmp/shellgym-itest/made-it.txt", "x")
	if !te.waitEvent("unit", "m/make-file", "completed", 15*time.Second) {
		t.Fatal("unit did not complete")
	}
	te.eng.Shutdown()

	// Reopen state from the same dir: the unit must still be completed.
	st, err := state.Open(filepath.Join(te.dir, "state"), "itest")
	if err != nil {
		t.Fatal(err)
	}
	st.View(func(d *state.Data) {
		if d.Unit("m/make-file").Status != state.UnitCompleted {
			t.Error("completion not persisted")
		}
		if d.CurrentUnit != "m/make-file" {
			t.Errorf("current unit = %q", d.CurrentUnit)
		}
	})
}

const levelUnit = `---
title: Level check
init:
  - name: prep
    run: |
      mkdir -p /tmp/shellgym-level
tasks:
  edge_part:
    check: |
      wait_file /tmp/shellgym-level/edge.txt
  level_part:
    mode: level
    check: |
      wait_cwd --now /tmp/shellgym-level
---
::task{name="edge_part"}
e
::
::task{name="level_part"}
l
::
`

func TestLevelTaskGatesUnitCompletion(t *testing.T) {
	te := newTestEnv(t, map[string]string{"010.m/010.level/unit.md": levelUnit})
	defer os.RemoveAll("/tmp/shellgym-level")

	if err := te.eng.ActivateUnit("m/level"); err != nil {
		t.Fatal(err)
	}
	sh := newStudentShell(t)
	// Satisfy the edge task while the level condition is FALSE.
	sh.Type("touch /tmp/shellgym-level/edge.txt")
	if te.waitEvent("unit", "m/level", "completed", 4*time.Second) {
		t.Fatal("unit completed although level task unsatisfied")
	}
	// Now make the level condition true.
	sh.Type("cd /tmp/shellgym-level")
	if !te.waitEvent("unit", "m/level", "completed", 15*time.Second) {
		t.Fatal("unit did not complete once level task satisfied")
	}
}

const showHintUnit = `---
title: Show hint
tasks:
  gated:
    check: |
      wait_file --timeout 2 /tmp/shellgym-showhint/ok.txt || hint_exit "Create ok.txt to proceed!"
---
::task{name="gated"}
w
::
`

func TestHintExitBuiltin(t *testing.T) {
	te := newTestEnv(t, map[string]string{"010.m/010.show-hint/unit.md": showHintUnit})
	defer os.RemoveAll("/tmp/shellgym-showhint")
	_ = os.MkdirAll("/tmp/shellgym-showhint", 0o755)

	if err := te.eng.ActivateUnit("m/show-hint"); err != nil {
		t.Fatal(err)
	}
	// The check fails fast (2s timeout) and hint_exit must deliver the
	// message as a hint event for THIS task.
	deadline := time.After(20 * time.Second)
	for {
		select {
		case ev := <-te.events:
			if d, ok := ev.Data.(HintEvent); ok && ev.Type == "hint" {
				if d.Task != "gated" || d.Hint != "Create ok.txt to proceed!" {
					t.Fatalf("wrong hint event: %+v", d)
				}
				// hint persisted too
				var stored string
				te.eng.Store.View(func(dd *state.Data) {
					stored = dd.Unit("m/show-hint").Task("gated").Hint
				})
				if stored != d.Hint {
					t.Fatalf("hint not persisted: %q", stored)
				}
				// finish the unit
				mustWrite(t, "/tmp/shellgym-showhint/ok.txt", "x")
				if !te.waitEvent("unit", "m/show-hint", "completed", 15*time.Second) {
					t.Fatal("unit did not complete after fixing the condition")
				}
				return
			}
		case <-deadline:
			t.Fatal("hint_exit hint event never arrived")
		}
	}
}

const multiHintUnit = `---
title: Two tasks two hints
init:
  - name: prep
    run: |
      mkdir -p /tmp/shellgym-multi
tasks:
  alpha:
    check: |
      wait_file --timeout 2 /tmp/shellgym-multi/alpha.txt || hint_exit "alpha-hint"
  beta:
    check: |
      wait_file --timeout 2 /tmp/shellgym-multi/beta.txt || hint_exit beta "beta-hint"
---
::task{name="alpha"}
a
::
::task{name="beta"}
b
::
`

func TestMultiTaskHintRouting(t *testing.T) {
	te := newTestEnv(t, map[string]string{"010.m/010.multi/unit.md": multiHintUnit})
	defer os.RemoveAll("/tmp/shellgym-multi")

	if err := te.eng.ActivateUnit("m/multi"); err != nil {
		t.Fatal(err)
	}
	// Both tasks fail fast and emit hints; each must land on its own task.
	got := map[string]string{}
	deadline := time.After(25 * time.Second)
	for len(got) < 2 {
		select {
		case ev := <-te.events:
			if d, ok := ev.Data.(HintEvent); ok && ev.Type == "hint" {
				got[d.Task] = d.Hint
			}
		case <-deadline:
			t.Fatalf("hints missing, got %v", got)
		}
	}
	if got["alpha"] != "alpha-hint" || got["beta"] != "beta-hint" {
		t.Fatalf("hints landed on wrong tasks: %v", got)
	}
	mustWrite(t, "/tmp/shellgym-multi/alpha.txt", "x")
	mustWrite(t, "/tmp/shellgym-multi/beta.txt", "x")
	if !te.waitEvent("unit", "m/multi", "completed", 20*time.Second) {
		t.Fatal("unit did not complete")
	}
}

const hintExitTerminatesUnit = `---
title: hint_exit terminates
init:
  - name: prep
    run: |
      mkdir -p /tmp/shellgym-he
tasks:
  gated:
    check: |
      wait_file --timeout 2 /tmp/shellgym-he/present.txt || hint_exit "not yet"
      wait_file_gone /tmp/shellgym-he/present.txt
---
::task{name="gated"}
w
::
`

func TestHintExitTerminatesScript(t *testing.T) {
	// present.txt never exists, so the baseline fails and hint_exit fires.
	// If hint_exit did NOT terminate the script, the next line
	// (wait_file_gone on a missing file) would succeed immediately and the
	// unit would complete - which must not happen.
	te := newTestEnv(t, map[string]string{"010.m/010.he/unit.md": hintExitTerminatesUnit})
	defer os.RemoveAll("/tmp/shellgym-he")

	if err := te.eng.ActivateUnit("m/he"); err != nil {
		t.Fatal(err)
	}
	if te.waitEvent("unit", "m/he", "completed", 8*time.Second) {
		t.Fatal("unit completed - hint_exit did not terminate the check script")
	}
	// And the hint itself must have been delivered with exit code 42
	// recorded for the attempt.
	runs := te.eng.Store.Runs("m/he", "gated")
	found42 := false
	for _, r := range runs {
		if r.ExitCode == 42 {
			found42 = true
		}
	}
	if !found42 {
		t.Fatalf("no run with the hint_exit code 42; runs: %+v", runs)
	}
}

const seqFirstUnit = `---
title: First
init:
  - name: prep
    run: |
      mkdir -p /tmp/shellgym-seqtest
tasks:
  made:
    check: |
      wait_file /tmp/shellgym-seqtest/a.txt
---
::task{name="made"}
Waiting...
::
`

const seqSecondUnit = `---
title: Second
needs: [first]
init:
  - name: seed
    run: |
      mkdir -p /tmp/shellgym-seqtest
      touch /tmp/shellgym-seqtest/b-init.txt
tasks:
  made:
    check: |
      wait_file /tmp/shellgym-seqtest/b.txt
---
::task{name="made"}
Waiting...
::
`

const seqThirdUnit = `---
title: Third
tasks:
  made:
    check: |
      wait_file /tmp/shellgym-seqtest/c.txt
---
::task{name="made"}
Waiting...
::
`

func TestUnitDependencyGating(t *testing.T) {
	te := newTestEnv(t, map[string]string{
		"010.m/010.first/unit.md":  seqFirstUnit,
		"010.m/020.second/unit.md": seqSecondUnit,
		"010.m/030.third/unit.md":  seqThirdUnit,
	})
	defer os.RemoveAll("/tmp/shellgym-seqtest")

	// A unit whose needs: are unsolved is locked: activation is rejected
	// and, crucially, its init must not run.
	if err := te.eng.ActivateUnit("m/second"); !errors.Is(err, ErrUnitLocked) {
		t.Fatalf("activate dependent unit: got %v, want ErrUnitLocked", err)
	}
	if _, err := os.Stat("/tmp/shellgym-seqtest/b-init.txt"); !os.IsNotExist(err) {
		t.Fatal("locked unit's init ran")
	}
	if !te.eng.UnitLocked("m/second") || te.eng.UnitLocked("m/first") {
		t.Fatal("lock flags: want second locked, first unlocked")
	}

	// A unit with no deps activates fine out of order.
	if te.eng.UnitLocked("m/third") {
		t.Fatal("independent unit reported locked")
	}
	if err := te.eng.ActivateUnit("m/third"); err != nil {
		t.Fatalf("activate independent unit out of order: %v", err)
	}

	// Solve the dependency...
	if err := te.eng.ActivateUnit("m/first"); err != nil {
		t.Fatal(err)
	}
	sh := newStudentShell(t)
	sh.Type("touch /tmp/shellgym-seqtest/a.txt")
	if !te.waitEvent("unit", "m/first", "completed", 15*time.Second) {
		t.Fatal("first unit did not complete")
	}

	// ...and the dependent unit unlocks: it activates and its init runs.
	if err := te.eng.ActivateUnit("m/second"); err != nil {
		t.Fatalf("activate unlocked unit: %v", err)
	}
	deadline := time.Now().Add(10 * time.Second)
	for {
		if _, err := os.Stat("/tmp/shellgym-seqtest/b-init.txt"); err == nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("unlocked unit's init did not run")
		}
		time.Sleep(100 * time.Millisecond)
	}
}
