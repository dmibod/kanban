package card

import (
	"context"
	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/persistence/models"
	"github.com/dmibod/kanban/shared/tools/db/mongo"
)

// Repository type
type Repository struct {
	repository *mongo.Repository
}

// CreateRepository method
func CreateRepository(r *mongo.Repository) Repository {
	return Repository{repository: r}
}

// FindByID method
func (r Repository) FindByID(ctx context.Context, id kernel.MemberID, visitor func(*models.Card) error) error {
	query := OneQuery{BoardID: id.SetID.String(), ID: id.ID.String()}

	return r.repository.Execute(ctx, query.Operation(ctx, visitor))
}

// FindByParent method
func (r Repository) FindByParent(ctx context.Context, id kernel.MemberID, visitor func(*models.Card) error) error {
	query := ListQuery{BoardID: id.SetID.String(), LaneID: id.ID.String()}

	return r.repository.Execute(ctx, query.Operation(ctx, visitor))
}

// Create card
func (r Repository) Create(ctx context.Context, boardID string, card *models.Card) error {
	command := CreateCommand{BoardID: boardID, Card: card}

	return r.repository.Execute(ctx, command.Operation(ctx))
}

// Remove card
func (r Repository) Remove(ctx context.Context, boardID string, id string) error {
	command := RemoveCommand{BoardID: boardID, ID: id}

	return r.repository.Execute(ctx, command.Operation(ctx))
}

// Update card
func (r Repository) Update(ctx context.Context, boardID string, id string, field string, value interface{}) error {
	command := UpdateCommand{BoardID: boardID, ID: id, Field: field, Value: value}

	return r.repository.Execute(ctx, command.Operation(ctx))
}
