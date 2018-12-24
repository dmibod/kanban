package console

import (
	"log"
	"os"

	api "github.com/dmibod/kanban/shared/tools/logger"
)

var _ api.Logger = (*Logger)(nil)

// Logger defines logger instance
type Logger struct {
	debug bool
	out   *log.Logger
	err   *log.Logger
}

// New creates Logger
func New(opts ...Option) *Logger {
	var options Options
	for _, o := range opts {
		o(&options)
	}

	return &Logger{
		debug: options.Debug,
		out:   log.New(os.Stdout, options.Prefix, log.LstdFlags),
		err:   log.New(os.Stderr, options.Prefix, log.LstdFlags),
	}
}

func (l *Logger) Info(v ...interface{}) {
	l.out.Print(v...)
}

func (l *Logger) Debug(v ...interface{}) {
	if l.debug {
		l.out.Print(v...)
	}
}

func (l *Logger) Error(v ...interface{}) {
	l.err.Print(v...)
}

func (l *Logger) Infoln(v ...interface{}) {
	l.out.Println(v...)
}

func (l *Logger) Debugln(v ...interface{}) {
	if l.debug {
		l.out.Println(v...)
	}
}

func (l *Logger) Errorln(v ...interface{}) {
	l.err.Println(v...)
}

func (l *Logger) Infof(f string, v ...interface{}) {
	l.out.Printf(f, v...)
}

func (l *Logger) Debugf(f string, v ...interface{}) {
	if l.debug {
		l.out.Printf(f, v...)
	}
}

func (l *Logger) Errorf(f string, v ...interface{}) {
	l.err.Printf(f, v...)
}
