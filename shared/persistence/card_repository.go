package persistence

import (
	"context"
	"github.com/dmibod/kanban/shared/domain/card"

	err "github.com/dmibod/kanban/shared/domain/error"
	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"gopkg.in/mgo.v2/bson"
)

// CardEntity maps card to/from mongo db
type CardEntity struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Name        string        `bson:"name"`
	Description string        `bson:"description"`
}

// CardVisitor type
type CardVisitor func(*CardEntity) error

// CardRepository type
type CardRepository struct {
	Repository *mongo.Repository
}

// FindCardByID method
func (r *CardRepository) FindCardByID(ctx context.Context, id kernel.ID) (*CardEntity, error) {
	if !id.IsValid() {
		return nil, err.ErrInvalidID
	}

	entity, err := r.Repository.FindByID(ctx, string(id), &CardEntity{})
	if err != nil {
		return nil, err
	}

	card, ok := entity.(*CardEntity)
	if !ok {
		return nil, ErrInvalidType
	}

	return card, nil
}

// FindCards method
func (r *CardRepository) FindCards(ctx context.Context, criteria interface{}, visitor CardVisitor) error {
	return r.Repository.Find(ctx, criteria, &CardEntity{}, func(entity interface{}) error {
		if card, ok := entity.(*CardEntity); ok {
			return visitor(card)
		}
		return ErrInvalidType
	})
}

// GetRepository - returns domain repository
func (r *CardRepository) GetRepository(ctx context.Context) *CardDomainRepository {
	return &CardDomainRepository{Context: ctx, Repository: r.Repository}
}

// CreateCardRepository creates new cards repository
func CreateCardRepository(f *mongo.RepositoryFactory) *CardRepository {
	return &CardRepository{
		Repository: f.CreateRepository("cards"),
	}
}

// CardDomainRepository type
type CardDomainRepository struct {
	context.Context
	Repository *mongo.Repository
}

// Create card
func (r *CardDomainRepository) Create(entity *card.Entity) error {
	return r.Repository.ExecuteCommands(r.Context, []mongo.Command{mongo.Insert(entity.ID.String(), entity)})
}

// Update card
func (r *CardDomainRepository) Update(entity *card.Entity) error {
	return nil //r.Repository.Update(r.Context, entity.ID.String(), entity)
}

// Delete card
func (r *CardDomainRepository) Delete(entity *card.Entity) error {
	return r.Repository.Remove(r.Context, entity.ID.String())
}
