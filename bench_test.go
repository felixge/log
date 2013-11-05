package log

import (
	"io/ioutil"
	"testing"
)

// BenchmarkDiscardLineLogger benchmarks the overhead of the logging
// abstractions provided by this library. The result can be seen as an estimate
// of maximal theoretical log throughput.
func BenchmarkDiscardLineLogger(b *testing.B) {
	l := NewLogger()
	l.Handle(Debug, NewLineWriter(ioutil.Discard))

	for i := 0; i < b.N; i++ {
		l.Debug("Hello %s", "World")
	}
}
