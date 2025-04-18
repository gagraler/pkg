package log

import (
	"context"
	"go.uber.org/zap"
)

func Info(args ...interface{}) {
	sugar.Info(args...)
}

func Infof(format string, args ...interface{}) {
	sugar.Infof(format, args...)
}

func WithContext(ctx context.Context) *zap.SugaredLogger {
	return sugar.With(ctx)
}

func Debug(args ...interface{}) {
	sugar.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	sugar.Debugf(format, args...)
}

func Warn(args ...interface{}) {
	sugar.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	sugar.Warnf(format, args...)
}

func Error(args ...interface{}) {
	sugar.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	sugar.Errorf(format, args...)
}

func Fatal(args ...interface{}) {
	sugar.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	sugar.Fatalf(format, args...)
}
