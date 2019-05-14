package hystrix

import (
	"errors"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/dmibod/kanban/shared/tools/circuit"
	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/logger/noop"
)

var (
	ErrRecoverFromPanic = errors.New("circuit: recover from panic")
)

// Breaker implements circuit breaker
type Breaker struct {
	logger.Logger
	name    string
	timeout int
}

// New creates new circuit breaker instance
func New(opts ...Option) *Breaker {

	var options Options

	for _, o := range opts {
		o(&options)
	}

	l := options.Logger

	if l == nil {
		l = &noop.Logger{}
	}

	return &Breaker{
		Logger:  l,
		name:    options.name,
		timeout: options.timeout,
	}
}

// ExecuteAsync executes handler within circuit breaker in async way
func (b *Breaker) ExecuteAsync(h circuit.Handler) error {
	output := make(chan bool, 1)

	hystrix.ConfigureCommand(b.name, hystrix.CommandConfig{Timeout: b.timeout})

	errors := hystrix.Go(b.name, func() error {

		err := b.recoverOnPanicHandler(h)()

		if err == nil {
			output <- true
		}

		return err
	}, nil)

	select {
	case <-output:
		return nil
	case err := <-errors:
		b.Errorln(err)
		return err
	}
}

// Execute executes handler within circuit breaker
func (b *Breaker) Execute(h circuit.Handler) error {
	hystrix.ConfigureCommand(b.name, hystrix.CommandConfig{Timeout: b.timeout})

	if err := hystrix.Do(b.name, func() error { return b.recoverOnPanicHandler(h)() }, nil); err != nil {
		b.Errorln(err)
		return err
	}

	return nil
}

func (b *Breaker) recoverOnPanicHandler(h circuit.Handler) circuit.Handler {
	return func() error {
		var err error
		defer func() {
			if e := recover(); e != nil {
				b.Errorln(e)
				err = ErrRecoverFromPanic
			}
		}()

		err = h()

		return err
	}
}
