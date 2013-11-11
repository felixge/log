package log

import (
	"io"
	"strings"
	"sync"
)

// NewLineHandler returns a Handler that writes newline separated log entries
// to the given io.Writer w using the provided format and style.
func NewLineHandler(w io.Writer, format string, style map[Level]TermStyle) *LineHandler {
	l := &LineHandler{
		w:        w,
		format:   NewFormat(format),
		style:    style,
		entries:  make(chan Entry, 1024),
		flushReq: make(chan chan struct{}),
	}
	go l.loop()
	return l
}

// LineHandler is a Handler that provides newline separated logging.
type LineHandler struct {
	w         io.Writer
	format    *Format
	style     map[Level]TermStyle
	entries   chan Entry // @TODO rename to entries
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

		line := strings.Replace(l.format.Format(e), "\n", "", -1)
		if style, ok := l.style[e.Level]; ok {
			line = style.Format(line)
		}
		io.WriteString(l.w, line+"\n")
	}
}
