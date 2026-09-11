package logger

import (
	"golf-booking-go/pkg/setting"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LoggerZap struct {
	*zap.Logger
}

func NewLogger(config setting.LoggerSetting) *LoggerZap {
	var level zapcore.Level

	switch config.LogLevel {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	case "fatal":
		level = zapcore.FatalLevel
	case "panic":
		level = zapcore.PanicLevel
	default:
		level = zapcore.InfoLevel
	}

	hook := &lumberjack.Logger{
		Filename:   config.Filename,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	}

	consoleCore := zapcore.NewCore(
		getConsoleEncoder(),
		zapcore.AddSync(os.Stdout),
		level,
	)

	fileCore := zapcore.NewCore(
		getJSONEncoder(),
		zapcore.AddSync(hook),
		level,
	)

	core := zapcore.NewTee(consoleCore, fileCore)

	logger := zap.New(
		core,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	return &LoggerZap{
		Logger: logger,
	}
}

func getConsoleEncoder() zapcore.Encoder {
	config := zap.NewDevelopmentEncoderConfig()

	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncodeCaller = zapcore.ShortCallerEncoder

	return zapcore.NewConsoleEncoder(config)
}

func getJSONEncoder() zapcore.Encoder {
	config := zap.NewProductionEncoderConfig()

	config.TimeKey = "time"
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncodeLevel = zapcore.CapitalLevelEncoder
	config.EncodeCaller = zapcore.ShortCallerEncoder

	return zapcore.NewJSONEncoder(config)
}
