package engine

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/sys/unix"
)

// ExecEvent is one observed process execution.
type ExecEvent struct {
	Seq  uint64    `json:"seq"`
	Time time.Time `json:"time"`
	PID  int       `json:"pid"`
	PPID int       `json:"ppid"`
	UID  int       `json:"uid"`
	// TTYNr is the controlling terminal (0 = none). Only tty-attached
	// processes count as student activity: the daemon's own task scripts
	// run in a fresh session with no controlling tty, which prevents a
	// check's own argv (which contains the searched pattern) from matching.
	TTYNr int      `json:"ttyNr"`
	Argv  []string `json:"argv"`
	// Env is captured eagerly for tty-attached processes (student
	// commands are low-rate; fast ones die before a lazy read could
	// happen). Empty for tty-less processes.
	Env []string `json:"-"`
}

// ExecWatcher records exec events into a bounded ring buffer, sourced from
// the kernel proc connector (netlink). If the connector is unavailable
// (missing CONFIG_PROC_EVENTS or CAP_NET_ADMIN), exec watching is simply
// disabled - exec-based checks (wait_exec, wait_env) then never fire.
type ExecWatcher struct {
	mu     sync.Mutex
	cond   *sync.Cond
	ring   []ExecEvent
	seq    uint64
	closed bool

	// Source is "netlink" when the connector is active, "" when exec
	// watching is unavailable.
	Source string
}

const ringSize = 4096

func NewExecWatcher() *ExecWatcher {
	w := &ExecWatcher{}
	w.cond = sync.NewCond(&w.mu)
	return w
}

// Start begins watching via the kernel proc connector. It returns an error
// if the connector cannot be opened (missing CONFIG_PROC_EVENTS or, most
// commonly, no CAP_NET_ADMIN because the daemon is not root) - exec
// watching is the only mechanism, so the daemon treats this as fatal rather
// than running with silently-broken wait_exec/wait_env checks.
func (w *ExecWatcher) Start() error {
	sock, err := openProcConnector()
	if err != nil {
		return fmt.Errorf("proc connector unavailable (run shellgym as root): %w", err)
	}
	w.Source = "netlink"
	go w.netlinkLoop(sock)
	return nil
}

func (w *ExecWatcher) Close() {
	w.mu.Lock()
	w.closed = true
	w.cond.Broadcast()
	w.mu.Unlock()
}

// Seq returns the current sequence number; events published after a given
// point have Seq greater than this.
func (w *ExecWatcher) Seq() uint64 {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.seq
}

// WaitMatch blocks until an event with Seq > after matches fn, returning it.
// Returns false when the deadline passes or the watcher closes.
func (w *ExecWatcher) WaitMatch(after uint64, deadline time.Time, fn func(ExecEvent) bool) (ExecEvent, bool) {
	timer := time.AfterFunc(time.Until(deadline), func() {
		w.mu.Lock()
		w.cond.Broadcast()
		w.mu.Unlock()
	})
	defer timer.Stop()

	w.mu.Lock()
	defer w.mu.Unlock()
	scanned := after
	for {
		for _, ev := range w.ring {
			if ev.Seq > scanned && fn(ev) {
				return ev, true
			}
		}
		scanned = w.seq
		if w.closed || time.Now().After(deadline) {
			return ExecEvent{}, false
		}
		w.cond.Wait()
	}
}

// Snapshot returns events with Seq > after (for debugging APIs).
func (w *ExecWatcher) Snapshot(after uint64, limit int) []ExecEvent {
	w.mu.Lock()
	defer w.mu.Unlock()
	var out []ExecEvent
	for _, ev := range w.ring {
		if ev.Seq > after {
			out = append(out, ev)
		}
	}
	if limit > 0 && len(out) > limit {
		out = out[len(out)-limit:]
	}
	return out
}

func (w *ExecWatcher) publish(ev ExecEvent) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.seq++
	ev.Seq = w.seq
	ev.Time = time.Now()
	w.ring = append(w.ring, ev)
	if len(w.ring) > ringSize {
		w.ring = w.ring[len(w.ring)-ringSize:]
	}
	w.cond.Broadcast()
}

// --- netlink proc connector -------------------------------------------------

const (
	cnIdxProc         = 1
	cnValProc         = 1
	procCnMcastListen = 1
	procEventExec     = 0x00000002
	nlMsgDone         = 0x3
)

func openProcConnector() (int, error) {
	sock, err := unix.Socket(unix.AF_NETLINK, unix.SOCK_DGRAM, unix.NETLINK_CONNECTOR)
	if err != nil {
		return -1, err
	}
	addr := &unix.SockaddrNetlink{Family: unix.AF_NETLINK, Groups: cnIdxProc, Pid: uint32(os.Getpid())}
	if err := unix.Bind(sock, addr); err != nil {
		unix.Close(sock)
		return -1, err
	}
	// Fork bursts (builds, package installs, parallel scripts) can outpace
	// the reader; a roomy receive buffer makes overflow (ENOBUFS) rare.
	// FORCE ignores rmem_max but needs CAP_NET_ADMIN - same capability the
	// connector itself needs, so the fallback is mostly for tests.
	if err := unix.SetsockoptInt(sock, unix.SOL_SOCKET, unix.SO_RCVBUFFORCE, 8<<20); err != nil {
		_ = unix.SetsockoptInt(sock, unix.SOL_SOCKET, unix.SO_RCVBUF, 8<<20)
	}
	if err := sendProcListen(sock, procCnMcastListen); err != nil {
		unix.Close(sock)
		return -1, err
	}
	return sock, nil
}

