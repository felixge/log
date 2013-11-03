package log

import (
	"runtime"
	"testing"
)

// BenchmarkRuntimeCaller determines the overhead of calling runtime.Caller to
// determine if it's reasonable to invoke it for every log call.
func BenchmarkRuntimeCaller(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runtime.Caller(2)
	}
}
