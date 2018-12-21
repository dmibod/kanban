package hystrix

import (
	"github.com/afex/hystrix-go/hystrix"
	"github.com/dmibod/kanban/shared/tools/circuit"
	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/logger/noop"
)

// Breaker implements circuit breaker
type Breaker struct {
	logger  logger.Logger
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
		log = &noop.Logger{}
	}

	return &Breaker{
		logger:  log,
		name:    options.name,
		timeout: options.timeout,
	}
}

// ExecuteAsync executes handler within circuit breaker in async way
func (b *Breaker) ExecuteAsync(h circuit.Handler) error {
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

// Execute executes handler within circuit breaker
func (b *Breaker) Execute(h circuit.Handler) error {
	hystrix.ConfigureCommand(b.name, hystrix.CommandConfig{Timeout: b.timeout})

	if err := hystrix.Do(b.name, func() error { return h() }, nil); err != nil {
		b.logger.Debugln(err)
		return err
	}

	b.logger.Debugln("success")
	return nil
}
