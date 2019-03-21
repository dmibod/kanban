package board

import (
	"github.com/dmibod/kanban/shared/domain/board"
	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"context"
)

// DomainRepository type
type DomainRepository struct {
	context.Context
	Repository *mongo.Repository
}

// Create board
func (r *DomainRepository) Create(entity *board.Entity) error {
	return r.Repository.ExecuteCommands(r.Context, []mongo.Command{mongo.Insert(entity.ID.String(), entity)})
}

// Update board
func (r *DomainRepository) Update(entity *board.Entity) error {
	return nil //r.Repository.Update(r.Context, entity.ID.String(), entity)
}

// Delete board
func (r *DomainRepository) Delete(entity *board.Entity) error {
	return r.Repository.Remove(r.Context, entity.ID.String())
}
