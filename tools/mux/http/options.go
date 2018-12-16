package http

import (
	"strconv"
	"os"
)

const muxPortEnvVar = "MUX_PORT"

// Options can be used to create a customized mux.
type Options struct {
	Port int
}

// Option is a function on the options for a http mux.
type Option func(*Options)

// WithPort initializes Port option
func WithPort(p int) Option {
	return func(o *Options) {
		o.Port = p
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

