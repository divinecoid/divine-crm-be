package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Logger is a custom logger
type Logger struct {
	level  string
	format string
}

// NewLogger creates a new logger
func NewLogger(level, format string) *Logger {
	return &Logger{
		level:  level,
		format: format,
	}
}

// Debug logs debug messages
func (l *Logger) Debug(msg string, args ...interface{}) {
	if l.shouldLog("debug") {
		l.log("DEBUG", msg, args...)
	}
}

// Info logs info messages
func (l *Logger) Info(msg string, args ...interface{}) {
	if l.shouldLog("info") {
		l.log("INFO", msg, args...)
	}
}

// Warn logs warning messages
func (l *Logger) Warn(msg string, args ...interface{}) {
	if l.shouldLog("warn") {
		l.log("WARN", msg, args...)
	}
}

// Error logs error messages
func (l *Logger) Error(msg string, args ...interface{}) {
	if l.shouldLog("error") {
		l.log("ERROR", msg, args...)
	}
}

// Fatal logs fatal messages and exits
func (l *Logger) Fatal(msg string, args ...interface{}) {
	l.log("FATAL", msg, args...)
	os.Exit(1)
}

func (l *Logger) log(level, msg string, args ...interface{}) {
	timestamp := time.Now().Format(time.RFC3339)

	if l.format == "json" {
		// JSON format
		logMsg := fmt.Sprintf(`{"time":"%s","level":"%s","msg":"%s"`, timestamp, level, msg)
		for i := 0; i < len(args); i += 2 {
			if i+1 < len(args) {
				logMsg += fmt.Sprintf(`,"%v":"%v"`, args[i], args[i+1])
			}
		}
		logMsg += "}"
		log.Println(logMsg)
	} else {
		// Text format
		logMsg := fmt.Sprintf("[%s] %s - %s", timestamp, level, msg)
		for i := 0; i < len(args); i += 2 {
			if i+1 < len(args) {
				logMsg += fmt.Sprintf(" %v=%v", args[i], args[i+1])
			}
		}
		log.Println(logMsg)
	}
}

func (l *Logger) shouldLog(level string) bool {
	levels := map[string]int{
		"debug": 0,
		"info":  1,
		"warn":  2,
		"error": 3,
		"fatal": 4,
	}

	configLevel := levels[l.level]
	msgLevel := levels[level]

	return msgLevel >= configLevel
}
