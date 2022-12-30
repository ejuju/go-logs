package logs

import "encoding/json"

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
