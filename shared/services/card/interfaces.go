package card

import (
	"github.com/dmibod/kanban/shared/kernel"
	"context"
)

// Reader interface
type Reader interface {
	// GetByID gets card by id
	GetByID(context.Context, kernel.ID) (*Model, error)
	// GetAll cards
	GetAll(context.Context) ([]*Model, error)
	// GetByLaneID gets cards by lane id
	GetByLaneID(context.Context, kernel.ID) ([]*Model, error)
}

// Writer interface
type Writer interface {
	// Create card
	Create(context.Context, *CreateModel) (kernel.ID, error)
	// Name card
	Name(context.Context, kernel.ID, string) error
	// Describe card
	Describe(context.Context, kernel.ID, string) error
	// Remove card
	Remove(context.Context, kernel.ID) error
}

// Service interface
type Service interface {
	Reader
	Writer
}
