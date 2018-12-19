package hystrix

import (
	"github.com/afex/hystrix-go/hystrix"
	"github.com/dmibod/kanban/shared/tools/circuit"
	"github.com/dmibod/kanban/shared/tools/log"
	"github.com/dmibod/kanban/shared/tools/log/logger"
)

// Breaker implements circuit breaker
type Breaker struct {
	logger  log.Logger
	name    string
	timeout int
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
		logger:  log,
		name:    options.name,
		timeout: options.timeout,
	}
}

// Execute executes handler within circuit breaker
func (b *Breaker) Execute(h circuit.Handler) error {
	output := make(chan bool, 1)

	hystrix.ConfigureCommand(b.name, hystrix.CommandConfig{Timeout: b.timeout})

	errors := hystrix.Go(b.name, func() error {

		err := h()

		if err == nil {
			output <- true
		}

		return err
	}, nil)

	select {
	case <-output:
		b.logger.Debugln("success")
		return nil
	case err := <-errors:
		b.logger.Debugln(err)
		return err
	}
}
