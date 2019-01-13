package persistence

import (
	"context"

	"github.com/dmibod/kanban/shared/domain/event"

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
		board, ok := entity.(*BoardEntity)
		if !ok {
			return ErrInvalidType
		}

		return visitor(board)
	})
}

func (r *BoardRepository) UpdateBoard(ctx context.Context, manager *event.Manager, operation func() error) error {
	h := &eventHandler{
		commands: []mongo.Command{},
	}

	manager.Listen(h)

	err := operation()
	if err != nil {
		return err
	}

	return r.Repository.ExecuteCommands(ctx, h.commands)
}

// CreateBoardRepository creates new repository
func CreateBoardRepository(f *mongo.RepositoryFactory) *BoardRepository {
	return &BoardRepository{
		Repository: f.CreateRepository("boards"),
	}
}

type eventHandler struct {
	commands []mongo.Command
}

func (h *eventHandler) Handle(event interface{}) {
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

	h.commands = append(h.commands, command)
}
