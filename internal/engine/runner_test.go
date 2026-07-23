package engine

import (
	"context"
	"strings"
	"testing"
	"time"
)

func testRunner(t *testing.T) *scriptRunner {
	t.Helper()
	dir := t.TempDir()
	return &scriptRunner{checksDir: dir, sockPath: dir + "/none.sock"}
}

func TestRunCapturesStreams(t *testing.T) {
	r := testRunner(t)
	res := r.Run(context.Background(), "echo out; echo err >&2; exit 3",
		map[string]string{"FOO": "bar"}, 5*time.Second)
	if res.ExitCode != 3 {
		t.Errorf("exit = %d", res.ExitCode)
	}
	if strings.TrimSpace(res.Stdout) != "out" || strings.TrimSpace(res.Stderr) != "err" {
		t.Errorf("streams: %q %q", res.Stdout, res.Stderr)
	}
}

func TestRunEnvAndPath(t *testing.T) {
	r := testRunner(t)
	res := r.Run(context.Background(), `echo "$FOO $GYM_SOCK"; echo "$PATH"`,
		map[string]string{"FOO": "bar"}, 5*time.Second)
	if res.ExitCode != 0 {
		t.Fatalf("exit = %d, stderr: %s", res.ExitCode, res.Stderr)
	}
	if !strings.Contains(res.Stdout, "bar") || !strings.Contains(res.Stdout, "none.sock") {
		t.Errorf("env not passed: %q", res.Stdout)
	}
	if !strings.Contains(res.Stdout, r.checksDir) {
		t.Errorf("checks dir not on PATH: %q", res.Stdout)
	}
}

func TestRunTimeoutKillsProcessTree(t *testing.T) {
	r := testRunner(t)
	start := time.Now()
	res := r.Run(context.Background(), "sleep 30 & wait", nil, 500*time.Millisecond)
	if elapsed := time.Since(start); elapsed > 3*time.Second {
		t.Fatalf("timeout not enforced: %v", elapsed)
	}
	if !res.TimedOut || res.ExitCode != 124 {
		t.Errorf("res: %+v", res)
	}
}

func TestRunContextCancel(t *testing.T) {
	r := testRunner(t)
	ctx, cancel := context.WithCancel(context.Background())
	go func() { time.Sleep(200 * time.Millisecond); cancel() }()
	start := time.Now()
	res := r.Run(ctx, "sleep 30", nil, 0)
	if time.Since(start) > 3*time.Second {
		t.Fatal("cancel not honored")
	}
	if res.TimedOut {
		t.Error("cancel must not be reported as timeout")
	}
}

func TestBackgroundChildDoesNotBlockWait(t *testing.T) {
	// A script that leaves a detached child inheriting stdout/stderr must
	// return as soon as the script itself exits.
	r := testRunner(t)
	start := time.Now()
	res := r.Run(context.Background(),
		"setsid sleep 20 & echo started", nil, 10*time.Second)
	if elapsed := time.Since(start); elapsed > 3*time.Second {
		t.Fatalf("Wait blocked on background child: %v", elapsed)
	}
	if res.ExitCode != 0 || !strings.Contains(res.Stdout, "started") {
		t.Errorf("res: %+v", res)
	}
}

func TestPipefail(t *testing.T) {
	r := testRunner(t)
	res := r.Run(context.Background(), "false | cat", nil, 5*time.Second)
	if res.ExitCode == 0 {
		t.Error("pipefail not active")
	}
}
