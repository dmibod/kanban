package persistence

import (
	"context"

	"github.com/dmibod/kanban/shared/domain"
	"github.com/dmibod/kanban/shared/kernel"
	"gopkg.in/mgo.v2/bson"

	"github.com/dmibod/kanban/shared/tools/db"
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

// LaneRepository interface
type LaneRepository interface {
	db.RepositoryEntity
	db.Repository
	DomainRepository(context.Context) domain.Repository
	FindLaneByID(context.Context, kernel.ID) (*LaneEntity, error)
	FindLanes(context.Context, interface{}, LaneVisitor) error
}

type laneRepository struct {
	db.Repository
}

func (r *laneRepository) CreateInstance() interface{} {
	return &LaneEntity{}
}

func (r *laneRepository) GetID(entity interface{}) string {
	return entity.(*LaneEntity).ID.Hex()
}

func (r *laneRepository) DomainRepository(ctx context.Context) domain.Repository {
	return &laneDomainRepository{
		ctx:        ctx,
		Repository: r.Repository,
	}
}

func (r *laneRepository) FindLaneByID(ctx context.Context, id kernel.ID) (*LaneEntity, error) {
	if !id.IsValid() {
		return nil, domain.ErrInvalidID
	}

	entity, err := r.Repository.FindByID(ctx, string(id))
	if err != nil {
		return nil, err
	}

	lane, ok := entity.(*LaneEntity)
	if !ok {
		return nil, ErrInvalidType
	}

	return lane, nil
}

func (r *laneRepository) FindLanes(ctx context.Context, criteria interface{}, visitor LaneVisitor) error {
	return r.Repository.Find(ctx, criteria, func(entity interface{}) error {
		lane, ok := entity.(*LaneEntity)
		if !ok {
			return ErrInvalidType
		}

		return visitor(lane)
	})
}

// CreateLaneRepository creates repository
func CreateLaneRepository(f db.RepositoryFactory) LaneRepository {
	r := &laneRepository{}
	r.Repository = f.CreateRepository("lanes", r)
	return r
}

type laneDomainRepository struct {
	ctx context.Context
	db.Repository
}

func (r *laneDomainRepository) Fetch(id kernel.ID) (interface{}, error) {
	if !id.IsValid() {
		return nil, domain.ErrInvalidID
	}

	persistent, err := r.Repository.FindByID(r.ctx, string(id))
	if err != nil {
		return nil, err
	}

	entity, ok := persistent.(*LaneEntity)
	if !ok {
		return nil, domain.ErrInvalidType
	}

	return r.mapEntityToDomain(entity), nil
}

func (r *laneDomainRepository) Delete(id kernel.ID) (interface{}, error) {
	entity, err := r.Fetch(id)
	if err == nil {
		err = r.Repository.Remove(r.ctx, string(id))
	}

	return entity, err
}

func (r *laneDomainRepository) Persist(entity interface{}) (kernel.ID, error) {
	lane, ok := entity.(domain.LaneEntity)
	if !ok {
		return kernel.EmptyID, domain.ErrInvalidType
	}

	persistent := r.mapDomainToEntity(&lane)

	if lane.ID.IsValid() {
		return lane.ID, r.Repository.Update(r.ctx, persistent)
	}

	id, err := r.Repository.Create(r.ctx, persistent)
	if err != nil {
		return kernel.EmptyID, err
	}

	return kernel.ID(id), nil
}

func (r *laneDomainRepository) mapEntityToDomain(entity *LaneEntity) *domain.LaneEntity {
	children := []kernel.ID{}
	for _, id := range entity.Children {
		children = append(children, kernel.ID(id))
	}

	return &domain.LaneEntity{
		ID:          kernel.ID(entity.ID.Hex()),
		Kind:        entity.Kind,
		Name:        entity.Name,
		Description: entity.Description,
		Layout:      entity.Layout,
		Children:    children,
	}
}

func (r *laneDomainRepository) mapDomainToEntity(domainEntity *domain.LaneEntity) *LaneEntity {
	children := []string{}
	for _, id := range domainEntity.Children {
		children = append(children, string(id))
	}

	entity := &LaneEntity{
		Kind:        domainEntity.Kind,
		Name:        domainEntity.Name,
		Description: domainEntity.Description,
		Layout:      domainEntity.Layout,
		Children:    children,
	}

	if domainEntity.ID.IsValid() {
		entity.ID = bson.ObjectIdHex(string(domainEntity.ID))
	}

	return entity
}
