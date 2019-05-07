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
	Lanes       []LaneModel
}

// LaneModel type
type LaneModel struct {
	ID          kernel.ID
	Type        string
	Name        string
	Description string
	Layout      string
	Lanes       []LaneModel
	Cards       []CardModel
}

// CardModel type
type CardModel struct {
	ID          kernel.ID
	Name        string
	Description string
}
