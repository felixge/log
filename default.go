package log

import (
	"fmt"
	"os"
	"runtime"
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
	DefaultFormatter        = NewLineFormatter(DefaultLayout, nil)
	DefaultColorFormatter   = NewLineFormatter(DefaultLayout, DefaultTermStyle)
	DefaultMessageFormatter = NewLineFormatter("message", nil)
	DefaultConfig           = Config{
		FlushTimeout: 30 * time.Second,
		FatalExit:    true,
	}
	DefaultErrorHandler = func(err error) {
		e := NewEntry(ERROR, "%s", err)
		fmt.Fprint(os.Stderr, DefaultFormatter.Format(e))
	}
	DefaultBufSize          = 4096
	DefaultFileWriterConfig = FileWriterConfig{
		Perm:         0600,
		Formatter:    DefaultFormatter,
		RotateSignal: syscall.SIGUSR1,
		ErrorHandler: DefaultErrorHandler,
		BufSize:      DefaultBufSize,
		Blocking:     false,
		Capacity:     1024,
		GoRoutines:   runtime.NumCPU(),
	}
	DefaultTermConfig = FileWriterConfig{
		Writer:       os.Stdout,
		Formatter:    DefaultColorFormatter,
		ErrorHandler: DefaultErrorHandler,
		BufSize:      DefaultBufSize,
		Blocking:     true,
		Capacity:     0,
		GoRoutines:   1,
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
