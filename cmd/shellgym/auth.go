package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	// authTokenEnv is the alternative to --auth-token-file. There is no flag
	// taking the token itself: a command line is readable by every user on
	// the box.
	authTokenEnv       = "SHELLGYM_AUTH_TOKEN"
	authTokenFileUsage = "file holding the web UI auth token (alternative: the " + authTokenEnv + " environment variable)"
)

// loadAuthToken returns the daemon's auth token: the contents of file when
// one is given, the environment variable otherwise. "" means no auth.
func loadAuthToken(file string) (string, error) {
	if file == "" {
		return strings.TrimSpace(os.Getenv(authTokenEnv)), nil
	}
	data, err := os.ReadFile(file)
	if err != nil {
		return "", fmt.Errorf("auth token file: %w", err)
	}
	token := strings.TrimSpace(string(data))
	if token == "" {
		return "", fmt.Errorf("auth token file %s is empty", file)
	}
	return token, nil
}

// isLoopbackAddr reports whether a listen address ("127.0.0.1:63636",
// ":63636") is reachable from this host only.
func isLoopbackAddr(addr string) bool {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return false
	}
	if host == "localhost" {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}
