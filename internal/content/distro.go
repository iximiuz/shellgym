package content

import (
	"os"
	"os/exec"
	"strings"
)

// DetectDistro reads /etc/os-release and returns the distro ID (e.g.
// "ubuntu", "rocky") and its ID_LIKE list (e.g. ["debian"], ["rhel",
// "centos", "fedora"]).
func DetectDistro() (id string, like []string) {
	raw, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return "linux", nil
	}
	for _, line := range strings.Split(string(raw), "\n") {
		if v, ok := strings.CutPrefix(line, "ID="); ok {
			id = strings.Trim(v, `"`)
		} else if v, ok := strings.CutPrefix(line, "ID_LIKE="); ok {
			like = strings.Fields(strings.Trim(v, `"`))
		}
	}
	if id == "" {
		id = "linux"
	}
	return id, like
}

// DetectCaps returns host capability tags used by unit `requires:` filters.
// Currently detected: "systemd" - a reachable system systemd instance
// (absent when there is no reachable systemd instance).
func DetectCaps() []string {
	var caps []string
	out, err := exec.Command("systemctl", "is-system-running").Output()
	state := strings.TrimSpace(string(out))
	// Non-zero exit is fine (e.g. "degraded"); only no-output/offline means
	// there is no systemd to talk to.
	if (err == nil || state != "") && state != "offline" {
		caps = append(caps, "systemd")
	}
	return caps
}
