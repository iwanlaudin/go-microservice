package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger adalah interface untuk logger kustom kita
type Logger interface {
	Debug(msg string, fields ...zapcore.Field)
	Info(msg string, fields ...zapcore.Field)
	Warn(msg string, fields ...zapcore.Field)
	Error(msg string, fields ...zapcore.Field)
	Fatal(msg string, fields ...zapcore.Field)
}

type zapLogger struct {
	sugaredLogger *zap.SugaredLogger
}

// New membuat instance baru dari logger
func New(logLevel string) Logger {
	config := zap.NewProductionConfig()

	// Set level berdasarkan konfigurasi
	switch logLevel {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.StacktraceKey = "" // Disable stacktrace

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	return &zapLogger{
		sugaredLogger: logger.Sugar(),
	}
}

func (l *zapLogger) Debug(msg string, fields ...zapcore.Field) {
	l.sugaredLogger.Debugw(msg, fieldsToInterface(fields)...)
}

func (l *zapLogger) Info(msg string, fields ...zapcore.Field) {
	l.sugaredLogger.Infow(msg, fieldsToInterface(fields)...)
}

func (l *zapLogger) Warn(msg string, fields ...zapcore.Field) {
	l.sugaredLogger.Warnw(msg, fieldsToInterface(fields)...)
}

func (l *zapLogger) Error(msg string, fields ...zapcore.Field) {
	l.sugaredLogger.Errorw(msg, fieldsToInterface(fields)...)
}

func (l *zapLogger) Fatal(msg string, fields ...zapcore.Field) {
	l.sugaredLogger.Fatalw(msg, fieldsToInterface(fields)...)
	os.Exit(1)
}

// fieldsToInterface mengkonversi zap.Field ke interface{} untuk SugaredLogger
func fieldsToInterface(fields []zapcore.Field) []interface{} {
	result := make([]interface{}, len(fields)*2)
	for i, field := range fields {
		result[i*2] = field.Key
		result[i*2+1] = field.Interface
	}
	return result
}

// Helper untuk membuat field String
func String(key, value string) zapcore.Field {
	return zap.String(key, value)
}

// Helper untuk membuat field Error
func Error(err error) zapcore.Field {
	return zap.Error(err)
}
