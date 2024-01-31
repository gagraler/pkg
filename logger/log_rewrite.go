package logger

/**
 * @author: x.gallagher.anderson@gmail.com
 * @time: 2024/1/31 22:24
 * @file: log_rewrite.go
 * @description: 重写zap日志方法
 */

func Debug(msg string) {
	sugarLogger.Debug(msg)
}

func Debugf(msg string, args ...interface{}) {
	sugarLogger.Debugf(msg, args...)
}

func Info(msg string) {
	sugarLogger.Info(msg)
}

func Infof(msg string, args ...interface{}) {
	sugarLogger.Infof(msg, args...)
}

func Warn(msg string) {
	sugarLogger.Warn(msg)
}

func Warnf(msg string, args ...interface{}) {
	sugarLogger.Warnf(msg, args...)
}

func Error(msg string) {
	sugarLogger.Error(msg)
}

func Errorf(msg string, args ...interface{}) {
	sugarLogger.Errorf(msg, args...)
}

func Fatal(msg string) {
	sugarLogger.Fatal(msg)
}

func Fatalf(msg string, args ...interface{}) {
	sugarLogger.Fatalf(msg, args...)
}

func Panic(msg string) {
	sugarLogger.Panic(msg)
}

func Panicf(msg string, args ...interface{}) {
	sugarLogger.Panicf(msg, args...)
}
