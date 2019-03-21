package lane

import (
	"github.com/dmibod/kanban/shared/persistence"
	"context"

	"github.com/dmibod/kanban/shared/domain/lane"

	err "github.com/dmibod/kanban/shared/domain/error"
	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/tools/db/mongo"
)

// Visitor type
type Visitor func(*LaneEntity) error

// Repository type
type Repository struct {
	Repository *mongo.Repository
}

// CreateRepository creates repository
func CreateRepository(f *mongo.RepositoryFactory) *Repository {
	return &Repository{
		Repository: f.CreateRepository("lanes"),
	}
}

// FindLaneByID method
func (r *Repository) FindLaneByID(ctx context.Context, id kernel.ID) (*LaneEntity, error) {
	if !id.IsValid() {
		return nil, err.ErrInvalidID
	}

	entity, err := r.Repository.FindByID(ctx, string(id), &LaneEntity{})
	if err != nil {
		return nil, err
	}

	lane, ok := entity.(*LaneEntity)
	if !ok {
		return nil, persistence.ErrInvalidType
	}

	return lane, nil
}

// FindLanes method
func (r *Repository) FindLanes(ctx context.Context, criteria interface{}, visitor Visitor) error {
	return r.Repository.Find(ctx, criteria, &LaneEntity{}, func(entity interface{}) error {
		lane, ok := entity.(*LaneEntity)
		if !ok {
			return persistence.ErrInvalidType
		}

		return visitor(lane)
	})
}

// GetRepository - returns domain repository
func (r *Repository) GetRepository(ctx context.Context) lane.Repository {
	return &DomainRepository{Context: ctx, Repository: r.Repository}
}
