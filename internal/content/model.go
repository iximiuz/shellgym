// Package content implements the content engine: it defines the on-disk
// format of learning paths (modules, units, tasks, vars) and renders unit
// markdown to HTML. It knows nothing about how tasks are executed.
package content

import (
	"fmt"
	"regexp"
)

// TaskMode determines how the validation engine treats a task.
type TaskMode string

const (
	// ModeEdge tasks run (typically blocking on a wait_* check) until they
	// exit 0 once; after that they are completed forever.
	ModeEdge TaskMode = "edge"
	// ModeLevel tasks are polled; their status reflects the current system
	// state and may flip back and forth until the unit completes.
	ModeLevel TaskMode = "level"
)

// VarSpec defines a unit parameter. Exactly one field must be set.
type VarSpec struct {
	Value string   `yaml:"value"` // fixed value
	Pick  []string `yaml:"pick"`  // random choice from a list
	Shell string   `yaml:"shell"` // output of a shell command
	// From inherits a var from a preceding unit in the same module
	// ("unit-name.VAR"), so dependent units share randomized state.
	From string `yaml:"from"`
}

// InitTask prepares the system before the unit is presented.
type InitTask struct {
	Name string `yaml:"name"`
	Run  string `yaml:"run"`
}

// Task is a verification task.
type Task struct {
	Name    string   // map key in frontmatter
	Mode    TaskMode `yaml:"mode"`
	Needs   []string `yaml:"needs"`
	Check   string   `yaml:"check"`
	Hint    string   `yaml:"hint"`    // optional dynamic-hint script
	Timeout int      `yaml:"timeout"` // per-attempt timeout, seconds (0 = engine default)
	// Solve holds the reference solution: shell lines typed into the
	// student shell by `shellgym solve`. Stripped from on-disk files in
	// --live mode and never exposed by the UI.
	Solve string `yaml:"solve"`
}

// Frontmatter is the YAML preamble of unit.md.
type Frontmatter struct {
	Title  string   `yaml:"title"`
	Labels []string `yaml:"labels"`
	Needs  []string `yaml:"needs"`
	// Requires lists host capabilities the unit depends on (currently:
	// "systemd"). Units whose requirements the runtime lacks are dropped
	// at load time.
	Requires []string           `yaml:"requires"`
	Vars     map[string]VarSpec `yaml:"vars"`
	Init     []InitTask         `yaml:"init"`
	Tasks    map[string]*Task   `yaml:"tasks"`
}

// Unit is a single scene with one or more tasks.
type Unit struct {
	ID       string // "module/unit", prefixes stripped
	Name     string // folder name without prefix
	ModuleID string
	Dir      string // absolute folder path
	Order    int    // numeric prefix
	Front    Frontmatter
	Body     string  // raw markdown body
	Tasks    []*Task // frontmatter tasks in stable (name-sorted, deps-checked) order
}

// Module groups units; may have an intro scene (module.md).
type Module struct {
	ID    string
	Name  string
	Dir   string
	Order int
	Title string // first heading of module.md, or derived from name
	Intro string // raw markdown of module.md ("" = none)
	Units []*Unit
}

// Path is a whole learning path.
type Path struct {
	ID          string `yaml:"id"`
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	// ShellUser is the login user whose shells are observed (e.g. laborant).
	ShellUser string `yaml:"shellUser"`
	Modules   []*Module
}

// Scene is one screen in the horizontal path: a module intro or a unit.
type Scene struct {
	Kind   string // "module" | "unit"
	Module *Module
	Unit   *Unit
}

// Scenes returns the linear scene sequence of the path.
func (p *Path) Scenes() []Scene {
	var out []Scene
	for _, m := range p.Modules {
		if m.Intro != "" {
			out = append(out, Scene{Kind: "module", Module: m})
		}
		for _, u := range m.Units {
			out = append(out, Scene{Kind: "unit", Module: m, Unit: u})
		}
	}
	return out
}

// Unit looks a unit up by its "module/unit" id.
func (p *Path) Unit(id string) *Unit {
	for _, m := range p.Modules {
		for _, u := range m.Units {
			if u.ID == id {
				return u
			}
		}
	}
	return nil
}

// Module looks a module up by id.
func (p *Path) Module(id string) *Module {
	for _, m := range p.Modules {
		if m.ID == id {
			return m
		}
	}
	return nil
}

var prefixRe = regexp.MustCompile(`^(\d+)\.(.+)$`)

func splitPrefix(folder string) (order int, name string, err error) {
	m := prefixRe.FindStringSubmatch(folder)
	if m == nil {
		return 0, "", fmt.Errorf("folder %q: want NNN.name format", folder)
	}
	fmt.Sscanf(m[1], "%d", &order)
	return order, m[2], nil
}
