package nlog

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"
)

// level is controlled by LOG_LEVEL env var
// valid values: DEBUG, INFO, ERROR (default: INFO)
var level = getLogLevel()

// logger is the internal logger with mutex protection and buffered writes
var logger = struct {
	sync.Mutex
	buf *bufio.Writer
}{
	buf: bufio.NewWriter(os.Stdout),
}

func getLogLevel() string {
	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		return "INFO"
	}
	return level
}

// Info logs informational messages
func Info(args ...any) {
	if level == "INFO" || level == "DEBUG" {
		log("[INFO] ", args...)
	}
}

// Debug logs debug messages
func Debug(args ...any) {
	if level == "DEBUG" {
		log("[DEBUG]", args...)
	}
}

// Error logs error messages (always logged)
func Error(args ...any) {
	log("[ERROR]", args...)
}

// log is the internal logging function
func log(level string, args ...any) {
	logger.Lock()
	defer logger.Unlock()

	// pre-allocate slice with capacity for prefix + args
	out := make([]any, 0, 2+len(args))
	out = append(out, level, time.Now().Format(time.DateTime))
	out = append(out, args...)

	fmt.Fprintln(logger.buf, out...)

	// ensure logs are written immediately
	logger.buf.Flush()
}
