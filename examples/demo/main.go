package main

import (
	"os"

	"github.com/ejuju/go-logs"
)

func main() {
	logger, err := logs.NewJSONLogger("./your_log_dir")
	if err != nil {
		panic(err)
	}

	// Log a simple message
	_ = logger.Log(logs.New(logs.LevelInfo, "hey, i'm a log"))

	// Log a message with additional data
	_ = logger.Log(logs.New(
		logs.LevelError,                         // Defines the log severity level
		"i'm the log message",                   // Defines the log message
		logs.WithSourceCodeLocation("src", 0),   // Adds source code location in the log
		logs.WithFS("fs", os.DirFS(".")),        // Adds info about a file system in the log
		logs.WithData("more_data", "some data"), // Adds some more data to the log
		logs.WithData("again", os.Args),         // Accepts any type that can be serialized.
	))
}
