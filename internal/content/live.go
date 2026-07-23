package content

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// StripSolveScripts removes `solve:` blocks from every unit.md under dir,
// rewriting the files in place. Used by `serve --live` so students cannot
// read the reference solutions off the disk. The transformation is textual
// (YAML formatting is preserved); it removes both block-scalar
// (`solve: |`) and single-line (`solve: ...`) forms.
func StripSolveScripts(dir string) (int, error) {
	stripped := 0
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || filepath.Base(path) != "unit.md" {
			return err
		}
		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		front, body, err := splitFrontmatter(string(raw))
		if err != nil || front == "" {
			return err
		}
		newFront, n := stripSolveFromYAML(front)
		if n == 0 {
			return nil
		}
		out := "---\n" + newFront + "\n---\n" + body
		if err := os.WriteFile(path, []byte(out), info.Mode()); err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		stripped += n
		return nil
	})
	return stripped, err
}

var solveKeyRe = regexp.MustCompile(`^(\s*)solve:\s*(\|[-+]?|>[-+]?)?\s*(\S.*)?$`)

// stripSolveFromYAML removes solve: entries (and their indented block
// bodies) from frontmatter text, returning the new text and the number of
// entries removed.
func stripSolveFromYAML(front string) (string, int) {
	lines := strings.Split(front, "\n")
	var out []string
	removed := 0
	for i := 0; i < len(lines); i++ {
		m := solveKeyRe.FindStringSubmatch(lines[i])
		if m == nil {
			out = append(out, lines[i])
			continue
		}
		removed++
		indent := len(m[1])
		if m[2] != "" || m[3] == "" {
			// Block scalar (or empty value): swallow more-indented lines.
			for i+1 < len(lines) {
				next := lines[i+1]
				if strings.TrimSpace(next) == "" || lineIndent(next) > indent {
					i++
					continue
				}
				break
			}
		}
	}
	return strings.Join(out, "\n"), removed
}

func lineIndent(s string) int {
	return len(s) - len(strings.TrimLeft(s, " "))
}
