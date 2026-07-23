// Package skills holds the authoring guides that get embedded into the
// shellgym binary; `shellgym skills <name>` prints them.
package skills

import "embed"

//go:embed *.md
var FS embed.FS
