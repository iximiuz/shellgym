package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
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
		pathDir    string
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
	// promptCount is the number of prompts the shell has printed so far,
	// counted from the output stream as it arrives (see syncMarker).
	promptCount int
}

// syncMarker is printed by the solve shell's PROMPT_COMMAND right before
// every prompt. Counting its occurrences tells when a typed line has
// finished executing, without altering the line itself - the daemon's
// line watcher sees exactly what a student would have typed, so
// $-anchored wait_line checks pass under solve too. The marker is never
// typed, so keystroke echo cannot inflate the count.
const syncMarker = "__SHELLGYM_SYNC__"

func newStudentShell() (*studentShell, error) {
	cmd := exec.Command("bash", "--norc", "-i")
	cmd.Dir = os.Getenv("HOME")
	cmd.Env = append(os.Environ(), "PS1=$ ", "TERM=dumb", "PROMPT_COMMAND=echo "+syncMarker)
	f, err := pty.Start(cmd)
	if err != nil {
		return nil, err
	}
	s := &studentShell{f: f, cmd: cmd}
	go func() {
		marker := []byte(syncMarker)
		buf := make([]byte, 4096)
		var carry []byte // tail of the previous chunk, for markers split across reads
		for {
			n, err := f.Read(buf)
			if n > 0 {
				data := append(carry, buf[:n]...)
				if c := bytes.Count(data, marker); c > 0 {
					s.mu.Lock()
					s.promptCount += c
					s.mu.Unlock()
				}
				if len(data) > len(marker)-1 {
					data = data[len(data)-(len(marker)-1):]
				}
				carry = append([]byte(nil), data...)
			}
			if err != nil {
				return
			}
		}
	}()
	// Wait for the first prompt: TypeSync counts prompts, so typing must
	// only ever start with the shell idle at one.
	for deadline := time.Now().Add(10 * time.Second); s.prompts() == 0; {
		if time.Now().After(deadline) {
			s.Close()
			return nil, fmt.Errorf("the shell printed no prompt within 10s")
		}
		time.Sleep(50 * time.Millisecond)
	}
	return s, nil
}

func (s *studentShell) Close() {
	_ = s.cmd.Process.Kill()
	_, _ = s.cmd.Process.Wait()
	_ = s.f.Close()
}

// prompts returns the number of prompts the shell has printed so far.
func (s *studentShell) prompts() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.promptCount
}

