package tlog

import (
	"os"

	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var defaultOptions = &Options{
	Level:    "info",
	LogFile:  "stdout",
	ErrFile:  "stderr",
	FileSize: 100, // 100M
	FileAge:  0,   // forever
}

type Options struct {
	Level    string
	ErrFile  string
	LogFile  string
	FileSize int
	FileAge  int
}

func (o *Options) newWriteSyncer() (zapcore.WriteSyncer, zapcore.WriteSyncer) {
	var (
		writer    zapcore.WriteSyncer = os.Stdout
		errWriter zapcore.WriteSyncer = os.Stderr
	)
	if o.LogFile != "" {
		writer = o.open(o.LogFile)
	}
	if o.ErrFile != "" {
		errWriter = o.open(o.ErrFile)
	}
	return writer, errWriter
}

func (o *Options) open(path string) zapcore.WriteSyncer {
	switch path {
	case "stdout":
		return os.Stdout
	case "stderr":
		return os.Stderr
	default:
		return zapcore.AddSync(&lumberjack.Logger{
			Filename:  path,
			MaxSize:   o.FileSize,
			LocalTime: true,
			MaxAge:    o.FileAge,
		})
	}
}
