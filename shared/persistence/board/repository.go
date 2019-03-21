package board

import (
	"github.com/dmibod/kanban/shared/persistence"
	"context"

	"github.com/dmibod/kanban/shared/domain/board"
	"github.com/dmibod/kanban/shared/tools/db/mongo"

	err "github.com/dmibod/kanban/shared/domain/error"
	"github.com/dmibod/kanban/shared/kernel"
)

// Visitor type
type Visitor func(*BoardEntity) error

// Repository type
type Repository struct {
	Repository *mongo.Repository
}

// CreateRepository creates new repository
func CreateRepository(f *mongo.RepositoryFactory) *Repository {
	return &Repository{
		Repository: f.CreateRepository("boards"),
	}
}

// FindBoardByID method
func (r *Repository) FindBoardByID(ctx context.Context, id kernel.ID) (*BoardEntity, error) {
	if !id.IsValid() {
		return nil, err.ErrInvalidID
	}

	entity, err := r.Repository.FindByID(ctx, string(id), &BoardEntity{})
	if err != nil {
		return nil, err
	}

	if board, ok := entity.(*BoardEntity); ok {
		return board, nil
	}

	return nil, persistence.ErrInvalidType
}

// FindBoards method
func (r *Repository) FindBoards(ctx context.Context, criteria interface{}, visitor Visitor) error {
	return r.Repository.Find(ctx, criteria, &BoardEntity{}, func(entity interface{}) error {
		if board, ok := entity.(*BoardEntity); ok {
			return visitor(board)
		}

		return persistence.ErrInvalidType
	})
}

// Handle domain event
func (r *Repository) Handle(ctx context.Context, event interface{}) {
	if event == nil {
		return
	}

	var command mongo.Command

	switch e := event.(type) {
	case board.CreatedEvent:
		command = mongo.Insert(string(e.ID), e.Entity)
	case board.DeletedEvent:
		command = mongo.Remove(string(e.ID))
	case board.NameChangedEvent:
		command = mongo.Update(string(e.ID), "name", e.NewValue)
	case board.DescriptionChangedEvent:
		command = mongo.Update(string(e.ID), "description", e.NewValue)
	case board.LayoutChangedEvent:
		command = mongo.Update(string(e.ID), "layout", e.NewValue)
	case board.SharedChangedEvent:
		command = mongo.Update(string(e.ID), "shared", e.NewValue)
	case board.ChildAppendedEvent:
		command = mongo.CustomUpdate(string(e.ID), mongo.AddToSet("children", e.ChildID))
	case board.ChildRemovedEvent:
		command = mongo.CustomUpdate(string(e.ID), mongo.PullFromSet("children", e.ChildID))
	default:
		return
	}

	r.Repository.ExecuteCommands(ctx, []mongo.Command{command})
}

// GetRepository - returns domain repository
func (r *Repository) GetRepository(ctx context.Context) board.Repository {
	return &DomainRepository{Context: ctx, Repository: r.Repository}
}
