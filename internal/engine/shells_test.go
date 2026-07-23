package engine

import (
	"fmt"
	"os"
	"os/user"
	"testing"
)

func TestStatFields(t *testing.T) {
	stat := "1234 (my shell) S 1 1234 1234 34816 1234 4194304 1 2 3 4 5 6 7 8 20 0 1 0 9999 1 2 3"
	f := statFields(stat)
	if len(f) < 22 {
		t.Fatalf("fields: %d %v", len(f), f)
	}
	if f[0] != "1234" || f[1] != "(my shell)" || f[2] != "S" {
		t.Errorf("head fields: %v", f[:3])
	}
	if f[6] != "34816" {
		t.Errorf("tty_nr field: %q", f[6])
	}
	if f[21] != "9999" {
		t.Errorf("starttime field: %q", f[21])
	}
}

func TestFindShellsSelf(t *testing.T) {
	// The test process itself is not a shell; just assert the scan runs
	// clean for the current user and returns well-formed entries.
	me, err := user.Current()
	if err != nil {
		t.Skip("no current user")
	}
	shells, err := FindShells(me.Username)
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range shells {
		if s.PID <= 0 || s.Cwd == "" {
			t.Errorf("malformed shell entry: %+v", s)
		}
		if _, err := os.Stat(fmt.Sprintf("/proc/%d", s.PID)); err != nil {
			t.Errorf("reported shell pid %d does not exist", s.PID)
		}
	}
}

func TestLookupUIDUnknown(t *testing.T) {
	if _, err := lookupUID("definitely-no-such-user-xyz"); err == nil {
		t.Error("want error for unknown user")
	}
}
