package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/creack/pty"
	"github.com/spf13/cobra"

	"github.com/iximiuz/labs-content/tools/shellgym/internal/content"
)

// solve is the acceptance-test driver: it spawns a REAL interactive shell
// on a pty (indistinguishable from a student's terminal to the daemon),
// walks the learning path, types each task's solve: script into the shell,
// and tracks unit completion through the daemon's API.
func newSolveCmd() *cobra.Command {
	var (
		api        string
		pathDir string
		unitFilter string
		timeout    time.Duration
	)
	cmd := &cobra.Command{
		Use:   "solve",
		Short: "Auto-solve a learning path through a real pty shell (acceptance test)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSolve(api, pathDir, unitFilter, timeout)
		},
	}
	cmd.Flags().StringVar(&api, "api", "http://localhost:63636", "daemon API base URL")
	cmd.Flags().StringVar(&pathDir, "path", "", "learning path directory (source of solve scripts; required)")
	cmd.Flags().StringVar(&unitFilter, "unit", "", "solve only this unit id (module/unit)")
	cmd.Flags().DurationVar(&timeout, "timeout", 2*time.Minute, "per-unit completion timeout")
	_ = cmd.MarkFlagRequired("path")
	return cmd
}

// studentShell is an interactive bash on a pty.
type studentShell struct {
	f   *os.File
	cmd *exec.Cmd
	mu  sync.Mutex
	out bytes.Buffer
	seq int
}

func newStudentShell() (*studentShell, error) {
	cmd := exec.Command("bash", "--norc", "-i")
	cmd.Dir = os.Getenv("HOME")
	cmd.Env = append(os.Environ(), "PS1=$ ", "TERM=dumb")
	f, err := pty.Start(cmd)
	if err != nil {
		return nil, err
	}
	s := &studentShell{f: f, cmd: cmd}
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := f.Read(buf)
			if n > 0 {
				s.mu.Lock()
				s.out.Write(buf[:n])
				if s.out.Len() > 1<<20 {
					// keep the tail only
					b := s.out.Bytes()
					s.out = *bytes.NewBuffer(append([]byte(nil), b[len(b)-1<<19:]...))
				}
				s.mu.Unlock()
			}
			if err != nil {
				return
			}
		}
	}()
	return s, nil
}

func (s *studentShell) Close() {
	_ = s.cmd.Process.Kill()
	_, _ = s.cmd.Process.Wait()
	_ = s.f.Close()
}

func (s *studentShell) contains(marker string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return bytes.Contains(s.out.Bytes(), []byte(marker))
}

// TypeSync types one command line and waits until the shell has finished
// executing it. A sync marker is chained onto the same line; the marker
// string is split in the sent text so the terminal's echo of the
// keystrokes cannot match the scan.
func (s *studentShell) TypeSync(line string, timeout time.Duration) error {
	s.seq++
	marker := fmt.Sprintf("__SYNC_%d__", s.seq)
	sep := "; "
	if strings.HasSuffix(strings.TrimSpace(line), "&") {
		sep = " "
	}
	typed := fmt.Sprintf("%s%secho __SY''NC_%d__\n", line, sep, s.seq)
	if _, err := s.f.WriteString(typed); err != nil {
		return err
	}
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if s.contains("\n"+marker) || s.contains("\r"+marker) {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("sync timeout after typing %q", line)
}

// --- daemon API client ------------------------------------------------------

type apiClient struct{ base string }

func (c *apiClient) get(path string, out any) error {
	resp, err := http.Get(c.base + path)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("GET %s: %s", path, resp.Status)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

func (c *apiClient) post(path string) error {
	resp, err := http.Post(c.base+path, "application/json", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("POST %s: %s: %s", path, resp.Status, strings.TrimSpace(string(body)))
	}
	return nil
}

type apiScene struct {
	Kind   string `json:"kind"`
	ID     string `json:"id"`
	Status string `json:"status"`
}

type apiPath struct {
	Scenes []apiScene `json:"scenes"`
}

type apiUnit struct {
	Vars map[string]string `json:"vars"`
}

// --- the walk ---------------------------------------------------------------

func runSolve(api, pathDir, unitFilter string, timeout time.Duration) error {
	// Load WITHOUT distro/capability filtering: the daemon's /api/path is
	// the authority on which units exist (its environment may differ from
	// where solve runs, e.g. a containerized daemon).
	path, err := content.Load(pathDir, "*", nil, []string{"*"})
	if err != nil {
		return err
	}

	sh, err := newStudentShell()
	if err != nil {
		return fmt.Errorf("start student shell: %w", err)
	}
	defer sh.Close()
	if err := sh.TypeSync("echo shell-ready", 10*time.Second); err != nil {
		return fmt.Errorf("student shell not responding: %w", err)
	}

	c := &apiClient{base: strings.TrimRight(api, "/")}
	var pathState apiPath
	if err := c.get("/api/path", &pathState); err != nil {
		return err
	}

	failed := 0
	for _, scene := range pathState.Scenes {
		if scene.Kind == "module" {
			_ = c.post("/api/module-seen/" + scene.ID)
			continue
		}
		if unitFilter != "" && scene.ID != unitFilter {
			continue
		}
		if scene.Status == "completed" {
			fmt.Printf("SKIP  %s (already completed)\n", scene.ID)
			continue
		}
		// Note: the daemon rejects activation of units whose needs: deps
		// are not solved, so a dependent of a failed (or filtered-out)
		// unit fails here too - solving it alone is meaningless anyway.
		if err := solveUnit(c, sh, path.Unit(scene.ID), timeout); err != nil {
			fmt.Printf("FAIL  %s (%v)\n", scene.ID, err)
			failed++
		} else {
			fmt.Printf("PASS  %s\n", scene.ID)
		}
	}
	if failed > 0 {
		return fmt.Errorf("%d unit(s) failed", failed)
	}
	return nil
}

func solveUnit(c *apiClient, sh *studentShell, u *content.Unit, timeout time.Duration) error {
	if u == nil {
		return fmt.Errorf("unit not present in local content")
	}
	if err := c.post("/api/activate/" + u.ID); err != nil {
		return err
	}
	var au apiUnit
	if err := c.get("/api/unit/"+u.ID, &au); err != nil {
		return err
	}

	// Export the unit's resolved vars into the student shell so solve
	// lines can reference them.
	names := make([]string, 0, len(au.Vars))
	for k := range au.Vars {
		names = append(names, k)
	}
	sort.Strings(names)
	for _, k := range names {
		if err := sh.TypeSync(fmt.Sprintf("export %s='%s'", k, au.Vars[k]), 30*time.Second); err != nil {
			return err
		}
	}

	time.Sleep(1 * time.Second) // give init tasks a moment

	// Type each task's solve script in task order.
	typed := false
	for _, t := range u.Tasks {
		for _, line := range strings.Split(t.Solve, "\n") {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			typed = true
			if err := sh.TypeSync(line, 5*time.Minute); err != nil {
				return err
			}
		}
	}
	if !typed {
		return fmt.Errorf("no solve script")
	}

	// Wait for unit completion via the API.
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		var p apiPath
		if err := c.get("/api/path", &p); err == nil {
			for _, s := range p.Scenes {
				if s.ID == u.ID && s.Status == "completed" {
					return nil
				}
			}
		}
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("not completed within %s", timeout)
}
