package stan

import (
	"time"

	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/nats-io/go-nats"
	"github.com/nats-io/go-nats-streaming"
)

type options struct {
	url       string
	clusterID string
	clientID  string
	durable   string
	stanOpts  []stan.Option
	natsOpts  []nats.Option
	logger.Logger
}

// Option initializes Options properties
type Option func(*options)

// WithName initializes name option
func WithName(n string) Option {
	return func(o *options) {
		o.natsOpts = append(o.natsOpts, nats.Name(n))
	}
}

// WithDurable initializes durable option
func WithDurable(n string) Option {
	return func(o *options) {
		o.durable = n
	}
}

// WithURL initializes url option
func WithURL(u string) Option {
	return func(o *options) {
		o.url = u
	}
}

// WithClusterID initializes clusterID option
func WithClusterID(id string) Option {
	return func(o *options) {
		o.clusterID = id
	}
}

// WithClientID initializes clientID option
func WithClientID(id string) Option {
	return func(o *options) {
		o.clientID = id
	}
}

// WithConnectionLostHandler initializes connectionLostHandler option
func WithConnectionLostHandler(h stan.ConnectionLostHandler) Option {
	return func(o *options) {
		o.stanOpts = append(o.stanOpts, stan.SetConnectionLostHandler(h))
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
