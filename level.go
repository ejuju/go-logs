package logs

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
