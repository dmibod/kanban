package persistence

import (
	"github.com/dmibod/kanban/shared/tools/circuit"
	"github.com/dmibod/kanban/shared/tools/circuit/hystrix"
	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"github.com/dmibod/kanban/shared/tools/logger"
)

type serviceWithCircuitBreaker struct {
	executor mongo.OperationExecutor
	breaker  interface {
		Execute(circuit.Handler) error
	}
}

// CreateService creates database service with circuit breaker
func CreateService(l logger.Logger) (mongo.OperationExecutor, mongo.SessionProvider) {
	e, p := mongo.CreateExecutor(mongo.WithLogger(l))
	return &serviceWithCircuitBreaker{
		executor: e,
		breaker:  hystrix.New(hystrix.WithLogger(l), hystrix.WithName("MONGO"), hystrix.WithTimeout(100)),
	}, p
}

// Execute executes database service operation within circuit breaker
func (s *serviceWithCircuitBreaker) Execute(ctx *mongo.OperationContext, op mongo.Operation) error {
	return s.breaker.Execute(func() error {
		return s.executor.Execute(ctx, op)
	})
}
