// Package checkclient implements the built-in checks that task scripts call
// (via PATH shims like `wait_cwd`). Each check is a small command: exit 0
// when the condition is met, non-zero otherwise. `wait_*` checks block until
// the condition is met (or --timeout / --now); state-reading helpers like
// `shell_cwd` print to stdout.
package checkclient

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Names lists every check command (used to generate PATH shims).
var Names = []string{
	"shell_cwd", "shells", "hint_exit", "set_var",
	"wait_cwd", "wait_exec", "wait_env",
	"wait_file", "wait_file_gone", "wait_file_contains",
	"wait_proc", "wait_proc_gone",
	"wait_port", "wait_port_free",
}

// HintExitCode is the distinct exit code of the hint_exit built-in. The
// runner injects a shell-function wrapper so that calling `hint_exit`
// inside a check script posts the hint AND terminates the script with
// this code.
const HintExitCode = 42

// Main runs check `name` with args; returns process exit code.
func Main(name string, args []string) int {
	if name == "hint_exit" {
		return hintExit(args)
	}
	if name == "set_var" {
		return setVar(args)
	}
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	timeout := fs.Float64("timeout", 0, "give up after this many seconds (0 = wait forever)")
	now := fs.Bool("now", false, "single instant check, no waiting")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	args = fs.Args()

	c := &client{sock: os.Getenv("GYM_SOCK"), since: sinceSeq()}
	deadline := time.Now().Add(365 * 24 * time.Hour)
	if *timeout > 0 {
		deadline = time.Now().Add(time.Duration(*timeout * float64(time.Second)))
	}
	if *now {
		deadline = time.Time{} // signals single-shot
	}

	ok, err := c.run(name, args, deadline)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", name, err)
		return 2
	}
	if !ok {
		return 1
	}
	return 0
}

func sinceSeq() uint64 {
	v, _ := strconv.ParseUint(os.Getenv("GYM_SINCE_SEQ"), 10, 64)
	return v
}

type client struct {
	sock  string
	since uint64
}

func (c *client) http() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
				var d net.Dialer
				return d.DialContext(ctx, "unix", c.sock)
			},
		},
	}
}

func (c *client) run(name string, args []string, deadline time.Time) (bool, error) {
	oneShot := deadline.IsZero()
	if oneShot {
		deadline = time.Now()
	}
	switch name {
	case "shell_cwd":
		if len(args) > 1 {
			return false, fmt.Errorf("usage: shell_cwd [shell-pid]")
		}
		shells, err := c.getShells()
		if err != nil {
			return false, err
		}
		if len(args) == 1 {
			pid, err := parsePID(args[0])
			if err != nil {
				return false, err
			}
			for _, s := range shells {
				if s.PID == pid {
					fmt.Println(s.Cwd)
					return true, nil
				}
			}
			return false, fmt.Errorf("no interactive shell with pid %d", pid)
		}
		if len(shells) == 0 {
			return false, fmt.Errorf("no interactive shells found")
		}
		fmt.Println(shells[0].Cwd)
		return true, nil
	case "shells":
		shells, err := c.getShells()
		if err != nil {
			return false, err
		}
		for _, s := range shells {
			fmt.Printf("%d\t%s\t%s\t%s\n", s.PID, s.Exe, s.TTY, s.Cwd)
		}
		return true, nil
	case "wait_cwd":
		if len(args) < 1 || len(args) > 2 {
			return false, fmt.Errorf("usage: wait_cwd [shell-pid] <abs-path-or-regex>")
		}
		onlyPID := 0
		pattern := args[0]
		if len(args) == 2 {
			pid, err := parsePID(args[0])
			if err != nil {
				return false, err
			}
			onlyPID = pid
			pattern = args[1]
		}
		return c.poll(oneShot, deadline, func() (bool, error) {
			shells, err := c.getShells()
			if err != nil {
				return false, err
			}
			for _, s := range shells {
				if onlyPID != 0 && s.PID != onlyPID {
					continue
				}
				if pathMatch(pattern, s.Cwd) {
					// Report which shell matched, so checks can capture it
					// (typically into a task var via set_var).
					fmt.Println(s.PID)
					return true, nil
				}
			}
			return false, nil
		})
	case "wait_exec":
		if len(args) != 1 {
			return false, fmt.Errorf("usage: wait_exec <regex>")
		}
		return c.execWait(execWaitRequest{Regex: args[0]}, oneShot, deadline)
	case "wait_env":
		if len(args) < 1 || len(args) > 2 {
			return false, fmt.Errorf("usage: wait_env <NAME> [regex]")
		}
		req := execWaitRequest{EnvName: args[0], EnvRegex: ".*"}
		if len(args) == 2 {
			req.EnvRegex = args[1]
		}
		return c.execWait(req, oneShot, deadline)
	case "wait_file":
		if len(args) != 1 {
			return false, fmt.Errorf("usage: wait_file <path-or-glob>")
		}
		return c.poll(oneShot, deadline, func() (bool, error) {
			m, err := filepath.Glob(args[0])
			return err == nil && len(m) > 0, nil
		})
	case "wait_file_gone":
		if len(args) != 1 {
			return false, fmt.Errorf("usage: wait_file_gone <path-or-glob>")
		}
		return c.poll(oneShot, deadline, func() (bool, error) {
			m, err := filepath.Glob(args[0])
			return err == nil && len(m) == 0, nil
		})
	case "wait_file_contains":
		if len(args) != 2 {
			return false, fmt.Errorf("usage: wait_file_contains <path> <regex>")
		}
		// (?m): ^ and $ anchor to lines - the natural reading for
		// "the file contains a line matching ...".
		re, err := regexp.Compile("(?m)" + args[1])
		if err != nil {
			return false, err
		}
		return c.poll(oneShot, deadline, func() (bool, error) {
			raw, err := os.ReadFile(args[0])
			return err == nil && re.Match(raw), nil
		})
	case "wait_proc":
		if len(args) != 1 {
			return false, fmt.Errorf("usage: wait_proc <regex>")
		}
		re, err := regexp.Compile(args[0])
		if err != nil {
			return false, err
		}
		return c.poll(oneShot, deadline, func() (bool, error) {
			return procMatch(re), nil
		})
	case "wait_proc_gone":
		if len(args) != 1 {
			return false, fmt.Errorf("usage: wait_proc_gone <regex>")
		}
		re, err := regexp.Compile(args[0])
		if err != nil {
			return false, err
		}
		return c.poll(oneShot, deadline, func() (bool, error) {
			return !procMatch(re), nil
		})
	case "wait_port":
		return c.portCheck(args, oneShot, deadline, true)
	case "wait_port_free":
		return c.portCheck(args, oneShot, deadline, false)
	}
	return false, fmt.Errorf("unknown check %q", name)
}

