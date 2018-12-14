package http

import (
	"github.com/dmibod/kanban/tools/mux"

	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	defaultPort   = 3000
	muxPortEnvVar = "MUX_PORT"
)

var _ mux.Mux = (*Mux)(nil)

// Mux defines mux instance
type Mux struct {
	port int
}

// New - creates mux
func New(opts ...Option) *Mux {
	var options Options

	options.Port = GetPortOrDefault(defaultPort)

	for _, o := range opts {
		o(&options)
	}

	return &Mux{
		port: options.Port,
	}
}

// Start - starts mux
func (m *Mux) Start() {
	log.Printf("Starting mux at port %v...\n", m.port)
	http.ListenAndServe(fmt.Sprintf(":%v", m.port), nil)
}

// Handle - attaches url handler to mux
func (m *Mux) Handle(pattern string, handler http.Handler) {
	http.Handle(pattern, handler)
}

// GetPortOrDefault - gets port from environment variable or fallbacks to default one
func GetPortOrDefault(defPort int) int {
	env := os.Getenv(muxPortEnvVar)

	if port, err := strconv.Atoi(env); err != nil {
		return defPort
	} else {
		return port
	}
}
