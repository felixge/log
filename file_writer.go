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
}

type FileWriter struct {
	config FileWriterConfig
	file   *os.File
	writer io.Writer
	opCh   chan interface{}
}

type flusher interface {
	Flush() error
}

type flushReq chan struct{}

type rotateReq struct{}

func NewFileWriterConfig(config FileWriterConfig) *FileWriter {
	w := &FileWriter{
		config: config,
		opCh:   make(chan interface{}, config.Capacity),
	}

	if config.Writer != nil {
		w.setWriter(config.Writer, config.BufSize)
	} else {
		if config.RotateSignal != nil {
			rotateCh := make(chan os.Signal, 1)
			signal.Notify(rotateCh, config.RotateSignal)
			go w.rotateLoop(rotateCh)
		}
		w.open()
	}

	go w.opLoop()
	return w
}

func NewFileWriter(path string) *FileWriter {
	config := DefaultFileWriterConfig
	config.Path = path
	return NewFileWriterConfig(config)
}

func (w *FileWriter) Log(entry Entry) {
	message := w.config.Formatter.Format(entry)

	if w.config.Blocking {
		w.opCh <- message
		return
	}

	select {
	case w.opCh <- message:
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
		w.writer = nil
		return
	}
	w.file = file
	w.setWriter(file, w.config.BufSize)
}

func (w *FileWriter) opLoop() {
	for op := range w.opCh {
		switch t := op.(type) {
		case string:
			w.log(t)
		case flushReq:
			w.flush()
			t <- struct{}{}
		case rotateReq:
			w.rotate()
		}
	}
}

func (w *FileWriter) setWriter(writer io.Writer, bufSize int) {
	if bufSize > 0 {
		w.writer = bufio.NewWriterSize(writer, bufSize)
		return
	}
	w.writer = writer
}

func (w *FileWriter) log(message string) {
	if _, err := io.WriteString(w.writer, message); err != nil {
		w.error(err)
	}
}

func (w *FileWriter) flush() {
	if flusher, ok := w.writer.(flusher); ok {
		if err := flusher.Flush(); err != nil {
			w.error(err)
		}
	}
}

func (w *FileWriter) rotateLoop(rotateCh <-chan os.Signal) {
	for _ = range rotateCh {
		w.opCh <- rotateReq{}
	}
}

func (w *FileWriter) rotate() {
	w.flush()
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
