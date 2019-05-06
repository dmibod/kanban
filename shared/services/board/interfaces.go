package board

import (
	"context"
	"github.com/dmibod/kanban/shared/kernel"
)

// Reader interface
type Reader interface {
	// GetByID get by id
	GetByID(context.Context, kernel.ID) (*Model, error)
	// GetByOwner boards
	GetByOwner(context.Context, string) ([]*ListModel, error)
}

// Writer interface
type Writer interface {
	// Create board
	Create(context.Context, *CreateModel) (kernel.ID, error)
	// Layout board
	Layout(context.Context, kernel.ID, string) error
	// Name board
	Name(context.Context, kernel.ID, string) error
	// Describe board
	Describe(context.Context, kernel.ID, string) error
	// Share board
	Share(context.Context, kernel.ID, bool) error
	// Remove board by id
	Remove(context.Context, kernel.ID) error
	// AppendLane to board
	AppendLane(context.Context, kernel.MemberID) error
	// ExcludeLane from board
	ExcludeLane(context.Context, kernel.MemberID) error
}

// Service interface
type Service interface {
	Reader
	Writer
}
