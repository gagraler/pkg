package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strings"
)

/**
 * @author: x.gallagher.anderson@gmail.com
 * @time: 2024/1/31 22:24
 * @file: log.go
 * @description: 基于zap sugared logger的日志库
 */

type Logger struct {
	l     *zap.SugaredLogger
	level zapcore.Level
}

type Field = zap.Field

type Config struct {
	OutPut       string
	FileName     string
	Level        string
	RotateSize   int
	RotateNum    int
	RotateAgeNum int
}

var sugarLogger *zap.SugaredLogger

func NewLogger(l *Config) *Logger {
	zapCfg := zap.NewDevelopmentConfig()

	zapCfg.Level = zap.NewAtomicLevelAt(LevelFromString(l.Level))
	zapCfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
	zapCfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	var core []zapcore.Core

	// Set up the non-colored level encoder for file output.
	consoleEncoderConfig := zapCfg.EncoderConfig
	consoleEncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder // Standard non-colored level encoder.
	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(consoleEncoderConfig),
		zapcore.AddSync(zapcore.Lock(os.Stdout)),
		zapCfg.Level,
	)
	core = append(core, consoleCore)

	if l.OutPut != "stdout" {
		fileEncoderConfig := zapCfg.EncoderConfig
		fileEncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder // Standard non-colored level encoder.

		fileWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   l.FileName,
			MaxSize:    l.RotateSize,
			MaxBackups: l.RotateNum,
			MaxAge:     l.RotateAgeNum,
			Compress:   true,
		})

		fileCore := zapcore.NewCore(
			zapcore.NewConsoleEncoder(fileEncoderConfig),
			fileWriter,
			zapCfg.Level,
		)

		core = append(core, fileCore)
	}

	multiCore := zapcore.NewTee(core...)

	sugarLogger = zap.New(multiCore, zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()

	return &Logger{
		l:     sugarLogger,
		level: zapCfg.Level.Level(),
	}
}

func LevelFromString(levelStr string) zapcore.Level {
	switch strings.ToLower(levelStr) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	case "panic":
		return zapcore.PanicLevel
	default:
		return zapcore.DebugLevel
	}
}
