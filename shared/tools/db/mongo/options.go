package mongo

import "github.com/dmibod/kanban/shared/tools/log"

// Options declares repository factory options
type Options struct {
	executor OperationExecutor
	logger   log.Logger
	db       string
}

// Option is a closure which should initialize specific Options properties
type Option func(*Options)

// WithLogger initializes logger option
func WithLogger(l log.Logger) Option {
	return func(o *Options) {
		o.logger = l
	}
}

// WithDatabase initializes db option
func WithDatabase(db string) Option {
	return func(o *Options) {
		o.db = db
	}
}

// WithExecutor initializes executor option
func WithExecutor(e OperationExecutor) Option {
	return func(o *Options) {
		o.executor = e
	}
}
