// Package webui is the built-in web UI: a small HTTP+WebSocket server with
// embedded static assets. It implements ui.UI.
package webui

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/websocket"

	"github.com/iximiuz/labs-content/tools/shellgym/internal/bus"
	"github.com/iximiuz/labs-content/tools/shellgym/internal/content"
	"github.com/iximiuz/labs-content/tools/shellgym/internal/engine"
	"github.com/iximiuz/labs-content/tools/shellgym/internal/state"
)

//go:embed assets
var assets embed.FS

// Options tunes the server.
type Options struct {
	// Live disables author-facing debug endpoints (student deployments).
	Live bool
}

type Server struct {
	Addr   string
	Engine *engine.Engine
	Bus    *bus.Bus
	Opts   Options

	mux *http.ServeMux
}

func New(addr string, eng *engine.Engine, b *bus.Bus, opts ...Options) *Server {
	var o Options
	if len(opts) > 0 {
		o = opts[0]
	}
	s := &Server{Addr: addr, Engine: eng, Bus: b, Opts: o}
	m := http.NewServeMux()
	sub, _ := fs.Sub(assets, "assets")
	m.Handle("/", http.FileServer(http.FS(sub)))
	m.HandleFunc("GET /api/path", s.handlePath)
	m.HandleFunc("GET /api/unit/{id...}", s.handleUnit)
	m.HandleFunc("POST /api/activate/{id...}", s.handleActivate)
	m.HandleFunc("POST /api/reset/{id...}", s.handleReset)
	m.HandleFunc("POST /api/module-seen/{id...}", s.handleModuleSeen)
	m.HandleFunc("GET /api/module/{id...}", s.handleModule)
	m.HandleFunc("GET /api/debug/{id...}", s.handleDebug)
	m.HandleFunc("GET /api/status", s.handleStatus)
	m.HandleFunc("GET /api/events", s.handleEvents)
	m.HandleFunc("GET /unit-assets/{rest...}", s.handleUnitAsset)
	s.mux = m
	return s
}

