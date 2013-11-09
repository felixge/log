package log

import (
	"bytes"
	"fmt"
	"runtime"
	"sync"
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

func TestLogger_Flush(t *testing.T) {
	// Configure l as a *Logger that writes to a *LineWriter that outputs to a
	// slow io.Writer which sleep for dt on every Write call.
	var (
		wg    sync.WaitGroup
		b     = &bytes.Buffer{}
		dt    = 10 * time.Millisecond
		count = 10
		w     = NewLineWriter(NewSlowWriter(b, dt), DefaultFormat, DefaultTermStyle)
		l     = NewLogger(w)
	)

	start := time.Now()
	for i := 1; i <= count; i++ {
		l.Debug("Message %s", i)
	}
	// The above log statements should have been async. To verify this, we check
	// that the total duration for writing them out did not exceed half the time
	// a single sync operation would have taken.
	if duration := time.Since(start); duration > dt/2 {
		t.Fatalf("Expected async logging, but detected sync behavior. %s", duration)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(dt / 2) // Make sure l.Flush() is in progress before we continue.
		l.Debug("Log during flush") // Try to write a log entry while flushing happens
		// Verify that the above call was blocked until Flush() finished. Flush
		// takes at least dt*count, so verify that we waited at least that long.
		if duration := time.Since(start); duration < dt*time.Duration(count) {
			t.Fatalf("Logging did not seem to block during Flush. %s", duration)
		}
	}()
	l.Flush()

	// Flush will take at least dt*count (because of the slow io.Writer), so
	// report an error if it finishes faster for some reason.
	if duration := time.Since(start); duration < dt*time.Duration(count) {
		t.Fatalf("Flush seems to have dropped messages. %s", duration)
	}
	// Make sure the goroutine we spawned finishes before we terminate the test
	wg.Wait()
}

func TestEntryFormat_defaultFormat(t *testing.T) {
	e := Entry{
		Time:     time.Now(),
		Level:    Info,
		Message:  "foo",
		File:     "bar.go",
		Line:     23,
		Function: "foo.bar",
	}

	str := e.Format(DefaultFormat)
	expected := fmt.Sprintf(
		"[%s UTC] [%s] %s (%s:%d)",
		e.Time.UTC().Format("2006-01-02 15:04:05.000"),
		Info,
		e.Message,
		e.Function,
		e.Line,
	)
	if str != expected {
		t.Errorf("Bad result: %q != %q", str, expected)
	}
}

func TestEntryFormat_customFormat(t *testing.T) {
	e := Entry{
		Time:     time.Now(),
		Level:    Info,
		Message:  "foo",
		File:     "bar.go",
		Line:     23,
		Function: "foo.bar",
	}

	str := e.Format("2006/01/02 15:04:05.000 level message file/line/function")
	expected := fmt.Sprintf(
		"%s %s %s %s/%d/%s",
		e.Time.Format("2006/01/02 15:04:05.000"),
		Info,
		e.Message,
		e.File,
		e.Line,
		e.Function,
	)
	if str != expected {
		t.Errorf("Bad result: %q != %q", str, expected)
	}
}

func TestLogger_Panic(t *testing.T) {
	var (
		wg   sync.WaitGroup
		w    = NewTestWriter()
		l    = NewLogger(w)
		file string
		line int
	)

	l.SetExit(false)

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer l.Panic()
		_, file, line, _ = runtime.Caller(0)
		panic("oh no")
	}()
	wg.Wait()

	if !w.MatchLevel("panic: oh no", Fatal) {
		t.Error("Panic was not logged.")
	}
	e := w.Entries[0]

	if e.File != file {
		t.Errorf("Bad line: %s != %s", e.File, file)
	}
	if e.Line != line+1 {
		t.Errorf("Bad file: %d != %d", e.Line, line+1)
	}
}

func TestNewEntry(t *testing.T) {
	pc, file, line, _ := runtime.Caller(0)
	e := NewEntry(Debug, "Hello %s", "World")
	fn := runtime.FuncForPC(pc).Name()

	t.Logf("fn: %s", fn)
	if e.File != file {
		t.Errorf("Bad file: %s != %s", e.File, file)
	}
	if e.Line != line+1 {
		t.Errorf("Bad line: %d != %d", e.Line, line+1)
	}
	if e.Function != fn {
		t.Errorf("Bad function: %d != %d", e.Function, fn)
	}
}
