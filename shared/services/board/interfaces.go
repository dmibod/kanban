package board

import (
	"github.com/dmibod/kanban/shared/kernel"
	"context"
)

// Reader interface
type Reader interface {
	// GetByID get by id
	GetByID(context.Context, kernel.ID) (*Model, error)
	// GetByOwner boards
	GetByOwner(context.Context, string) ([]*ListModel, error)
	// GetAll boards
	GetAll(context.Context) ([]*ListModel, error)
}

// Writer interface
type Writer interface {
	// Create by payload
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
	// AppendChild to lane
	AppendChild(context.Context, kernel.ID, kernel.ID) error
	// ExcludeChild from lane
	ExcludeChild(context.Context, kernel.ID, kernel.ID) error
}

// Service interface
type Service interface {
	Reader
	Writer
}
