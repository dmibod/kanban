package persistence

import (
	"github.com/dmibod/kanban/shared/tools/circuit"
	"github.com/dmibod/kanban/shared/tools/circuit/hystrix"
	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"github.com/dmibod/kanban/shared/tools/logger"
)

type executorWithCircuitBreaker struct {
	executor mongo.OperationExecutor
	breaker  interface {
		Execute(circuit.Handler) error
	}
}

// CreateExecutor with circuit breaker
func CreateExecutor(p mongo.SessionProvider, l logger.Logger) mongo.OperationExecutor {
	return &executorWithCircuitBreaker{
		executor: mongo.CreateExecutor(p, l),
		breaker:  hystrix.New(hystrix.WithLogger(l), hystrix.WithName("MONGO"), hystrix.WithTimeout(100)),
	}
}

// Execute executes database service operation within circuit breaker
func (e *executorWithCircuitBreaker) Execute(ctx *mongo.OperationContext, op mongo.Operation) error {
	return e.breaker.Execute(func() error {
		return e.executor.Execute(ctx, op)
	})
}
