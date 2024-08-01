package logger

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap"
)

/**
 * @author: gagral.x@gmail.com
 * @time: 2024/7/17 21:23
 * @file: logger_kratos.go
 * @description: 实现kratos的log接口
 */

var _ log.Logger = (*Log)(nil)

type Log struct {
	log    *zap.Logger
	msgKey string
}

type Option func(*Log)

func WithMessageKey(msgKey string) Option {
	return func(l *Log) {
		l.msgKey = msgKey
	}
}

func NewLog(zaplog *zap.Logger) *Log {
	return &Log{
		log:    zaplog,
		msgKey: log.DefaultMessageKey,
	}
}

func (l *Log) Log(level log.Level, values ...interface{}) error {
	var (
		msg       = ""
		valuesLen = len(values)
	)
	if valuesLen == 0 || valuesLen%2 != 0 {
		l.log.Warn(fmt.Sprint("Key values must appear in pairs: ", values))
		return nil
	}

	data := make([]zap.Field, 0, (valuesLen/2)+1)
	for i := 0; i < valuesLen; i += 2 {
		if values[i].(string) == l.msgKey {
			msg, _ = values[i+1].(string)
			continue
		}
		data = append(data, zap.Any(fmt.Sprint(values[i]), values[i+1]))
	}

	switch level {
	case log.LevelDebug:
		l.log.Debug(msg, data...)
	case log.LevelInfo:
		l.log.Info(msg, data...)
	case log.LevelWarn:
		l.log.Warn(msg, data...)
	case log.LevelError:
		l.log.Error(msg, data...)
	case log.LevelFatal:
		l.log.Fatal(msg, data...)
	}
	return nil
}

func (l *Log) Sync() error {
	return l.log.Sync()
}

func (l *Log) Close() error {
	return l.Sync()
}
