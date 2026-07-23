package content

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStripSolveScripts(t *testing.T) {
	dir := t.TempDir()
	unit := `---
title: A unit
tasks:
  t1:
    check: |
      wait_file /tmp/x
    solve: |
      touch /tmp/x
      echo done
  t2:
    solve: touch /tmp/y
    check: |
      true
---
Body stays.
`
	path := filepath.Join(dir, "010.m", "010.u", "unit.md")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(unit), 0o644); err != nil {
		t.Fatal(err)
	}
	n, err := StripSolveScripts(dir)
	if err != nil {
		t.Fatal(err)
	}
	if n != 2 {
		t.Fatalf("stripped = %d", n)
	}
	raw, _ := os.ReadFile(path)
	out := string(raw)
	if strings.Contains(out, "solve") || strings.Contains(out, "touch /tmp/x") || strings.Contains(out, "touch /tmp/y") {
		t.Fatalf("solve not stripped:\n%s", out)
	}
	if !strings.Contains(out, "wait_file /tmp/x") || !strings.Contains(out, "Body stays.") {
		t.Fatalf("over-stripped:\n%s", out)
	}
	// The stripped file must still load.
	writeFile(t, filepath.Join(dir, "path.yaml"), "id: x\ntitle: X\n")
	if _, err := Load(dir, "ubuntu", nil, nil); err != nil {
		t.Fatalf("stripped file does not load: %v", err)
	}
}
