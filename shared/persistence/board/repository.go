package board

import (
	"context"
	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/persistence/models"
	"github.com/dmibod/kanban/shared/tools/db/mongo"
)

// Repository for boards
type Repository struct {
	repository *mongo.Repository
}

// CreateRepository method
func CreateRepository(r *mongo.Repository) Repository {
	return Repository{repository: r}
}

// FindByID method
func (r Repository) FindByID(ctx context.Context, id kernel.ID, visitor func(*models.Board) error) error {
	query := OneQuery{ID: id.String()}

	return r.repository.Execute(ctx, query.Operation(ctx, visitor))
}

// FindByOwner method
func (r Repository) FindByOwner(ctx context.Context, owner string, visitor func(*models.BoardListModel) error) error {
	query := ListQuery{Owner: owner}

	return r.repository.Execute(ctx, query.Operation(ctx, visitor))
}

// Create board
func (r Repository) Create(ctx context.Context, board *models.Board) error {
	command := CreateCommand{Board: board}

	return r.repository.Execute(ctx, command.Operation(ctx))
}

// Remove board
func (r Repository) Remove(ctx context.Context, id string) error {
	command := RemoveCommand{ID: id}

	return r.repository.Execute(ctx, command.Operation(ctx))
}

// Update board
func (r Repository) Update(ctx context.Context, id string, field string, value interface{}) error {
	command := UpdateCommand{ID: id, Field: field, Value: value}

	return r.repository.Execute(ctx, command.Operation(ctx))
}

// Attach child to board
func (r Repository) Attach(ctx context.Context, id string, childID string) error {
	command := AttachCommand{ID: id, ChildID: childID}

	return r.repository.Execute(ctx, command.Operation(ctx))
}

// Detach child from board
func (r Repository) Detach(ctx context.Context, id string, childID string) error {
	command := DetachCommand{ID: id, ChildID: childID}

	return r.repository.Execute(ctx, command.Operation(ctx))
}
