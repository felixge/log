package log

import (
	"runtime"
	"testing"
)

func TestNewEntry(t *testing.T) {
	pc, file, line, _ := runtime.Caller(1)
	e := NewEntry(DEBUG, "Hello %s", "World")
	fn := runtime.FuncForPC(pc).Name()

	if e.File != file {
		t.Errorf("Bad file: %s != %s", e.File, file)
	}
	if e.Line != line {
		t.Errorf("Bad line: %d != %d", e.Line, line)
	}
	if e.Function != fn {
		t.Errorf("Bad function: %d != %d", e.Function, fn)
	}
}
