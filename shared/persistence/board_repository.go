package persistence

import (
	"context"
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

// CreateBoardRepository creates new repository
func CreateBoardRepository(f *mongo.RepositoryFactory) *BoardRepository {
	return &BoardRepository{
		Repository: f.CreateRepository("boards"),
	}
}