func (c *client) portCheck(args []string, oneShot bool, deadline time.Time, wantListen bool) (bool, error) {
	if len(args) != 1 {
		return false, fmt.Errorf("usage: wait_port[_free] <port>")
	}
	port, err := strconv.Atoi(args[0])
	if err != nil {
		return false, err
	}
	return c.poll(oneShot, deadline, func() (bool, error) {
		return portListening(port) == wantListen, nil
	})
}

func (c *client) poll(oneShot bool, deadline time.Time, fn func() (bool, error)) (bool, error) {
	for {
		ok, err := fn()
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
		if oneShot || time.Now().After(deadline) {
			return false, nil
		}
		time.Sleep(200 * time.Millisecond)
	}
}

func parsePID(s string) (int, error) {
	pid, err := strconv.Atoi(s)
	if err != nil || pid <= 0 {
		return 0, fmt.Errorf("invalid shell pid %q", s)
	}
	return pid, nil
}

// pathMatch: exact path match, or regex when the pattern contains regex
// metacharacters and compiles.
func pathMatch(pattern, cwd string) bool {
	if pattern == cwd {
		return true
	}
	if strings.ContainsAny(pattern, "?*[](|^$+") {
		if re, err := regexp.Compile("^(" + pattern + ")$"); err == nil {
			return re.MatchString(cwd)
		}
	}
	return false
}

func procMatch(re *regexp.Regexp) bool {
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return false
	}
	self := os.Getpid()
	for _, e := range entries {
		pid, err := strconv.Atoi(e.Name())
		if err != nil || pid == self {
			continue
		}
		raw, err := os.ReadFile(fmt.Sprintf("/proc/%d/cmdline", pid))
		if err != nil || len(raw) == 0 {
			continue
		}
		cmdline := strings.ReplaceAll(strings.TrimRight(string(raw), "\x00"), "\x00", " ")
		if re.MatchString(cmdline) {
			return true
		}
	}
	return false
}

// portListening reports whether any local TCP socket listens on port
// (state 0A) per /proc/net/tcp and tcp6.
func portListening(port int) bool {
	for _, f := range []string{"/proc/net/tcp", "/proc/net/tcp6"} {
		raw, err := os.ReadFile(f)
		if err != nil {
			continue
		}
		for _, line := range strings.Split(string(raw), "\n")[1:] {
			fields := strings.Fields(line)
			if len(fields) < 4 || fields[3] != "0A" {
				continue
			}
			addr := fields[1]
			idx := strings.LastIndex(addr, ":")
			if idx < 0 {
				continue
			}
			p, err := strconv.ParseInt(addr[idx+1:], 16, 32)
			if err == nil && int(p) == port {
				return true
			}
		}
	}
	return false
}

