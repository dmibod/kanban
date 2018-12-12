package http

import(
	"github.com/dmibod/kanban/tools/mux"

	"fmt"
	"net/http"
)

var _ mux.Mux = (*Mux)(nil)

type Mux struct{
	port int
}

func New(opts ...Option) *Mux {
	var options Options
	for _, o := range opts {
		o(&options)
	}

	return &Mux{
		port: options.Port,
	}
}

func (m *Mux)Start(){
	http.ListenAndServe(fmt.Sprintf(":%v", m.port), nil)
}

func (m *Mux)Handle(pattern string, handler http.Handler){
  http.Handle(pattern, handler)
}