package lane

import (
	"context"
	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/persistence/models"
	"github.com/dmibod/kanban/shared/tools/db/mongo"
)

// Repository for lanes
type Repository struct {
	repository *mongo.Repository
}

// CreateRepository method
func CreateRepository(r *mongo.Repository) Repository {
	return Repository{repository: r}
}

// FindByID method
func (r Repository) FindByID(ctx context.Context, id kernel.MemberID, visitor func(*models.Lane) error) error {
	query := OneQuery{BoardID: id.SetID.String(), ID: id.ID.String()}

	return r.repository.Execute(ctx, query.Operation(ctx, visitor))
}

// FindByParent method
func (r Repository) FindByParent(ctx context.Context, id kernel.MemberID, visitor func(*models.LaneListModel) error) error {
	query := ListQuery{BoardID: id.SetID.String(), ParentID: id.ID.String()}

	return r.repository.Execute(ctx, query.Operation(ctx, visitor))
}

// Create lane
func (r Repository) Create(ctx context.Context, boardID string, lane *models.Lane) error {
	command := CreateCommand{BoardID: boardID, Lane: lane}

	return r.repository.Execute(ctx, command.Operation(ctx))
}

// Remove lane
func (r Repository) Remove(ctx context.Context, boardID string, id string) error {
	command := RemoveCommand{BoardID: boardID, ID: id}

	return r.repository.Execute(ctx, command.Operation(ctx))
}

// Update lane
func (r Repository) Update(ctx context.Context, boardID string, id string, field string, value interface{}) error {
	command := UpdateCommand{BoardID: boardID, ID: id, Field: field, Value: value}

	return r.repository.Execute(ctx, command.Operation(ctx))
}

// Attach child to lane
func (r Repository) Attach(ctx context.Context, boardID string, id string, childID string) error {
	command := AttachCommand{BoardID: boardID, ID: id, ChildID: childID}

	return r.repository.Execute(ctx, command.Operation(ctx))
}

// Detach child from lane
func (r Repository) Detach(ctx context.Context, boardID string, id string, childID string) error {
	command := DetachCommand{BoardID: boardID, ID: id, ChildID: childID}

	return r.repository.Execute(ctx, command.Operation(ctx))
}
