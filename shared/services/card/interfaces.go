package card

import (
	"github.com/dmibod/kanban/shared/kernel"
	"context"
)

// Reader interface
type Reader interface {
	// GetByID gets card by id
	GetByID(context.Context, kernel.MemberID) (*Model, error)
	// GetByLaneID gets cards by lane id
	GetByLaneID(context.Context, kernel.MemberID) ([]*Model, error)
}

// Writer interface
type Writer interface {
	// Create card
	Create(context.Context, kernel.ID, *CreateModel) (kernel.ID, error)
	// Name card
	Name(context.Context, kernel.MemberID, string) error
	// Describe card
	Describe(context.Context, kernel.MemberID, string) error
	// Remove card
	Remove(context.Context, kernel.MemberID) error
}

// Service interface
type Service interface {
	Reader
	Writer
}
