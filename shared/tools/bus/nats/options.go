package nats

import (
	"context"
	"time"

	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/nats-io/go-nats"
)

type options struct {
	url      string
	natsOpts []nats.Option
	ctx      context.Context
	logger.Logger
}

// Option initializes Options properties
type Option func(*options)

// WithContext initializes ctx option
func WithContext(c context.Context) Option {
	return func(o *options) {
		o.ctx = c
	}
}

// WithName initializes name option
func WithName(n string) Option {
	return func(o *options) {
		o.natsOpts = append(o.natsOpts, nats.Name(n))
	}
}

// WithURL initializes url option
func WithURL(u string) Option {
	return func(o *options) {
		o.url = u
	}
}

// WithReconnectDelay initializes reconnectDelay option
func WithReconnectDelay(t time.Duration) Option {
	return func(o *options) {
		o.natsOpts = append(o.natsOpts, nats.ReconnectWait(t))
	}
}

// WithDisconnectHandler initializes disconnectHandler option
func WithDisconnectHandler(h nats.ConnHandler) Option {
	return func(o *options) {
		o.natsOpts = append(o.natsOpts, nats.DisconnectHandler(h))
	}
}

// WithReconnectHandler initializes reconnectHandler option
func WithReconnectHandler(h nats.ConnHandler) Option {
	return func(o *options) {
		o.natsOpts = append(o.natsOpts, nats.ReconnectHandler(h))
	}
}

// WithCloseHandler initializes closeHandler option
func WithCloseHandler(h nats.ConnHandler) Option {
	return func(o *options) {
		o.natsOpts = append(o.natsOpts, nats.ClosedHandler(h))
	}
}

// WithLogger initializes logger option
func WithLogger(l logger.Logger) Option {
	return func(o *options) {
		o.Logger = l
	}
}
