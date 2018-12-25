package message

import (
	"github.com/dmibod/kanban/shared/tools/circuit"
	"github.com/dmibod/kanban/shared/tools/circuit/hystrix"
	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/msg/nats"
	natz "github.com/nats-io/go-nats"
	"time"
)

type serviceWithCircuitBreaker struct {
	executor nats.OperationExecutor
	breaker  interface {
		Execute(circuit.Handler) error
	}
}

// CreateService creates database service with circuit breaker
func CreateService(l logger.Logger) nats.OperationExecutor {
	e := nats.CreateExecutor(
		nats.WithLogger(l),
		nats.WithReconnectDelay(time.Second),
		nats.WithName("KANBAN"),
		nats.WithReconnectHandler(func(c *natz.Conn) { l.Debugln("reconnect...") }),
		nats.WithDisconnectHandler(func(c *natz.Conn) { l.Debugln("disconnect...") }),
		nats.WithCloseHandler(func(c *natz.Conn) { l.Debugln("close...") }))

	return &serviceWithCircuitBreaker{
		executor: e,
		breaker:  hystrix.New(hystrix.WithLogger(l), hystrix.WithName("NATS"), hystrix.WithTimeout(100)),
	}
}

// Execute executes database service operation within circuit breaker
func (s *serviceWithCircuitBreaker) Execute(ctx *nats.OperationContext, op nats.Operation) error {
	return s.breaker.Execute(func() error {
		return s.executor.Execute(ctx, op)
	})
}
