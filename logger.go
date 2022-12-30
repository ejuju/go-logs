package logs

import (
	"bytes"
	"io"
	"sync"
)

type LoggerFunc func(*Log) error

type DefaultLogger struct {
	Writers     []io.Writer // For ex: stdout and/or file
	Serializer  Serializer  // For ex: As JSON
	BaseOptions []LogOption // For ex: creation timestamp, source code location
	LogPrefix   string      // For ex: "HTTP" or "Server Name"
	LogSuffix   string      // For ex: ",\n" to seperate JSON logs by commas and line breaks
}

func (dl *DefaultLogger) LoggerFunc() (LoggerFunc, error) {
	// Init writer with stdout and logfile
	w := newWriterWrapper(dl.Writers...)

	// Init mutex
	mu := &sync.Mutex{}

	return func(l *Log) error {
		mu.Lock()
		defer mu.Unlock()

		// Apply base options to log
		for _, opt := range dl.BaseOptions {
			opt(l)
		}

		// Get log bytes
		b := bytes.Join([][]byte{
			[]byte(dl.LogPrefix),
			dl.Serializer(l),
			[]byte(dl.LogSuffix + "\n"),
		}, nil)

		// Write log
		_, err := w.Write(b)
		return err
	}, nil
}

// writerWrapper is a utility type that implements io.Writer by wrapping one or more io.Writers
type writerWrapper []io.Writer

// newWriterWrapper instanciates a new WriterWrapper.
func newWriterWrapper(w ...io.Writer) writerWrapper { return w }

// Does not fail if one of the underlying writers returns an error.
func (ww writerWrapper) Write(b []byte) (int, error) {
	var numBytesWritten = 0
	var errs errWrapper
	for _, w := range ww {
		n, err := w.Write(b)
		if err != nil {
			errs = append(errs, err)
		}
		numBytesWritten += n
	}
	if errs != nil {
		return numBytesWritten, errs
	}
	return numBytesWritten, nil
}

// errWrapper is a utility type that wraps one or more errors.
type errWrapper []error

// Error is the implementation of the error interface.
// It joins error messages with ", ".
func (ew errWrapper) Error() string {
	out := ""
	for i, err := range ew {
		if i > 0 {
			out += ", "
		}
		out += err.Error()
	}
	return out
}
