package http

import (
	"github.com/dmibod/kanban/shared/tools/log"
	"github.com/dmibod/kanban/shared/tools/log/logger"
	"github.com/dmibod/kanban/shared/tools/mux"

	"fmt"
	"net/http"
	"sync"
)

const (
	defaultPort = 3000
	anyMethod   = "*"
)

var _ mux.Mux = (*Mux)(nil)

// Mux defines mux instance
type Mux struct {
	sync.Mutex
	port     int
	logger   log.Logger
	handlers map[string]*methodHandler
}

// New - creates mux
func New(opts ...Option) *Mux {
	var options Options

	options.Port = GetPortOrDefault(defaultPort)

	for _, o := range opts {
		o(&options)
	}

	if options.Logger == nil {
		options.Logger = logger.New(logger.WithPrefix("[MUX] "), logger.WithDebug(true))
	}

	return &Mux{
		port:     options.Port,
		logger:   options.Logger,
		handlers: make(map[string]*methodHandler),
	}
}

// Start - starts mux
func (m *Mux) Start() {
	m.logger.Debugf("Starting mux at port %v...\n", m.port)
	http.ListenAndServe(fmt.Sprintf(":%v", m.port), nil)
}

// All serves any-method request
func (m *Mux) All(pattern string, h http.Handler) {
	m.handle(anyMethod, pattern, h)
}

// Get serves GET request
func (m *Mux) Get(pattern string, h http.Handler) {
	m.handle(http.MethodGet, pattern, h)
}

// Post serves POST request
func (m *Mux) Post(pattern string, h http.Handler) {
	m.handle(http.MethodPost, pattern, h)
}

func (m *Mux) handle(method string, pattern string, h http.Handler) {
	m.Lock()
	defer m.Unlock()
	mh, ok := m.handlers[pattern]
	if !ok {
		mh = &methodHandler{logger: m.logger, methods: make(map[string]http.Handler)}
		http.Handle(pattern, mh)
		m.handlers[pattern] = mh
	}
	mh.methods[method] = h
}

type methodHandler struct {
	logger   log.Logger
	methods map[string]http.Handler
}

func (mh *methodHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mh.logger.Debugf("%v request received\n", r.Method)
	if h, ok := mh.methods[r.Method]; ok {
		h.ServeHTTP(w, r)
	} else if h, ok := mh.methods[anyMethod]; ok {
		h.ServeHTTP(w, r)
	} else {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		mh.logger.Errorln("Wrong HTTP method")
	}
}
