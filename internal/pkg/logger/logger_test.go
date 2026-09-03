package logger_test

import (
	"testing"

	"erm-dokter/internal/pkg/logger"
)

func TestLogger(t *testing.T) {
	l := logger.New()
	if l == nil {
		t.Fatal("Expected non-nil logger")
	}

	// Should execute without panicking
	l.Info("Test info message: %s", "hello")
	l.Warn("Test warn message: %d", 42)
	l.Error("Test error message: %v", "something wrong")
}
