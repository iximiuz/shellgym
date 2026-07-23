package content

import (
	"strings"
	"testing"
)

func TestResolveVars(t *testing.T) {
	vars, err := ResolveVars(map[string]VarSpec{
		"FIXED": {Value: "hello"},
		"PICK":  {Pick: []string{"a", "b", "c"}},
		"SHELL": {Shell: "echo -n out-$((1+1))"},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if vars["FIXED"] != "hello" {
		t.Errorf("FIXED = %q", vars["FIXED"])
	}
	if vars["PICK"] != "a" && vars["PICK"] != "b" && vars["PICK"] != "c" {
		t.Errorf("PICK = %q", vars["PICK"])
	}
	if vars["SHELL"] != "out-2" {
		t.Errorf("SHELL = %q", vars["SHELL"])
	}
}

func TestResolveVarsFrom(t *testing.T) {
	vars, err := ResolveVars(map[string]VarSpec{
		"INHERITED": {From: "other-unit.TOKEN"},
	}, func(unit, name string) (string, error) {
		if unit != "other-unit" || name != "TOKEN" {
			t.Fatalf("lookup called with %s.%s", unit, name)
		}
		return "xyz", nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if vars["INHERITED"] != "xyz" {
		t.Errorf("INHERITED = %q", vars["INHERITED"])
	}
}

func TestResolveVarsErrors(t *testing.T) {
	if _, err := ResolveVars(map[string]VarSpec{"bad-name": {Value: "x"}}, nil); err == nil {
		t.Error("want error for lowercase var name")
	}
	if _, err := ResolveVars(map[string]VarSpec{"EMPTY": {}}, nil); err == nil {
		t.Error("want error for empty spec")
	}
	if _, err := ResolveVars(map[string]VarSpec{"F": {From: "no-dot"}}, func(string, string) (string, error) {
		return "", nil
	}); err == nil || !strings.Contains(err.Error(), "unit-name.VAR") {
		t.Errorf("want from-format error, got %v", err)
	}
	if _, err := ResolveVars(map[string]VarSpec{"S": {Shell: "exit 3"}}, nil); err == nil {
		t.Error("want error for failing shell var")
	}
}

func TestInterpolate(t *testing.T) {
	vars := map[string]string{"NAME": "alpha", "N": "7"}
	got := Interpolate("dir ${NAME} has ${N} files; ${UNKNOWN} stays; $NAME too", vars)
	want := "dir alpha has 7 files; ${UNKNOWN} stays; $NAME too"
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
