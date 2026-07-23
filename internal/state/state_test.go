package state

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestStorePersistence(t *testing.T) {
	dir := t.TempDir()
	s, err := Open(dir, "p1")
	if err != nil {
		t.Fatal(err)
	}
	err = s.Update(func(d *Data) {
		d.CurrentUnit = "m/u"
		u := d.Unit("m/u")
		u.Status = UnitActive
		u.Vars = map[string]string{"X": "1"}
		u.Task("t").Status = "running"
	})
	if err != nil {
		t.Fatal(err)
	}

	s2, err := Open(dir, "p1")
	if err != nil {
		t.Fatal(err)
	}
	s2.View(func(d *Data) {
		if d.CurrentUnit != "m/u" {
			t.Errorf("current = %q", d.CurrentUnit)
		}
		u := d.Unit("m/u")
		if u.Status != UnitActive || u.Vars["X"] != "1" || u.Task("t").Status != "running" {
			t.Errorf("unit state lost: %+v", u)
		}
	})
}

func TestStoreCorruptRecovery(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "p1"), 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "p1", "progress.json")
	if err := os.WriteFile(path, []byte("{not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	s, err := Open(dir, "p1")
	if err != nil {
		t.Fatal(err)
	}
	s.View(func(d *Data) {
		if len(d.Units) != 0 {
			t.Error("expected fresh state")
		}
	})
	if _, err := os.Stat(path + ".corrupt"); err != nil {
		t.Error("corrupt backup not kept")
	}
}

func TestRunHistoryJSONL(t *testing.T) {
	dir := t.TempDir()
	s, err := Open(dir, "p1")
	if err != nil {
		t.Fatal(err)
	}
	for i := range 50 {
		if err := s.AddRun("m/u", "t1", TaskRun{StartedAt: time.Now(), ExitCode: i, Kind: "run"}); err != nil {
			t.Fatal(err)
		}
	}
	runs := s.Runs("m/u", "t1")
	if len(runs) != maxRunsPerTask {
		t.Fatalf("runs = %d", len(runs))
	}
	if runs[len(runs)-1].ExitCode != 49 {
		t.Error("kept wrong end of history")
	}
	// streams get truncated
	if err := s.AddRun("m/u", "t1", TaskRun{Stdout: strings.Repeat("x", 10000)}); err != nil {
		t.Fatal(err)
	}
	runs = s.Runs("m/u", "t1")
	last := runs[len(runs)-1]
	if len(last.Stdout) > maxStreamLen+100 || !strings.Contains(last.Stdout, "truncated") {
		t.Errorf("stdout not truncated: %d bytes", len(last.Stdout))
	}
	// task listing + unit drop
	if tasks := s.RunTasks("m/u"); len(tasks) != 1 || tasks[0] != "t1" {
		t.Errorf("RunTasks = %v", tasks)
	}
	if err := s.DropUnit("m/u"); err != nil {
		t.Fatal(err)
	}
	if runs := s.Runs("m/u", "t1"); len(runs) != 0 {
		t.Error("runs survived DropUnit")
	}
	// progress.json stays small even with lots of run output
	raw, err := os.ReadFile(filepath.Join(dir, "p1", "progress.json"))
	if err != nil {
		t.Fatal(err)
	}
	if len(raw) > 4096 {
		t.Errorf("progress.json unexpectedly large: %d bytes", len(raw))
	}
}
