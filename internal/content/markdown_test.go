package content

import (
	"strings"
	"testing"
)

func TestRenderUnitPlain(t *testing.T) {
	html, err := RenderUnit("# Title\n\nSome *text*.", "")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, "<h1") || !strings.Contains(html, "<em>text</em>") {
		t.Errorf("html: %s", html)
	}
}

func TestRenderTaskComponent(t *testing.T) {
	src := `Intro.

::task{name="check_it"}
#active
Waiting for **it**...
#completed
Done with ` + "`it`" + `.
::

Outro.`
	html, err := RenderUnit(src, "")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		`data-task="check_it"`,
		`data-status="pending"`,
		`task-section-active`,
		`task-section-completed`,
		`<strong>it</strong>`,
		`task-section-hint`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %q in:\n%s", want, html)
		}
	}
	if !strings.Contains(html, "Intro.") || !strings.Contains(html, "Outro.") {
		t.Error("surrounding markdown lost")
	}
}

func TestRenderTaskDefaultBodyIsActive(t *testing.T) {
	html, err := RenderUnit("::task{name=\"t\"}\nJust waiting.\n::", "")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, "task-section-active") || !strings.Contains(html, "Just waiting.") {
		t.Errorf("default body not treated as active: %s", html)
	}
}

func TestRenderHintComponent(t *testing.T) {
	html, err := RenderUnit("::hint{title=\"Stuck?\"}\nTry `ls`.\n::", "")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, "<details") || !strings.Contains(html, "Stuck?") ||
		!strings.Contains(html, "<code>ls</code>") {
		t.Errorf("hint html: %s", html)
	}
}

func TestRenderCalloutComponents(t *testing.T) {
	html, err := RenderUnit("::tip{title=\"Tab completion\"}\nPress `Tab` to complete.\n::", "")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		`class="callout callout-tip"`,
		`callout-title`,
		`Tab completion`,
		`<code>Tab</code>`,
		`callout-icon`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %q in tip html:\n%s", want, html)
		}
	}
	// details/summary is the hint's shape; callouts are always visible
	if strings.Contains(html, "<details") {
		t.Errorf("tip must not fold: %s", html)
	}
	// each kind gets its own modifier class and default title
	for kind, title := range map[string]string{"tip": "Tip", "note": "Note", "warn": "Warning"} {
		html, err = RenderUnit("::"+kind+"\nSomething useful.\n::", "")
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(html, `class="callout callout-`+kind+`"`) {
			t.Errorf("%s: modifier class missing: %s", kind, html)
		}
		if !strings.Contains(html, ">"+title+"</span>") {
			t.Errorf("%s: default title missing: %s", kind, html)
		}
	}
}

func TestRenderImageComponent(t *testing.T) {
	html, err := RenderUnit("::image{src=\"pic.png\" alt=\"A pic\"}\n::", "/unit-assets/m/u/")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, `src="/unit-assets/m/u/pic.png"`) {
		t.Errorf("image src not prefixed: %s", html)
	}
}

func TestRelativeImgRewrite(t *testing.T) {
	html, err := RenderUnit("![alt](diagram.png) and ![ext](https://x/y.png)", "/unit-assets/m/u/")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, `src="/unit-assets/m/u/diagram.png"`) {
		t.Errorf("relative img not rewritten: %s", html)
	}
	if !strings.Contains(html, `src="https://x/y.png"`) {
		t.Errorf("absolute img mangled: %s", html)
	}
}

func TestLinksOpenInNewTab(t *testing.T) {
	src := "See [docs](https://example.com/x), https://auto.example.com and\n\n::hint{title=\"H\"}\n[inner](https://in.example.com)\n::\n"
	html, err := RenderUnit(src, "")
	if err != nil {
		t.Fatal(err)
	}
	for _, href := range []string{"https://example.com/x", "https://auto.example.com", "https://in.example.com"} {
		if !strings.Contains(html, `href="`+href+`" target="_blank" rel="noopener"`) {
			t.Errorf("link %s not opened in a new tab: %s", href, html)
		}
	}
}

func TestUnknownComponent(t *testing.T) {
	if _, err := RenderUnit("::nope{}\n::", ""); err == nil {
		t.Error("want error for unknown component")
	}
}

func TestUnclosedComponent(t *testing.T) {
	if _, err := RenderUnit("::task{name=\"x\"}\nno close", ""); err == nil {
		t.Error("want error for unclosed component")
	}
}

func TestRenderTaskComponentDefaultName(t *testing.T) {
	html, err := RenderUnit("::task\n#active\nWaiting...\n::", "", "only_task")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, `data-task="only_task"`) {
		t.Errorf("default task name not applied: %s", html)
	}
	// Without a default, a name-less task must error.
	if _, err := RenderUnit("::task\nWaiting...\n::", ""); err == nil {
		t.Error("want error for name-less task without default")
	}
}

