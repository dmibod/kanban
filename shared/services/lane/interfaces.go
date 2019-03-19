package lane

import (
	"github.com/dmibod/kanban/shared/kernel"
	"context"
)

// Reader interface
type Reader interface {
	// GetByID get by id
	GetByID(context.Context, kernel.ID) (*Model, error)
	// GetAll lanes
	GetAll(context.Context) ([]*ListModel, error)
	// GetByLaneID gets lanes by lane id
	GetByLaneID(context.Context, kernel.ID) ([]*ListModel, error)
	// GetByBoardID gets lanes by board id
	GetByBoardID(context.Context, kernel.ID) ([]*ListModel, error)
}

// Writer interface
type Writer interface {
	// Create lane
	Create(context.Context, *CreateModel) (kernel.ID, error)
	// Layout lane
	Layout(context.Context, kernel.ID, string) error
	// Name lane
	Name(context.Context, kernel.ID, string) error
	// Describe lane
	Describe(context.Context, kernel.ID, string) error
	// Remove lane
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
