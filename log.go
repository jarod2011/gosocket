package gosocket

import (
	log "github.com/jarod2011/go-log"
)

type Logger interface {
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Log(v ...interface{})
	Logf(format string, v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
}

type logger struct{}

func (l *logger) Debug(v ...interface{}) {
	log.D().Log(v...)
}

func (l *logger) Debugf(format string, v ...interface{}) {
	log.D().Logf(format, v...)
}

func (l *logger) Log(v ...interface{}) {
	log.I().Log(v...)
}

func (l *logger) Logf(format string, v ...interface{}) {
	log.I().Logf(format, v...)
}

func (l *logger) Warn(v ...interface{}) {
	log.W().Log(v...)
}

func (l *logger) Warnf(format string, v ...interface{}) {
	log.W().Logf(format, v...)
}

func (l *logger) Error(v ...interface{}) {
	log.E().Log(v...)
}

func (l *logger) Errorf(format string, v ...interface{}) {
	log.E().Logf(format, v...)
}

func (l *logger) Fatal(v ...interface{}) {
	log.F().Log(v...)
}

func (l *logger) Fatalf(format string, v ...interface{}) {
	log.F().Logf(format, v...)
}

func newLogger() Logger {
	return &logger{}
}