func TestRenderMDCBlockAttrs(t *testing.T) {
	src := "::hint\n---\ntitle: Block Title\n---\nThe hint *body*.\n::"
	html, err := RenderUnit(src, "")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, "Block Title") || !strings.Contains(html, "<em>body</em>") {
		t.Errorf("MDC attrs/body wrong: %s", html)
	}
	if strings.Contains(html, "---") {
		t.Errorf("yaml fence leaked into output: %s", html)
	}
	// colon-prefixed keys (labs convention) are accepted too
	html, err = RenderUnit("::image\n---\n:src: pic.png\n:alt: A pic\n---\n::", "/a/")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, `src="/a/pic.png"`) {
		t.Errorf("colon-key attrs not parsed: %s", html)
	}
}

func TestAuthorHTMLCannotRunScript(t *testing.T) {
	src := "Intro.\n\n" +
		"<script>alert('s1')</script>\n\n" +
		`<button class="gym-reset" onclick="fetch('/api/reset/m/u',{method:'POST'})">Reset</button>` + "\n\n" +
		`<img src="x.png" onerror="alert('s2')">` + "\n\n" +
		"[plain](javascript:alert('s3')) and <a href=\"javascript:alert('s4')\">raw</a>\n\n" +
		`<iframe src="https://example.com/"></iframe><object data="x.swf"></object>` + "\n\n" +
		`<svg onload="alert('s5')"><script>alert('s6')</script></svg>` + "\n\n" +
		`<form action="https://example.com/"><input name="q"></form>` + "\n\n" +
		`<a href="https://example.com/" target="_blank" rel="opener">tab</a>` + "\n\n" +
		"<style>.x { color: red; }</style><script>alert('s7')</script>\n\n" +
		"::hint{title=\"Stuck?\"}\n<div onmouseover=\"alert('s8')\">hover</div>\n::\n"
	html, err := RenderUnit(src, "/unit-assets/m/u/")
	if err != nil {
		t.Fatal(err)
	}
	for _, bad := range []string{
		"alert(", "<script", "onclick", "onerror", "onload", "onmouseover",
		"javascript:", "<iframe", "<object", "<form", `rel="opener"`,
	} {
		if strings.Contains(html, bad) {
			t.Errorf("sanitized output contains %q:\n%s", bad, html)
		}
	}
	// The hint box's own bulb icon is the only SVG left.
	if n := strings.Count(html, "<svg"); n != strings.Count(html, `<svg class="hint-bulb"`) {
		t.Errorf("author SVG survived:\n%s", html)
	}
	for _, want := range []string{"Intro.", "Reset", "hover", "tab", "hint-box"} {
		if !strings.Contains(html, want) {
			t.Errorf("sanitized output lost %q:\n%s", want, html)
		}
	}
}

// Published paths use raw HTML for formatting: <pre><code> blocks and
// <span> highlights, plus <style> blocks that show the solution matching
// the drawn variant (data-* attributes). All of it must survive.
func TestAuthorHTMLKeepsFormatting(t *testing.T) {
	src := "<style>\n.sol { display: none; }\n.scene-body:has(.sol) > .sol[data-draw=\"labels\"] { display: block; }\n</style>\n\n" +
		"<pre class=\"sol\" data-draw=\"labels\"><code>kubectl apply -f - &lt;&lt;'EOF'\n" +
		"  name: <span class=\"sol-v\">web</span>\nEOF\n</code></pre>\n\n" +
		"```bash\nls -l\n```\n\n" +
		"| Left | Right |\n|:-----|------:|\n| a | b |\n\n" +
		"- [x] done\n- [ ] todo\n\n" +
		"A <kbd>Ctrl</kbd>-<kbd>C</kbd> and a [link](https://example.com/).\n\n" +
		"![pic](pic.png)\n"
	html, err := RenderUnit(src, "/unit-assets/m/u/")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"<style>\n.sol { display: none; }\n.scene-body:has(.sol) > .sol[data-draw=\"labels\"] { display: block; }\n</style>",
		`<pre class="sol" data-draw="labels"><code>kubectl apply -f - &lt;&lt;`,
		`<span class="sol-v">web</span>`,
		`<code class="language-bash">`,
		`style="text-align: left"`,
		`style="text-align: right"`,
		`type="checkbox"`,
		`<kbd>Ctrl</kbd>`,
		`href="https://example.com/" target="_blank" rel="noopener"`,
		`src="/unit-assets/m/u/pic.png"`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("sanitized output lost %q:\n%s", want, html)
		}
	}
}
