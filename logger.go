package logs

import (
	"time"
)

// Logger can write and format a log.
type Logger interface {
	Log(*Log) error
}

// Creates a new log
func New(lvl Level, msg string, opts ...LogOption) *Log {
	l := &Log{
		CreatedAt: time.Now(),
		Level:     LevelLabels[lvl],
		Message:   msg,
		Data:      map[string]any{},
	}
	for _, opt := range opts {
		opt(l)
	}
	return l
}

// Level represents the severity of a log.
// Severity of logs could range anywhere between simple debug info to critical errors.
type Level int

// Represents the level of severity of a log.
const (
	LevelUnknown Level = iota // Only for temporary use, like context.TODO()
	LevelDebug                // Debug (usually not meant to be kept in production)
	LevelInfo                 // Informative data
	LevelWarn                 // Warnings
	LevelError                // Internal errors
	LevelPanic                // Panics / fatal errors
)

// Holds textual representations of the log levels.
var LevelLabels = [...]string{
	LevelUnknown: "UNKNOWN",
	LevelDebug:   "DEBUG",
	LevelInfo:    "INFO",
	LevelWarn:    "WARN",
	LevelError:   "ERROR",
	LevelPanic:   "PANIC",
}

// Log holds logging data, it has a timestamp, a level of severity and a message.
// It can also include additional data fields.
type Log struct {
	CreatedAt time.Time      `json:"created_at"`
	Level     string         `json:"level"`
	Message   string         `json:"message"`
	Data      map[string]any `json:"data,omitempty"`
}
