package content

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Load reads a learning path from dir. distro is the running distro id
// (e.g. "ubuntu"); units whose labels don't match it are dropped.
// distroLike lists ID_LIKE values ("debian" for ubuntu, etc.). caps lists
// host capabilities (e.g. "systemd"); units with unmet `requires:` are
// dropped.
func Load(dir string, distro string, distroLike []string, caps []string) (*Path, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	p := &Path{}
	rawPath, err := os.ReadFile(filepath.Join(dir, "path.yaml"))
	if err != nil {
		return nil, fmt.Errorf("read path.yaml: %w", err)
	}
	if err := yaml.Unmarshal(rawPath, p); err != nil {
		return nil, fmt.Errorf("parse path.yaml: %w", err)
	}
	if p.ID == "" {
		p.ID = filepath.Base(dir)
	}
	if p.ShellUser == "" {
		p.ShellUser = "laborant"
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		if !e.IsDir() || strings.HasPrefix(e.Name(), ".") {
			continue
		}
		mod, err := loadModule(filepath.Join(dir, e.Name()), e.Name(), distro, distroLike, caps)
		if err != nil {
			return nil, err
		}
		if len(mod.Units) > 0 || mod.Intro != "" {
			p.Modules = append(p.Modules, mod)
		}
	}
	sort.Slice(p.Modules, func(i, j int) bool { return p.Modules[i].Order < p.Modules[j].Order })

	if err := validate(p); err != nil {
		return nil, err
	}
	return p, nil
}

func loadModule(dir, folder, distro string, distroLike []string, caps []string) (*Module, error) {
	order, name, err := splitPrefix(folder)
	if err != nil {
		return nil, fmt.Errorf("module %s: %w", folder, err)
	}
	mod := &Module{ID: name, Name: name, Dir: dir, Order: order, Title: titleFromName(name)}

	if raw, err := os.ReadFile(filepath.Join(dir, "module.md")); err == nil {
		mod.Intro = string(raw)
		if t := firstHeading(mod.Intro); t != "" {
			mod.Title = t
		}
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		if !e.IsDir() || strings.HasPrefix(e.Name(), ".") {
			continue
		}
		u, err := loadUnit(filepath.Join(dir, e.Name()), e.Name(), mod.ID)
		if err != nil {
			return nil, err
		}
		if !distroMatch(u.Front.Labels, distro, distroLike) || !capsMatch(u.Front.Requires, caps) {
			continue
		}
		mod.Units = append(mod.Units, u)
	}
	sort.Slice(mod.Units, func(i, j int) bool { return mod.Units[i].Order < mod.Units[j].Order })
	return mod, nil
}

func loadUnit(dir, folder, moduleID string) (*Unit, error) {
	order, name, err := splitPrefix(folder)
	if err != nil {
		return nil, fmt.Errorf("unit %s: %w", folder, err)
	}
	raw, err := os.ReadFile(filepath.Join(dir, "unit.md"))
	if err != nil {
		return nil, fmt.Errorf("unit %s: %w", folder, err)
	}
	front, body, err := splitFrontmatter(string(raw))
	if err != nil {
		return nil, fmt.Errorf("unit %s: %w", folder, err)
	}

	u := &Unit{
		ID:       moduleID + "/" + name,
		Name:     name,
		ModuleID: moduleID,
		Dir:      dir,
		Order:    order,
		Body:     body,
	}
	if err := yaml.Unmarshal([]byte(front), &u.Front); err != nil {
		return nil, fmt.Errorf("unit %s: frontmatter: %w", folder, err)
	}
	if u.Front.Title == "" {
		return nil, fmt.Errorf("unit %s: missing title", folder)
	}

	names := make([]string, 0, len(u.Front.Tasks))
	for tname, t := range u.Front.Tasks {
		t.Name = tname
		if t.Mode == "" {
			t.Mode = ModeEdge
		}
		if t.Mode != ModeEdge && t.Mode != ModeLevel {
			return nil, fmt.Errorf("unit %s: task %s: bad mode %q", folder, tname, t.Mode)
		}
		names = append(names, tname)
	}
	sort.Strings(names)
	// Dependency-first (topological) order, name-sorted within ties: the
	// UI and `shellgym solve` both walk tasks in this order.
	added := map[string]bool{}
	var add func(n string)
	add = func(n string) {
		if added[n] {
			return
		}
		added[n] = true
		t := u.Front.Tasks[n]
		deps := append([]string(nil), t.Needs...)
		sort.Strings(deps)
		for _, d := range deps {
			if _, ok := u.Front.Tasks[d]; ok {
				add(d)
			}
		}
		u.Tasks = append(u.Tasks, t)
	}
	for _, tname := range names {
		add(tname)
	}
	return u, nil
}

