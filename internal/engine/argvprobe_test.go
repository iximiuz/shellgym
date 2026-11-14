package engine

import (
	"context"
	"fmt"
	"os/exec"
	"reflect"
	"strings"
	"testing"
	"time"
)

// argvRecord builds a trace_pipe record of the shell_execve probe the way
// tracefs prints it: quoted strings, (fault) for unreadable entries.
func argvRecord(pid int, cmd string, argv ...string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "            date-%d    [002] DNZff    29.289614: execve: (0x5627288dc530) cmd=%q", pid, cmd)
	for i := 0; i < argvProbeMaxArgs; i++ {
		if i < len(argv) {
			fmt.Fprintf(&b, " a%d=\"%s\"", i, argv[i])
		} else {
			fmt.Fprintf(&b, " a%d=(fault)", i)
		}
	}
	return b.String()
}

func TestParseArgvTraceLine(t *testing.T) {
	for _, tc := range []struct {
		name string
		line string
		pid  int
		argv []string
	}{
		{"simple", argvRecord(7035, "/usr/bin/date", "date", "--moonshine"), 7035, []string{"date", "--moonshine"}},
		{"one word", argvRecord(12, "/usr/bin/whoami", "whoami"), 12, []string{"whoami"}},
		{"spaces and quotes inside", argvRecord(9, "/usr/bin/echo", "echo", `say "hi" there`, "+%A %d %B"), 9, []string{"echo", `say "hi" there`, "+%A %d %B"}},
		{"empty argument", argvRecord(9, "/usr/bin/printf", "printf", "", "x"), 9, []string{"printf", "", "x"}},
		{"dashed comm", strings.Replace(argvRecord(42, "/usr/bin/kubectl", "kubectl", "get", "po"), "date-42", "my-tool-x-42", 1), 42, []string{"kubectl", "get", "po"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			pid, argv, ok := parseArgvTraceLine(tc.line)
			if !ok || pid != tc.pid || !reflect.DeepEqual(argv, tc.argv) {
				t.Fatalf("got pid=%d argv=%q ok=%v, want pid=%d argv=%q", pid, argv, ok, tc.pid, tc.argv)
			}
		})
	}
}

func TestParseArgvTraceLineRejects(t *testing.T) {
	full := make([]string, argvProbeMaxArgs)
	for i := range full {
		full[i] = fmt.Sprintf("x%d", i)
	}
	for name, line := range map[string]string{
		"readline record": `bash-1103  [002] DNZff  39.302011: readline: (0x55c5 <- 0x55c5) line="ls"`,
		"lost events":     "CPU:2 [LOST 3 EVENTS]",
		"no argv":         argvRecord(5, "/usr/bin/x"),
		"maybe truncated": argvRecord(5, "/usr/bin/x", full...),
		"cut short":       strings.SplitN(argvRecord(5, "/usr/bin/date", "date", "-u"), " a3=", 2)[0],
	} {
		if _, _, ok := parseArgvTraceLine(line); ok {
			t.Errorf("%s: parsed, want rejected", name)
		}
	}
}

// TestArgvProbeCatchesFastCommands runs commands that exit within about a
// millisecond (a date with an unknown option) through bash and expects
// every one of them as an exec event with its full argv. Needs root and
// tracefs with uprobes; run the suite as root on a playground.
func TestArgvProbeCatchesFastCommands(t *testing.T) {
	w := startWatcher(t)
	if err := w.StartArgvProbe("/bin/bash"); err != nil {
		t.Skipf("argv probe unavailable: %v", err)
	}
	after := w.Seq()

	const n = 200
	marker := fmt.Sprintf("--shellgym-probe-%d", time.Now().UnixNano())
	script := fmt.Sprintf("for i in $(seq 1 %d); do date %s-$i 2>/dev/null; done", n, marker)
	if out, err := exec.Command("bash", "-c", script).CombinedOutput(); err != nil && len(out) > 0 {
		t.Logf("bash: %v: %s", err, out)
	}

	seen := map[string]bool{}
	deadline := time.Now().Add(5 * time.Second)
	for len(seen) < n && time.Now().Before(deadline) {
		_, _ = w.WaitMatch(context.Background(), after, time.Now().Add(200*time.Millisecond), func(ev ExecEvent) bool {
			if len(ev.Argv) == 2 && ev.Argv[0] == "date" && strings.HasPrefix(ev.Argv[1], marker+"-") {
				seen[ev.Argv[1]] = true
			}
			return false
		})
	}
	if len(seen) != n {
		t.Fatalf("saw %d of %d fast commands", len(seen), n)
	}
}
