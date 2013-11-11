package log

import (
	"path/filepath"
	"runtime"
	"time"
)

func NewEntry(lvl Level, args ...interface{}) Entry {
	fn, file, line := getCaller()
	return Entry{
		Time:     time.Now(),
		Level:    lvl,
		Message:  formatMessage(args),
		File:     file,
		Line:     line,
		Function: fn,
	}
}

func getCaller() (fn, file string, line int) {
	var (
		skip   = 0
		pkgDir string
		ok     bool
		pc     uintptr
	)

	for ; ; skip++ {
		pc, file, line, ok = runtime.Caller(skip)
		if !ok {
			break
		} else if skip == 0 {
			pkgDir = filepath.Dir(file)
			continue
		} else if dir := filepath.Dir(file); dir != pkgDir {
			fn = runtime.FuncForPC(pc).Name()
			if fn != "runtime.panic" {
				break
			}
		}
	}
	return fn, file, line
}

type Entry struct {
	Time     time.Time
	Level    Level
	Message  string
	File     string
	Function string
	Line     int
}
