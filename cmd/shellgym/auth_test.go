package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsLoopbackAddr(t *testing.T) {
	for addr, want := range map[string]bool{
		"127.0.0.1:63636": true,
		"127.0.0.2:63636": true,
		"localhost:63636": true,
		"[::1]:63636":     true,
		":63636":          false,
		"0.0.0.0:63636":   false,
		"[::]:63636":      false,
		"10.0.0.5:63636":  false,
		"example.com:80":  false,
		"63636":           false,
	} {
		if got := isLoopbackAddr(addr); got != want {
			t.Errorf("isLoopbackAddr(%q) = %v, want %v", addr, got, want)
		}
	}
}

func TestLoadAuthToken(t *testing.T) {
	t.Setenv(authTokenEnv, "")
	if tok, err := loadAuthToken(""); err != nil || tok != "" {
		t.Errorf("no token configured: %q, %v", tok, err)
	}

	t.Setenv(authTokenEnv, " from-env\n")
	if tok, err := loadAuthToken(""); err != nil || tok != "from-env" {
		t.Errorf("env token: %q, %v", tok, err)
	}

	// The file wins over the environment.
	file := filepath.Join(t.TempDir(), "token")
	if err := os.WriteFile(file, []byte("from-file\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if tok, err := loadAuthToken(file); err != nil || tok != "from-file" {
		t.Errorf("file token: %q, %v", tok, err)
	}

	if err := os.WriteFile(file, []byte("\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := loadAuthToken(file); err == nil {
		t.Error("an empty token file was accepted")
	}
	if _, err := loadAuthToken(filepath.Join(t.TempDir(), "missing")); err == nil {
		t.Error("a missing token file was accepted")
	}
}
