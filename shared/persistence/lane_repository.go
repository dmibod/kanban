package persistence

import (
	"context"

	err "github.com/dmibod/kanban/shared/domain/error"
	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"gopkg.in/mgo.v2/bson"
)

// LaneEntity definition
type LaneEntity struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Kind        string        `bson:"kind"`
	Name        string        `bson:"name"`
	Description string        `bson:"description"`
	Layout      string        `bson:"layout"`
	Children    []string      `bson:"children"`
}

// LaneVisitor type
type LaneVisitor func(*LaneEntity) error

// LaneRepository type
type LaneRepository struct {
	Repository *mongo.Repository
}

// FindLaneByID method
func (r *LaneRepository) FindLaneByID(ctx context.Context, id kernel.ID) (*LaneEntity, error) {
	if !id.IsValid() {
		return nil, err.ErrInvalidID
	}

	entity, err := r.Repository.FindByID(ctx, string(id), &LaneEntity{})
	if err != nil {
		return nil, err
	}

	lane, ok := entity.(*LaneEntity)
	if !ok {
		return nil, ErrInvalidType
	}

	return lane, nil
}

// FindLanes method
func (r *LaneRepository) FindLanes(ctx context.Context, criteria interface{}, visitor LaneVisitor) error {
	return r.Repository.Find(ctx, criteria, &LaneEntity{}, func(entity interface{}) error {
		lane, ok := entity.(*LaneEntity)
		if !ok {
			return ErrInvalidType
		}

		return visitor(lane)
	})
}

// CreateLaneRepository creates repository
func CreateLaneRepository(f *mongo.RepositoryFactory) *LaneRepository {
	return &LaneRepository{
		Repository: f.CreateRepository("lanes"),
	}
}
