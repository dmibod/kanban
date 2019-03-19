package board

import (
	"github.com/dmibod/kanban/shared/kernel"
)


// CreateModel type
type CreateModel struct {
	Owner       string
	Name        string
	Description string
	Layout      string
}

// ListModel type
type ListModel struct {
	ID          kernel.ID
	Owner       string
	Name        string
	Description string
	Shared      bool
	Layout      string
}

// Model type
type Model struct {
	ID          kernel.ID
	Owner       string
	Name        string
	Description string
	Shared      bool
	Layout      string
	Lanes       []interface{}
}
