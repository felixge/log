package log

import (
	"os"
	"time"
)

var (
	DefaultLogger = NewLogger(NewLineHandler(os.Stdout, DefaultFormat, DefaultTermStyle))
	// DefaultFormat defines the default log format used by NewTermLogger.
	DefaultFormat = "[2006-01-02 15:04:05.000 UTC] [level] message (function:line)"
	// DefaultTermStyle defines the default colors/style used by NewTermLogger
	DefaultTermStyle = map[Level]TermStyle{
		DEBUG: DarkGrey,
		INFO:  0,
		WARN:  Yellow,
		ERROR: Red,
		FATAL: White | BgRed,
	}
	DefaultExit = true
	DefaultFlushTimeout = 30 * time.Second
)

func Debug(args ...interface{}) {
	DefaultLogger.Debug(args...)
}

func Info(args ...interface{}) {
	DefaultLogger.Info(args...)
}

func Warn(args ...interface{}) {
	DefaultLogger.Warn(args...)
}

func Error(args ...interface{}) error {
	return DefaultLogger.Error(args...)
}

func Fatal(args ...interface{}) {
	DefaultLogger.Fatal(args...)
}
