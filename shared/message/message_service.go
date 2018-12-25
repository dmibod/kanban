package message

import (
	"github.com/dmibod/kanban/shared/tools/circuit"
	"github.com/dmibod/kanban/shared/tools/circuit/hystrix"
	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/msg/nats"
	natz "github.com/nats-io/go-nats"
	"github.com/nats-io/go-nats-streaming"
	"time"
)

type serviceWithCircuitBreaker struct {
	executor nats.OperationExecutor
	breaker  interface {
		Execute(circuit.Handler) error
	}
}

// CreateService creates message service with circuit breaker
func CreateService(l logger.Logger) nats.OperationExecutor {
	e := nats.CreateExecutor(
		nats.WithLogger(l),
		nats.WithReconnectDelay(time.Second),
		nats.WithName("KANBAN"),
		nats.WithClientID("KANBAN_CLIENT"),
		nats.WithConnectionLostHandler(func(c stan.Conn, reason error) { l.Debugf("connection lost, reason %v...", reason) }),
		nats.WithReconnectHandler(func(c *natz.Conn) { l.Debugln("reconnect...") }),
		nats.WithDisconnectHandler(func(c *natz.Conn) { l.Debugln("disconnect...") }),
		nats.WithCloseHandler(func(c *natz.Conn) { l.Debugln("close...") }))

	return &serviceWithCircuitBreaker{
		executor: e,
		breaker:  hystrix.New(hystrix.WithLogger(l), hystrix.WithName("NATS"), hystrix.WithTimeout(1000)),
	}
}

// Execute executes message service operation within circuit breaker
func (s *serviceWithCircuitBreaker) Execute(ctx *nats.OperationContext, op nats.Operation) error {
	return s.breaker.Execute(func() error {
		return s.executor.Execute(ctx, op)
	})
}

// Status signalling connection up/down transitions
func (s *serviceWithCircuitBreaker) Status() <-chan bool {
	return s.executor.Status()
}
