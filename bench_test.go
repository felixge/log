package log

import (
	"fmt"
	"io/ioutil"
	"runtime"
	"sync"
	"testing"
	"time"
)

// BenchmarkDiscardLineLogger benchmarks the overhead of the logging
// abstractions provided by this library. The result can be seen as an estimate
// of maximal theoretical log throughput.
func BenchmarkDiscardLineLogger(b *testing.B) {
	var (
		cpus   = runtime.NumCPU()
		before = runtime.GOMAXPROCS(cpus)
		wg     sync.WaitGroup
	)
	defer runtime.GOMAXPROCS(before)


	config := DefaultFileWriterConfig
	config.Blocking = true
	config.Writer = ioutil.Discard
	l := NewLogger(DefaultConfig)
	l.Handle(DEBUG, NewFileWriterConfig(config))

	start := time.Now()
	b.ResetTimer()
	for g := 0; g < cpus; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < b.N; i++ {
				l.Debug("Hello %s", "World")
			}
		}()
	}
	wg.Wait()
	if err := l.Flush(); err != nil {
		b.Fatalf("Flush error: %s", err)
	}
	b.StopTimer()
	duration := time.Since(start)

	total := b.N * cpus
	hz := NewHz(total, duration)
	b.Logf("%s (%d ops in %s)", hz, total, duration)
}

var prefixes = map[int]string{
	1000:    "k",
	1000000: "M",
}

func NewHz(count int, d time.Duration) Hz {
	return Hz(float64(count) / d.Seconds())
}

type Hz float64

func (h Hz) String() string {
	var (
		prefix string
		factor int
	)
	for f, p := range prefixes {
		if int(h) > f {
			prefix = p
			factor = f
		}
	}
	h = h / Hz(factor)
	return fmt.Sprintf("%.2f %sHZ", float64(h), prefix)
}
