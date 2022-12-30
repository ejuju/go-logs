package logs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// Log holds logging data, it has a timestamp, a level of severity and a message.
// It can also include additional data fields.
type Log struct {
	Message string         `json:"message"`
	Data    map[string]any `json:"data,omitempty"`
}

// Creates a new log with the timestamp set to the current time.
func NewLog(msg string, opts ...LogOption) *Log {
	l := &Log{Message: msg, Data: map[string]any{}}
	for _, opt := range opts {
		opt(l)
	}
	return l
}

// LogOption modifies a log.
// LogOptions are typically used to add more data to a log.
type LogOption func(*Log)

const (
	dataKeyPrefix      = "__"
	DataKeyTimestamp   = dataKeyPrefix + "created_at"
	DataKeyLevel       = dataKeyPrefix + "level"
	DataKeySrcFunction = dataKeyPrefix + "src_function"
	DataKeySrcFileLine = dataKeyPrefix + "src_file_line"
	DataKeyFSys        = dataKeyPrefix + "fsys"
)

// WithData adds more data to a log.
// The value should serializable in order to be writable to the logger output.
func WithData(key string, value any) LogOption {
	return func(l *Log) { l.Data[key] = value }
}

// WithTimestamp adds a creation datetime to the log.
func WithTimestamp() LogOption { return func(l *Log) { l.Data[DataKeyTimestamp] = time.Now() } }

// WithLevel adds a severity level to the log.
func WithLevel(lvl string) LogOption { return func(l *Log) { l.Data[DataKeyLevel] = lvl } }

// WithSrc stores the location where the log was created in the source code.
func WithSrc() LogOption {
	return func(l *Log) {
		pc, file, line, ok := runtime.Caller(2)
		if !ok {
			return
		}
		l.Data[DataKeySrcFunction] = runtime.FuncForPC(pc).Name()
		l.Data[DataKeySrcFileLine] = file + ":" + strconv.Itoa(line)
	}
}

// WithFS adds info about a file system (name and size of files) to the log.
func WithFSys(fsys fs.FS) LogOption {
	type FileInfo struct {
		Path string `json:"path"`
		Size int    `json:"size"`
	}

	return func(l *Log) {
		files := []FileInfo{}
		err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			info, err := d.Info()
			if err != nil {
				return err
			}
			files = append(files, FileInfo{
				Path: path,
				Size: int(info.Size()),
			})
			return nil
		})
		if err != nil {
			l.Data[DataKeyFSys] = err.Error()
		}
		l.Data[DataKeyFSys] = files
	}
}

// Serializer can convert a Log to bytes so that it can be written.
type Serializer func(*Log) []byte

// Returns the JSON representation of a log.
// This function will panic if the JSON marshalling of the log returns an error.
func AsJSON(l *Log) []byte {
	b, err := json.MarshalIndent(l, "", "\t")
	if err != nil {
		panic(err)
	}
	return b
}

// Returns a single line text representation of a log.
func AsSingleLine(l *Log) []byte {
	return []byte(fmt.Sprintf("%q %#v", l.Message, l.Data))
}

// LogLevel represents the severity of a log.
// Severity of logs could range anywhere between simple debug info to critical errors.
type LogLevel int

// String returns the textual representation of a level.
func (lvl LogLevel) String() string { return levelLabels[lvl] }

// Represents the level of severity of a log.
const (
	LevelUnknown LogLevel = iota // Only for temporary use, like context.TODO()
	LevelDebug                   // Debug (usually not meant to be kept in production)
	LevelInfo                    // Informative data
	LevelWarn                    // Warnings
	LevelError                   // Internal errors
	LevelPanic                   // Panics / fatal errors
)

// Holds textual representations of the log levels.
var levelLabels = [...]string{
	LevelUnknown: "UNKNOWN",
	LevelDebug:   "DEBUG",
	LevelInfo:    "INFO",
	LevelWarn:    "WARN",
	LevelError:   "ERROR",
	LevelPanic:   "PANIC",
}

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
