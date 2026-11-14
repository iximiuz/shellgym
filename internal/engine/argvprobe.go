package engine

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"golang.org/x/sys/unix"
)

// ArgvProbe captures the argv of every command bash executes, sourced from
// a kernel uprobe on bash's shell_execve() - the function bash calls right
// before execve(2), in the process that is about to become the command,
// with the command's path and argv. The kernel copies the strings into the
// trace buffer at that moment, so the argv survives however short-lived
// the command is.
//
// It backs up the proc connector: an exec event carries only a pid, and
// the argv must be read from /proc/<pid>/cmdline, which is empty once the
// process has exited. A sub-millisecond command (`whoami`, a `date` with a
// bad option) can exit before the daemon gets to that read; its event is
// then completed with the argv this probe recorded for the same pid.
//
// Like the line watcher it is optional: it needs root, tracefs with uprobe
// events, a bash that exports shell_execve, and a known calling convention
// (amd64, arm64). Without it, the daemon behaves as before - very fast
// commands may go unseen.
type ArgvProbe struct {
	// Target is the probed binary and the shell_execve offset in it (for logs).
	Target string

	tracefs  string
	instance string
	fd       int
	done     chan struct{} // closed when readLoop has exited
	closed   atomic.Bool
	onArgv   func(pid int, argv []string)
	onDrop   func(n int)
	onTick   func()
}

const (
	argvProbeName     = "execve"
	argvProbeInstance = "shellgym-exec"
	// argvProbeMaxArgs is how many argv entries the probe fetches. A record
	// with all of them present may be truncated and is dropped rather than
	// matched with a partial argv.
	argvProbeMaxArgs = 24
	// argvProbeBufferKB is the per-CPU trace buffer of the probe's instance:
	// room for a few seconds of a bash exec storm while the reader catches up.
	argvProbeBufferKB = 4096
)

// argvFieldNames are the probe's fetch arguments in record order: the
// executed path, then argv[0..argvProbeMaxArgs-1].
var argvFieldNames = func() []string {
	names := []string{"cmd"}
	for i := 0; i < argvProbeMaxArgs; i++ {
		names = append(names, fmt.Sprintf("a%d", i))
	}
	return names
}()

// argvProbeRegs returns the registers holding shell_execve's first two
// arguments (char *command, char **args) for the running architecture.
func argvProbeRegs() (cmd, args string, err error) {
	switch runtime.GOARCH {
	case "amd64":
		return "%di", "%si", nil
	case "arm64":
		return "%x0", "%x1", nil
	}
	return "", "", fmt.Errorf("no shell_execve calling convention for %s", runtime.GOARCH)
}

// startArgvProbe installs the shell_execve uprobe for the given shell binary
// and streams the recorded argvs to onArgv. onDrop gets the number of
// records the kernel dropped (trace buffer full), and onTick runs on every
// reader wakeup (at least every 200 ms) for housekeeping. It returns an
// error, leaving nothing behind in tracefs, when any prerequisite is
// missing.
func startArgvProbe(shellPath string, onArgv func(pid int, argv []string), onDrop func(n int), onTick func()) (*ArgvProbe, error) {
	regCmd, regArgs, err := argvProbeRegs()
	if err != nil {
		return nil, err
	}
	tracefs, err := findTracefs()
	if err != nil {
		return nil, err
	}
	shellPath, err = filepath.EvalSymlinks(shellPath)
	if err != nil {
		return nil, fmt.Errorf("shell: %w", err)
	}
	if filepath.Base(shellPath) != "bash" {
		return nil, fmt.Errorf("shell %s is not bash", shellPath)
	}
	offset, ok, err := elfSymbolOffset(shellPath, "shell_execve")
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("%s exports no shell_execve symbol", shellPath)
	}

	p := &ArgvProbe{
		Target:   fmt.Sprintf("%s+0x%x", shellPath, offset),
		tracefs:  tracefs,
		instance: filepath.Join(tracefs, "instances", argvProbeInstance),
		fd:       -1,
		onArgv:   onArgv,
		onDrop:   onDrop,
		onTick:   onTick,
	}
	// A previous daemon may have died without cleaning up.
	p.teardown()

	var spec strings.Builder
	fmt.Fprintf(&spec, "p:%s/%s %s:0x%x cmd=+0(%s):string", probeGroup, argvProbeName, shellPath, offset, regCmd)
	for i := 0; i < argvProbeMaxArgs; i++ {
		// args[i] is a pointer at args+8*i; the string it points to is the argument.
		fmt.Fprintf(&spec, " a%d=+0(+%d(%s)):string", i, 8*i, regArgs)
	}
	spec.WriteString("\n")
	if err := appendFile(filepath.Join(tracefs, "uprobe_events"), spec.String()); err != nil {
		return nil, fmt.Errorf("register uprobe: %w", err)
	}
	if err := os.Mkdir(p.instance, 0o755); err != nil && !errors.Is(err, os.ErrExist) {
		p.teardown()
		return nil, fmt.Errorf("create tracefs instance: %w", err)
	}
	// Best effort: a larger buffer only makes drops under exec storms rarer.
	_ = writeFile(filepath.Join(p.instance, "buffer_size_kb"), fmt.Sprintf("%d\n", argvProbeBufferKB))
	if err := writeFile(filepath.Join(p.instance, "events", probeGroup, argvProbeName, "enable"), "1\n"); err != nil {
		p.teardown()
		return nil, fmt.Errorf("enable uprobe: %w", err)
	}
	fd, err := unix.Open(filepath.Join(p.instance, "trace_pipe"), unix.O_RDONLY|unix.O_NONBLOCK|unix.O_CLOEXEC, 0)
	if err != nil {
		p.teardown()
		return nil, fmt.Errorf("open trace_pipe: %w", err)
	}
	p.fd = fd
	p.done = make(chan struct{})
	go p.readLoop()
	return p, nil
}

