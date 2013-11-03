package log

import (
	"testing"
)

func TestLogger(t *testing.T) {
	l := NewLogger()
	w := NewTestWriter()
	l.Handle(Debug, w)

	l.Debug("Test A")
	l.Info("Test B")
	l.Warn("Test C")
	err := l.Error("Test D")
	if err.Error() != "Test D" {
		t.Errorf("Bad error return: %s", err)
	}

	if entries := len(w.Entries); entries != 4 {
		t.Errorf("Bad #entries: %d", entries)
	}

	if !w.MatchLevel("A$", Debug) {
		t.Errorf("Missing entry: A")
	}
	if !w.MatchLevel("B$", Info) {
		t.Errorf("Missing entry: B")
	}
	if !w.MatchLevel("C$", Warn) {
		t.Errorf("Missing entry: C")
	}
	if !w.MatchLevel("D$", Error) {
		t.Errorf("Missing entry: D")
	}
}
