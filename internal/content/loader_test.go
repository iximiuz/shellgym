package content

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeFile(t *testing.T, path, data string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
}

func scaffold(t *testing.T, units map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "path.yaml"), "id: test-path\ntitle: Test\n")
	for rel, content := range units {
		writeFile(t, filepath.Join(dir, rel), content)
	}
	return dir
}

const minimalUnit = `---
title: A unit
tasks:
  t1:
    check: |
      true
---
Body text.
`

func TestLoadBasic(t *testing.T) {
	dir := scaffold(t, map[string]string{
		"010.mod-a/module.md":              "# Module A\nIntro.",
		"010.mod-a/010.unit-one/unit.md":   minimalUnit,
		"010.mod-a/020.unit-two/unit.md":   minimalUnit,
		"020.mod-b/010.unit-three/unit.md": minimalUnit,
	})
	p, err := Load(dir, "ubuntu", []string{"debian"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if p.ID != "test-path" || len(p.Modules) != 2 {
		t.Fatalf("unexpected path: %+v", p)
	}
	if p.Modules[0].ID != "mod-a" || p.Modules[0].Title != "Module A" {
		t.Fatalf("module: %+v", p.Modules[0])
	}
	if got := p.Modules[0].Units[0].ID; got != "mod-a/unit-one" {
		t.Fatalf("unit id: %s", got)
	}
	scenes := p.Scenes()
	if len(scenes) != 4 { // 1 module intro + 3 units
		t.Fatalf("scenes: %d", len(scenes))
	}
	if scenes[0].Kind != "module" || scenes[1].Unit.Name != "unit-one" {
		t.Fatalf("scene order broken")
	}
	if u := p.Unit("mod-b/unit-three"); u == nil || u.ModuleID != "mod-b" {
		t.Fatalf("Unit() lookup failed")
	}
}

func TestLoadDistroFilter(t *testing.T) {
	dir := scaffold(t, map[string]string{
		"010.m/010.everywhere/unit.md": minimalUnit,
		"010.m/020.debian-only/unit.md": strings.Replace(minimalUnit,
			"title: A unit", "title: A unit\nlabels: [ubuntu, debian]", 1),
		"010.m/030.rpm-only/unit.md": strings.Replace(minimalUnit,
			"title: A unit", "title: A unit\nlabels: [rocky, fedora]", 1),
	})
	p, err := Load(dir, "ubuntu", []string{"debian"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if n := len(p.Modules[0].Units); n != 2 {
		t.Fatalf("want 2 units on ubuntu, got %d", n)
	}
	p, err = Load(dir, "rocky", []string{"rhel", "centos", "fedora"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, u := range p.Modules[0].Units {
		names = append(names, u.Name)
	}
	if len(names) != 2 || names[1] != "rpm-only" {
		t.Fatalf("rocky units: %v", names)
	}
}

func TestValidateEdgeOnLevel(t *testing.T) {
	bad := `---
title: Bad
tasks:
  lvl:
    mode: level
    check: |
      true
  edge:
    needs: [lvl]
    check: |
      true
---
x
`
	dir := scaffold(t, map[string]string{"010.m/010.u/unit.md": bad})
	if _, err := Load(dir, "ubuntu", nil, nil); err == nil ||
		!strings.Contains(err.Error(), "may not depend on level") {
		t.Fatalf("want edge-on-level error, got %v", err)
	}
}

func TestValidateUnknownNeed(t *testing.T) {
	bad := strings.Replace(minimalUnit, "title: A unit", "title: A unit\nneeds: [nope]", 1)
	dir := scaffold(t, map[string]string{"010.m/010.u/unit.md": bad})
	if _, err := Load(dir, "ubuntu", nil, nil); err == nil ||
		!strings.Contains(err.Error(), "unknown unit") {
		t.Fatalf("want unknown-unit error, got %v", err)
	}
}

func TestValidateTaskCycle(t *testing.T) {
	bad := `---
title: Cycle
tasks:
  a:
    needs: [b]
    check: |
      true
  b:
    needs: [a]
    check: |
      true
---
x
`
	dir := scaffold(t, map[string]string{"010.m/010.u/unit.md": bad})
	if _, err := Load(dir, "ubuntu", nil, nil); err == nil ||
		!strings.Contains(err.Error(), "cycle") {
		t.Fatalf("want cycle error, got %v", err)
	}
}

func TestValidateFromRefs(t *testing.T) {
	first := `---
title: First
vars:
  TOKEN: { value: abc }
tasks:
  t:
    check: |
      true
---
x
`
	second := `---
title: Second
needs: [first]
vars:
  TOKEN: { from: first.TOKEN }
tasks:
  t:
    check: |
      true
---
x
`
	dir := scaffold(t, map[string]string{
		"010.m/010.first/unit.md":  first,
		"010.m/020.second/unit.md": second,
	})
	if _, err := Load(dir, "ubuntu", nil, nil); err != nil {
		t.Fatalf("valid from ref rejected: %v", err)
	}
	// referencing a missing var must fail
	dir2 := scaffold(t, map[string]string{
		"010.m/010.first/unit.md": first,
		"010.m/020.second/unit.md": strings.Replace(second,
			"from: first.TOKEN", "from: first.MISSING", 1),
	})
	if _, err := Load(dir2, "ubuntu", nil, nil); err == nil ||
		!strings.Contains(err.Error(), "no var") {
		t.Fatalf("want missing-var error, got %v", err)
	}
}

func TestSplitFrontmatter(t *testing.T) {
	fm, body, err := splitFrontmatter("---\na: 1\n---\nhello")
	if err != nil || fm != "a: 1" || body != "hello" {
		t.Fatalf("got %q %q %v", fm, body, err)
	}
	if _, _, err := splitFrontmatter("---\nunterminated"); err == nil {
		t.Fatal("want error for unterminated frontmatter")
	}
	fm, body, err = splitFrontmatter("no front")
	if err != nil || fm != "" || body != "no front" {
		t.Fatalf("got %q %q %v", fm, body, err)
	}
}

func TestRealLinux101Path(t *testing.T) {
	// The reference path shipped with the tool must always load and render
	// on both distro families.
	root := filepath.Join("..", "..", "paths", "sample-linux-101")
	if _, err := os.Stat(root); err != nil {
		t.Skip("sample-linux-101 path not present")
	}
	for _, d := range []struct {
		id   string
		like []string
	}{{"ubuntu", []string{"debian"}}, {"rocky", []string{"rhel", "centos", "fedora"}}} {
		p, err := Load(root, d.id, d.like, []string{"systemd", "python3"})
		if err != nil {
			t.Fatalf("%s: %v", d.id, err)
		}
		units := 0
		for _, m := range p.Modules {
			units += len(m.Units)
			for _, u := range m.Units {
				vars := map[string]string{}
				for name := range u.Front.Vars {
					vars[name] = "SAMPLE"
				}
				defTask := ""
				if len(u.Tasks) == 1 {
					defTask = u.Tasks[0].Name
				}
				if _, err := RenderUnit(Interpolate(u.Body, vars), "/x/", defTask); err != nil {
					t.Errorf("%s: render: %v", u.ID, err)
				}
			}
		}
		if units < 20 {
			t.Fatalf("%s: suspiciously few units: %d", d.id, units)
		}
	}

	// Without python3 the whole http.server chain in networking-basics must
	// be marked unsupported (serve-yourself directly, stop-serving by
	// cascade) but still present.
	p, err := Load(root, "ubuntu", []string{"debian"}, []string{"systemd"})
	if err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{
		"networking-basics/who-is-listening",
		"networking-basics/serve-yourself",
		"networking-basics/stop-serving",
	} {
		u := p.Unit(id)
		if u == nil {
			t.Fatalf("unit %s missing from python3-less load", id)
		}
		if !u.Unsupported {
			t.Errorf("unit %s: want unsupported without python3", id)
		}
	}
}

func TestRequiresCapabilityFilter(t *testing.T) {
	gated := strings.Replace(minimalUnit, "title: A unit", "title: A unit\nrequires: [systemd]", 1)
	// Depends on the gated unit via needs: - unsupported by cascade.
	dependent := strings.Replace(minimalUnit, "title: A unit", "title: A unit\nneeds: [gated]", 1)
	// References a var of the dependent - unsupported transitively, via from:.
	varRef := strings.Replace(minimalUnit, "title: A unit",
		"title: A unit\nvars:\n  V: { from: dependent.W }", 1)
	dependent = strings.Replace(dependent, "title: A unit\n", "title: A unit\nvars:\n  W: { value: x }\n", 1)
	dir := scaffold(t, map[string]string{
		"010.m/010.plain/unit.md":     minimalUnit,
		"010.m/020.gated/unit.md":     gated,
		"010.m/030.dependent/unit.md": dependent,
		"010.m/040.var-ref/unit.md":   varRef,
	})
	p, err := Load(dir, "ubuntu", nil, nil) // no capabilities
	if err != nil {
		t.Fatal(err)
	}
	if n := len(p.Modules[0].Units); n != 4 {
		t.Fatalf("want all units kept, got %d", n)
	}
	for name, want := range map[string]bool{"plain": false, "gated": true, "dependent": true, "var-ref": true} {
		u := p.Unit("m/" + name)
		if u.Unsupported != want {
			t.Errorf("unit %s: Unsupported = %v, want %v", name, u.Unsupported, want)
		}
	}
	if got := p.Unit("m/gated").MissingCaps; len(got) != 1 || got[0] != "systemd" {
		t.Errorf("gated MissingCaps = %v, want [systemd]", got)
	}
	if got := p.Unit("m/dependent").MissingCaps; len(got) != 0 {
		t.Errorf("dependent MissingCaps = %v, want none (cascaded, not direct)", got)
	}

	p, err = Load(dir, "ubuntu", nil, []string{"systemd"})
	if err != nil {
		t.Fatal(err)
	}
	for _, u := range p.Modules[0].Units {
		if u.Unsupported {
			t.Errorf("unit %s: unsupported despite systemd cap", u.ID)
		}
	}
}

func TestTopologicalTaskOrder(t *testing.T) {
	unit := `---
title: Ordered
tasks:
  z_first:
    check: |
      true
  a_second:
    needs: [z_first]
    check: |
      true
---
x
`
	dir := scaffold(t, map[string]string{"010.m/010.u/unit.md": unit})
	p, err := Load(dir, "ubuntu", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	tasks := p.Modules[0].Units[0].Tasks
	if tasks[0].Name != "z_first" || tasks[1].Name != "a_second" {
		t.Fatalf("want dependency-first order, got %s, %s", tasks[0].Name, tasks[1].Name)
	}
}

func TestSolveAttributeParsed(t *testing.T) {
	unit := `---
title: With solve
tasks:
  t1:
    check: |
      true
    solve: |
      echo hello
---
x
`
	dir := scaffold(t, map[string]string{"010.m/010.u/unit.md": unit})
	p, err := Load(dir, "ubuntu", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.TrimSpace(p.Modules[0].Units[0].Tasks[0].Solve); got != "echo hello" {
		t.Fatalf("solve = %q", got)
	}
}