func (s *Server) Run(ctx context.Context) error {
	srv := &http.Server{Addr: s.Addr, Handler: s.mux}
	go func() {
		<-ctx.Done()
		shCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = srv.Shutdown(shCtx)
	}()
	log.Printf("webui: listening on %s", s.Addr)
	err := srv.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

// --- path/scene handlers ----------------------------------------------------

type sceneJSON struct {
	Kind     string `json:"kind"` // module | unit
	ID       string `json:"id"`   // module id or unit id
	Title    string `json:"title"`
	ModuleID string `json:"moduleId"`
	Status   string `json:"status"` // pending|active|completed|seen
	// Locked marks units whose needs: dependencies are not all completed:
	// browsable, but not activatable (no init, no checks) until they are.
	Locked bool `json:"locked"`
}

type pathJSON struct {
	ID          string      `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Scenes      []sceneJSON `json:"scenes"`
	Current     string      `json:"current"` // current unit id ("" = none yet)
	Completed   int         `json:"completed"`
	Total       int         `json:"total"` // units only
}

func (s *Server) handlePath(w http.ResponseWriter, r *http.Request) {
	p := s.Engine.Path
	out := pathJSON{ID: p.ID, Title: p.Title, Description: p.Description}
	s.Engine.Store.View(func(d *state.Data) {
		out.Current = d.CurrentUnit
		for _, sc := range p.Scenes() {
			switch sc.Kind {
			case "module":
				st := "pending"
				if d.SeenModules[sc.Module.ID] {
					st = "seen"
				}
				out.Scenes = append(out.Scenes, sceneJSON{
					Kind: "module", ID: sc.Module.ID, ModuleID: sc.Module.ID,
					Title: sc.Module.Title, Status: st,
				})
			case "unit":
				out.Total++
				st := "pending"
				if us, ok := d.Units[sc.Unit.ID]; ok {
					st = string(us.Status)
				}
				if st == string(state.UnitCompleted) {
					out.Completed++
				}
				out.Scenes = append(out.Scenes, sceneJSON{
					Kind: "unit", ID: sc.Unit.ID, ModuleID: sc.Unit.ModuleID,
					Title: sc.Unit.Front.Title, Status: st,
					Locked: engine.UnitLockedIn(p, d, sc.Unit.ID),
				})
			}
		}
	})
	writeJSON(w, out)
}

type taskJSON struct {
	Name   string `json:"name"`
	Mode   string `json:"mode"`
	Status string `json:"status"`
	Hint   string `json:"hint,omitempty"`
	// Needs lists task names that must complete before this one goes live;
	// the UI shows dependent tasks as blocked until then.
	Needs []string `json:"needs,omitempty"`
}

type unitJSON struct {
	ID     string            `json:"id"`
	Title  string            `json:"title"`
	Module string            `json:"moduleId"`
	HTML   string            `json:"html"`
	Tasks  []taskJSON        `json:"tasks"`
	Status string            `json:"status"`
	Locked bool              `json:"locked"`
	Vars   map[string]string `json:"vars,omitempty"` // consumed by `shellgym solve`
}

func (s *Server) handleUnit(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	u := s.Engine.Path.Unit(id)
	if u == nil {
		http.Error(w, "unknown unit", 404)
		return
	}
	vars, err := s.Engine.EnsureVars(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	body := content.Interpolate(u.Body, vars)
	title := content.Interpolate(u.Front.Title, vars)
	defTask := ""
	if len(u.Tasks) == 1 {
		defTask = u.Tasks[0].Name
	}
	html, err := content.RenderUnit(body, "/unit-assets/"+u.ID+"/", defTask)
	if err != nil {
		http.Error(w, fmt.Sprintf("render: %v", err), 500)
		return
	}
	out := unitJSON{ID: u.ID, Title: title, Module: u.ModuleID, HTML: html, Status: "pending",
		Locked: s.Engine.UnitLocked(id), Vars: vars}
	s.Engine.Store.View(func(d *state.Data) {
		us := d.Unit(id)
		out.Status = string(us.Status)
		for _, t := range u.Tasks {
			tj := taskJSON{Name: t.Name, Mode: string(t.Mode), Status: "pending", Needs: t.Needs}
			if ts, ok := us.Tasks[t.Name]; ok {
				tj.Status = ts.Status
				tj.Hint = ts.Hint
			}
			out.Tasks = append(out.Tasks, tj)
		}
	})
	writeJSON(w, out)
}

type moduleJSON struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	HTML  string `json:"html"`
}

func (s *Server) handleModule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	m := s.Engine.Path.Module(id)
	if m == nil || m.Intro == "" {
		http.Error(w, "unknown module", 404)
		return
	}
	html, err := content.RenderUnit(m.Intro, "")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	writeJSON(w, moduleJSON{ID: m.ID, Title: m.Title, HTML: html})
}

func (s *Server) handleActivate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.Engine.ActivateUnit(id); err != nil {
		code := 400
		if errors.Is(err, engine.ErrUnitLocked) {
			code = 409
		}
		http.Error(w, err.Error(), code)
		return
	}
	writeJSON(w, map[string]string{"ok": "true"})
}

func (s *Server) handleReset(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.Engine.ResetUnit(id); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	writeJSON(w, map[string]string{"ok": "true"})
}

func (s *Server) handleModuleSeen(w http.ResponseWriter, r *http.Request) {
	if err := s.Engine.MarkModuleSeen(r.PathValue("id")); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	writeJSON(w, map[string]string{"ok": "true"})
}

type debugTaskJSON struct {
	Name string          `json:"name"`
	Runs []state.TaskRun `json:"runs"`
}

func (s *Server) handleDebug(w http.ResponseWriter, r *http.Request) {
	if s.Opts.Live {
		http.Error(w, "debug is disabled in live mode", 404)
		return
	}
	id := r.PathValue("id")
	names := s.Engine.Store.RunTasks(id)
	sortTaskNames(names)
	out := make([]debugTaskJSON, 0, len(names))
	for _, n := range names {
		out = append(out, debugTaskJSON{Name: n, Runs: s.Engine.Store.Runs(id, n)})
	}
	writeJSON(w, out)
}

func sortTaskNames(names []string) {
	isInit := func(n string) bool { return strings.HasPrefix(n, "init:") }
	for i := 0; i < len(names); i++ {
		for j := i + 1; j < len(names); j++ {
			a, b := names[i], names[j]
			if (isInit(b) && !isInit(a)) || (isInit(a) == isInit(b) && b < a) {
				names[i], names[j] = names[j], names[i]
			}
		}
	}
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	distro, _ := content.DetectDistro()
	src := ""
	if s.Engine.Watcher != nil {
		src = s.Engine.Watcher.Source
	}
	shells, _ := engine.FindShells(s.Engine.Path.ShellUser)
	writeJSON(w, map[string]any{
		"distro":     distro,
		"execSource": src,
		"shellUser":  s.Engine.Path.ShellUser,
		"shells":     shells,
		"live":       s.Opts.Live,
	})
}

// --- websocket --------------------------------------------------------------

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // playground proxies
}

func (s *Server) handleEvents(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()
	ch, unsub := s.Bus.Subscribe()
	defer unsub()

	// Reader goroutine: only to detect close.
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}()

	ping := time.NewTicker(20 * time.Second)
	defer ping.Stop()
	for {
		select {
		case ev, ok := <-ch:
			if !ok {
				return
			}
			if err := conn.WriteJSON(ev); err != nil {
				return
			}
		case <-ping.C:
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case <-done:
			return
		}
	}
}

// --- unit assets ------------------------------------------------------------

func (s *Server) handleUnitAsset(w http.ResponseWriter, r *http.Request) {
	rest := r.PathValue("rest") // "<module>/<unit>/<relpath...>"
	parts := strings.SplitN(rest, "/", 3)
	if len(parts) != 3 {
		http.Error(w, "bad path", 400)
		return
	}
	u := s.Engine.Path.Unit(parts[0] + "/" + parts[1])
	if u == nil {
		http.Error(w, "unknown unit", 404)
		return
	}
	clean := filepath.Clean(parts[2])
	if strings.HasPrefix(clean, "..") {
		http.Error(w, "bad path", 400)
		return
	}
	full := filepath.Join(u.Dir, clean)
	if _, err := os.Stat(full); err != nil {
		http.Error(w, "not found", 404)
		return
	}
	http.ServeFile(w, r, full)
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