func sendProcListen(sock int, op uint32) error {
	// nlmsghdr + cn_msg + op
	buf := make([]byte, 16+20+4)
	le := binary.LittleEndian
	le.PutUint32(buf[0:], uint32(len(buf))) // nlmsg_len
	le.PutUint16(buf[4:], nlMsgDone)        // nlmsg_type
	le.PutUint32(buf[12:], uint32(os.Getpid()))
	// cn_msg: idx, val, seq, ack, len, flags
	le.PutUint32(buf[16:], cnIdxProc)
	le.PutUint32(buf[20:], cnValProc)
	le.PutUint16(buf[32:], 4) // data len
	le.PutUint32(buf[36:], op)
	addr := &unix.SockaddrNetlink{Family: unix.AF_NETLINK}
	return unix.Sendto(sock, buf, 0, addr)
}

func (w *ExecWatcher) netlinkLoop(sock int) {
	defer unix.Close(sock)
	buf := make([]byte, 65536)
	le := binary.LittleEndian
	var lastOverflowLog time.Time
	for {
		w.mu.Lock()
		closed := w.closed
		w.mu.Unlock()
		if closed {
			return
		}
		n, _, err := unix.Recvfrom(sock, buf, 0)
		if err != nil {
			if err == unix.EINTR {
				continue
			}
			if err == unix.ENOBUFS {
				// The kernel dropped events because our receive buffer
				// overflowed (fork burst). The socket is still healthy;
				// events in the gap are lost, which wait_* checks absorb
				// by re-polling - so keep reading, never disable.
				if time.Since(lastOverflowLog) > time.Minute {
					lastOverflowLog = time.Now()
					log.Printf("execwatch: netlink overflow (ENOBUFS): some exec events were lost")
				}
				continue
			}
			log.Printf("execwatch: netlink read error: %v; exec-based checks disabled", err)
			w.Source = ""
			return
		}
		// walk netlink messages
		for off := 0; off+16 <= n; {
			msgLen := int(le.Uint32(buf[off:]))
			if msgLen < 16 || off+msgLen > n {
				break
			}
			// cn_msg payload at off+16, proc_event at off+16+20
			pe := off + 16 + 20
			if pe+16 <= n {
				what := le.Uint32(buf[pe:])
				if what == procEventExec {
					// exec event: process pid at pe+16 (after what, cpu, timestamp[8])
					pid := int(le.Uint32(buf[pe+16:]))
					// Harvest concurrently: a burst of execs (shell pipelines,
					// login scripts) harvested serially would delay the later
					// /proc reads past the lifetime of short-lived commands
					// like `ss`, silently dropping their events.
					go func(pid int) {
						if ev, ok := harvestProc(pid); ok {
							w.publish(ev)
						}
					}(pid)
				}
			}
			off += nlmsgAlign(msgLen)
		}
	}
}

func nlmsgAlign(n int) int { return (n + 3) &^ 3 }

// harvestProc reads argv/uid/ppid of a freshly-exec'ed pid from /proc. The
// proc connector fires AFTER the new mm is installed, so /proc reflects the
// executed command. Very short-lived processes may be gone already - that
// is fine for reps where the interesting commands run at human speed.
func harvestProc(pid int) (ExecEvent, bool) {
	raw, err := os.ReadFile(fmt.Sprintf("/proc/%d/cmdline", pid))
	if err != nil || len(raw) == 0 {
		return ExecEvent{}, false
	}
	argv := splitNul(string(raw))
	status, err := os.ReadFile(fmt.Sprintf("/proc/%d/status", pid))
	if err != nil {
		return ExecEvent{}, false
	}
	// -1 = unknown: the process may exit before we can read its stat (very
	// short-lived commands like `ls`). Consumers must only reject a
	// CONFIRMED 0 - daemon-spawned scripts live long enough for the read
	// to succeed, so unknown means "probably a fast interactive command".
	ttyNr := -1
	if stat, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", pid)); err == nil {
		if f := statFields(string(stat)); len(f) > 6 {
			if v, err := strconv.Atoi(f[6]); err == nil {
				ttyNr = v
			}
		}
	}
	uid, ppid := -1, -1
	for _, line := range strings.Split(string(status), "\n") {
		if v, ok := strings.CutPrefix(line, "Uid:"); ok {
			f := strings.Fields(v)
			if len(f) > 0 {
				uid, _ = strconv.Atoi(f[0])
			}
		} else if v, ok := strings.CutPrefix(line, "PPid:"); ok {
			ppid, _ = strconv.Atoi(strings.TrimSpace(v))
		}
	}
	ev := ExecEvent{PID: pid, PPID: ppid, UID: uid, TTYNr: ttyNr, Argv: argv}
	if ttyNr != 0 {
		// Eager env capture (bounded): fast interactive commands are gone
		// before wait_env could read /proc lazily.
		if raw, err := os.ReadFile(fmt.Sprintf("/proc/%d/environ", pid)); err == nil && len(raw) > 0 {
			if len(raw) > 32*1024 {
				raw = raw[:32*1024]
			}
			ev.Env = splitNul(string(raw))
		}
	}
	return ev, true
}

// splitNul splits a NUL-separated /proc string (cmdline, environ).
func splitNul(s string) []string {
	return strings.Split(strings.TrimRight(s, "\x00"), "\x00")
}