// hintExit posts a hint message for the current task to the daemon:
//
//	hint_exit <message>          (task inferred from $GYM_TASK)
//	hint_exit <task> <message>   (explicit task, for multi-task scripts)
//
// The binary itself exits with HintExitCode; the shell-function wrapper
// the runner injects then terminates the WHOLE check script with that
// code, so `wait_x ... || hint_exit "msg"` both reports and stops.
func hintExit(args []string) int {
	unit := os.Getenv("GYM_UNIT")
	task := os.Getenv("GYM_TASK")
	var msg string
	switch len(args) {
	case 1:
		msg = args[0]
	case 2:
		task = args[0]
		msg = args[1]
	default:
		fmt.Fprintln(os.Stderr, "usage: hint_exit [task] <message>")
		return HintExitCode
	}
	if unit == "" || task == "" {
		fmt.Fprintln(os.Stderr, "hint_exit: GYM_UNIT/GYM_TASK not set (must run inside a task script)")
		return HintExitCode
	}
	c := &client{sock: os.Getenv("GYM_SOCK")}
	if c.sock == "" {
		fmt.Fprintln(os.Stderr, "hint_exit: GYM_SOCK not set")
		return HintExitCode
	}
	body, _ := json.Marshal(map[string]string{"unit": unit, "task": task, "message": msg})
	resp, err := c.http().Post("http://gym/hint", "application/json", bytes.NewReader(body))
	if err != nil {
		fmt.Fprintf(os.Stderr, "hint_exit: %v\n", err)
		return HintExitCode
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Fprintf(os.Stderr, "hint_exit: daemon returned %s\n", resp.Status)
	}
	return HintExitCode
}

// setVar persists a task var on the current unit:
//
//	set_var <NAME> <value>
//
// The var joins the unit's vars: exported into subsequent runs of the
// unit's own scripts and into scripts of units that declare this unit in
// `needs:`. This is the way to pass small values between tasks (and on
// to dependent units) - use files only for BLOB-like data.
func setVar(args []string) int {
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: set_var <NAME> <value>")
		return 2
	}
	unit := os.Getenv("GYM_UNIT")
	if unit == "" {
		fmt.Fprintln(os.Stderr, "set_var: GYM_UNIT not set (must run inside a task script)")
		return 2
	}
	c := &client{sock: os.Getenv("GYM_SOCK")}
	if c.sock == "" {
		fmt.Fprintln(os.Stderr, "set_var: GYM_SOCK not set")
		return 2
	}
	body, _ := json.Marshal(map[string]string{"unit": unit, "name": args[0], "value": args[1]})
	resp, err := c.http().Post("http://gym/vars", "application/json", bytes.NewReader(body))
	if err != nil {
		fmt.Fprintf(os.Stderr, "set_var: %v\n", err)
		return 2
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		msg, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		fmt.Fprintf(os.Stderr, "set_var: daemon returned %s: %s\n", resp.Status, strings.TrimSpace(string(msg)))
		return 2
	}
	return 0
}

// --- daemon API calls -------------------------------------------------------

type shellInfo struct {
	PID int    `json:"pid"`
	Exe string `json:"exe"`
	TTY string `json:"tty"`
	Cwd string `json:"cwd"`
}

func (c *client) getShells() ([]shellInfo, error) {
	if c.sock == "" {
		return nil, fmt.Errorf("GYM_SOCK not set (check must run inside a task script)")
	}
	resp, err := c.http().Get("http://gym/shells")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var shells []shellInfo
	return shells, json.NewDecoder(resp.Body).Decode(&shells)
}

type execWaitRequest struct {
	After      uint64  `json:"after"`
	Regex      string  `json:"regex"`
	EnvName    string  `json:"envName"`
	EnvRegex   string  `json:"envRegex"`
	TimeoutSec float64 `json:"timeoutSec"`
}

type execWaitResponse struct {
	Matched bool `json:"matched"`
}

func (c *client) execWait(req execWaitRequest, oneShot bool, deadline time.Time) (bool, error) {
	if c.sock == "" {
		return false, fmt.Errorf("GYM_SOCK not set (check must run inside a task script)")
	}
	req.After = c.since
	if oneShot {
		req.TimeoutSec = 0.05
	} else {
		req.TimeoutSec = time.Until(deadline).Seconds()
	}
	body, _ := json.Marshal(req)
	resp, err := c.http().Post("http://gym/exec/wait", "application/json", bytes.NewReader(body))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return false, fmt.Errorf("daemon returned %s", resp.Status)
	}
	var out execWaitResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return false, err
	}
	return out.Matched, nil
}
