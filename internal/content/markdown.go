package content

import (
	"bytes"
	"fmt"
	"html"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	ghtml "github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
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

// Raw HTML in the markdown passes goldmark (WithUnsafe) and is then filtered
// by authorHTML - see RenderMarkdown.
var md = goldmark.New(
	goldmark.WithExtensions(extension.GFM),
	goldmark.WithParserOptions(parser.WithASTTransformers(
		util.Prioritized(newTabLinks{}, 100),
	)),
	goldmark.WithRendererOptions(ghtml.WithUnsafe()),
)

// authorHTML is the allowlist for HTML rendered from an author's markdown.
// The UI inserts it into the student's gym page, which talks to the daemon's
// unauthenticated API, so nothing that can run script survives: <script>,
// on* handlers, javascript: URLs, <iframe>/<object>/<form>, inline SVG.
// Formatting HTML stays - paths use <pre><code>, <span class>, data-*
// attributes, and <style> blocks (CSS cannot run script) for variant-specific
// text.
var authorHTML = func() *bluemonday.Policy {
	p := bluemonday.UGCPolicy()
	p.AllowElements("kbd")
	p.AllowAttrs("class").Globally()
	p.AllowDataAttributes()
	p.AllowDataURIImages()
	// newTabLinks output; rel=opener would hand window.opener back.
	p.RequireNoFollowOnLinks(false)
	p.AllowAttrs("target").Matching(regexp.MustCompile(`^_blank$`)).OnElements("a")
	p.AllowAttrs("rel").Matching(regexp.MustCompile(`^noopener$`)).OnElements("a")
	// GFM task list items and table cell alignment.
	p.AllowAttrs("type").Matching(regexp.MustCompile(`^checkbox$`)).OnElements("input")
	p.AllowAttrs("checked", "disabled").OnElements("input")
	p.AllowStyles("text-align").OnElements("th", "td")
	// AllowUnsafe only lets <style> be declared below; <script> stays out of
	// the allowlist, so it is dropped together with its content.
	p.AllowUnsafe(true)
	p.AllowElements("style")
	p.AllowElementsContent("style")
	return p
}()

// newTabLinks makes every link (and GFM autolink) open in a new tab: the unit
// page is the student's workspace (terminal + live task state), and navigating
// away from it mid-rep would be disruptive. rel=noopener keeps the opened page
// from reaching back into ours via window.opener.
type newTabLinks struct{}

func (newTabLinks) Transform(doc *ast.Document, _ text.Reader, _ parser.Context) {
	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		switch n.Kind() {
		case ast.KindLink, ast.KindAutoLink:
			n.SetAttributeString("target", []byte("_blank"))
			n.SetAttributeString("rel", []byte("noopener"))
		}
		return ast.WalkContinue, nil
	})
}

// RenderMarkdown renders plain markdown (no components) to HTML, sanitized
// with authorHTML. Every piece of author markdown (plain text, component
// sections, module intros) goes through here; the component wrappers around
// it are built by this package with escaped attributes.
func RenderMarkdown(src string) (string, error) {
	var buf bytes.Buffer
	if err := md.Convert([]byte(src), &buf); err != nil {
		return "", err
	}
	return authorHTML.Sanitize(buf.String()), nil
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
	case "tip", "note", "warn":
		return renderCallout(c)
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

// callouts are the always-visible boxes (::tip, ::note, ::warn) that sit
// alongside tasks. They share one markup shape and stylesheet family and
// differ only in default title, icon and color (via the callout-<kind>
// modifier class).
var callouts = map[string]struct {
	title string
	icon  string
}{
	"tip":  {"Tip", tipBulbSVG},
	"note": {"Note", noteInfoSVG},
	"warn": {"Warning", warnTriangleSVG},
}

func renderCallout(c *Component) (string, error) {
	kind := callouts[c.Name]
	title := c.Attrs["title"]
	if title == "" {
		title = kind.title
	}
	inner, err := RenderMarkdown(c.Sections[""])
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(
		`<div class="callout callout-%s"><div class="callout-head">%s<span class="callout-title">%s</span></div><div class="callout-body">%s</div></div>`,
		c.Name, kind.icon, html.EscapeString(title), inner), nil
}

// tipBulbSVG is mdi:lightbulb-outline for ::tip callouts - technique notes
// (keystrokes, tab completion, history tricks). Deliberately the
// plain-outline cousin of the hint box's filled "lightbulb-on" bulb: same
// family, calmer voice.
const tipBulbSVG = `<svg class="callout-icon" viewBox="0 0 24 24" aria-hidden="true"><path fill="currentColor" d="M12,2A7,7 0 0,1 19,9C19,11.38 17.81,13.47 16,14.74V17A1,1 0 0,1 15,18H9A1,1 0 0,1 8,17V14.74C6.19,13.47 5,11.38 5,9A7,7 0 0,1 12,2M9,21V20H15V21A1,1 0 0,1 14,22H10A1,1 0 0,1 9,21M12,4A5,5 0 0,0 7,9C7,11.05 8.23,12.81 10,13.58V16H14V13.58C15.77,12.81 17,11.05 17,9A5,5 0 0,0 12,4Z"/></svg>`

// noteInfoSVG is mdi:information-outline for ::note callouts.
const noteInfoSVG = `<svg class="callout-icon" viewBox="0 0 24 24" aria-hidden="true"><path fill="currentColor" d="M11,9H13V7H11M12,20C7.59,20 4,16.41 4,12C4,7.59 7.59,4 12,4C16.41,4 20,7.59 20,12C20,16.41 16.41,20 12,20M12,2A10,10 0 0,0 2,12A10,10 0 0,0 12,22A10,10 0 0,0 22,12A10,10 0 0,0 12,2M11,17H13V11H11V17Z"/></svg>`

// warnTriangleSVG is mdi:alert-outline for ::warn callouts.
const warnTriangleSVG = `<svg class="callout-icon" viewBox="0 0 24 24" aria-hidden="true"><path fill="currentColor" d="M12,2L1,21H23M12,6L19.53,19H4.47M11,10V14H13V10M11,16V18H13V16"/></svg>`

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
