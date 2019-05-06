package lane

import (
	"github.com/dmibod/kanban/shared/kernel"
)

// CreateModel type
type CreateModel struct {
	Type        string
	Name        string
	Description string
	Layout      string
}

// ListModel type
type ListModel struct {
	ID          kernel.ID
	Type        string
	Name        string
	Description string
	Layout      string
}

// Model type
type Model struct {
	ID          kernel.ID
	Type        string
	Name        string
	Description string
	Layout      string
}
