package nats

import (
	"sync"
	"time"

	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/logger/noop"
	"github.com/nats-io/go-nats"
)

const (
	defaultURL            = nats.DefaultURL
	defaultReconnectDelay = time.Second
)

// OperationExecutor executes operation
type OperationExecutor interface {
	Execute(*OperationContext, Operation) error
}

type executor struct {
	sync.Mutex
	url    string
	opts   []nats.Option
	conn   *nats.Conn
	logger logger.Logger
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

	url := o.url
	if url == "" {
		url = defaultURL
	}

	return &executor{
		logger: l,
		url:    url,
		opts:   o.opts,
	}
}

// Execute operation
func (e *executor) Execute(c *OperationContext, o Operation) error {
	err := e.ensureConnection(c)
	if err != nil {
		e.logger.Errorln("cannot open connection")
		return err
	}

	err = o(c.ctx, e.conn)
	if err != nil {
		e.dropDeadConnection()
	}

	return err
}

func (e *executor) ensureConnection(ctx *OperationContext) error {
	e.Lock()
	defer e.Unlock()

	if e.conn == nil {
		conn, err := nats.Connect(e.url, e.opts...)
		if err != nil {
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
