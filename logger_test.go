package log

import (
	"bytes"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	l, w := NewTestLogger()

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

	if !w.MatchLevel("A$", DEBUG) {
		t.Errorf("Missing entry: A")
	}
	if !w.MatchLevel("B$", INFO) {
		t.Errorf("Missing entry: B")
	}
	if !w.MatchLevel("C$", WARN) {
		t.Errorf("Missing entry: C")
	}
	if !w.MatchLevel("D$", ERROR) {
		t.Errorf("Missing entry: D")
	}
}

func TestLogger_Flush(t *testing.T) {
	// Configure l as a *Logger that writes to a *LineHandler that outputs to a
	// slow io.Writer which sleep for dt on every Write call.
	var (
		wg    sync.WaitGroup
		b     = &bytes.Buffer{}
		dt    = 10 * time.Millisecond
		count = 10
		w     = NewLineHandler(NewSlowWriter(b, dt), DefaultFormat, DefaultTermStyle)
		l     = NewLogger(DefaultConfig, w)
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
		time.Sleep(dt / 2)          // Make sure l.Flush() is in progress before we continue.
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

func TestLogger_Panic(t *testing.T) {
	var (
		wg   sync.WaitGroup
		w    = NewTestHandler()
		l    = NewLogger(Config{FatalExit: false}, w)
		file string
		line int
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer l.Panic()
		_, file, line, _ = runtime.Caller(1)
		panic("oh no")
	}()
	wg.Wait()

	if !w.MatchLevel("panic: oh no", FATAL) {
		t.Error("Panic was not logged.")
	}
	e := w.Entries[0]

	if e.File != file {
		t.Errorf("Bad file: %s != %s", e.File, file)
	}
	if e.Line != line {
		t.Errorf("Bad line: %d != %d", e.Line, line)
	}
}
