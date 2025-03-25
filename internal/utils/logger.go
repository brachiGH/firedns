package utils

import (
	"log"
	"os"
)

// Logger structure
type Logger struct {
	Info  *log.Logger
	Warn  *log.Logger
	Error *log.Logger
}

// NewLogger initializes and returns a logger instance
func NewLogger() *Logger {
	return &Logger{
		Info:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		Warn:  log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile),
		Error: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}
