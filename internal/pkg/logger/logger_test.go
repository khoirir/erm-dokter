package logger_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"strings"
	"testing"

	"erm-dokter/internal/pkg/logger"
)

func TestLogger_DefaultAndOutputJSON(t *testing.T) {
	var buf bytes.Buffer
	l := logger.NewWithOptions(&buf, "json", slog.LevelDebug, "erm-dokter-test")
	if l == nil {
		t.Fatal("Expected non-nil logger")
	}

	l.Info("Dokter %s berhasil login dengan id %d", "Budi", 101)

	output := buf.String()
	if output == "" {
		t.Fatal("Expected log output, got empty string")
	}

	var parsed map[string]any
	if err := json.Unmarshal([]byte(output), &parsed); err != nil {
		t.Fatalf("Log output is not valid JSON: %v, raw output: %s", err, output)
	}

	if parsed["level"] != "INFO" {
		t.Errorf("Expected level 'INFO', got '%v'", parsed["level"])
	}

	expectedMsg := "Dokter Budi berhasil login dengan id 101"
	if parsed["msg"] != expectedMsg {
		t.Errorf("Expected msg '%s', got '%v'", expectedMsg, parsed["msg"])
	}

	if parsed["app"] != "erm-dokter-test" {
		t.Errorf("Expected app 'erm-dokter-test', got '%v'", parsed["app"])
	}

	source, ok := parsed["source"].(map[string]any)
	if !ok {
		t.Fatalf("Expected source map in JSON log, got %v", parsed["source"])
	}
	file, _ := source["file"].(string)
	if !strings.Contains(file, "logger_test.go") {
		t.Errorf("Expected source file to contain 'logger_test.go', got '%s'", file)
	}
}

func TestLogger_LevelsAndCallerSource(t *testing.T) {
	var buf bytes.Buffer
	l := logger.NewWithOptions(&buf, "json", slog.LevelDebug, "erm-dokter")

	l.Debug("Debug test: %s", "detail")
	l.Warn("Warn test: %d", 404)
	l.Error("Error test: %s", "fatal issue")

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("Expected 3 log lines, got %d", len(lines))
	}

	expectedLevels := []string{"DEBUG", "WARN", "ERROR"}
	for i, line := range lines {
		var item map[string]any
		if err := json.Unmarshal([]byte(line), &item); err != nil {
			t.Fatalf("Line %d is not valid JSON: %v", i, err)
		}
		if item["level"] != expectedLevels[i] {
			t.Errorf("Line %d expected level %s, got %v", i, expectedLevels[i], item["level"])
		}
		source, ok := item["source"].(map[string]any)
		if !ok || !strings.Contains(source["file"].(string), "logger_test.go") {
			t.Errorf("Line %d expected source file 'logger_test.go', got %v", i, source)
		}
	}
}

func TestLogger_WithAttributes(t *testing.T) {
	var buf bytes.Buffer
	baseLogger := logger.NewWithOptions(&buf, "json", slog.LevelInfo, "erm-dokter")
	childLogger := baseLogger.With("dokter_id", "DR001", "unit", "ralan")

	childLogger.Info("Memulai pemeriksaan pasien")

	var parsed map[string]any
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if parsed["dokter_id"] != "DR001" {
		t.Errorf("Expected dokter_id 'DR001', got '%v'", parsed["dokter_id"])
	}
	if parsed["unit"] != "ralan" {
		t.Errorf("Expected unit 'ralan', got '%v'", parsed["unit"])
	}
}

func TestLogger_ContextMethods(t *testing.T) {
	var buf bytes.Buffer
	l := logger.NewWithOptions(&buf, "json", slog.LevelInfo, "erm-dokter")

	l.InfoContext(context.Background(), "Order resep dibuat", "no_rawat", "2026/09/10/0001", "jumlah_obat", 3)

	var parsed map[string]any
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if parsed["no_rawat"] != "2026/09/10/0001" {
		t.Errorf("Expected no_rawat '2026/09/10/0001', got '%v'", parsed["no_rawat"])
	}
	if parsed["jumlah_obat"] != float64(3) {
		t.Errorf("Expected jumlah_obat 3, got '%v'", parsed["jumlah_obat"])
	}
}

func TestLogger_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	l := logger.NewWithOptions(&buf, "text", slog.LevelInfo, "erm-dokter")

	l.Info("Pesan text format")

	output := buf.String()
	if !strings.Contains(output, "level=INFO") || !strings.Contains(output, "msg=\"Pesan text format\"") {
		t.Errorf("Unexpected text log output: %s", output)
	}
}

func TestLogger_NewFromEnv(t *testing.T) {
	os.Setenv("LOG_FORMAT", "text")
	os.Setenv("LOG_LEVEL", "warn")
	os.Setenv("APP_NAME", "custom-app")
	defer func() {
		os.Unsetenv("LOG_FORMAT")
		os.Unsetenv("LOG_LEVEL")
		os.Unsetenv("APP_NAME")
	}()

	l := logger.New()
	if l == nil {
		t.Fatal("Expected non-nil logger from New()")
	}
	if l.Handler() == nil {
		t.Fatal("Expected non-nil slog.Handler")
	}
}
