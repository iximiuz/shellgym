package engine

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"sync"
	"syscall"
	"time"

	"github.com/iximiuz/labs-content/tools/shellgym/internal/checkclient"
	"github.com/iximiuz/labs-content/tools/shellgym/internal/state"
)

// RunResult is the outcome of one script execution.
type RunResult struct {
	ExitCode  int
	Stdout    string
	Stderr    string
	Duration  time.Duration
	TimedOut  bool
	StartedAt time.Time
}

// scriptRunner executes task scripts with the gym environment: built-in
// checks on PATH, unit vars exported, daemon socket available.
type scriptRunner struct {
	checksDir string
	sockPath  string
}

// WriteCheckShims materializes the built-in check commands as tiny shell
// shims in dir, which is prepended to every task script's PATH. selfExe is
// the shellgym binary path.
func WriteCheckShims(dir, selfExe string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	for _, name := range checkclient.Names {
		shim := fmt.Sprintf("#!/bin/sh\nexec %q check %s \"$@\"\n", selfExe, name)
		if err := os.WriteFile(filepath.Join(dir, name), []byte(shim), 0o755); err != nil {
			return err
		}
	}
	return nil
}

// Run executes script with the given extra env. A timeout of 0 means no
// deadline (the context still allows cancellation on unit switch/shutdown).
func (r *scriptRunner) Run(ctx context.Context, script string, env map[string]string, timeout time.Duration) RunResult {
	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}
	// The prelude turns hint_exit into a script-terminating construct: the
	// PATH shim (a child process) posts the hint and returns HintExitCode;
	// the function wrapper then exits the calling script with that code.
	prelude := fmt.Sprintf("hint_exit() { command hint_exit \"$@\"; exit %d; }\n",
		checkclient.HintExitCode)
	cmd := exec.Command("bash", "-o", "pipefail", "-c", prelude+script)
	cmd.Env = append(os.Environ(),
		"PATH="+r.checksDir+":"+os.Getenv("PATH"),
		"GYM_SOCK="+r.sockPath,
	)
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		cmd.Env = append(cmd.Env, k+"="+env[k])
	}
	// Real pipe FDs (not exec.Cmd's copy goroutines): a script may leave a
	// long-lived background child that inherits stdout/stderr, and Wait()
	// must NOT block until that grandchild exits.
	outR, outW, err := os.Pipe()
	if err != nil {
		return RunResult{StartedAt: time.Now(), ExitCode: 127, Stderr: err.Error()}
	}
	errR, errW, err := os.Pipe()
	if err != nil {
		outR.Close()
		outW.Close()
		return RunResult{StartedAt: time.Now(), ExitCode: 127, Stderr: err.Error()}
	}
	cmd.Stdout = outW
	cmd.Stderr = errW
	// Own session (implies own process group, and - crucially - NO
	// controlling terminal, so daemon-run scripts are never mistaken for
	// student activity by the exec watcher) + killable as a tree.
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}

	res := RunResult{StartedAt: time.Now()}
	if err := cmd.Start(); err != nil {
		outR.Close()
		outW.Close()
		errR.Close()
		errW.Close()
		res.ExitCode = 127
		res.Stderr = err.Error()
		return res
	}
	outW.Close()
	errW.Close()

	var mu sync.Mutex
	var stdout, stderr bytes.Buffer
	drain := func(dst *bytes.Buffer, src *os.File) {
		buf := make([]byte, 4096)
		for {
			n, err := src.Read(buf)
			if n > 0 {
				mu.Lock()
				dst.Write(buf[:n])
				mu.Unlock()
			}
			if err != nil {
				return
			}
		}
	}
	go drain(&stdout, outR)
	go drain(&stderr, errR)

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case err := <-done:
		res.ExitCode = exitCode(err)
	case <-ctx.Done():
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		<-done
		res.TimedOut = ctx.Err() == context.DeadlineExceeded
		res.ExitCode = 124
	}
	// Give drains a moment to pick up trailing output, then snapshot.
	// Lingering background writers keep the pipes open - close our read
	// ends so the drain goroutines terminate instead of leaking.
	time.Sleep(50 * time.Millisecond)
	mu.Lock()
	res.Stdout = stdout.String()
	res.Stderr = stderr.String()
	mu.Unlock()
	outR.Close()
	errR.Close()

	res.Duration = time.Since(res.StartedAt)
	return res
}

func exitCode(err error) int {
	if err == nil {
		return 0
	}
	if ee, ok := err.(*exec.ExitError); ok {
		return ee.ExitCode()
	}
	return 126
}

func toTaskRun(res RunResult, kind string) state.TaskRun {
	return state.TaskRun{
		StartedAt: res.StartedAt,
		Duration:  res.Duration.Seconds(),
		ExitCode:  res.ExitCode,
		Stdout:    res.Stdout,
		Stderr:    res.Stderr,
		Kind:      kind,
		TimedOut:  res.TimedOut,
	}
}
