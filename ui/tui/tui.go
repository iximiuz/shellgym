// Package ui defines the pluggable UI contract. The daemon can host any
// number of UIs (web, future TUI); each consumes the same engine, content
// snapshot, and event bus.
package tui

import "context"

// UI is a frontend implementation (web server, TUI, ...).
type UI interface {
	// Run serves the UI until ctx is cancelled.
	Run(ctx context.Context) error
}
