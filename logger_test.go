package log

import (
	"testing"
	"time"
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

func TestEntryFormat(t *testing.T) {
	e := Entry{
		Time:    time.Now(),
		Level:   Info,
		Message: "foo",
		File:    "bar.go",
		Line:    23,
	}

	str := e.Format(DefaultFormat)
	t.Log(str)
}

func BenchmarkEntryFormat(b *testing.B) {
	e := Entry{
		Time:    time.Now(),
		Level:   Info,
		Message: "foo",
		File:    "bar.go",
		Line:    23,
	}

	for i := 0; i < b.N; i++ {
		e.Format("15:04:05.000 [level] message (file:line)")
	}
}
