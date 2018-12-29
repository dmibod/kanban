package bus

import (
	"github.com/dmibod/kanban/shared/tools/bus"
	"github.com/dmibod/kanban/shared/tools/circuit"
	"github.com/dmibod/kanban/shared/tools/circuit/hystrix"
	"github.com/dmibod/kanban/shared/tools/logger"
)

type transportWithCircuitBreaker struct {
	bus.Transport
	breaker interface {
		Execute(circuit.Handler) error
	}
}

// CreateTransport creates transport with circuit breaker
func CreateTransport(t bus.Transport, l logger.Logger) bus.Transport {
	return &transportWithCircuitBreaker{
		Transport: t,
		breaker: hystrix.New(
			hystrix.WithLogger(l),
			hystrix.WithName("BUS"),
			hystrix.WithTimeout(100)),
	}
}

// Publish message
func (t *transportWithCircuitBreaker) Publish(topic string, message []byte) error {
	return t.breaker.Execute(func() error {
		return t.Transport.Publish(topic, message)
	})
}

// Subscribe for messages
func (t *transportWithCircuitBreaker) Subscribe(topic string, queue string, handler bus.MessageHandler) (interface{}, error) {
	var sub interface{}

	err := t.breaker.Execute(func() error {
		var trErr error
		sub, trErr = t.Transport.Subscribe(topic, queue, handler)
		return trErr
	})

	return sub, err
}

// Unsubscribe subscription
func (t *transportWithCircuitBreaker) Unsubscribe(handle interface{}) error {
	return t.breaker.Execute(func() error {
		return t.Transport.Unsubscribe(handle)
	})
}
