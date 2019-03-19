package card

import (
	"github.com/dmibod/kanban/shared/kernel"
)

// CreateModel type
type CreateModel struct {
	Name        string
	Description string
}

// Model type
type Model struct {
	ID          kernel.ID
	Name        string
	Description string
}
