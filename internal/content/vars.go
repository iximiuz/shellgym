package content

import (
	"fmt"
	"math/rand"
	"os/exec"
	"regexp"
	"sort"
	"strings"
	"time"
)

// VarLookup resolves a `from:` reference: returns the value of varName in
// unitName's (already resolved) vars.
type VarLookup func(unitName, varName string) (string, error)

// ResolveVars materializes a unit's var specs into concrete values.
// Called once per unit activation; the result is persisted so the values
// stay stable across daemon restarts. lookup may be nil if no spec uses
// `from:`.
func ResolveVars(specs map[string]VarSpec, lookup VarLookup) (map[string]string, error) {
	out := make(map[string]string, len(specs))
	names := make([]string, 0, len(specs))
	for n := range specs {
		if !varNameRe.MatchString(n) {
			return nil, fmt.Errorf("var %q: name must match %s", n, varNameRe)
		}
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		spec := specs[n]
		switch {
		case spec.Value != "":
			out[n] = spec.Value
		case len(spec.Pick) > 0:
			out[n] = spec.Pick[rand.Intn(len(spec.Pick))]
		case spec.From != "":
			unitName, varName, ok := strings.Cut(spec.From, ".")
			if !ok {
				return nil, fmt.Errorf("var %q: from must be \"unit-name.VAR\"", n)
			}
			if lookup == nil {
				return nil, fmt.Errorf("var %q: from-references not supported here", n)
			}
			v, err := lookup(unitName, varName)
			if err != nil {
				return nil, fmt.Errorf("var %q: %w", n, err)
			}
			out[n] = v
		case spec.Shell != "":
			cmd := exec.Command("bash", "-o", "pipefail", "-c", spec.Shell)
			done := make(chan struct{})
			var b []byte
			var err error
			go func() { b, err = cmd.Output(); close(done) }()
			select {
			case <-done:
			case <-time.After(10 * time.Second):
				_ = cmd.Process.Kill()
				<-done
				return nil, fmt.Errorf("var %q: shell command timed out", n)
			}
			if err != nil {
				return nil, fmt.Errorf("var %q: %w", n, err)
			}
			out[n] = strings.TrimSpace(string(b))
		default:
			return nil, fmt.Errorf("var %q: one of value/pick/shell/from required", n)
		}
	}
	return out, nil
}

var (
	varNameRe = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)
	interpRe  = regexp.MustCompile(`\$\{([A-Z][A-Z0-9_]*)\}`)
)

// Interpolate replaces ${NAME} references with var values. Unknown names are
// left untouched (they may be shell syntax or task-field placeholders that
// the UI resolves live).
func Interpolate(text string, vars map[string]string) string {
	return interpRe.ReplaceAllStringFunc(text, func(m string) string {
		name := interpRe.FindStringSubmatch(m)[1]
		if v, ok := vars[name]; ok {
			return v
		}
		return m
	})
}
