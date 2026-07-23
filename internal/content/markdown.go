package content

import (
	"bytes"
	"fmt"
	"html"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	ghtml "github.com/yuin/goldmark/renderer/html"
)

// The markdown body supports MDC-style block components:
//
//   ::task{name="chdir"}
//   #active
//   Waiting...
//   #completed
//   Done!
//   ::
//
//   ::hint{title="Stuck?"}
//   Folded hint markdown.
//   ::
//
//   ::image{src="tree.png" alt="File tree"}
//   ::
//
// Components are pre-processed into HTML wrappers before goldmark runs, with
// their inner section markdown rendered recursively.

var md = goldmark.New(
	goldmark.WithExtensions(extension.GFM),
	goldmark.WithRendererOptions(ghtml.WithUnsafe()),
)

// RenderMarkdown renders plain markdown (no components) to HTML.
func RenderMarkdown(src string) (string, error) {
	var buf bytes.Buffer
	if err := md.Convert([]byte(src), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

var (
	compOpenRe = regexp.MustCompile(`^::([a-z][a-z-]*)(\{[^}]*\})?\s*$`)
	attrRe     = regexp.MustCompile(`([a-zA-Z][a-zA-Z0-9-]*)="([^"]*)"`)
	sectionRe  = regexp.MustCompile(`^#([a-z][a-z-]*)\s*$`)
)

// Component is a parsed block component.
type Component struct {
	Name     string
	Attrs    map[string]string
	Sections map[string]string // section name -> raw markdown ("" section = default body)
}

// RenderUnit renders a unit body (with vars already interpolated) to HTML.
// assetPrefix is prepended to relative image/src paths (unit-local assets).
// defaultTask, when non-empty, is the task name a name-less ::task
// component binds to (single-task units may omit the name).
func RenderUnit(body, assetPrefix string, defaultTask ...string) (string, error) {
	defTask := ""
	if len(defaultTask) > 0 {
		defTask = defaultTask[0]
	}
	lines := strings.Split(body, "\n")
	var out strings.Builder
	var plain []string

	flushPlain := func() error {
		if len(plain) == 0 {
			return nil
		}
		h, err := RenderMarkdown(strings.Join(plain, "\n"))
		if err != nil {
			return err
		}
		out.WriteString(h)
		plain = plain[:0]
		return nil
	}

	for i := 0; i < len(lines); i++ {
		m := compOpenRe.FindStringSubmatch(lines[i])
		if m == nil {
			plain = append(plain, lines[i])
			continue
		}
		// find the closing "::"
		end := -1
		for j := i + 1; j < len(lines); j++ {
			if strings.TrimSpace(lines[j]) == "::" {
				end = j
				break
			}
		}
		if end < 0 {
			return "", fmt.Errorf("component %q at line %d: missing closing '::'", m[1], i+1)
		}
		if err := flushPlain(); err != nil {
			return "", err
		}
		comp := parseComponent(m[1], m[2], lines[i+1:end])
		h, err := renderComponent(comp, assetPrefix, defTask)
		if err != nil {
			return "", err
		}
		out.WriteString(h)
		i = end
	}
	if err := flushPlain(); err != nil {
		return "", err
	}
	return rewriteAssetURLs(out.String(), assetPrefix), nil
}

func parseComponent(name, rawAttrs string, inner []string) *Component {
	c := &Component{Name: name, Attrs: map[string]string{}, Sections: map[string]string{}}
	for _, am := range attrRe.FindAllStringSubmatch(rawAttrs, -1) {
		c.Attrs[am[1]] = am[2]
	}
	// MDC block form: a leading `---\n<yaml>\n---` inside the component
	// carries the attributes instead of the inline {...} list.
	consumed := 0
	if len(inner) > 0 && strings.TrimSpace(inner[0]) == "---" {
		for j := 1; j < len(inner); j++ {
			if strings.TrimSpace(inner[j]) == "---" {
				var attrs map[string]string
				raw := strings.Join(inner[1:j], "\n")
				if err := yaml.Unmarshal([]byte(raw), &attrs); err == nil {
					for k, v := range attrs {
						c.Attrs[strings.TrimPrefix(k, ":")] = v
					}
					consumed = j + 1
				}
				break
			}
		}
	}
	inner = inner[consumed:]
	section := ""
	var buf []string
	flush := func() {
		if len(buf) > 0 || section != "" {
			c.Sections[section] = strings.TrimSpace(strings.Join(buf, "\n"))
		}
		buf = nil
	}
	for _, line := range inner {
		if sm := sectionRe.FindStringSubmatch(line); sm != nil {
			flush()
			section = sm[1]
			continue
		}
		buf = append(buf, line)
	}
	flush()
	return c
}

func renderComponent(c *Component, assetPrefix, defaultTask string) (string, error) {
	switch c.Name {
	case "task":
		return renderTaskComponent(c, defaultTask)
	case "hint":
		title := c.Attrs["title"]
		if title == "" {
			title = "Hint"
		}
		inner, err := RenderMarkdown(c.Sections[""])
		if err != nil {
			return "", err
		}
		return fmt.Sprintf(
			`<details class="hint-box"><summary><span class="hint-title">%s</span>%s</summary><div class="hint-body">%s</div></details>`,
			html.EscapeString(title), hintBulbSVG, inner), nil
	case "image":
		src := c.Attrs["src"]
		if !strings.Contains(src, "://") && !strings.HasPrefix(src, "/") {
			src = assetPrefix + src
		}
		return fmt.Sprintf(`<figure class="image-box"><img src="%s" alt="%s"></figure>`,
			html.EscapeString(src), html.EscapeString(c.Attrs["alt"])), nil
	default:
		return "", fmt.Errorf("unknown component %q", c.Name)
	}
}

// hintBulbSVG is mdi:lightbulb-on-40 - the lightbulb the labs platform
// shows on hint boxes (rendered on the right edge of the summary row).
const hintBulbSVG = `<svg class="hint-bulb" viewBox="0 0 24 24" aria-hidden="true"><path fill="currentColor" d="M1 11h3v2H1zM13 1h-2v3h2zM4.9 3.5L3.5 4.9L5.6 7L7 5.6zm14.2 0L17 5.6L18.4 7l2.1-2.1zM10 22c0 .6.4 1 1 1h2c.6 0 1-.4 1-1v-1h-4zm10-11v2h3v-2zm-2 1c0 2.2-1.2 4.2-3 5.2V19c0 .6-.4 1-1 1h-4c-.6 0-1-.4-1-1v-1.8c-1.8-1-3-3-3-5.2c0-3.3 2.7-6 6-6s6 2.7 6 6m-2 0c0-2.21-1.79-4-4-4s-4 1.79-4 4c0 .74.22 1.41.57 2h6.86c.35-.59.57-1.26.57-2"/></svg>`

func renderTaskComponent(c *Component, defaultTask string) (string, error) {
	name := c.Attrs["name"]
	if name == "" {
		name = defaultTask
	}
	if name == "" {
		return "", fmt.Errorf("task component: missing name attribute (only single-task units may omit it)")
	}
	sections := map[string]string{}
	for _, s := range []string{"active", "completed", "failed"} {
		src := c.Sections[s]
		if src == "" && s == "active" {
			src = c.Sections[""] // default body doubles as the active text
		}
		if src == "" {
			continue
		}
		h, err := RenderMarkdown(src)
		if err != nil {
			return "", err
		}
		sections[s] = h
	}
	var b strings.Builder
	fmt.Fprintf(&b, `<div class="task-box" data-task="%s" data-status="pending">`, html.EscapeString(name))
	b.WriteString(`<div class="task-status-icon"></div><div class="task-text">`)
	for _, s := range []string{"active", "completed", "failed"} {
		if sections[s] != "" {
			fmt.Fprintf(&b, `<div class="task-section task-section-%s">%s</div>`, s, sections[s])
		}
	}
	b.WriteString(`<div class="task-section task-section-hint" hidden></div>`)
	b.WriteString(`</div></div>`)
	return b.String(), nil
}

var imgSrcRe = regexp.MustCompile(`(<img[^>]+src=")([^"]+)(")`)

// rewriteAssetURLs prefixes relative <img> srcs produced by plain markdown
// (![alt](img.png)) with the unit asset route.
func rewriteAssetURLs(htmlSrc, assetPrefix string) string {
	if assetPrefix == "" {
		return htmlSrc
	}
	return imgSrcRe.ReplaceAllStringFunc(htmlSrc, func(m string) string {
		parts := imgSrcRe.FindStringSubmatch(m)
		src := parts[2]
		if strings.Contains(src, "://") || strings.HasPrefix(src, "/") || strings.HasPrefix(src, "data:") {
			return m
		}
		return parts[1] + assetPrefix + src + parts[3]
	})
}
