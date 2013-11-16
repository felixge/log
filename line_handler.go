package log

import (
	"io"
	"sync"
)

// NewLineHandler returns a Handler that writes newline separated log entries
// to the given io.Writer w using the provided format and style.
func NewLineHandler(w io.Writer, formatter Formatter) *LineHandler {
	l := &LineHandler{
		w:         w,
		formatter: formatter,
		flushReq:  make(chan chan struct{}),
		entries:   make(chan Entry),
	}
	go l.loop()
	return l
}

// LineHandler is a Handler that provides newline separated logging.
type LineHandler struct {
	w         io.Writer
	formatter Formatter
	entries   chan Entry
	flushReq  chan chan struct{}
	flushLock sync.Mutex
}

// HandleLog writes the given log entry to a new line.
// @TODO Process entries in another goroutine.
func (l *LineHandler) HandleLog(e Entry) {
	l.flushLock.Lock()
	defer l.flushLock.Unlock()

	l.entries <- e
}

// Flush waits for any buffered log Entries to be written out.
// @TODO Make this block any HandleLog
func (l *LineHandler) Flush() {
	l.flushLock.Lock()
	defer l.flushLock.Unlock()

	flushReq := make(chan struct{})
	l.flushReq <- flushReq
	<-flushReq
}

func (l *LineHandler) loop() {
	var flushReq chan struct{}
	for {
		var e Entry
		if flushReq == nil {
			select {
			case e = <-l.entries:
			case flushReq = <-l.flushReq:
				continue
			}
		} else {
			select {
			case e = <-l.entries:
			default:
				flushReq <- struct{}{}
				flushReq = nil
				continue
			}
		}

		line := l.formatter.Format(e)
		io.WriteString(l.w, line)
	}
}
