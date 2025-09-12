// internal/jsonlog/jsonlog.go
package jsonlog

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

// Severity levels for log entries.
type Level int8

const (
	LevelInfo Level = iota
	LevelError
	LevelFatal
	LevelOff
)

func (l Level) String() string {
	switch l {
	case LevelInfo:
		return "INFO"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}

// Logger struct with a mutex for safe concurrent writes.
type Logger struct {
	out      io.Writer
	minLevel Level
	mu       sync.Mutex
}

// New returns a new Logger that writes to stdout (or another writer).
func New(out io.Writer, minLevel Level) *Logger {
	return &Logger{out: out, minLevel: minLevel}
}

// PrintInfo logs an informational message.
func (l *Logger) PrintInfo(message string, properties map[string]string) {
	l.print(LevelInfo, message, properties)
}

// PrintError logs an error message.
func (l *Logger) PrintError(err error, properties map[string]string) {
	l.print(LevelError, err.Error(), properties)
}

// PrintFatal logs a fatal error and exits the app.
func (l *Logger) PrintFatal(err error, properties map[string]string) {
	l.print(LevelFatal, err.Error(), properties)
	os.Exit(1)
}

// print handles writing the JSON log entry.
func (l *Logger) print(level Level, message string, properties map[string]string) error {
	if level < l.minLevel {
		return nil
	}

	// Core log entry data
	entry := struct {
		Level      string            `json:"level"`
		Time       string            `json:"time"`
		Message    string            `json:"message"`
		Properties map[string]string `json:"properties,omitempty"`
		Trace      string            `json:"trace,omitempty"`
	}{
		Level:      level.String(),
		Time:       time.Now().UTC().Format(time.RFC3339),
		Message:    message,
		Properties: properties,
	}

	// For errors and fatals, include stack trace
	if level >= LevelError {
		entry.Trace = string(debug.Stack())
	}

	js, err := json.Marshal(entry)
	if err != nil {
		js = []byte(fmt.Sprintf(`{"level": %q, "time": %q, "message": "unable to marshal log message: %s"}`,
			level.String(), time.Now().UTC().Format(time.RFC3339), err.Error()))
	}

	// Ensure each log entry is on its own line
	l.mu.Lock()
	defer l.mu.Unlock()
	_, err = l.out.Write(append(js, '\n'))
	return err
}

// Write makes Logger satisfy io.Writer (useful for HTTP server error logs).
func (l *Logger) Write(p []byte) (n int, err error) {
	return len(p), l.print(LevelError, string(p), nil)
}
