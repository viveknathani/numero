package nlog

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

func TestLogLevels(t *testing.T) {
	// save original stdout and env
	oldStdout := os.Stdout
	oldLevel := os.Getenv("LOG_LEVEL")
	defer func() {
		// restore original stdout and env
		os.Stdout = oldStdout
		os.Setenv("LOG_LEVEL", oldLevel)
		// reset logger
		logger.buf = bufio.NewWriter(os.Stdout)
	}()

	tests := []struct {
		name      string
		logLevel  string
		logFunc   func(...any)
		message   string
		expected  bool // should message be logged
		levelText string
	}{
		{"info with info level", "INFO", Info, "test info", true, "[INFO] "},
		{"debug with info level", "INFO", Debug, "test debug", false, "[DEBUG]"},
		{"error with info level", "INFO", Error, "test error", true, "[ERROR]"},
		{"info with debug level", "DEBUG", Info, "test info", true, "[INFO] "},
		{"debug with debug level", "DEBUG", Debug, "test debug", true, "[DEBUG]"},
		{"error with debug level", "DEBUG", Error, "test error", true, "[ERROR]"},
		{"info without level", "", Info, "test info", true, "[INFO] "},
		{"debug without level", "", Debug, "test debug", false, "[DEBUG]"},
		{"error without level", "", Error, "test error", true, "[ERROR]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create pipe to capture stdout
			r, w, _ := os.Pipe()
			os.Stdout = w
			logger.buf = bufio.NewWriter(w)

			// set log level
			if tt.logLevel != "" {
				os.Setenv("LOG_LEVEL", tt.logLevel)
			} else {
				os.Unsetenv("LOG_LEVEL")
			}

			// reset cached level
			level = getLogLevel()

			// log message
			tt.logFunc(tt.message)

			// read captured output
			w.Close()
			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			// verify output
			if tt.expected {
				if !strings.Contains(output, tt.message) {
					t.Errorf("expected output to contain message %q, got %q", tt.message, output)
				}
				if !strings.Contains(output, tt.levelText) {
					t.Errorf("expected output to contain level %q, got %q", tt.levelText, output)
				}
				if !strings.Contains(output, time.Now().Format(time.DateTime)[:10]) {
					t.Errorf("expected output to contain today's date, got %q", output)
				}
			} else {
				if output != "" {
					t.Errorf("expected no output, got %q", output)
				}
			}
		})
	}
}

func TestConcurrentLogging(t *testing.T) {
	// save and restore stdout
	oldStdout := os.Stdout
	defer func() {
		os.Stdout = oldStdout
		logger.buf = bufio.NewWriter(os.Stdout)
	}()

	// create pipe to capture stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	logger.buf = bufio.NewWriter(w)

	// log concurrently
	const n = 100
	done := make(chan bool)
	for i := 0; i < n; i++ {
		go func(i int) {
			Info("concurrent log", i)
			done <- true
		}(i)
	}

	// wait for all goroutines
	for i := 0; i < n; i++ {
		<-done
	}

	// read captured output
	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// verify we have n complete lines
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != n {
		t.Errorf("expected %d lines, got %d", n, len(lines))
	}

	// verify each line is properly formatted
	for _, line := range lines {
		if !strings.Contains(line, "[INFO]") {
			t.Errorf("malformed line: %q", line)
		}
		if !strings.Contains(line, "concurrent log") {
			t.Errorf("malformed line: %q", line)
		}
	}
}
