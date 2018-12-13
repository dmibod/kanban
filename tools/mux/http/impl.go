package http

import(
	"github.com/dmibod/kanban/tools/mux"

	"fmt"
	"log"
	"net/http"
)

var _ mux.Mux = (*Mux)(nil)

// Mux defines mux instance
type Mux struct{
	port int
}

// New - creates mux
func New(opts ...Option) *Mux {
	var options Options
	for _, o := range opts {
		o(&options)
	}

	return &Mux{
		port: options.Port,
	}
}

// Start - starts mux
func (m *Mux)Start(){
	log.Printf("Starting mux at port %v...\n", m.port)
	http.ListenAndServe(fmt.Sprintf(":%v", m.port), nil)
}

// Handle - attaches url handler to mux
func (m *Mux)Handle(pattern string, handler http.Handler){
  http.Handle(pattern, handler)
}