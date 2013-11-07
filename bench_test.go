package log

import (
	"io/ioutil"
	"testing"
	"time"
)

// BenchmarkDiscardLineLogger benchmarks the overhead of the logging
// abstractions provided by this library. The result can be seen as an estimate
// of maximal theoretical log throughput.
func BenchmarkDiscardLineLogger(b *testing.B) {
	l := NewLogger()
	l.Handle(Debug, NewLineWriter(ioutil.Discard, DefaultFormat, nil))

	for i := 0; i < b.N; i++ {
		l.Debug("Hello %s", "World")
	}
	if err := l.Flush(); err != nil {
		panic(err)
	}
}

// BenchmarkEntryFormat tests the performance of the entry formatting function.
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
