package logger

import (
	"log/slog"
	"os"
	"strings"
)

var logConfig = &slog.HandlerOptions{
	AddSource: false,
}

// global logger
var slogger *slog.Logger

func InitLogger(levelStr string) {
	level := getLogLevelFromString(levelStr)
	logConfig.Level = level
	slogger = slog.New(slog.NewJSONHandler(os.Stdout, logConfig))
}

func Debug(msg string, args ...interface{}) {
	slogger.Debug(msg, args...)
}

func Info(msg string, args ...interface{}) {
	slogger.Info(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	slogger.Warn(msg, args...)
}

func Error(msg string, args ...interface{}) {
	slogger.Error(msg, args...)
}

func getLogLevelFromString(levelStr string) slog.Level {
	levelStr = strings.ToLower(strings.TrimSpace(levelStr))

	switch levelStr {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo // Default to info
	}
}
