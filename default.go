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
		PANIC: White | BgRed,
	}
	DefaultFormatter        = NewLineFormatter(DefaultLayout, nil)
	DefaultColorFormatter   = NewLineFormatter(DefaultLayout, DefaultTermStyle)
	DefaultMessageFormatter = NewLineFormatter("message", nil)
	DefaultConfig           = Config{
		FlushTimeout: 30 * time.Second,
	}
	DefaultErrorHandler = func(err error) {
		e := NewEntry(ERROR, "%s", err)
		fmt.Fprint(os.Stderr, DefaultFormatter.Format(e))
	}
	DefaultFileWriterConfig = FileWriterConfig{
		Perm:          0600,
		Formatter:     DefaultFormatter,
		RotateSignal:  syscall.SIGUSR1,
		ErrorHandler:  DefaultErrorHandler,
		Blocking:      false,
		Capacity:      1024,
		BufSize:       4096,
		FlushInterval: time.Second,
	}
	DefaultTermConfig = FileWriterConfig{
		Writer:       os.Stdout,
		Formatter:    DefaultColorFormatter,
		ErrorHandler: DefaultErrorHandler,
		Blocking:     true,
	}
	DefaultWriter = NewFileWriterConfig(DefaultTermConfig)

	DefaultLogger = NewLogger(DefaultConfig, DefaultWriter)
)

func Debug(args ...interface{}) {
	DefaultLogger.Log(NewEntryWithStack(DEBUG, 3, 1, args...))
}

func Info(args ...interface{}) {
	DefaultLogger.Log(NewEntryWithStack(INFO, 3, 1, args...))
}

func Warn(args ...interface{}) {
	DefaultLogger.Log(NewEntryWithStack(WARN, 3, 1, args...))
}

func Error(args ...interface{}) error {
	e := NewEntryWithStack(ERROR, 3, 1, args...)
	DefaultLogger.Log(e)
	return NewError(e)
}

func Panic(args ...interface{}) {
	DefaultLogger.Log(NewEntryWithStack(PANIC, 3, 1, args...))
}

// @TODO Panic level
