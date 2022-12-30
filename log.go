package logs

import (
	"io/fs"
	"runtime"
	"strconv"
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
