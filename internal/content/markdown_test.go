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
