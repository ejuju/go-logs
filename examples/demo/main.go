package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/ejuju/go-logs"
)

func main() {
	// Create log file
	f, err := os.Create(fmt.Sprintf("logs.%d.txt", time.Now().Unix()))
	if err != nil {
		panic(err)
	}

	// Configure default logger
	config := &logs.DefaultLogger{
		Writers:    []io.Writer{f, os.Stdout}, // write to both file and stdout
		Serializer: logs.AsJSON,               // serialize as JSON
		LogSuffix:  ",",                       // append comma to each log
		BaseOptions: []logs.LogOption{
			logs.WithTimestamp(), // store timestamp in logs
			logs.WithSrc(),       // store source code location in logs
		},
	}

	// Get logging func
	log, err := config.LoggerFunc()
	if err != nil {
		panic(err)
	}

	// Write a simple log
	log(logs.NewLog("hey, i'm a log"))

	// Write a log with additional data
	log(logs.NewLog(
		"hey, i'm another log",
		logs.WithLevel(logs.LevelInfo.String()),
		logs.WithFSys(os.DirFS(".")),
		logs.WithData("some_key", "some data"),
		logs.WithData("other_key", os.Args),
	))
}
