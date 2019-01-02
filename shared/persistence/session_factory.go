package persistence

import (
	"github.com/dmibod/kanban/shared/tools/circuit"
	"github.com/dmibod/kanban/shared/tools/circuit/hystrix"
	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"github.com/dmibod/kanban/shared/tools/logger"
	"gopkg.in/mgo.v2"
)

type factoryWithCircuitBreaker struct {
	factory mongo.SessionFactory
	breaker interface {
		Execute(circuit.Handler) error
	}
}

// CreateSessionFactory with circuit breaker
func CreateSessionFactory(f mongo.SessionFactory, l logger.Logger) mongo.SessionFactory {
	return &factoryWithCircuitBreaker{
		factory: f,
		breaker: hystrix.New(hystrix.WithLogger(l), hystrix.WithName("MONGO"), hystrix.WithTimeout(100)),
	}
}

// Execute executes database service operation within circuit breaker
func (f *factoryWithCircuitBreaker) Session() (*mgo.Session, error) {
	var session *mgo.Session
	err := f.breaker.Execute(func() error {
		s, e := f.factory.Session()
		session = s
		return e
	})
	return session, err
}
