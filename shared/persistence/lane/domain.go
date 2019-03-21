package lane

import (
	"context"

	"github.com/dmibod/kanban/shared/domain/lane"
	"github.com/dmibod/kanban/shared/tools/db/mongo"
)

// DomainRepository type
type DomainRepository struct {
	context.Context
	Repository *mongo.Repository
}

// Create lane
func (r *DomainRepository) Create(entity *lane.Entity) error {
	return r.Repository.ExecuteCommands(r.Context, []mongo.Command{mongo.Insert(entity.ID.String(), entity)})
}

// Update lane
func (r *DomainRepository) Update(entity *lane.Entity) error {
	return nil //r.Repository.Update(r.Context, entity.ID.String(), entity)
}

// Delete lane
func (r *DomainRepository) Delete(entity *lane.Entity) error {
	return r.Repository.Remove(r.Context, entity.ID.String())
}
