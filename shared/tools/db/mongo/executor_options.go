package mongo

import (
	"github.com/dmibod/kanban/shared/tools/logger"
	"time"
)

type options struct {
	url      string
	timeout  time.Duration
	authdb   string
	user     string
	password string
	logger.Logger
}

// Option initializes Options properties
type Option func(*options)

// WithURL initializes url option
func WithURL(u string) Option {
	return func(o *options) {
		o.url = u
	}
}

// WithTimeout initializes timeout
func WithTimeout(t time.Duration) Option {
	return func(o *options) {
		o.timeout = t
	}
}

// WithAuthDb initializes authdb option
func WithAuthDb(db string) Option {
	return func(o *options) {
		o.authdb = db
	}
}

// WithUser initializes user option
func WithUser(u string) Option {
	return func(o *options) {
		o.user = u
	}
}

// WithPassword initializes password option
func WithPassword(p string) Option {
	return func(o *options) {
		o.password = p
	}
}

// WithLogger initializes logger option
func WithLogger(l logger.Logger) Option {
	return func(o *options) {
		o.Logger = l
	}
}
