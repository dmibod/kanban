package lane

import (
	"context"
	"github.com/dmibod/kanban/shared/kernel"
)

// Reader interface
type Reader interface {
	// GetByID get by id
	GetByID(context.Context, kernel.MemberID) (*Model, error)
	// GetByBoardID gets lanes by board id
	GetByBoardID(context.Context, kernel.ID) ([]*ListModel, error)
	// GetByLaneID gets lanes by lane id
	GetByLaneID(context.Context, kernel.MemberID) ([]*ListModel, error)
}

// Writer interface
type Writer interface {
	// Create lane
	Create(context.Context, kernel.ID, *CreateModel) (kernel.ID, error)
	// Layout lane
	Layout(context.Context, kernel.MemberID, string) error
	// Name lane
	Name(context.Context, kernel.MemberID, string) error
	// Describe lane
	Describe(context.Context, kernel.MemberID, string) error
	// Remove lane
	Remove(context.Context, kernel.MemberID) error
	// AppendChild to lane
	AppendChild(context.Context, kernel.MemberID, kernel.ID) error
	// ExcludeChild from lane
	ExcludeChild(context.Context, kernel.MemberID, kernel.ID) error
}

// Service interface
type Service interface {
	Reader
	Writer
}
