package persistence

import (
	"context"

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
		card, ok := entity.(*CardEntity)
		if !ok {
			return ErrInvalidType
		}

		return visitor(card)
	})
}

// CreateCardRepository creates new cards repository
func CreateCardRepository(f *mongo.RepositoryFactory) *CardRepository {
	return &CardRepository{
		Repository: f.CreateRepository("cards"),
	}
}
