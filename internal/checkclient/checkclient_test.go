package checkclient

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"
)

// oneShot invokes a check in --now form (zero deadline = single evaluation),
// bypassing the daemon socket - both checks under test are pure fs reads.
func oneShot(t *testing.T, name string, args ...string) bool {
	t.Helper()
	ok, err := (&client{}).run(name, args, time.Time{})
	if err != nil {
		t.Fatalf("%s %v: %v", name, args, err)
	}
	return ok
}

func TestWaitFileMode(t *testing.T) {
	f := filepath.Join(t.TempDir(), "secret.txt")
	if err := os.WriteFile(f, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if oneShot(t, "wait_file_mode", f, "600") {
		t.Error("644 file must not match 600")
	}
	if !oneShot(t, "wait_file_mode", f, "644") {
		t.Error("644 file must match 644")
	}
	if err := syscall.Chmod(f, 0o4755); err != nil {
		t.Fatal(err)
	}
	if !oneShot(t, "wait_file_mode", f, "4755") {
		t.Error("setuid bits must compare against 4-digit octal")
	}
	if oneShot(t, "wait_file_mode", f, "755") {
		t.Error("3-digit octal must not ignore special bits")
	}
	if oneShot(t, "wait_file_mode", filepath.Join(t.TempDir(), "gone"), "600") {
		t.Error("missing file must not match")
	}
	if _, err := (&client{}).run("wait_file_mode", []string{f, "9x9"}, time.Time{}); err == nil {
		t.Error("want error for invalid octal mode")
	}
}

func TestWaitDir(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "box")
	if err := os.WriteFile(filepath.Join(dir, "box.txt"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if oneShot(t, "wait_dir", sub) {
		t.Error("missing dir must not match")
	}
	if oneShot(t, "wait_dir", filepath.Join(dir, "box.txt")) {
		t.Error("plain file must not match")
	}
	if err := os.Mkdir(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	if !oneShot(t, "wait_dir", sub) {
		t.Error("existing dir must match")
	}
	if !oneShot(t, "wait_dir", filepath.Join(dir, "b*")) {
		t.Error("glob matching a dir must match")
	}
}

func TestWaitFileNewer(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "target")
	ref := filepath.Join(dir, "ref")
	old := time.Now().Add(-48 * time.Hour)
	for _, f := range []string{target, ref} {
		if err := os.WriteFile(f, nil, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Chtimes(target, old, old); err != nil {
		t.Fatal(err)
	}
	if oneShot(t, "wait_file_newer", target, ref) {
		t.Error("older target must not match")
	}
	now := time.Now()
	if err := os.Chtimes(target, now, now); err != nil {
		t.Fatal(err)
	}
	if !oneShot(t, "wait_file_newer", target, ref) {
		t.Error("freshly touched target must match")
	}
	if oneShot(t, "wait_file_newer", filepath.Join(dir, "gone"), ref) {
		t.Error("missing target must not match")
	}
	if oneShot(t, "wait_file_newer", target, filepath.Join(dir, "gone")) {
		t.Error("missing reference must not match")
	}
}
