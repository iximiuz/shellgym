package engine

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
)

// ShellInfo describes one interactive shell of the observed user.
type ShellInfo struct {
	PID     int    `json:"pid"`
	Exe     string `json:"exe"`
	TTY     string `json:"tty"`
	Cwd     string `json:"cwd"`
	StartNs int64  `json:"-"` // process start time (jiffies), for "most recent"
}

var shellNames = map[string]bool{
	"bash": true, "zsh": true, "sh": true, "fish": true, "dash": true, "ash": true,
}

// FindShells scans /proc for interactive shells owned by username: a shell
// binary with a controlling terminal. Sorted by most recently started first.
func FindShells(username string) ([]ShellInfo, error) {
	uid, err := lookupUID(username)
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return nil, err
	}
	var out []ShellInfo
	for _, e := range entries {
		pid, err := strconv.Atoi(e.Name())
		if err != nil {
			continue
		}
		info, ok := inspectShell(pid, uid)
		if ok {
			out = append(out, info)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].StartNs > out[j].StartNs })
	return out, nil
}

func inspectShell(pid, uid int) (ShellInfo, bool) {
	var info ShellInfo
	st, err := os.Stat(fmt.Sprintf("/proc/%d", pid))
	if err != nil {
		return info, false
	}
	sys, ok := st.Sys().(*syscall.Stat_t)
	if !ok || int(sys.Uid) != uid {
		return info, false
	}
	comm, err := os.ReadFile(fmt.Sprintf("/proc/%d/comm", pid))
	if err != nil {
		return info, false
	}
	name := strings.TrimSpace(string(comm))
	if !shellNames[strings.TrimPrefix(name, "-")] {
		return info, false
	}
	// Interactive shells have a controlling tty: field 7 of /proc/pid/stat
	// (tty_nr) is non-zero.
	stat, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
	if err != nil {
		return info, false
	}
	fields := statFields(string(stat))
	if len(fields) < 22 {
		return info, false
	}
	ttyNr, _ := strconv.Atoi(fields[6])
	if ttyNr == 0 {
		return info, false
	}
	start, _ := strconv.ParseInt(fields[21], 10, 64)
	cwd, err := os.Readlink(fmt.Sprintf("/proc/%d/cwd", pid))
	if err != nil {
		return info, false
	}
	tty, _ := os.Readlink(fmt.Sprintf("/proc/%d/fd/0", pid))
	return ShellInfo{PID: pid, Exe: name, TTY: filepath.Base(tty), Cwd: cwd, StartNs: start}, true
}

// statFields splits /proc/pid/stat handling the parenthesized comm field
// (which may contain spaces).
func statFields(stat string) []string {
	close := strings.LastIndex(stat, ")")
	if close < 0 {
		return nil
	}
	head := strings.Fields(stat[:strings.Index(stat, "(")])
	rest := strings.Fields(stat[close+1:])
	out := append([]string{}, head...)              // pid
	out = append(out, stat[len(head[0])+1:close+1]) // (comm)
	return append(out, rest...)
}

func lookupUID(username string) (int, error) {
	u, err := user.Lookup(username)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(u.Uid)
}

func lookupHome(username string) (string, error) {
	u, err := user.Lookup(username)
	if err != nil {
		return "", err
	}
	return u.HomeDir, nil
}

// ShellEnvOf extracts an environment variable observed in the *initial* env
// of a process (note: /proc/pid/environ does not reflect later exports; use
// exec-event children to observe those).
func ShellEnvOf(pid int, name string) (string, bool) {
	raw, err := os.ReadFile(fmt.Sprintf("/proc/%d/environ", pid))
	if err != nil {
		return "", false
	}
	for _, kv := range strings.Split(string(raw), "\x00") {
		if v, ok := strings.CutPrefix(kv, name+"="); ok {
			return v, true
		}
	}
	return "", false
}
