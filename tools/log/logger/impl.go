package logger

import (
	"log"

	api "github.com/dmibod/kanban/tools/log"
)

var _ api.Logger = (*Logger)(nil)

// Logger defines logger instance
type Logger struct {
	module string
}

// New creates Logger
func New(opts ...Option) *Logger {
	var options Options
	for _, o := range opts {
		o(&options)
	}

	return &Logger{
		module: options.Module,
	}
}

func (l *Logger) Info(m string) {
	log.Printf("[%v] %v\n", l.module, m)
}

func (l *Logger) Debug(m string) {
	log.Printf("[%v] %v\n", l.module, m)
}

func (l *Logger) Error(m string) {
	log.Printf("[%v] %v\n", l.module, m)
}
