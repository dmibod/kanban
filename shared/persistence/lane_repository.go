package persistence

import (
	"context"

	"github.com/dmibod/kanban/shared/domain/lane"

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

// GetRepository - returns domain repository
func (r *LaneRepository) GetRepository(ctx context.Context) *LaneDomainRepository {
	return &LaneDomainRepository{Context: ctx, Repository: r.Repository}
}

// CreateLaneRepository creates repository
func CreateLaneRepository(f *mongo.RepositoryFactory) *LaneRepository {
	return &LaneRepository{
		Repository: f.CreateRepository("lanes"),
	}
}

// LaneDomainRepository type
type LaneDomainRepository struct {
	context.Context
	Repository *mongo.Repository
}

// Create lane
func (r *LaneDomainRepository) Create(entity *lane.Entity) error {
	return r.Repository.ExecuteCommands(r.Context, []mongo.Command{mongo.Insert(entity.ID.String(), entity)})
}

// Update lane
func (r *LaneDomainRepository) Update(entity *lane.Entity) error {
	return nil //r.Repository.Update(r.Context, entity.ID.String(), entity)
}

// Delete lane
func (r *LaneDomainRepository) Delete(entity *lane.Entity) error {
	return r.Repository.Remove(r.Context, entity.ID.String())
}
