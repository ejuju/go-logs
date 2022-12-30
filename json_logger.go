package logs

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// JSONLogger is a logger that writes logs to a file,
// in JSON format.
type JSONLogger struct {
	out   io.Writer
	tries int
	mu    sync.RWMutex
}

// NewJSONLogger allocates a new JSONLogger.
// It handles creating the logs directory (if needed) and log file.
func NewJSONLogger(dirname string) (*JSONLogger, error) {
	// Create directory for log files if needed
	err := os.Mkdir(dirname, os.ModePerm)
	if errors.Is(err, os.ErrExist) {
		err = nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to make directory: %w", err)
	}

	// Create log file
	f, err := os.Create(fmt.Sprintf("%s/logs.%d.txt", dirname, time.Now().Unix()))
	if err != nil {
		return nil, fmt.Errorf("failed to make create log file: %w", err)
	}

	return &JSONLogger{out: f, tries: 5}, nil
}

// Log writes log data in a JSON format to the underlying log file.
func (jfl *JSONLogger) Log(l *Log) error {
	jfl.mu.Lock()
	defer jfl.mu.Unlock()

	for i := 0; i < jfl.tries; i++ {
		b, err := json.MarshalIndent(l, "", "\t")
		if err != nil {
			continue
		}
		_, err = jfl.out.Write(append(b, ',', '\n'))
		if err != nil {
			continue
		}
		break
	}
	return nil
}
