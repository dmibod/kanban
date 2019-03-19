package persistence

import (
	"context"

	"github.com/dmibod/kanban/shared/domain/board"
	"github.com/dmibod/kanban/shared/tools/db/mongo"

	err "github.com/dmibod/kanban/shared/domain/error"
	"github.com/dmibod/kanban/shared/kernel"
	"gopkg.in/mgo.v2/bson"
)

// BoardEntity entity
type BoardEntity struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Owner       string        `bson:"owner"`
	Name        string        `bson:"name"`
	Description string        `bson:"description"`
	Layout      string        `bson:"layout"`
	Shared      bool          `bson:"shared"`
	Children    []string      `bson:"children"`
}

// BoardVisitor type
type BoardVisitor func(*BoardEntity) error

// BoardRepository type
type BoardRepository struct {
	Repository *mongo.Repository
}

// FindBoardByID method
func (r *BoardRepository) FindBoardByID(ctx context.Context, id kernel.ID) (*BoardEntity, error) {
	if !id.IsValid() {
		return nil, err.ErrInvalidID
	}

	entity, err := r.Repository.FindByID(ctx, string(id), &BoardEntity{})
	if err != nil {
		return nil, err
	}

	board, ok := entity.(*BoardEntity)
	if !ok {
		return nil, ErrInvalidType
	}

	return board, nil
}

// FindBoards method
func (r *BoardRepository) FindBoards(ctx context.Context, criteria interface{}, visitor BoardVisitor) error {
	return r.Repository.Find(ctx, criteria, &BoardEntity{}, func(entity interface{}) error {
		if board, ok := entity.(*BoardEntity); ok {
			return visitor(board)
		}

		return ErrInvalidType
	})
}

// Handle domain event
func (r *BoardRepository) Handle(ctx context.Context, event interface{}) {
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
func (r *BoardRepository) GetRepository(ctx context.Context) *BoardDomainRepository {
	return &BoardDomainRepository{Context: ctx, Repository: r.Repository}
}

// CreateBoardRepository creates new repository
func CreateBoardRepository(f *mongo.RepositoryFactory) *BoardRepository {
	return &BoardRepository{
		Repository: f.CreateRepository("boards"),
	}
}

// BoardDomainRepository type
type BoardDomainRepository struct {
	context.Context
	Repository *mongo.Repository
}

// Create board
func (r *BoardDomainRepository) Create(entity *board.Entity) error {
	return r.Repository.ExecuteCommands(r.Context, []mongo.Command{mongo.Insert(entity.ID.String(), entity)})
}

// Update board
func (r *BoardDomainRepository) Update(entity *board.Entity) error {
	return nil //r.Repository.Update(r.Context, entity.ID.String(), entity)
}

// Delete board
func (r *BoardDomainRepository) Delete(entity *board.Entity) error {
	return r.Repository.Remove(r.Context, entity.ID.String())
}
