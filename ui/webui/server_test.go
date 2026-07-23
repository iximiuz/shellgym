package webui

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"

	"github.com/iximiuz/labs-content/tools/shellgym/internal/bus"
	"github.com/iximiuz/labs-content/tools/shellgym/internal/content"
	"github.com/iximiuz/labs-content/tools/shellgym/internal/engine"
	"github.com/iximiuz/labs-content/tools/shellgym/internal/state"
)

func testServer(t *testing.T) (*Server, *bus.Bus) {
	t.Helper()
	dir := t.TempDir()
	write := func(rel, data string) {
		p := filepath.Join(dir, "content", rel)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(data), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("path.yaml", "id: apitest\ntitle: API Test\n")
	write("010.mod/module.md", "# The module\nHello.")
	write("010.mod/010.unit-a/unit.md", `---
title: Unit A ${TOKEN}
vars:
  TOKEN: { value: tok123 }
tasks:
  t1:
    check: |
      true
---
Do the thing with ${TOKEN}.

::task{name="t1"}
Waiting...
::
`)
	write("010.mod/010.unit-a/pic.png", "PNGDATA")

	p, err := content.Load(filepath.Join(dir, "content"), "ubuntu", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	st, err := state.Open(filepath.Join(dir, "state"), p.ID)
	if err != nil {
		t.Fatal(err)
	}
	b := bus.New()
	// Note: no watcher and no check API needed for API-level tests; tasks
	// are never actually started because we don't activate units here.
	eng := engine.New(p, st, b, nil, engine.Options{ChecksDir: dir, SockPath: filepath.Join(dir, "x.sock")})
	t.Cleanup(eng.Shutdown)
	return New(":0", eng, b), b
}

func TestPathEndpoint(t *testing.T) {
	s, _ := testServer(t)
	ts := httptest.NewServer(s.mux)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/path")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var out struct {
		Title  string `json:"title"`
		Total  int    `json:"total"`
		Scenes []struct {
			Kind, ID, Status string
		} `json:"scenes"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}
	if out.Title != "API Test" || out.Total != 1 || len(out.Scenes) != 2 {
		t.Fatalf("path payload: %+v", out)
	}
	if out.Scenes[0].Kind != "module" || out.Scenes[1].ID != "mod/unit-a" {
		t.Fatalf("scenes: %+v", out.Scenes)
	}
}

func TestUnitEndpointInterpolatesVars(t *testing.T) {
	s, _ := testServer(t)
	ts := httptest.NewServer(s.mux)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/unit/mod/unit-a")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var out struct {
		Title string                          `json:"title"`
		HTML  string                          `json:"html"`
		Tasks []struct{ Name, Status string } `json:"tasks"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}
	if out.Title != "Unit A tok123" {
		t.Errorf("title: %q", out.Title)
	}
	if !strings.Contains(out.HTML, "Do the thing with tok123") {
		t.Errorf("vars not interpolated: %s", out.HTML)
	}
	if len(out.Tasks) != 1 || out.Tasks[0].Name != "t1" || out.Tasks[0].Status != "pending" {
		t.Errorf("tasks: %+v", out.Tasks)
	}
}

func TestModuleEndpoint(t *testing.T) {
	s, _ := testServer(t)
	ts := httptest.NewServer(s.mux)
	defer ts.Close()
	resp, err := http.Get(ts.URL + "/api/module/mod")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var out struct{ Title, HTML string }
	_ = json.NewDecoder(resp.Body).Decode(&out)
	if out.Title != "The module" || !strings.Contains(out.HTML, "Hello.") {
		t.Errorf("module payload: %+v", out)
	}
}

func TestUnitAssetServing(t *testing.T) {
	s, _ := testServer(t)
	ts := httptest.NewServer(s.mux)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/unit-assets/mod/unit-a/pic.png")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("asset status: %d", resp.StatusCode)
	}
	// Path traversal must be rejected.
	req, _ := http.NewRequest("GET", ts.URL+"/unit-assets/mod/unit-a/"+strings.ReplaceAll("../../../etc/passwd", "/", "%2f"), nil)
	resp2, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp2.Body.Close()
	if resp2.StatusCode == 200 {
		t.Error("path traversal served")
	}
}

func TestWebSocketStreamsBusEvents(t *testing.T) {
	s, b := testServer(t)
	ts := httptest.NewServer(s.mux)
	defer ts.Close()

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http") + "/api/events"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	time.Sleep(100 * time.Millisecond) // let subscription settle
	b.Publish(bus.Event{Type: "task", Data: engine.TaskEvent{Unit: "u", Task: "t", Status: "completed"}})

	_ = conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	var ev struct {
		Type string                              `json:"type"`
		Data struct{ Unit, Task, Status string } `json:"data"`
	}
	if err := conn.ReadJSON(&ev); err != nil {
		t.Fatal(err)
	}
	if ev.Type != "task" || ev.Data.Status != "completed" {
		t.Errorf("event: %+v", ev)
	}
}

func TestStatusEndpoint(t *testing.T) {
	s, _ := testServer(t)
	ts := httptest.NewServer(s.mux)
	defer ts.Close()
	resp, err := http.Get(ts.URL + "/api/status")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var out map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&out)
	if out["distro"] == "" {
		t.Errorf("status: %+v", out)
	}
}

func TestDebugDisabledInLiveMode(t *testing.T) {
	s, _ := testServer(t)
	s.Opts.Live = true
	ts := httptest.NewServer(s.mux)
	defer ts.Close()
	resp, err := http.Get(ts.URL + "/api/debug/mod/unit-a")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 404 {
		t.Fatalf("debug endpoint alive in live mode: %d", resp.StatusCode)
	}
	// status advertises live so the UI hides the debug button
	resp2, err := http.Get(ts.URL + "/api/status")
	if err != nil {
		t.Fatal(err)
	}
	defer resp2.Body.Close()
	var st map[string]any
	_ = json.NewDecoder(resp2.Body).Decode(&st)
	if st["live"] != true {
		t.Errorf("status.live = %v", st["live"])
	}
}
