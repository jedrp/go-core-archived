package pllog

import (
	"log"

	"github.com/sirupsen/logrus"
)

type DefaultLogger struct {
	Level logrus.Level
}

func NewDefaultLogger(l logrus.Level) *DefaultLogger {
	return &DefaultLogger{
		Level: l,
	}
}

func (l *DefaultLogger) Trace(args ...interface{}) {
	l.Log(logrus.TraceLevel, args...)
}
func (l *DefaultLogger) Debug(args ...interface{}) {
	l.Log(logrus.DebugLevel, args...)
}
func (l *DefaultLogger) Info(args ...interface{}) {
	l.Log(logrus.InfoLevel, args...)
}
func (l *DefaultLogger) Warn(args ...interface{}) {
	l.Log(logrus.WarnLevel, args...)
}
func (l *DefaultLogger) Error(args ...interface{}) {
	l.Log(logrus.ErrorLevel, args...)
}
func (l *DefaultLogger) Fatal(args ...interface{}) {
	l.Log(logrus.FatalLevel, args...)
}
func (l *DefaultLogger) Panic(args ...interface{}) {
	l.Log(logrus.PanicLevel, args...)
}
func (l *DefaultLogger) Tracef(format string, args ...interface{}) {
	l.Logf(logrus.TraceLevel, format, args...)
}
func (l *DefaultLogger) Debugf(format string, args ...interface{}) {
	l.Logf(logrus.DebugLevel, format, args...)
}
func (l *DefaultLogger) Infof(format string, args ...interface{}) {
	l.Logf(logrus.InfoLevel, format, args...)
}
func (l *DefaultLogger) Warnf(format string, args ...interface{}) {
	l.Logf(logrus.WarnLevel, format, args...)
}
func (l *DefaultLogger) Errorf(format string, args ...interface{}) {
	l.Logf(logrus.ErrorLevel, format, args...)
}
func (l *DefaultLogger) Fatalf(format string, args ...interface{}) {
	l.Logf(logrus.FatalLevel, format, args...)
}
func (l *DefaultLogger) Panicf(format string, args ...interface{}) {
	l.Logf(logrus.PanicLevel, format, args...)
}
func (l *DefaultLogger) WithFields(fields map[string]interface{}) PlLogentry {
	return l
}

func (l *DefaultLogger) Logf(level logrus.Level, format string, args ...interface{}) {
	if l.IsLevelEnabled(level) {
		log.Printf(format, args...)
		if level <= logrus.PanicLevel {
			panic(l)
		}
	}
}

func (l *DefaultLogger) Log(level logrus.Level, args ...interface{}) {
	if l.IsLevelEnabled(level) {
		log.Println(args...)
		if level <= logrus.PanicLevel {
			panic(l)
		}
	}
}

func (l *DefaultLogger) IsLevelEnabled(level logrus.Level) bool {
	return l.Level >= level
}
