package log

import (
	"errors"
	"net/http"

	"github.com/yrss1/workout/internal/pkg/errs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"moul.io/chizap"
)

type Field zapcore.Field

func String(key, value string) Field {
	return Field(zap.String(key, value))
}

func Any(key string, value interface{}) Field {
	return Field(zap.Any(key, value))
}

func Int64(key string, value int64) Field {
	return Field(zap.Int64(key, value))
}

func Uint64(key string, value uint64) Field {
	return Field(zap.Uint64(key, value))
}

func Error(err error) Field {
	var domainError *errs.Error
	if errors.As(err, &domainError) {
		return Field(zap.Object("error", domainError))
	}
	return Field(zap.Error(err))
}

type Log struct {
	logger *zap.Logger
}

func (l *Log) Logger() *zap.Logger {
	return l.logger
}

func (l *Log) Debug(msg string, fields ...Field) {
	zf := make([]zap.Field, len(fields))
	for i, f := range fields {
		zf[i] = zap.Field(f)
	}
	l.logger.Debug(msg, zf...)
}

func (l *Log) Info(msg string, fields ...Field) {
	zf := make([]zap.Field, len(fields))
	for i, f := range fields {
		zf[i] = zap.Field(f)
	}
	l.logger.Info(msg, zf...)
}

func (l *Log) Print(msg string, fields ...Field) {
	zf := make([]zap.Field, len(fields))
	for i, f := range fields {
		zf[i] = zap.Field(f)
	}
	l.logger.Info(msg, zf...)
}

func (l *Log) Warn(msg string, fields ...Field) {
	zf := make([]zap.Field, len(fields))
	for i, f := range fields {
		zf[i] = zap.Field(f)
	}
	l.logger.Warn(msg, zf...)
}

func (l *Log) Warning(msg string, fields ...Field) {
	zf := make([]zap.Field, len(fields))
	for i, f := range fields {
		zf[i] = zap.Field(f)
	}
	l.logger.Warn(msg, zf...)
}

func (l *Log) Error(msg string, fields ...Field) {
	zf := make([]zap.Field, len(fields))
	for i, f := range fields {
		zf[i] = zap.Field(f)
	}
	l.logger.Error(msg, zf...)
}

func (l *Log) Fatal(msg string, fields ...Field) {
	zf := make([]zap.Field, len(fields))
	for i, f := range fields {
		zf[i] = zap.Field(f)
	}
	l.logger.Fatal(msg, zf...)
}

func (l *Log) Panic(msg string, fields ...Field) {
	zf := make([]zap.Field, len(fields))
	for i, f := range fields {
		zf[i] = zap.Field(f)
	}
	l.logger.Panic(msg, zf...)
}

func (l *Log) With(fields ...Field) *Log {
	if len(fields) == 0 {
		return l
	}
	zapFields := make([]zap.Field, len(fields))
	for i, field := range fields {
		zapFields[i] = zap.Field(field)
	}
	return &Log{
		logger: l.logger.With(zapFields...),
	}
}

func (l *Log) Middleware() func(next http.Handler) http.Handler {
	return chizap.New(l.logger, &chizap.Opts{
		WithReferer:   false,
		WithUserAgent: false,
	})
}

func NewLog(level string) (*Log, error) {
	config := zap.NewProductionConfig()
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, err
	}
	config.Level = lvl
	config.Development = config.Level.Level() == zapcore.DebugLevel
	config.DisableStacktrace = config.Level.Level() != zapcore.DebugLevel
	config.DisableCaller = config.Level.Level() != zapcore.DebugLevel
	config.EncoderConfig.MessageKey = "message"
	logger, err := config.Build()
	if err != nil {
		return nil, err
	}
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)
	return &Log{logger: logger}, nil
}
