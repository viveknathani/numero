package nlog

import (
	"fmt"
	"os"
	"time"
)

// Info should be used for informational logging.
func Info(args ...any) {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "INFO" || logLevel == "DEBUG" {
		// we add the extra space here so that the width becomes consistent with DEBUG and ERROR
		logInternal("[INFO]"+" ", args...)
	}
}

// Debug should be used for debugging level logging.
func Debug(args ...any) {
	if os.Getenv("LOG_LEVEL") == "DEBUG" {
		logInternal("[DEBUG]", args...)
	}
}

// Error should be used for error logging. Will be emitted always.
func Error(args ...any) {
	logInternal("[ERROR]", args...)
}

// logInternal is the internal implementation of numero's lightweight logger
func logInternal(level string, args ...any) {
	currentTime := time.Now()
	prefix := []any{level, currentTime.Format(time.DateTime)}
	fmt.Println(append(prefix, args...)...)
}
