package log

import (
	"fmt"
	"os"
	"syscall"
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
	DefaultFormatter      = NewLineFormatter(DefaultLayout, nil)
	DefaultColorFormatter = NewLineFormatter(DefaultLayout, DefaultTermStyle)
	DefaultConfig         = Config{
		FlushTimeout: 30 * time.Second,
		FatalExit:    true,
	}
	DefaultErrorHandler = func(err error) {
		e := NewEntry(ERROR, "%s", err)
		fmt.Fprint(os.Stderr, DefaultFormatter.Format(e))
	}
	DefaultFileWriterConfig = FileWriterConfig{
		Perm:         0600,
		Formatter:    DefaultFormatter,
		RotateSignal: syscall.SIGUSR1,
		ErrorHandler: DefaultErrorHandler,
		BufSize:      4096,
		Capacity:     1024,
		Blocking:     false,
	}
	DefaultTermConfig = FileWriterConfig{
		Writer:       os.Stdout,
		Formatter:    DefaultColorFormatter,
		ErrorHandler: DefaultErrorHandler,
		BufSize:      4096,
		Capacity:     1024,
		Blocking:     true,
	}
	DefaultWriter = NewFileWriterConfig(DefaultTermConfig)
	DefaultLogger = NewLogger(DefaultConfig, DefaultWriter)
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

// @TODO Panic level
