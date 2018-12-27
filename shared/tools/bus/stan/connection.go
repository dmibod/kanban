package stan

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/dmibod/kanban/shared/tools/logger/noop"

	"github.com/dmibod/kanban/shared/tools/bus"
	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/nats-io/go-nats"
	"github.com/nats-io/go-nats-streaming"
)

// Errors
var (
	ErrNotConnected = errors.New("stan: not connected")
)

const (
	defaultClusterID      = "test-cluster"
	defaultReconnectDelay = time.Second
)

var _ bus.Connection = (*Connection)(nil)

// Connection interface
type Connection struct {
	mu        sync.Mutex
	ctx       context.Context
	status    chan bool
	close     chan struct{}
	url       string
	clusterID string
	clientID  string
	logger    logger.Logger
	natsConn  *nats.Conn
	stanConn  stan.Conn
	stanOpts  []stan.Option
	natsOpts  []nats.Option
	subOpts   []stan.SubscriptionOption
}

// CreateConnection creates new connection
func CreateConnection(opts ...Option) *Connection {
	var o options

	o.natsOpts = []nats.Option{}
	o.stanOpts = []stan.Option{}

	for _, opt := range opts {
		opt(&o)
	}

	url := o.url
	if url == "" {
		url = stan.DefaultNatsURL
	}

	clusterID := o.clusterID
	if clusterID == "" {
		clusterID = defaultClusterID
	}

	l := o.Logger
	if l == nil {
		l = &noop.Logger{}
	}

	ctx := o.ctx
	if ctx == nil {
		ctx = context.TODO()
	}

	var conn *Connection

	o.natsOpts = append(o.natsOpts, nats.DisconnectHandler(func(nc *nats.Conn) { conn.status <- false }))
	o.natsOpts = append(o.natsOpts, nats.ReconnectHandler(func(nc *nats.Conn) { conn.status <- true }))
	o.stanOpts = append(o.stanOpts, stan.NatsURL(url))

	subOpts := []stan.SubscriptionOption{stan.StartWithLastReceived()}

	if o.durable != "" {
		subOpts = append(subOpts, stan.DurableName(o.durable))
	}

	conn = &Connection{
		logger:    l,
		url:       url,
		ctx:       ctx,
		status:    make(chan bool, 1),
		close:     make(chan struct{}, 1),
		clusterID: clusterID,
		clientID:  o.clientID,
		stanOpts:  o.stanOpts,
		natsOpts:  o.natsOpts,
		subOpts:   subOpts,
	}

	return conn
}

// Connect to broker
func (c *Connection) Connect() <-chan bool {
	if !c.IsConnected() {
		go c.connect()
	}
	return c.status
}

// IsConnected status
func (c *Connection) IsConnected() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.stanConn != nil && c.stanConn.NatsConn().IsConnected()
}

// Publish message
func (c *Connection) Publish(topic string, message []byte) error {
	if !c.IsConnected() {
		return ErrNotConnected
	}
	return c.stanConn.Publish(topic, message)
}

// Subscribe for messages
func (c *Connection) Subscribe(topic string, queue string, handler bus.MessageHandler) (interface{}, error) {
	if !c.IsConnected() {
		return nil, ErrNotConnected
	}
	return c.stanConn.QueueSubscribe(topic, queue, func(msg *stan.Msg) {
		handler.Handle(msg.Data)
	}, c.subOpts...)
}

// Unsubscribe subscription
func (c *Connection) Unsubscribe(handle interface{}) error {
	if !c.IsConnected() {
		return ErrNotConnected
	}
	if s, ok := handle.(stan.Subscription); ok {
		return s.Unsubscribe()
	}
	return nil
}

// Close connection
func (c *Connection) Close() <-chan struct{} {
	return c.close
}

func (c *Connection) connect() {
	timer := time.NewTicker(time.Second * 1)
	for {
		select {
		case <-c.ctx.Done():
			c.logger.Debugln("close signal")
			c.mu.Lock()
			c.disconnect()
			c.mu.Unlock()
			c.logger.Debugln("closed")
			return
		case <-timer.C:
			c.logger.Debugln("connect signal")
			c.mu.Lock()
			c.connectNats()
			if c.natsConn != nil {
				c.connectStan()
				if c.stanConn != nil {
					c.logger.Debugln("send up signal")
					c.status <- true
					timer.Stop()
				}
			}
			c.mu.Unlock()
		}
	}
}

func (c *Connection) connectNats() {
	if c.natsConn == nil {
		c.logger.Debugln("connect nats")
		natsConn, err := nats.Connect(c.url, c.natsOpts...)
		if err != nil {
			c.logger.Errorln(err)
			return
		}
		c.logger.Debugln("nats connected")
		c.natsConn = natsConn
		c.stanOpts = append(c.stanOpts, stan.NatsConn(natsConn))
	}
}

func (c *Connection) connectStan() {
	if c.stanConn == nil {
		c.logger.Debugln("connect stan")
		stanConn, err := stan.Connect(c.clusterID, c.clientID, c.stanOpts...)
		if err != nil {
			c.logger.Errorln(err)
			return
		}
		c.logger.Debugln("stan connected")
		c.stanConn = stanConn
	}
}

func (c *Connection) disconnect() {
	if c.stanConn != nil {
		c.logger.Debugln("close stan connection")
		err := c.stanConn.Close()
		c.stanConn = nil
		if err != nil {
			c.logger.Errorln(err)
		}
	}
	if c.natsConn != nil {
		c.logger.Debugln("close nats connection")
		c.natsConn.Close()
		c.natsConn = nil
	}
	c.status <- false
	c.close <- struct{}{}
}
