package http

import (
	"os"
	"strconv"

	"github.com/dmibod/kanban/shared/tools/logger"
)

const muxPortEnvVar = "MUX_PORT"

// Options can be used to create a customized mux.
type Options struct {
	Port   int
	Logger logger.Logger
}

// Option is a function on the options for a http mux.
type Option func(*Options)

// WithPort initializes Port option
func WithPort(p int) Option {
	return func(o *Options) {
		o.Port = p
	}
}

// WithLogger initializes Logger option
func WithLogger(l logger.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}

// GetPortOrDefault gets port from environment variable or fallbacks to default one
func GetPortOrDefault(defPort int) int {
	env := os.Getenv(muxPortEnvVar)

	port, err := strconv.Atoi(env)
	if err != nil {
		return defPort
	}

	return port
}
