package http

import (
	"github.com/dmibod/kanban/tools/mux"

	"fmt"
	"log"
	"net/http"
	"sync"
)

const (
	defaultPort   = 3000
)

var _ mux.Mux = (*Mux)(nil)

// Mux defines mux instance
type Mux struct {
	sync.Mutex
	port     int
	handlers map[string]*methodHandler
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

// Get serves GET request
func (m *Mux) Get(pattern string, h http.Handler) {
	m.Handle(http.MethodGet, pattern, h)
}

// Post serves POST request
func (m *Mux) Post(pattern string, h http.Handler) {
	m.Handle(http.MethodPost, pattern, h)
}

// Handle serves METHOD request
func (m *Mux) Handle(method string, pattern string, h http.Handler) {
	m.Lock()
	defer m.Unlock()
	mh, ok := m.handlers[pattern]; 
	if !ok {
		mh = &methodHandler{}
		http.Handle(pattern, mh)
		m.handlers[pattern] = mh
	}
	mh.methods[method] = h
}

type methodHandler struct {
	methods map[string]http.Handler
}

func (h *methodHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h, ok := h.methods[r.Method]; ok {
		h.ServeHTTP(w, r)
	} else {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		log.Println("Wrong HTTP method")
	}
}
