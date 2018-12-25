package nats

import (
	"sync"
	"time"

	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/logger/noop"
	"github.com/nats-io/go-nats"
	"github.com/nats-io/go-nats-streaming"
)

const (
	defaultClusterID      = "test-cluster"
	defaultReconnectDelay = time.Second
)

// OperationExecutor executes operation
type OperationExecutor interface {
	Execute(*OperationContext, Operation) error

	Notify(chan<- bool)
}

type executor struct {
	sync.Mutex
	url       string
	clusterID string
	clientID  string
	stanOpts  []stan.Option
	natsOpts  []nats.Option
	conn      Connection
	logger    logger.Logger
	listeners []chan<- bool
}

// CreateExecutor creates executor
func CreateExecutor(opts ...Option) OperationExecutor {
	var o options

	for _, opt := range opts {
		opt(&o)
	}

	l := o.logger
	if l == nil {
		l = &noop.Logger{}
	}

	clusterID := o.clusterID
	if clusterID == "" {
		clusterID = defaultClusterID
	}

	listeners := []chan<- bool{}

	exec := &executor{
		logger:    l,
		url:       o.url,
		clusterID: clusterID,
		clientID:  o.clientID,
		stanOpts:  o.stanOpts,
		natsOpts:  o.natsOpts,
		listeners: listeners,
	}

	o.natsOpts = append(o.natsOpts, nats.DisconnectHandler(func(nc *nats.Conn) { exec.notify(false) }))
	o.natsOpts = append(o.natsOpts, nats.ReconnectHandler(func(nc *nats.Conn) { exec.notify(true) }))

	return exec
}

// Execute operation
func (e *executor) Execute(c *OperationContext, o Operation) error {
	err := e.ensureConnection(c)
	if err != nil {
		return err
	}

	err = o(c.ctx, e.conn)
	if err != nil {
		e.dropDeadConnection()
	}

	return err
}

// Notify allows to subscribe for connection up/down transitions
func (e *executor) Notify(ch chan<- bool) {
	e.Lock()
	defer e.Unlock()

	e.listeners = append(e.listeners, ch)
}

func (e *executor) ensureConnection(ctx *OperationContext) error {
	e.Lock()
	defer e.Unlock()

	if e.conn == nil {
		conn, err := e.createConnection()
		if err != nil {
			e.logger.Errorln("cannot open connection")
			return err
		}

		e.logger.Debugln("new connection")
		e.conn = conn
	}

	return nil
}

func (e *executor) dropDeadConnection() {
	e.Lock()
	defer e.Unlock()

	if e.conn != nil {
		err := e.conn.Flush()
		if err == nil {
			e.logger.Debugln("flush ok")
			return
		}

		e.logger.Debugln("close connection")

		e.conn.Close()
		e.conn = nil
	}
}

func (e *executor) createConnection() (Connection, error) {
	url := e.url

	if url == "" {
		url = nats.DefaultURL
	}

	nc, err := CreateNatsConnection(url, e.natsOpts...)

	if e.clusterID == "" {
		return nc, err
	}

	if url == "" {
		url = stan.DefaultNatsURL
	}

	e.stanOpts = append(e.stanOpts, stan.NatsConn(nc.conn))

	return CreateStanConnection(url, e.clusterID, e.clientID, e.stanOpts...)
}

func (e *executor) notify(s bool) {
	if s {
		e.logger.Debugln("send recover signal")
	} else {
		e.logger.Debugln("send release signal")
	}
	for _, ch := range e.listeners {
		if len(ch) == 0 {
			ch <- s
		}
	}
	if s {
		e.logger.Debugln("recover signal sent")
	} else {
		e.logger.Debugln("release signal sent")
	}
}
