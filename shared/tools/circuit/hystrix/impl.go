package hystrix

import (
	"github.com/dmibod/kanban/shared/tools/log"
	"github.com/dmibod/kanban/shared/tools/log/logger"
)

// Breaker implements circuit breaker
type Breaker struct {
	logger log.Logger
}

// New creates new circuit breaker instance
func New(opts ...Option) *Breaker {

	var options Options

	for _, o := range opts {
		o(&options)
	}

	log := options.logger

	if log == nil {
		log = logger.New(logger.WithPrefix("[CIRCUIT] "), logger.WithDebug(true))
	}

	return &Breaker{
		logger: log,
	}
}

