package log

import (
	"bufio"
	"io"
	"os"
	"os/signal"
)

type FileWriterConfig struct {
	Path         string
	Perm         os.FileMode
	Writer       io.Writer
	Formatter    Formatter
	RotateSignal os.Signal
	ErrorHandler ErrorHandler
	BufSize      int
	Blocking     bool
	Capacity     int
	GoRoutines   int
}

type FileWriter struct {
	config  FileWriterConfig
	file    *os.File
	buf     *bufio.Writer
	opCh    chan interface{}
	entryCh chan Entry
}

type flushReq chan struct{}

type rotateReq struct{}

func NewFileWriterConfig(config FileWriterConfig) *FileWriter {
	w := &FileWriter{
		config:  config,
		opCh:    make(chan interface{}, config.Capacity),
		entryCh: make(chan Entry, config.Capacity),
	}

	if config.Writer != nil {
		w.buf = bufio.NewWriterSize(config.Writer, config.BufSize)
	} else {
		if config.RotateSignal != nil {
			rotateCh := make(chan os.Signal, 1)
			signal.Notify(rotateCh, config.RotateSignal)
			go w.rotateLoop(rotateCh)
		}
		w.open()
	}

	go w.opLoop()
	for i := 0; i < config.GoRoutines; i++ {
		go w.entryLoop()
	}
	return w
}

func NewFileWriter(path string) *FileWriter {
	config := DefaultFileWriterConfig
	config.Path = path
	return NewFileWriterConfig(config)
}

func (w *FileWriter) Log(entry Entry) {
	if w.config.Blocking {
		w.entryCh <- entry
		return
	}

	select {
	case w.entryCh <- entry:
	default:
		w.error(&ErrEntryDropped{entry})
		return
	}
}

func (w *FileWriter) Flush() {
	req := make(flushReq)
	w.opCh <- req
	<-req
}

func (w *FileWriter) open() {
	file, err := os.OpenFile(w.config.Path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, w.config.Perm)
	if err != nil {
		w.error(err)
		w.file = nil
		w.buf = nil
		return
	}
	w.file = file
	w.buf = bufio.NewWriterSize(w.file, w.config.BufSize)
}

func (w *FileWriter) opLoop() {
	for op := range w.opCh {
		switch t := op.(type) {
		case string:
			w.log(t)
		case flushReq:
			w.flush(t)
		case rotateReq:
			w.rotate()
		}
	}
}

func (w *FileWriter) entryLoop() {
	for {
		event := <-w.entryCh
		w.opCh <- w.config.Formatter.Format(event)
	}
}

func (w *FileWriter) log(message string) {
	if _, err := w.buf.WriteString(message); err != nil {
		w.error(err)
	}
}

func (w *FileWriter) flush(req flushReq) {
	if err := w.buf.Flush(); err != nil {
		w.error(err)
	}
	req <- struct{}{}
}

func (w *FileWriter) rotateLoop(rotateCh <-chan os.Signal) {
	for _ = range rotateCh {
		w.opCh <- rotateReq{}
	}
}

func (w *FileWriter) rotate() {
	if err := w.buf.Flush(); err != nil {
		w.error(err)
	}
	if err := w.file.Close(); err != nil {
		w.error(err)
	}
	w.open()
}

func (w *FileWriter) error(err error) {
	if w.config.ErrorHandler != nil {
		w.config.ErrorHandler(err)
	}
}