func splitFrontmatter(raw string) (front, body string, err error) {
	if !strings.HasPrefix(raw, "---\n") {
		return "", raw, nil
	}
	rest := raw[len("---\n"):]
	idx := strings.Index(rest, "\n---\n")
	if idx < 0 {
		return "", "", fmt.Errorf("unterminated frontmatter")
	}
	return rest[:idx], rest[idx+len("\n---\n"):], nil
}

// capsMatch reports whether every required capability is present.
func capsMatch(requires, caps []string) bool {
	for _, c := range caps {
		if c == "*" {
			return true
		}
	}
	for _, r := range requires {
		found := false
		for _, c := range caps {
			if strings.EqualFold(r, c) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func distroMatch(labels []string, distro string, distroLike []string) bool {
	if len(labels) == 0 || distro == "*" {
		return true
	}
	for _, l := range labels {
		l = strings.ToLower(strings.TrimSpace(l))
		if l == distro {
			return true
		}
		for _, like := range distroLike {
			if l == like {
				return true
			}
		}
	}
	return false
}

func titleFromName(name string) string {
	name = strings.ReplaceAll(name, "-", " ")
	if name == "" {
		return name
	}
	return strings.ToUpper(name[:1]) + name[1:]
}

func firstHeading(md string) string {
	for _, line := range strings.Split(md, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") {
			return strings.TrimSpace(strings.TrimLeft(line, "#"))
		}
	}
	return ""
}

// validate enforces structural rules: known deps, no edge-on-level task
// dependencies, no dependency cycles, short unit dep chains.
func validate(p *Path) error {
	for _, m := range p.Modules {
		unitByName := map[string]*Unit{}
		for _, u := range m.Units {
			unitByName[u.Name] = u
		}
		for _, u := range m.Units {
			for _, need := range u.Front.Needs {
				dep, ok := unitByName[need]
				if !ok {
					return fmt.Errorf("unit %s: needs unknown unit %q (unit deps must stay within the module)", u.ID, need)
				}
				if dep.Order >= u.Order {
					return fmt.Errorf("unit %s: needs %q which does not precede it", u.ID, need)
				}
			}
			if depth := unitDepDepth(u, unitByName, 0); depth > 5 {
				return fmt.Errorf("unit %s: dependency chain longer than 5", u.ID)
			}
			for vn, spec := range u.Front.Vars {
				if spec.From == "" {
					continue
				}
				refUnit, refVar, ok := strings.Cut(spec.From, ".")
				if !ok {
					return fmt.Errorf("unit %s: var %s: from must be \"unit-name.VAR\"", u.ID, vn)
				}
				dep, exists := unitByName[refUnit]
				if !exists {
					return fmt.Errorf("unit %s: var %s: from references unknown unit %q", u.ID, vn, refUnit)
				}
				if dep.Order >= u.Order {
					return fmt.Errorf("unit %s: var %s: from must reference a preceding unit", u.ID, vn)
				}
				if _, exists := dep.Front.Vars[refVar]; !exists {
					return fmt.Errorf("unit %s: var %s: unit %q has no var %q", u.ID, vn, refUnit, refVar)
				}
			}
			if err := validateTasks(u); err != nil {
				return err
			}
		}
	}
	return nil
}

func unitDepDepth(u *Unit, byName map[string]*Unit, depth int) int {
	if depth > 5 {
		return depth
	}
	max := depth
	for _, need := range u.Front.Needs {
		if dep := byName[need]; dep != nil {
			if d := unitDepDepth(dep, byName, depth+1); d > max {
				max = d
			}
		}
	}
	return max
}

func validateTasks(u *Unit) error {
	byName := u.Front.Tasks
	for _, t := range u.Tasks {
		if strings.TrimSpace(t.Check) == "" {
			return fmt.Errorf("unit %s: task %s: empty check script", u.ID, t.Name)
		}
		for _, need := range t.Needs {
			dep, ok := byName[need]
			if !ok {
				return fmt.Errorf("unit %s: task %s: needs unknown task %q", u.ID, t.Name, need)
			}
			if t.Mode == ModeEdge && dep.Mode == ModeLevel {
				return fmt.Errorf("unit %s: edge task %s may not depend on level task %s", u.ID, t.Name, need)
			}
		}
	}
	// cycle check via DFS
	const (
		white = 0
		grey  = 1
		black = 2
	)
	color := map[string]int{}
	var visit func(name string) error
	visit = func(name string) error {
		switch color[name] {
		case grey:
			return fmt.Errorf("unit %s: task dependency cycle involving %q", u.ID, name)
		case black:
			return nil
		}
		color[name] = grey
		for _, need := range byName[name].Needs {
			if _, ok := byName[need]; !ok {
				continue
			}
			if err := visit(need); err != nil {
				return err
			}
		}
		color[name] = black
		return nil
	}
	for name := range byName {
		if err := visit(name); err != nil {
			return err
		}
	}
	return nil
}
