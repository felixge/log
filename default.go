package log

import (
	"os"
	"time"
)

var (
	DefaultLayout    = "[2006-01-02 15:04:05.000 UTC] [level] message (function:line)"
	DefaultTermStyle = map[Level]TermStyle{
		DEBUG: DarkGrey,
		INFO:  0,
		WARN:  Yellow,
		ERROR: Red,
		FATAL: White | BgRed,
	}
	DefaultFormatter = NewLineFormatter(DefaultLayout, DefaultTermStyle)
	DefaultConfig    = Config{
		FlushTimeout: 30 * time.Second,
		FatalExit:    true,
	}
	DefaultLogger = NewLogger(DefaultConfig, NewLineHandler(os.Stdout, DefaultFormatter))
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
