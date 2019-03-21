package card

import (
	"context"
	"github.com/dmibod/kanban/shared/persistence"
	"github.com/dmibod/kanban/shared/domain/card"
	err "github.com/dmibod/kanban/shared/domain/error"
	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/tools/db/mongo"
)

// Visitor type
type Visitor func(*CardEntity) error

// Repository type
type Repository struct {
	Repository *mongo.Repository
}

// CreateRepository creates new cards repository
func CreateRepository(f *mongo.RepositoryFactory) *Repository {
	return &Repository{
		Repository: f.CreateRepository("cards"),
	}
}

// FindCardByID method
func (r *Repository) FindCardByID(ctx context.Context, id kernel.ID) (*CardEntity, error) {
	if !id.IsValid() {
		return nil, err.ErrInvalidID
	}

	entity, err := r.Repository.FindByID(ctx, string(id), &CardEntity{})
	if err != nil {
		return nil, err
	}

	card, ok := entity.(*CardEntity)
	if !ok {
		return nil, persistence.ErrInvalidType
	}

	return card, nil
}

// FindCards method
func (r *Repository) FindCards(ctx context.Context, criteria interface{}, visitor Visitor) error {
	return r.Repository.Find(ctx, criteria, &CardEntity{}, func(entity interface{}) error {
		if card, ok := entity.(*CardEntity); ok {
			return visitor(card)
		}
		return persistence.ErrInvalidType
	})
}

// GetRepository - returns domain repository
func (r *Repository) GetRepository(ctx context.Context) card.Repository {
	return &DomainRepository{Context: ctx, Repository: r.Repository}
}
