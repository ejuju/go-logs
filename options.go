package logs

import (
	"io/fs"
	"runtime"
	"strconv"
)

// LogOption modifies a log.
// LogOptions are typically used to add more data to a log.
type LogOption func(*Log)

// WithData adds more data to a log.
// The value should serializable in order to be writable to the logger output.
func WithData(key string, value any) LogOption {
	return func(l *Log) { l.Data[key] = value }
}

// WithSourceCodeLocation adds the information about where the log is
// coming from in the source code.
// Use offset 0 for the current position in the source code.
// Add 1 to the offset with result in logging the calling function's location.
//
// For ex: "main.foo at /main.go:45"
func WithSourceCodeLocation(key string, offset int) LogOption {
	return func(l *Log) {
		pc, file, line, ok := runtime.Caller(2 + offset)
		if !ok {
			return
		}
		l.Data[key] = runtime.FuncForPC(pc).Name() + " at " + file + ":" + strconv.Itoa(line)
	}
}

// WithFS adds info about a file system (sub-directories and files).
func WithFS(key string, fsys fs.FS) LogOption {
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
			l.Data[key] = err.Error()
		}
		l.Data[key] = files
	}
}
