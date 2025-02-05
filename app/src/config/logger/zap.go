package logger

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	h *zap.Logger
}

func (zlogger *zapLogger) InitLogger(levelStr string) {
	levelStr = strings.ToLower(strings.TrimSpace(levelStr))

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "time"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig:     encoderCfg,
		OutputPaths: []string{
			"stdout",
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
	}
	if levelStr != "info" {
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
		config.Development = true
	}

	l := zap.Must(config.Build())
	zlogger.h = l.WithOptions(zap.AddCallerSkip(1)) // Skip the logger.go file
}

func (zlogger *zapLogger) Debug(msg string, args ...any) {
	zapFields := zlogger.anySliceToZapFieldSlice(args)
	zlogger.h.Debug(msg, zapFields...)
}

func (zlogger *zapLogger) Info(msg string, args ...any) {
	zapFields := zlogger.anySliceToZapFieldSlice(args)
	zlogger.h.Info(msg, zapFields...)
}

func (zlogger *zapLogger) Warn(msg string, args ...any) {
	zapFields := zlogger.anySliceToZapFieldSlice(args)
	zlogger.h.Warn(msg, zapFields...)
}

func (zlogger *zapLogger) Error(msg string, args ...any) {
	zapFields := zlogger.anySliceToZapFieldSlice(args)
	zlogger.h.Error(msg, zapFields...)
}

func (zlogger *zapLogger) anySliceToZapFieldSlice(fields []any) []zap.Field {
	var zapFields []zap.Field
	n := len(fields)
	for i := 0; i < n; i += 2 {
		zapFields = append(zapFields, zap.Any(fields[i].(string), fields[i+1]))
	}
	return zapFields
}

func (zlogger *zapLogger) Any(key string, value any) zap.Field {
	return zap.Any(key, value)
}
