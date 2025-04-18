package log

import (
	"fmt"

	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const filename string = "log.LOG"

// getFileLogWriter returns the WriteSyncer for logging to a file.
func getFileLogWriter(config *Config) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/%s", config.Path, filename),
		MaxSize:    config.RotateSize,
		MaxBackups: config.RotateNum,
		MaxAge:     config.KeepHours,
		Compress:   true,
	}
	return zapcore.AddSync(lumberJackLogger)
}
