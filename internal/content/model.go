// Package content implements the content engine: it defines the on-disk
// format of learning paths (modules, units, tasks, vars) and renders unit
// markdown to HTML. It knows nothing about how tasks are executed.
package content

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
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
	// "systemd", "python3", and "readline"). Units
	// whose requirements the runtime lacks are marked Unsupported: still
	// shown, never run.
	Requires []string `yaml:"requires"`
	// Variant ("key=value") makes the unit part of a variant of the path:
	// for every key one value is drawn per attempt, and units carrying
	// another value of that key are hidden. See Path.ApplyVariants.
	Variant string `yaml:"variant"`
	// StartOver gives the student a "Start over" button on the unit while
	// it is in progress: it forgets the unit's progress and re-runs its
	// init with the same vars (see engine.StartOverUnit). For scenes the
	// student can break beyond repair.
	StartOver bool               `yaml:"startOver"`
	Vars      map[string]VarSpec `yaml:"vars"`
	Init      []InitTask         `yaml:"init"`
	Tasks     map[string]*Task   `yaml:"tasks"`
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
	// MissingCaps lists `requires:` capabilities absent on this host.
	MissingCaps []string
	// Unsupported marks a unit that cannot run in this environment: it has
	// missing capabilities itself, or (transitively) builds on a unit that
	// does. Unsupported units stay browsable but are never activated.
	Unsupported bool
	// Variant is the parsed `variant:` field (zero when unset).
	Variant Variant
	// Hidden marks a unit whose variant value was not drawn for this
	// attempt (see Path.ApplyVariants): it is left out of the student's
	// path, but stays in the model - loadable, activatable - so authoring
	// tools can still reach every unit.
	Hidden bool
}

// Variant is a unit's "key=value" membership in a path variant.
type Variant struct {
	Key   string
	Value string
}

// IsZero reports whether the unit belongs to no variant (always shown).
func (v Variant) IsZero() bool { return v.Key == "" }

func (v Variant) String() string {
	if v.IsZero() {
		return ""
	}
	return v.Key + "=" + v.Value
}

var variantRe = regexp.MustCompile(`^([a-z0-9][a-z0-9-]*)=([a-z0-9][a-z0-9-]*)$`)

// ParseVariant parses a "key=value" variant spec (lowercase letters,
// digits, and dashes on both sides).
func ParseVariant(spec string) (Variant, error) {
	m := variantRe.FindStringSubmatch(strings.TrimSpace(spec))
	if m == nil {
		return Variant{}, fmt.Errorf("variant %q: want key=value (lowercase letters, digits, and dashes)", spec)
	}
	return Variant{Key: m[1], Value: m[2]}, nil
}

// VariantKeys lists every variant key used in the path with its values
// (sorted, deduplicated) - the pool each per-attempt draw picks from.
func (p *Path) VariantKeys() map[string][]string {
	seen := map[string]map[string]bool{}
	for _, m := range p.Modules {
		for _, u := range m.Units {
			if u.Variant.IsZero() {
				continue
			}
			if seen[u.Variant.Key] == nil {
				seen[u.Variant.Key] = map[string]bool{}
			}
			seen[u.Variant.Key][u.Variant.Value] = true
		}
	}
	out := make(map[string][]string, len(seen))
	for k, vals := range seen {
		for v := range vals {
			out[k] = append(out[k], v)
		}
		sort.Strings(out[k])
	}
	return out
}

// ApplyVariants marks units Hidden according to the drawn values: a unit
// with variant key=value stays visible only if picks[key] == value. Units
// without a variant, and units whose key has no pick, are always visible.
func (p *Path) ApplyVariants(picks map[string]string) {
	for _, m := range p.Modules {
		for _, u := range m.Units {
			if u.Variant.IsZero() {
				u.Hidden = false
				continue
			}
			pick, ok := picks[u.Variant.Key]
			u.Hidden = ok && pick != u.Variant.Value
		}
	}
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

// Scenes returns the linear scene sequence of the path, hidden units
// included (callers presenting the path to a student skip Unit.Hidden).
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

var prefixRe = regexp.MustCompile(`^(\d+)\.([a-z0-9][a-z0-9-]*)$`)

func splitPrefix(folder string) (order int, name string, err error) {
	m := prefixRe.FindStringSubmatch(folder)
	if m == nil {
		return 0, "", fmt.Errorf("folder %q: want NNN.name format (name: lowercase letters, digits, and dashes)", folder)
	}
	fmt.Sscanf(m[1], "%d", &order)
	return order, m[2], nil
}
