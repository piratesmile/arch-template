package tlog

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ Logger = (*zapLogger)(nil)

var (
	once sync.Once
	std  Logger
)

func Init(opt *Options) {
	once.Do(func() {
		std = NewZapLogger(opt)
		// 使用全局 logger 记录日志，caller 需要再加一层
		stdLogger, ok := std.(*zapLogger)
		if ok {
			stdLogger.logger = stdLogger.logger.WithOptions(zap.AddCallerSkip(1))
		}
	})
}

type zapLogger struct {
	logger *zap.Logger
}

func NewZapLogger(opt *Options) Logger {
	if opt == nil {
		opt = defaultOptions
	}
	level, err := zap.ParseAtomicLevel(opt.Level)
	if err != nil {
		level = zap.NewAtomicLevel()
	}
	encoderConfig := zap.NewProductionEncoderConfig()

	encoderConfig.EncodeTime = func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(time.Format("2006-01-02 15:04:05.999"))
	}
	writeSyncer, errWriter := opt.newWriteSyncer()

	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), writeSyncer, level)

	logger := zap.New(core,
		zap.AddStacktrace(zap.PanicLevel),
		zap.AddCallerSkip(1),
		zap.ErrorOutput(errWriter),
		zap.WithCaller(true),
	)
	zap.RedirectStdLog(logger)

	return &zapLogger{logger: logger}
}

func (z *zapLogger) Debug(ctx context.Context, msg string, fields Fields) {
	z.logger.Debug(msg, getFields(ctx, fields)...)
}

func (z *zapLogger) Info(ctx context.Context, msg string, fields Fields) {
	z.logger.Info(msg, getFields(ctx, fields)...)
}

func (z *zapLogger) Warn(ctx context.Context, msg string, fields Fields) {
	z.logger.Warn(msg, getFields(ctx, fields)...)
}

func (z *zapLogger) Error(ctx context.Context, msg string, fields Fields) {
	z.logger.Error(msg, getFields(ctx, fields)...)
}

func (z *zapLogger) Panic(ctx context.Context, msg string, fields Fields) {
	z.logger.Panic(msg, getFields(ctx, fields)...)
}

func (z *zapLogger) Fatal(ctx context.Context, msg string, fields Fields) {
	z.logger.Fatal(msg, getFields(ctx, fields)...)
}

func (z *zapLogger) Sync() error {
	return z.logger.Sync()
}

func getFields(ctx context.Context, fields Fields) []zap.Field {
	if fields == nil {
		fields = Fields{}
	}
	fillFieldsFromContext(ctx, fields)

	var zapFields = make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}

	return zapFields
}

func fillFieldsFromContext(ctx context.Context, fields Fields) {
	if requestID, ok := ctx.Value("X-Request-ID").(string); ok {
		fields["requestID"] = requestID
	}
	if url, ok := ctx.Value("url").(string); ok {
		fields["url"] = url
	}
}
