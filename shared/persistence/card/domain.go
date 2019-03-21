package card

import (
	"github.com/dmibod/kanban/shared/domain/card"
	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"context"
)

// DomainRepository type
type DomainRepository struct {
	context.Context
	Repository *mongo.Repository
}

// Create card
func (r *DomainRepository) Create(entity *card.Entity) error {
	return r.Repository.ExecuteCommands(r.Context, []mongo.Command{mongo.Insert(entity.ID.String(), entity)})
}

// Update card
func (r *DomainRepository) Update(entity *card.Entity) error {
	return nil //r.Repository.Update(r.Context, entity.ID.String(), entity)
}

// Delete card
func (r *DomainRepository) Delete(entity *card.Entity) error {
	return r.Repository.Remove(r.Context, entity.ID.String())
}
