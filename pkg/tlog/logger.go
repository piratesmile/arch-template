package tlog

import (
	"context"
)

type Fields map[string]interface{}

type Logger interface {
	Debug(ctx context.Context, msg string, fields Fields)
	Info(ctx context.Context, msg string, fields Fields)
	Warn(ctx context.Context, msg string, fields Fields)
	Error(ctx context.Context, msg string, fields Fields)
	Panic(ctx context.Context, msg string, fields Fields)
	Fatal(ctx context.Context, msg string, fields Fields)
	Sync() error
}

func Global() Logger {
	return std
}

func Debug(ctx context.Context, msg string, fields Fields) {
	if std != nil {
		std.Debug(ctx, msg, fields)
	}
}

func Info(ctx context.Context, msg string, fields Fields) {
	if std != nil {
		std.Info(ctx, msg, fields)
	}
}

func Warn(ctx context.Context, msg string, fields Fields) {
	if std != nil {
		std.Warn(ctx, msg, fields)
	}
}

func Error(ctx context.Context, msg string, fields Fields) {
	if std != nil {
		std.Error(ctx, msg, fields)
	}
}

func Panic(ctx context.Context, msg string, fields Fields) {
	if std != nil {
		std.Panic(ctx, msg, fields)
	}
}

func Fatal(ctx context.Context, msg string, fields Fields) {
	if std != nil {
		std.Fatal(ctx, msg, fields)
	}
}

func Sync() error {
	if std != nil {
		return std.Sync()
	}
	return nil
}
