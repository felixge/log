package main

import (
	"fmt"
	"github.com/felixge/log"
	"time"
)

func main() {
	config := log.DefaultFileWriterConfig
	config.Path = "log.txt"
	config.ErrorHandler = nil
	config.Blocking = true
	file := log.NewFileWriterConfig(config)
	l := log.NewLogger(log.DefaultConfig, file)

	start := time.Now()
	for i := 0; i < 1000000; i++{
		l.Debug("Entry %d", i)
		//time.Sleep(time.Microsecond)
	}
	l.Flush()
	fmt.Printf("duration: %s\n", time.Since(start))
}
