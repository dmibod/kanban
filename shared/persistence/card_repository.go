package persistence

import (
	"context"
	"github.com/dmibod/kanban/shared/domain"
	"github.com/dmibod/kanban/shared/kernel"
	"gopkg.in/mgo.v2/bson"

	"github.com/dmibod/kanban/shared/tools/db"
)

// CardEntity maps card to/from mongo db
type CardEntity struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Name        string        `bson:"name"`
	Description string        `bson:"description"`
}

// CardVisitor type
type CardVisitor func(*CardEntity) error

// CardRepository interface
type CardRepository interface {
	db.RepositoryEntity
	db.Repository
	DomainRepository(context.Context) domain.Repository
	FindCardByID(context.Context, kernel.ID) (*CardEntity, error)
	FindCards(context.Context, interface{}, CardVisitor) error
}

type cardRepository struct {
	db.Repository
}

func (r *cardRepository) CreateInstance() interface{} {
	return &CardEntity{}
}

func (r *cardRepository) GetID(entity interface{}) string {
	return entity.(*CardEntity).ID.Hex()
}

func (r *cardRepository) DomainRepository(ctx context.Context) domain.Repository {
	return &cardDomainRepository{
		ctx:        ctx,
		Repository: r.Repository,
	}
}

func (r *cardRepository) FindCardByID(ctx context.Context, id kernel.ID) (*CardEntity, error) {
	if !id.IsValid() {
		return nil, domain.ErrInvalidID
	}

	entity, err := r.Repository.FindByID(ctx, string(id))
	if err != nil {
		return nil, err
	}

	card, ok := entity.(*CardEntity)
	if !ok {
		return nil, ErrInvalidType
	}

	return card, nil
}

func (r *cardRepository) FindCards(ctx context.Context, criteria interface{}, visitor CardVisitor) error {
	return r.Repository.Find(ctx, criteria, func(entity interface{}) error {
		card, ok := entity.(*CardEntity)
		if !ok {
			return ErrInvalidType
		}

		return visitor(card)
	})
}

// CreateCardRepository creates new cards repository
func CreateCardRepository(f db.RepositoryFactory) CardRepository {
	r := &cardRepository{}
	r.Repository = f.CreateRepository("cards", r)
	return r
}

type cardDomainRepository struct {
	ctx context.Context
	db.Repository
}

func (r *cardDomainRepository) Fetch(id kernel.ID) (interface{}, error) {
	if !id.IsValid() {
		return nil, domain.ErrInvalidID
	}

	persistent, err := r.Repository.FindByID(r.ctx, string(id))
	if err != nil {
		return nil, err
	}

	entity, ok := persistent.(*CardEntity)
	if !ok {
		return nil, domain.ErrInvalidType
	}

	return r.mapEntityToDomain(entity), nil
}

func (r *cardDomainRepository) Delete(id kernel.ID) (interface{}, error) {
	entity, err := r.Fetch(id)
	if err == nil {
		err = r.Repository.Remove(r.ctx, string(id))
	}

	return entity, err
}

func (r *cardDomainRepository) Persist(entity interface{}) (kernel.ID, error) {
	card, ok := entity.(domain.CardEntity)
	if !ok {
		return kernel.EmptyID, domain.ErrInvalidType
	}

	persistent := r.mapDomainToEntity(&card)

	if card.ID.IsValid() {
		return card.ID, r.Repository.Update(r.ctx, persistent)
	}

	id, err := r.Repository.Create(r.ctx, persistent)
	if err != nil {
		return kernel.EmptyID, err
	}

	return kernel.ID(id), nil
}

func (r *cardDomainRepository) mapEntityToDomain(entity *CardEntity) *domain.CardEntity {
	return &domain.CardEntity{
		ID:          kernel.ID(entity.ID.Hex()),
		Name:        entity.Name,
		Description: entity.Description,
	}
}

func (r *cardDomainRepository) mapDomainToEntity(domainEntity *domain.CardEntity) *CardEntity {
	entity := &CardEntity{
		Name:        domainEntity.Name,
		Description: domainEntity.Description,
	}

	if domainEntity.ID.IsValid() {
		entity.ID = bson.ObjectIdHex(string(domainEntity.ID))
	}

	return entity
}
