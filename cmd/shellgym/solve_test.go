package main

import "testing"

func TestExpandVars(t *testing.T) {
	vars := map[string]string{"PAUSE": "2", "PAUSE_MAX": "9", "BOGUS": "carrots", "LO": "1", "HI": "3", "REF": "$PAUSE"}
	cases := map[string]string{
		"sleep $PAUSE && hostname":         "sleep 2 && hostname",
		"sleep ${PAUSE}s; echo $PAUSE_MAX": "sleep 2s; echo 9",
		"date --$BOGUS || whoami":          "date --carrots || whoami",
		"echo $PAUSEX $? $HOME":            "echo $PAUSEX $? $HOME",
		"#!type sleep $PAUSE":              "#!type sleep 2",
		// Adjacent references resolve independently, in either spelling.
		"ls sales-q[$LO$HI].csv":     "ls sales-q[13].csv",
		"ls sales-q[${LO}${HI}].csv": "ls sales-q[13].csv",
		"echo $HI$LO$HI":             "echo 313",
		// A substituted value is typed verbatim, never expanded again.
		"echo $REF":        "echo $PAUSE",
		"echo $((PAUSE+1))": "echo $((PAUSE+1))",
		"echo ${PAUSE":      "echo ${PAUSE",
	}
	for in, want := range cases {
		if got := expandVars(in, vars); got != want {
			t.Errorf("%q: got %q, want %q", in, got, want)
		}
	}
}
