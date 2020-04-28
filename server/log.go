package server

import (
	log "github.com/jarod2011/go-log"
)

type Logger interface {
	Log(v ...interface{})
	Logf(format string, v ...interface{})
	EnableDebug()
}

const debugPrefix = "[DEBUG]"
const infoPrefix = "[INFOS]"
const warnPrefix = "[WARNN]"
const errPrefix = "[ERROR]"
const fatalPrefix = "[FATAL]"

type defaultLogger struct{}

var logMap = map[string]log.Logger{
	debugPrefix: log.D(),
	infoPrefix:  log.I(),
	warnPrefix:  log.W(),
	errPrefix:   log.E(),
	fatalPrefix: log.F(),
}

func getLogger(prefix string) log.Logger {
	if logger, ok := logMap[prefix]; ok {
		return logger
	}
	return log.I()
}

func (d defaultLogger) Log(v ...interface{}) {
	if len(v) > 0 {
		if m, ok := v[0].(string); ok {
			getLogger(m).Log(v[1:]...)
			return
		}
	}
	log.I().Log(v...)
}

func (d defaultLogger) Logf(format string, v ...interface{}) {
	if len(format) > 7 {
		getLogger(format[0:7]).Logf(format[7:], v...)
		return
	}
	getLogger("").Logf(format, v...)
}

func (d defaultLogger) EnableDebug() {
	log.SetLevel(log.Debug)
}
