package persistence

import (
	"github.com/dmibod/kanban/shared/tools/circuit"
	"github.com/dmibod/kanban/shared/tools/circuit/hystrix"
	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"github.com/dmibod/kanban/shared/tools/log"
)

// DatabaseServiceWithCircuitBreaker declares DatabaseService with CircuitBreaker
type databaseServiceWithCircuitBreaker struct {
	executor mongo.DatabaseCommandExecutor
	breaker  interface {
		Execute(circuit.Handler) error
	}
}

// CreateDatabaseService creates DatabaseService with CircuitBreaker
func CreateDatabaseService(l log.Logger) mongo.DatabaseCommandExecutor {
	return &databaseServiceWithCircuitBreaker{
		executor: mongo.CreateDatabaseService(l),
		breaker:  hystrix.New(hystrix.WithLogger(l), hystrix.WithName("MONGO"), hystrix.WithTimeout(100)),
	}
}

// Exec executes DatabaseService with CircuitBreaker
func (s *databaseServiceWithCircuitBreaker) Exec(c *mongo.DatabaseCommand, h mongo.DatabaseCommandHandler) error {
	return s.breaker.Execute(func() error {
		return s.executor.Exec(c, h)
	})
}
