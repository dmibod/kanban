package persistence

import (
	"github.com/dmibod/kanban/shared/domain/event"
)

// Service interface
type Service interface {
	Listen(event.Bus)
}
