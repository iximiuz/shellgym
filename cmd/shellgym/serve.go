package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strings"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/iximiuz/labs-content/tools/shellgym/internal/bus"
	"github.com/iximiuz/labs-content/tools/shellgym/internal/content"
	"github.com/iximiuz/labs-content/tools/shellgym/internal/engine"
	"github.com/iximiuz/labs-content/tools/shellgym/internal/state"
	"github.com/iximiuz/labs-content/tools/shellgym/ui/webui"
)

func newServeCmd() *cobra.Command {
	var (
		pathDir   string
		addr      string
		stateDir  string
		runDir    string
		shellUser string
		live      bool
	)
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Run the ShellGym daemon (web UI + validation engine)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return serve(pathDir, addr, stateDir, runDir, shellUser, live)
		},
	}
	cmd.Flags().StringVar(&pathDir, "path", "", "learning path directory (required)")
	cmd.Flags().StringVar(&addr, "addr", ":63636", "web UI listen address")
	cmd.Flags().StringVar(&stateDir, "state", "/var/lib/shellgym", "state directory")
	cmd.Flags().StringVar(&runDir, "run", "/run/shellgym", "runtime directory (socket, check shims)")
	cmd.Flags().StringVar(&shellUser, "user", "", "observed login user (default: from path.yaml)")
	cmd.Flags().BoolVar(&live, "live", false,
		"student-facing mode: strip solve scripts from on-disk unit files and disable the debug API")
	_ = cmd.MarkFlagRequired("path")
	return cmd
}

func serve(pathDir, addr, stateDir, runDir, shellUser string, live bool) error {
	if live {
		n, err := content.StripSolveScripts(pathDir)
		if err != nil {
			return err
		}
		log.Printf("live mode: stripped %d solve script(s) from %s", n, pathDir)
	}

	distro, like := content.DetectDistro()
	caps := content.DetectCaps()
	path, err := content.Load(pathDir, distro, like, caps)
	if err != nil {
		return err
	}
	if shellUser != "" {
		path.ShellUser = shellUser
	}

	// Command line watching (readline uprobe) is optional: when it cannot
	// start, the daemon runs without it and units that require it are
	// marked unsupported - browsable, never activated.
	lines := engine.NewLineWatcher()
	if err := lines.Start(engine.LoginShell(path.ShellUser)); err != nil {
		log.Printf("line watcher: unavailable, units requiring readline are marked unsupported: %v", err)
	} else {
		defer lines.Close()
		log.Printf("line watcher: uprobe on %s", lines.Target)
		caps = append(caps, "readline")
		path.ApplyCaps(caps)
	}
	log.Printf("loaded path %q: %d modules, distro=%s caps=%v", path.ID, len(path.Modules), distro, caps)

	st, err := state.Open(stateDir, path.ID)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(runDir, 0o755); err != nil {
		return err
	}
	selfExe, err := os.Executable()
	if err != nil {
		return err
	}
	checksDir := filepath.Join(runDir, "bin")
	if err := engine.WriteCheckShims(checksDir, selfExe); err != nil {
		return err
	}
	sockPath := filepath.Join(runDir, "gym.sock")

	watcher := engine.NewExecWatcher()
	if err := watcher.Start(); err != nil {
		return fmt.Errorf("exec watcher: %w", err)
	}
	defer watcher.Close()
	log.Printf("exec watcher: %s", watcher.Source)

	b := bus.New()
	eng := engine.New(path, st, b, watcher, engine.Options{
		ChecksDir: checksDir,
		SockPath:  sockPath,
		Lines:     lines,
	})

	if err := engine.ServeCheckAPI(sockPath, path.ShellUser, watcher, eng); err != nil {
		return err
	}

	eng.Resume()
	defer eng.Shutdown()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	srv := webui.New(addr, eng, b, webui.Options{Live: live})
	return srv.Run(ctx)
}

func newValidateCmd() *cobra.Command {
	var pathDir string
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Lint and render a learning path without running it",
		RunE: func(cmd *cobra.Command, args []string) error {
			distro, like := content.DetectDistro()
			// Validate with ALL capabilities assumed so requires-gated
			// units are checked too.
			path, err := content.Load(pathDir, distro, like, []string{"*"})
			if err != nil {
				return err
			}
			units := 0
			for _, m := range path.Modules {
				for _, u := range m.Units {
					units++
					vars := map[string]string{}
					for name := range u.Front.Vars {
						vars[name] = "SAMPLE"
					}
					defTask := ""
					if len(u.Tasks) == 1 {
						defTask = u.Tasks[0].Name
					}
					if _, err := content.RenderUnit(content.Interpolate(u.Body, vars), "/unit-assets/"+u.ID+"/", defTask); err != nil {
						return err
					}
				}
				if m.Intro != "" {
					if _, err := content.RenderUnit(m.Intro, ""); err != nil {
						return err
					}
				}
			}
			cmd.Printf("OK: %d modules, %d units\n", len(path.Modules), units)
			keys := path.VariantKeys()
			names := make([]string, 0, len(keys))
			for k := range keys {
				names = append(names, k)
			}
			sort.Strings(names)
			for _, k := range names {
				cmd.Printf("variant %s: %s\n", k, strings.Join(keys[k], ", "))
				if len(keys[k]) == 1 {
					cmd.Printf("  warning: a single value is always drawn - units in variant %s=%s are always shown\n", k, keys[k][0])
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&pathDir, "path", "", "learning path directory (required)")
	_ = cmd.MarkFlagRequired("path")
	return cmd
}
