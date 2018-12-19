package hystrix

import (
	"github.com/dmibod/kanban/shared/tools/log"
	"github.com/dmibod/kanban/shared/tools/cb"
)

// Options declares circuit breaker options
type Options struct {
	logger log.Logger
	subject circuit.Subject
}

// Option is a closure which should initialize specific Options properties
type Option func(*Options)

// WithLogger initializes logger option
func WithLogger(l log.Logger) Option {
	return func(o *Options) {
		o.logger = l
	}
}

// WithSubject initializes subject option
func WithSubject(s circuit.Subject) Option {
	return func(o *Options) {
		o.subject = s
	}
}
