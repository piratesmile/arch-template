package tlog

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const loggerKey = "ctx-logger"

var (
	// defaultLogger is the default logger. It is initialized once per package
	// include upon calling DefaultLogger.
	defaultLogger     *zap.SugaredLogger
	defaultLoggerOnce sync.Once
)

var encoderConfig = zapcore.EncoderConfig{
	TimeKey:       "timestamp",
	LevelKey:      "level",
	NameKey:       "logger",
	CallerKey:     "caller",
	MessageKey:    "message",
	StacktraceKey: "stacktrace",
	LineEnding:    zapcore.DefaultLineEnding,
	EncodeLevel:   zapcore.LowercaseLevelEncoder,
	EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(time.DateTime))
	},
	EncodeDuration: zapcore.SecondsDurationEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
}

func NewZapLogger(opt *Options) *zap.SugaredLogger {
	if opt == nil {
		opt = defaultOptions
	}
	level, err := zap.ParseAtomicLevel(opt.Level)
	if err != nil {
		level = zap.NewAtomicLevel()
	}
	writeSyncer, errWriter := opt.newWriteSyncer()

	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), writeSyncer, level)

	logger := zap.New(core,
		zap.AddStacktrace(zap.PanicLevel),
		zap.ErrorOutput(errWriter),
		zap.WithCaller(true),
	)
	zap.RedirectStdLog(logger)

	return logger.Sugar()
}

// DefaultLogger returns the default logger for the package.
func DefaultLogger() *zap.SugaredLogger {
	defaultLoggerOnce.Do(func() {
		defaultLogger = NewZapLogger(defaultOptions)
	})
	return defaultLogger
}

// WithLogger creates a new context with the provided logger attached.
func WithLogger(ctx context.Context, logger *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// FromContext returns the logger stored in the context. If no such logger
// exists, a default logger is returned.
func FromContext(ctx context.Context) *zap.SugaredLogger {
	if logger, ok := ctx.Value(loggerKey).(*zap.SugaredLogger); ok {
		return logger
	}
	return DefaultLogger()
}
