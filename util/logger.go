package util

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
)

var Log *Logger

type Logger struct {
	log   *logrus.Logger
	entry *logrus.Entry
}

var mapLogLevel = map[string]logrus.Level{
	"trace": logrus.TraceLevel,
	"debug": logrus.DebugLevel,
	"info":  logrus.InfoLevel,
	"warn":  logrus.WarnLevel,
	"error": logrus.ErrorLevel,
	"fatal": logrus.FatalLevel,
}

func NewLogger() *Logger {
	log := logrus.StandardLogger()

	level := os.Getenv("LOG_LEVEL")
	log.SetLevel(mapLogLevel[level])

	switch os.Getenv("LOG_FORMAT") {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{})
	default:
		logrus.SetFormatter(&logrus.TextFormatter{
			DisableColors:   false,
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}

	switch os.Getenv("LOG_OUTPUT") {
	default:
		log.Out = os.Stdout
	}

	return &Logger{
		log: log,
	}
}

func (l *Logger) WithContext(ctx context.Context) *Logger {
	if l.entry == nil {
		l.entry = l.log.WithContext(ctx)
		return l
	}
	l.entry = l.entry.WithContext(ctx)
	return l
}

func (l *Logger) WithField(key string, i interface{}) *Logger {
	if l.entry == nil {
		l.entry = l.log.WithField(key, i)
		return l
	}
	l.entry = l.entry.WithField(key, i)
	return l
}

func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	loggerFields := logrus.Fields{}
	for k, v := range fields {
		loggerFields[k] = v
	}
	if l.entry == nil {
		l.entry = l.log.WithFields(loggerFields)
		return l
	}
	l.entry = l.entry.WithFields(loggerFields)
	return l
}

func (l *Logger) WithError(err error) *Logger {
	if l.entry == nil {
		l.entry = l.log.WithError(err)
		return l
	}
	l.entry = l.entry.WithError(err)
	return l
}

func (l *Logger) Trace(args ...interface{}) *Logger {
	if l.entry == nil {
		l.log.Trace(args...)
		return l
	}
	l.entry.Trace(args...)
	return l
}

func (l *Logger) Tracef(format string, args ...interface{}) *Logger {
	if l.entry == nil {
		l.log.Tracef(format, args...)
		return l
	}
	l.entry.Tracef(format, args...)
	return l
}

func (l *Logger) Debug(args ...interface{}) *Logger {
	if l.entry == nil {
		l.log.Debug(args...)
		return l
	}
	l.entry.Debug(args...)
	return l
}

func (l *Logger) Debugf(format string, args ...interface{}) *Logger {
	if l.entry == nil {
		l.log.Debugf(format, args...)
		return l
	}
	l.entry.Debugf(format, args...)
	return l
}

func (l *Logger) Info(args ...interface{}) *Logger {
	if l.entry == nil {
		l.log.Info(args...)
		return l
	}
	l.entry.Info(args...)
	return l
}

func (l *Logger) Infof(format string, args ...interface{}) *Logger {
	if l.entry == nil {
		l.log.Infof(format, args...)
		return l
	}
	l.entry.Infof(format, args...)
	return l
}

func (l *Logger) Warn(args ...interface{}) *Logger {
	if l.entry == nil {
		l.log.Warn(args...)
		return l
	}
	l.entry.Warn(args...)
	return l
}

func (l *Logger) Warnf(format string, args ...interface{}) *Logger {
	if l.entry == nil {
		l.log.Warnf(format, args...)
		return l
	}
	l.entry.Warnf(format, args...)
	return l
}

func (l *Logger) Error(args ...interface{}) *Logger {
	if l.entry == nil {
		l.log.Error(args...)
		return l
	}
	l.entry.Error(args...)
	return l
}

func (l *Logger) Errorf(format string, args ...interface{}) *Logger {
	if l.entry == nil {
		l.log.Errorf(format, args...)
		return l
	}
	l.entry.Errorf(format, args...)
	return l
}

func (l *Logger) Fatal(args ...interface{}) *Logger {
	if l.entry == nil {
		l.log.Fatal(args...)
		return l
	}
	l.log.Fatal(args...)
	return l
}

func (l *Logger) Fatalf(format string, args ...interface{}) *Logger {
	if l.entry == nil {
		l.log.Fatalf(format, args...)
		return l
	}
	l.log.Fatalf(format, args...)
	return l
}

func (l *Logger) Panic(args ...interface{}) *Logger {
	if l.entry == nil {
		l.log.Panic(args...)
		return l
	}
	l.log.Panic(args...)
	return l
}

func (l *Logger) Panicf(format string, args ...interface{}) *Logger {
	if l.entry == nil {
		l.log.Panicf(format, args...)
		return l
	}
	l.log.Panicf(format, args...)
	return l
}
