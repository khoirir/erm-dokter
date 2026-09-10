package logger

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type Logger struct {
	handler slog.Handler
}

func New() *Logger {
	format := os.Getenv("LOG_FORMAT")
	levelStr := os.Getenv("LOG_LEVEL")
	appName := os.Getenv("APP_NAME")
	if appName == "" {
		appName = "erm-dokter"
	}

	var w io.Writer = os.Stdout

	if !isRunningInTest() {
		logFilePath := os.Getenv("LOG_FILE_PATH")
		if logFilePath == "" {
			logFilePath = "logs/app.log"
		}

		if !strings.EqualFold(logFilePath, "none") && !strings.EqualFold(logFilePath, "off") {
			if !filepath.IsAbs(logFilePath) {
				logFilePath = filepath.Join(findProjectRoot(), logFilePath)
			}
			if err := os.MkdirAll(filepath.Dir(logFilePath), 0755); err == nil {
				if file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644); err == nil {
					w = io.MultiWriter(os.Stdout, file)
				}
			}
		}
	}

	return NewWithOptions(w, format, parseLogLevel(levelStr), appName)
}

func NewWithOptions(w io.Writer, format string, level slog.Level, appName string) *Logger {
	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				if src, ok := a.Value.Any().(*slog.Source); ok {
					src.File = cleanSourcePath(src.File)
				}
			}
			return a
		},
	}

	var h slog.Handler
	if strings.EqualFold(format, "text") {
		h = slog.NewTextHandler(w, opts)
	} else {
		h = slog.NewJSONHandler(w, opts)
	}

	if appName != "" {
		h = h.WithAttrs([]slog.Attr{slog.String("app", appName)})
	}

	return &Logger{handler: h}
}

func (l *Logger) Handler() slog.Handler {
	return l.handler
}

func (l *Logger) With(args ...any) *Logger {
	attrs := argsToAttrs(args)
	return &Logger{
		handler: l.handler.WithAttrs(attrs),
	}
}

func (l *Logger) Debug(msg string, v ...any) {
	l.logWithDepth(context.Background(), 3, slog.LevelDebug, msg, v...)
}

func (l *Logger) Info(msg string, v ...any) {
	l.logWithDepth(context.Background(), 3, slog.LevelInfo, msg, v...)
}

func (l *Logger) Warn(msg string, v ...any) {
	l.logWithDepth(context.Background(), 3, slog.LevelWarn, msg, v...)
}
func (l *Logger) Error(msg string, v ...any) {
	l.logWithDepth(context.Background(), 3, slog.LevelError, msg, v...)
}

func (l *Logger) InfoContext(ctx context.Context, msg string, args ...any) {
	l.logAttrsWithDepth(ctx, 3, slog.LevelInfo, msg, args...)
}

func (l *Logger) WarnContext(ctx context.Context, msg string, args ...any) {
	l.logAttrsWithDepth(ctx, 3, slog.LevelWarn, msg, args...)
}

func (l *Logger) ErrorContext(ctx context.Context, msg string, args ...any) {
	l.logAttrsWithDepth(ctx, 3, slog.LevelError, msg, args...)
}

type contextKey string

const RequestIDKey contextKey = "request_id"

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

func GetRequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if val, ok := ctx.Value(RequestIDKey).(string); ok {
		return val
	}
	return ""
}

func (l *Logger) logWithDepth(ctx context.Context, depth int, level slog.Level, msg string, v ...any) {
	if !l.handler.Enabled(ctx, level) {
		return
	}
	if len(v) > 0 {
		msg = fmt.Sprintf(msg, v...)
	}
	var pcs [1]uintptr
	runtime.Callers(depth, pcs[:])
	r := slog.NewRecord(time.Now(), level, msg, pcs[0])
	if reqID := GetRequestID(ctx); reqID != "" {
		r.AddAttrs(slog.String("request_id", reqID))
	}
	_ = l.handler.Handle(ctx, r)
}

func (l *Logger) logAttrsWithDepth(ctx context.Context, depth int, level slog.Level, msg string, args ...any) {
	if !l.handler.Enabled(ctx, level) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(depth, pcs[:])
	r := slog.NewRecord(time.Now(), level, msg, pcs[0])
	attrs := argsToAttrs(args)
	if reqID := GetRequestID(ctx); reqID != "" {
		attrs = append(attrs, slog.String("request_id", reqID))
	}
	r.AddAttrs(attrs...)
	_ = l.handler.Handle(ctx, r)
}

func parseLogLevel(levelStr string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(levelStr)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func cleanSourcePath(file string) string {
	normalized := filepath.ToSlash(file)
	if idx := strings.Index(normalized, "erm-dokter/"); idx != -1 {
		return normalized[idx+len("erm-dokter/"):]
	}
	return filepath.Base(file)
}

func argsToAttrs(args []any) []slog.Attr {
	var attrs []slog.Attr
	for i := 0; i < len(args); i++ {
		switch v := args[i].(type) {
		case slog.Attr:
			attrs = append(attrs, v)
		case string:
			if i+1 < len(args) {
				attrs = append(attrs, slog.Any(v, args[i+1]))
				i++
			} else {
				attrs = append(attrs, slog.String("key", v))
			}
		default:
			attrs = append(attrs, slog.Any(fmt.Sprintf("attr_%d", i), v))
		}
	}
	return attrs
}

func isRunningInTest() bool {
	if flag.Lookup("test.v") != nil {
		return true
	}
	base := filepath.Base(os.Args[0])
	return strings.HasSuffix(base, ".test") || strings.HasSuffix(base, ".test.exe")
}

func findProjectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "."
}
