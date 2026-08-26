package logger

import (
	"fmt"
	"os"
	"testing"
)

func TestLoggerInit(t *testing.T) {
	tests := []struct {
		name     string
		mode     LogMode
		filePath string
	}{
		{name: "DEBUG mode", mode: DEBUG, filePath: "test_debug_init.md"},
		{name: "ERROR mode", mode: ERROR, filePath: "test_error_init.md"},
		{name: "PROD mode", mode: PROD, filePath: "test_prod_init.md"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.Remove(tt.filePath)
			defer resetForTests()

			logger := Init(tt.mode, tt.filePath)

			if logger.GetMode() != tt.mode {
				t.Errorf("expected mode %s, got %s", tt.mode, logger.GetMode())
			}

			if _, err := os.Stat(tt.filePath); os.IsNotExist(err) {
				t.Errorf("log file was not created")
			}
		})
	}
}

func TestLoggerDebugMode(t *testing.T) {
	defer os.Remove("test_debug_mode.md")
	defer resetForTests()

	logger := Init(DEBUG, "test_debug_mode.md")

	logger.Info("This is info")
	logger.Debug("This is debug")
	logger.Error("This is error")
	logger.Warn("This is warn")

	content, err := os.ReadFile("test_debug_mode.md")
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	logContent := string(content)
	if c := countLogEntries(logContent, "INFO"); c != 1 {
		t.Errorf("expected 1 INFO log, got %d", c)
	}
	if c := countLogEntries(logContent, "DEBUG"); c != 1 {
		t.Errorf("expected 1 DEBUG log, got %d", c)
	}
	if c := countLogEntries(logContent, "ERROR"); c != 1 {
		t.Errorf("expected 1 ERROR log, got %d", c)
	}
	if c := countLogEntries(logContent, "WARN"); c != 1 {
		t.Errorf("expected 1 WARN log, got %d", c)
	}
}

func TestLoggerErrorMode(t *testing.T) {
	defer os.Remove("test_error_mode.md")
	defer resetForTests()

	logger := Init(ERROR, "test_error_mode.md")

	logger.Info("This is info")
	logger.Debug("This is debug")
	logger.Error("This is error")
	logger.Warn("This is warn")

	content, err := os.ReadFile("test_error_mode.md")
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	logContent := string(content)
	if c := countLogEntries(logContent, "INFO"); c != 0 {
		t.Errorf("expected 0 INFO logs_test in ERROR mode, got %d", c)
	}
	if c := countLogEntries(logContent, "DEBUG"); c != 0 {
		t.Errorf("expected 0 DEBUG logs_test in ERROR mode, got %d", c)
	}
	if c := countLogEntries(logContent, "ERROR"); c != 1 {
		t.Errorf("expected 1 ERROR log, got %d", c)
	}
}

func TestLoggerMust(t *testing.T) {
	defer os.Remove("test_must.md")
	defer resetForTests()

	logger := Init(DEBUG, "test_must.md")

	// Should not panic
	logger.Must(nil, "no error")

	// Should panic
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, but didn't get one")
		}
	}()

	logger.Must(fmt.Errorf("test error"), "test context")
}

func TestLoggerMustNotNil(t *testing.T) {
	defer os.Remove("test_must_not_nil.md")
	defer resetForTests()

	logger := Init(DEBUG, "test_must_not_nil.md")

	// Should not panic
	logger.MustNotNil("not nil", "check")

	// Should panic
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, but didn't get one")
		}
	}()

	var nilValue interface{}
	logger.MustNotNil(nilValue, "nil check")
}

func countLogEntries(content, level string) int {
	pattern := fmt.Sprintf("] %s:", level)
	count := 0
	for i := 0; i < len(content); {
		if len(content) > i+len(pattern)-1 && content[i:i+len(pattern)] == pattern {
			count++
			i += len(pattern)
		} else {
			i++
		}
	}
	return count
}
