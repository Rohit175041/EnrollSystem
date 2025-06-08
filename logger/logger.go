package logger

import (
	"fmt"
	"os"
	"time"
)

const (
	INFO = iota
	WARN
	ERROR
	FATAL
)

var levelLabels = []string{"INFO", "WARN", "ERROR", "FATAL"}

type Logger struct {
	minLevel int
	file     *os.File
}

// NewLogger creates a new logger and opens a log file for writing.
// It truncates the file if it already exists to start fresh each time.
func NewLogger(minLevel int, filePath string) (*Logger, error) {
	// Open the file with O_TRUNC to clear previous logs
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return nil, err
	}

	return &Logger{
		minLevel: minLevel,
		file:     file,
	}, nil
}

func (l *Logger) Close() {
	if l.file != nil {
		_ = l.file.Close()
	}
}

func (l *Logger) log(level int, msg string, args ...interface{}) {
	if level < l.minLevel {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	label := levelLabels[level]
	formattedMsg := fmt.Sprintf(msg, args...)
	logLine := fmt.Sprintf("[%s] [%s] %s\n", timestamp, label, formattedMsg)

	// Print to console (optional, uncomment if needed)
	// fmt.Print(logLine)

	// Write to file
	if l.file != nil {
		_, _ = l.file.WriteString(logLine)
	}

	if level == FATAL {
		l.Close()
		os.Exit(1)
	}
}

func (l *Logger) Info(msg string, args ...interface{})  { l.log(INFO, msg, args...) }
func (l *Logger) Warn(msg string, args ...interface{})  { l.log(WARN, msg, args...) }
func (l *Logger) Error(msg string, args ...interface{}) { l.log(ERROR, msg, args...) }
func (l *Logger) Fatal(msg string, args ...interface{}) { l.log(FATAL, msg, args...) }