// TypeSync types one command line exactly as given and waits until the
// shell prints its next prompt, i.e. has finished executing the line (a
// trailing & returns to the prompt at once, as for a student).
func (s *studentShell) TypeSync(line string, timeout time.Duration) error {
	before := s.prompts()
	if _, err := s.f.WriteString(line + "\n"); err != nil {
		return err
	}
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if s.prompts() > before {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("sync timeout after typing %q", line)
}

// Directive executes a `#!` solve-script control line - the escape hatch for
// interactions a plain typed-and-synced line cannot express: signal keys
// (Ctrl-C, Ctrl-Z), pager keystrokes, and commands that hold the foreground
// so a chained sync marker would never print.
//
//	#!type TEXT     write TEXT to the pty verbatim (no Enter, no sync)
//	#!keys K K...   send named keys: enter, tab, space, esc, C-<letter>
//	                (or ctrl-<letter>), or any single literal character
//	#!wait SECONDS  pause to let the previous keystrokes take effect
func (s *studentShell) Directive(d string) error {
	verb, rest, _ := strings.Cut(d, " ")
	rest = strings.TrimSpace(rest)
	switch verb {
	case "type":
		if rest == "" {
			return fmt.Errorf("nothing to type")
		}
		_, err := s.f.WriteString(rest)
		return err
	case "keys":
		if rest == "" {
			return fmt.Errorf("no keys given")
		}
		for _, k := range strings.Fields(rest) {
			b, err := keyBytes(k)
			if err != nil {
				return err
			}
			if _, err := s.f.Write(b); err != nil {
				return err
			}
			time.Sleep(50 * time.Millisecond)
		}
		return nil
	case "wait":
		secs, err := strconv.ParseFloat(rest, 64)
		if err != nil || secs < 0 || secs > 60 {
			return fmt.Errorf("bad wait duration %q", rest)
		}
		time.Sleep(time.Duration(secs * float64(time.Second)))
		return nil
	}
	return fmt.Errorf("unknown directive %q", verb)
}

func keyBytes(name string) ([]byte, error) {
	switch strings.ToLower(name) {
	case "enter", "return":
		return []byte{'\r'}, nil
	case "tab":
		return []byte{'\t'}, nil
	case "space":
		return []byte{' '}, nil
	case "esc":
		return []byte{0x1b}, nil
	}
	lower := strings.ToLower(name)
	if ctrl, ok := strings.CutPrefix(lower, "c-"); ok || strings.HasPrefix(lower, "ctrl-") {
		if !ok {
			ctrl = strings.TrimPrefix(lower, "ctrl-")
		}
		if len(ctrl) == 1 && ctrl[0] >= 'a' && ctrl[0] <= 'z' {
			return []byte{ctrl[0] & 0x1f}, nil
		}
		return nil, fmt.Errorf("unknown control key %q", name)
	}
	if len([]rune(name)) == 1 {
		return []byte(name), nil
	}
	return nil, fmt.Errorf("unknown key %q", name)
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
	Kind        string `json:"kind"`
	ID          string `json:"id"`
	Status      string `json:"status"`
	Unsupported bool   `json:"unsupported"`
	Variant     string `json:"variant"`
	Hidden      bool   `json:"hidden"`
}

type apiPath struct {
	Scenes []apiScene `json:"scenes"`
}

type apiTask struct {
	Name   string `json:"name"`
	Mode   string `json:"mode"`
	Status string `json:"status"`
}

type apiUnit struct {
	Status string            `json:"status"`
	Vars   map[string]string `json:"vars"`
	Tasks  []apiTask         `json:"tasks"`
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
	// Ask for hidden units too: every variant of the path gets solved in
	// one walk, not just the one this attempt drew.
	var pathState apiPath
	if err := c.get("/api/path?hidden=1", &pathState); err != nil {
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
		label := scene.ID
		if scene.Variant != "" {
			label += " [" + scene.Variant
			if scene.Hidden {
				label += ", hidden"
			}
			label += "]"
		}
		if scene.Status == "completed" {
			fmt.Printf("SKIP  %s (already completed)\n", label)
			continue
		}
		if scene.Unsupported {
			fmt.Printf("SKIP  %s (not supported by the daemon's environment)\n", label)
			continue
		}
		// Note: the daemon rejects activation of units whose needs: deps
		// are not solved, so a dependent of a failed (or filtered-out)
		// unit fails here too - solving it alone is meaningless anyway.
		if err := solveUnit(c, sh, path.Unit(scene.ID), timeout); err != nil {
			fmt.Printf("FAIL  %s (%v)\n", label, err)
			failed++
		} else {
			fmt.Printf("PASS  %s\n", label)
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

	// Type each task's solve script in task order, waiting for the task to
	// pass before starting the next one - exactly like a student watching
	// the task box tick. Typing straight through would race the pollers:
	// a transient state (a cwd passed through, a file created then moved)
	// can flip faster than a ~200ms poll notices, and a needs-gated task
	// does not even start checking until its deps complete.
	deadline := time.Now().Add(timeout)
	typed := false
	for i, t := range u.Tasks {
		for _, line := range strings.Split(t.Solve, "\n") {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			line = expandVars(line, au.Vars)
			if d, ok := strings.CutPrefix(line, "#!"); ok {
				typed = true
				if err := sh.Directive(strings.TrimSpace(d)); err != nil {
					return fmt.Errorf("solve directive %q: %w", line, err)
				}
				continue
			}
			if strings.HasPrefix(line, "#") {
				continue
			}
			typed = true
			if err := sh.TypeSync(line, 5*time.Minute); err != nil {
				return err
			}
		}
		// Level tasks flip back and forth by design, and the LAST task of
		// the unit is covered by the unit-completion wait below.
		if t.Mode == content.ModeLevel || i == len(u.Tasks)-1 {
			continue
		}
		if err := waitTask(c, u.ID, t.Name, deadline); err != nil {
			return err
		}
	}
	if !typed {
		return fmt.Errorf("no solve script")
	}

	// Wait for unit completion via the API. Ask the unit itself, not the
	// path listing: a unit hidden by the variant draw is not listed there.
	for time.Now().Before(deadline) {
		var au apiUnit
		if err := c.get("/api/unit/"+u.ID, &au); err == nil && au.Status == "completed" {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("not completed within %s", timeout)
}

// expandVars substitutes the unit's vars ($NAME and ${NAME}) into a solve
// line before it is typed. A student types the literal value, and
// wait_line checks judge the line exactly as typed - a reference like
// `sleep $PAUSE && hostname` would never match what the unit asks for.
// Other `$` references are left for bash to expand (the vars are exported
// into the solve shell too, for uses like `$((PAUSE + 1))`).
//
// The line is scanned once, left to right, and every reference is
// resolved against the original text: a value that was just substituted
// is never rescanned, and adjacent references such as `$LO$HI` both
// resolve (a per-variable pass would turn `$LO$HI` into `$LO4` first, and
// the word-boundary match for `$LO` would then fail).
func expandVars(line string, vars map[string]string) string {
	return varRef.ReplaceAllStringFunc(line, func(ref string) string {
		m := varRef.FindStringSubmatch(ref)
		name := m[1]
		if name == "" {
			name = m[2]
		}
		if v, ok := vars[name]; ok {
			return v
		}
		return ref
	})
}

// varRef matches one `${NAME}` (group 1) or `$NAME` (group 2) reference.
// The greedy identifier match makes `$PAUSE_MAX` one reference, never
// `$PAUSE` followed by `_MAX`.
var varRef = regexp.MustCompile(`\$(?:\{([A-Za-z_][A-Za-z0-9_]*)\}|([A-Za-z_][A-Za-z0-9_]*))`)

// waitTask polls the unit API until the named task is completed (or, for
// level tasks, currently satisfied).
func waitTask(c *apiClient, unitID, task string, deadline time.Time) error {
	for time.Now().Before(deadline) {
		var au apiUnit
		if err := c.get("/api/unit/"+unitID, &au); err == nil {
			for _, t := range au.Tasks {
				if t.Name == task && (t.Status == "completed" || t.Status == "satisfied") {
					return nil
				}
			}
		}
		time.Sleep(300 * time.Millisecond)
	}
	return fmt.Errorf("task %s not completed before the unit deadline", task)
}
