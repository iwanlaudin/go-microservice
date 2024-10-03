package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

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

func New(logLevel string) Logger {
	config := zap.Config{
		Encoding:         "json", // Or “console” if you want terminal format
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

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

func fieldsToInterface(fields []zapcore.Field) []interface{} {
	result := make([]interface{}, len(fields)*2)
	for i, field := range fields {
		result[i*2] = field.Key

		switch field.Type {
		case zapcore.StringType:
			result[i*2+1] = field.String
		case zapcore.Int64Type:
			result[i*2+1] = field.Integer
		case zapcore.BoolType:
			result[i*2+1] = field.Integer == 1
		case zapcore.ErrorType:
			result[i*2+1] = field.Interface
		default:
			result[i*2+1] = "unknown"
		}
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
