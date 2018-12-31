package hystrix

import (
	"github.com/dmibod/kanban/shared/tools/logger"
)

// Options declares circuit breaker options
type Options struct {
	logger.Logger
	name    string
	timeout int
}

// Option is a closure which should initialize specific Options properties
type Option func(*Options)

// WithLogger initializes logger option
func WithLogger(l logger.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}

// WithName initializes name option
func WithName(n string) Option {
	return func(o *Options) {
		o.name = n
	}
}

// WithTimeout initializes timeout option
func WithTimeout(t int) Option {
	return func(o *Options) {
		o.timeout = t
	}
}
