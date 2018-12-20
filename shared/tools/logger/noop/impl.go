package noop

import (
	api "github.com/dmibod/kanban/shared/tools/logger"
)

var _ api.Logger = (*Logger)(nil)

// Logger defines logger instance
type Logger struct {
}

func (l *Logger) Info(v ...interface{}) {
}

func (l *Logger) Debug(v ...interface{}) {
}

func (l *Logger) Error(v ...interface{}) {
}

func (l *Logger) Infoln(v ...interface{}) {
}

func (l *Logger) Debugln(v ...interface{}) {
}

func (l *Logger) Errorln(v ...interface{}) {
}

func (l *Logger) Infof(f string, v ...interface{}) {
}

func (l *Logger) Debugf(f string, v ...interface{}) {
}

func (l *Logger) Errorf(f string, v ...interface{}) {
}

