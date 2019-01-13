package event

// Handler interface
type Handler interface {
	Handle(event interface{})
}

// HandleFunc type
type HandleFunc func(interface{})

// Handle event
func (h HandleFunc) Handle(event interface{}) {
	h(event)
}

// Source interface
type Source interface {
	Fire()
}

// Registry interface
type Registry interface {
	Register(event interface{})
}

// Manager type
type Manager struct {
	fireOnRegister bool
	events         []interface{}
	handlers       []Handler
}

// CreateEventManager instance
func CreateEventManager() *Manager {
	return &Manager{
		events:   []interface{}{},
		handlers: []Handler{},
	}
}

// CreateFireOnRegisterEventManager instance
func CreateFireOnRegisterEventManager() *Manager {
	return &Manager{
		fireOnRegister: true,
		events:         []interface{}{},
		handlers:       []Handler{},
	}
}

// Register event
func (m *Manager) Register(event interface{}) {
	if event != nil {
		if m.fireOnRegister {
			m.notify(event)
		} else {
			m.events = append(m.events, event)
		}
	}
}

// Listen events
func (m *Manager) Listen(handler Handler) {
	if handler != nil {
		m.handlers = append(m.handlers, handler)
	}
}

// Fire events
func (m *Manager) Fire() {
	for _, event := range m.events {
		m.notify(event)
	}
	m.events = m.events[:0]
}

func (m *Manager) notify(event interface{}) {
	for _, handler := range m.handlers {
		handler.Handle(event)
	}
}