// Close stops the reader and removes the probe and its tracefs instance
// (see LineWatcher.Close for why the reader must leave poll(2) first).
func (p *ArgvProbe) Close() {
	if !p.closed.CompareAndSwap(false, true) {
		return
	}
	select {
	case <-p.done:
	case <-time.After(2 * time.Second):
	}
	p.teardown()
}

func (p *ArgvProbe) teardown() {
	if p.fd >= 0 {
		_ = unix.Close(p.fd)
		p.fd = -1
	}
	enable := filepath.Join("events", probeGroup, argvProbeName, "enable")
	_ = writeFile(filepath.Join(p.instance, enable), "0\n")
	_ = writeFile(filepath.Join(p.tracefs, enable), "0\n")
	_ = os.Remove(p.instance)
	_ = appendFile(filepath.Join(p.tracefs, "uprobe_events"), fmt.Sprintf("-:%s/%s\n", probeGroup, argvProbeName))
}

func (p *ArgvProbe) readLoop() {
	defer close(p.done)
	buf := make([]byte, 256*1024)
	var pending []byte
	for !p.closed.Load() {
		if p.onTick != nil {
			p.onTick()
		}
		fds := []unix.PollFd{{Fd: int32(p.fd), Events: unix.POLLIN}}
		n, err := unix.Poll(fds, 200)
		if err != nil && err != unix.EINTR {
			log.Printf("execwatch: argv probe poll: %v; argv probe disabled", err)
			return
		}
		if n <= 0 {
			continue
		}
		n, err = unix.Read(p.fd, buf)
		if err != nil {
			if err == unix.EAGAIN || err == unix.EINTR {
				continue
			}
			if !p.closed.Load() {
				log.Printf("execwatch: argv probe read: %v; argv probe disabled", err)
			}
			return
		}
		pending = append(pending, buf[:n]...)
		for {
			i := bytes.IndexByte(pending, '\n')
			if i < 0 {
				break
			}
			line := string(pending[:i])
			pending = pending[i+1:]
			if pid, argv, ok := parseArgvTraceLine(line); ok {
				p.onArgv(pid, argv)
			} else if m := lostEventsRe.FindStringSubmatch(line); m != nil && p.onDrop != nil {
				n, _ := strconv.Atoi(m[1])
				p.onDrop(n)
			}
		}
	}
}

// lostEventsRe matches the marker trace_pipe prints where the kernel had to
// drop records, e.g. "CPU:2 [LOST 31 EVENTS]".
var lostEventsRe = regexp.MustCompile(`\[LOST (\d+) EVENTS\]`)

// argvTraceHeadRe matches the header of one trace_pipe record of the
// shell_execve probe, e.g.
//
//	date-7035  [002] DNZff  29.289614: execve: (0x5627288dc530) cmd="/usr/bin/date" a0="date" a1="--x" a2=(fault) ...
//
// The comm is at most 16 bytes and may itself contain dashes, so the pid is
// anchored on the "[cpu]" column that follows it.
var argvTraceHeadRe = regexp.MustCompile(`^\s*.{1,16}?-(\d+)\s+\[\d+\]\s+\S+\s+[\d.]+: ` + argvProbeName + `: \([^)]*\) `)

// parseArgvTraceLine extracts the pid and the argv from one record of the
// shell_execve probe. tracefs prints the strings in double quotes without
// escaping, so a value is delimited by the name of the field after it. The
// argv ends at the first unreadable entry - the NULL that terminates it
// shows up as (fault). A record whose every argv field is present may be
// truncated and yields false, as does any record that does not parse.
func parseArgvTraceLine(s string) (pid int, argv []string, ok bool) {
	m := argvTraceHeadRe.FindStringSubmatchIndex(s)
	if m == nil {
		return 0, nil, false
	}
	pid, err := strconv.Atoi(s[m[2]:m[3]])
	if err != nil {
		return 0, nil, false
	}
	rest := s[m[1]:]
	values := make([]*string, len(argvFieldNames))
	for i, name := range argvFieldNames {
		rest, ok = strings.CutPrefix(rest, name+"=")
		if !ok {
			return 0, nil, false
		}
		switch {
		case strings.HasPrefix(rest, "(fault)"):
			rest = rest[len("(fault)"):]
		case strings.HasPrefix(rest, `"`):
			rest = rest[1:]
			end := -1
			if i+1 < len(argvFieldNames) {
				end = strings.Index(rest, `" `+argvFieldNames[i+1]+"=")
			} else if strings.HasSuffix(rest, `"`) {
				end = len(rest) - 1
			}
			if end < 0 {
				return 0, nil, false
			}
			v := rest[:end]
			values[i] = &v
			rest = rest[end+1:]
		default:
			return 0, nil, false
		}
		rest = strings.TrimPrefix(rest, " ")
	}
	if values[len(values)-1] != nil {
		return 0, nil, false
	}
	for _, v := range values[1:] {
		if v == nil {
			break
		}
		argv = append(argv, *v)
	}
	if len(argv) == 0 {
		return 0, nil, false
	}
	return pid, argv, true
}
