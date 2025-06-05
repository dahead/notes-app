package logger

import (
	"io"
	"log"
	"os"
)

var (
	debugLogger *log.Logger
	isDebug     bool
)

func init() {
	// Check if debug logging is enabled via environment variable
	if os.Getenv("NOTESAPP_DEBUG") != "" {
		isDebug = true
		debugLogger = log.New(os.Stdout, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		// If debug is disabled, use discard writer to silence output
		debugLogger = log.New(io.Discard, "", 0)
	}
}

// Debug logs a debug message if debug logging is enabled
func Debug(format string, v ...interface{}) {
	if isDebug {
		debugLogger.Printf(format, v...)
	}
}
