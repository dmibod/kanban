package persistence

import (
	"context"
	b "github.com/dmibod/kanban/shared/persistence/board"
	c "github.com/dmibod/kanban/shared/persistence/card"
	"github.com/dmibod/kanban/shared/persistence/models"
	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"gopkg.in/mgo.v2/bson"

	"github.com/dmibod/kanban/shared/domain/board"
	"github.com/dmibod/kanban/shared/domain/card"
	"github.com/dmibod/kanban/shared/domain/lane"
	"github.com/dmibod/kanban/shared/kernel"
)

// Repository type
type Repository struct {
	repository *mongo.Repository
}

// FindBoardByID method
func (r Repository) FindBoardByID(ctx context.Context, id kernel.ID, visitor func(*models.Board) error) error {
	query := b.OneQuery{ID: id.String()}

	return r.repository.Execute(ctx, query.Operation(ctx, visitor))
}

// FindBoardsByOwner method
func (r Repository) FindBoardsByOwner(ctx context.Context, owner string, visitor func(*models.BoardListModel) error) error {
	query := b.ListQuery{Owner: owner}

	return r.repository.Execute(ctx, query.Operation(ctx, visitor))
}

// FindLaneByID method
func (r Repository) FindLaneByID(ctx context.Context, id kernel.MemberID, visitor func(*models.Lane) error) error {
	query := LaneQuery{BoardID: id.SetID.String(), ID: id.ID.String()}

	return r.repository.Execute(ctx, query.Operation(ctx, visitor))
}

// FindLanesByParent method
func (r Repository) FindLanesByParent(ctx context.Context, id kernel.MemberID, visitor func(*models.LaneListModel) error) error {
	query := LaneListQuery{BoardID: id.SetID.String(), ParentID: id.ID.String()}

	return r.repository.Execute(ctx, query.Operation(ctx, visitor))
}

// FindCardByID method
func (r Repository) FindCardByID(ctx context.Context, id kernel.MemberID, visitor func(*models.Card) error) error {
	query := c.OneQuery{BoardID: id.SetID.String(), ID: id.ID.String()}

	return r.repository.Execute(ctx, query.Operation(ctx, visitor))
}

// FindCardsByParent method
func (r Repository) FindCardsByParent(ctx context.Context, id kernel.MemberID, visitor func(*models.Card) error) error {
	query := c.ListQuery{BoardID: id.SetID.String(), LaneID: id.ID.String()}

	return r.repository.Execute(ctx, query.Operation(ctx, visitor))
}

// Handle domain event
func (r Repository) Handle(ctx context.Context, event interface{}) error {
	if event == nil {
		return nil
	}

	switch e := event.(type) {
	case board.CreatedEvent:
		return r.createBoard(ctx, r.mapBoard(&e.Entity))
	case board.DeletedEvent:
		return r.removeBoard(ctx, e.ID.String())
	case board.NameChangedEvent:
		return r.updateBoard(ctx, e.ID.String(), "name", e.NewValue)
	case board.DescriptionChangedEvent:
		return r.updateBoard(ctx, e.ID.String(), "description", e.NewValue)
	case board.LayoutChangedEvent:
		return r.updateBoard(ctx, e.ID.String(), "layout", e.NewValue)
	case board.SharedChangedEvent:
		return r.updateBoard(ctx, e.ID.String(), "shared", e.NewValue)
	case board.ChildAppendedEvent:
		return r.attachToBoard(ctx, e.ID.SetID.String(), e.ID.ID.String())
	case board.ChildRemovedEvent:
		return r.detachFromBoard(ctx, e.ID.SetID.String(), e.ID.ID.String())
	case lane.CreatedEvent:
		return r.createLane(ctx, e.ID.SetID.String(), r.mapLane(&e.Entity))
	case lane.DeletedEvent:
		return r.removeLane(ctx, e.ID.SetID.String(), e.ID.ID.String())
	case lane.NameChangedEvent:
		return r.updateLane(ctx, e.ID.SetID.String(), e.ID.ID.String(), "name", e.NewValue)
	case lane.DescriptionChangedEvent:
		return r.updateLane(ctx, e.ID.SetID.String(), e.ID.ID.String(), "description", e.NewValue)
	case lane.LayoutChangedEvent:
		return r.updateLane(ctx, e.ID.SetID.String(), e.ID.ID.String(), "layout", e.NewValue)
	case lane.ChildAppendedEvent:
		return r.attachToLane(ctx, e.ID.SetID.String(), e.ID.ID.String(), e.ChildID.String())
	case lane.ChildRemovedEvent:
		return r.detachFromLane(ctx, e.ID.SetID.String(), e.ID.ID.String(), e.ChildID.String())
	case card.CreatedEvent:
		return r.createCard(ctx, e.ID.SetID.String(), r.mapCard(&e.Entity))
	case card.DeletedEvent:
		return r.removeCard(ctx, e.ID.SetID.String(), e.ID.ID.String())
	case card.NameChangedEvent:
		return r.updateCard(ctx, e.ID.SetID.String(), e.ID.ID.String(), "name", e.NewValue)
	case card.DescriptionChangedEvent:
		return r.updateCard(ctx, e.ID.SetID.String(), e.ID.ID.String(), "description", e.NewValue)
	default:
		return nil
	}
}

func (r Repository) mapBoard(entity *board.Entity) *models.Board {
	return &models.Board{
		ID:          bson.ObjectIdHex(entity.ID.String()),
		Owner:       entity.Owner,
		Name:        entity.Name,
		Description: entity.Description,
		Layout:      entity.Layout,
		Shared:      entity.Shared,
		Children:    []bson.ObjectId{},
		Lanes:       []models.Lane{},
		Cards:       []models.Card{},
	}
}

func (r Repository) createBoard(ctx context.Context, board *models.Board) error {
	command := b.CreateCommand{Board: board}

	return r.repository.Execute(ctx, command.Operation(ctx))
}

func (r Repository) removeBoard(ctx context.Context, id string) error {
	command := b.RemoveCommand{ID: id}

	return r.repository.Execute(ctx, command.Operation(ctx))
}

func (r Repository) updateBoard(ctx context.Context, id string, field string, value interface{}) error {
	command := b.UpdateCommand{ID: id, Field: field, Value: value}

	return r.repository.Execute(ctx, command.Operation(ctx))
}

func (r Repository) attachToBoard(ctx context.Context, id string, childID string) error {
	command := b.AttachCommand{ID: id, ChildID: childID}

	return r.repository.Execute(ctx, command.Operation(ctx))
}

func (r Repository) detachFromBoard(ctx context.Context, id string, childID string) error {
	command := b.DetachCommand{ID: id, ChildID: childID}

	return r.repository.Execute(ctx, command.Operation(ctx))
}

func (r Repository) mapLane(entity *lane.Entity) *models.Lane {
	return &models.Lane{
		ID:          bson.ObjectIdHex(entity.ID.ID.String()),
		Kind:        entity.Kind,
		Name:        entity.Name,
		Description: entity.Description,
		Layout:      entity.Layout,
		Children:    []bson.ObjectId{},
	}
}

func (r Repository) createLane(ctx context.Context, boardID string, lane *models.Lane) error {
	command := CreateLaneCommand{BoardID: boardID, Lane: lane}

	return r.repository.Execute(ctx, command.Operation(ctx))
}

func (r Repository) removeLane(ctx context.Context, boardID string, id string) error {
	command := RemoveLaneCommand{BoardID: boardID, ID: id}

	return r.repository.Execute(ctx, command.Operation(ctx))
}

func (r Repository) updateLane(ctx context.Context, boardID string, id string, field string, value interface{}) error {
	command := UpdateLaneCommand{BoardID: boardID, ID: id, Field: field, Value: value}

	return r.repository.Execute(ctx, command.Operation(ctx))
}

func (r Repository) attachToLane(ctx context.Context, boardID string, id string, childID string) error {
	command := AttachToLaneCommand{BoardID: boardID, ID: id, ChildID: childID}

	return r.repository.Execute(ctx, command.Operation(ctx))
}

func (r Repository) detachFromLane(ctx context.Context, boardID string, id string, childID string) error {
	command := DetachFromLaneCommand{BoardID: boardID, ID: id, ChildID: childID}

	return r.repository.Execute(ctx, command.Operation(ctx))
}

func (r Repository) mapCard(entity *card.Entity) *models.Card {
	return &models.Card{
		ID:          bson.ObjectIdHex(entity.ID.ID.String()),
		Name:        entity.Name,
		Description: entity.Description,
	}
}

func (r Repository) createCard(ctx context.Context, boardID string, card *models.Card) error {
	command := c.CreateCommand{BoardID: boardID, Card: card}

	return r.repository.Execute(ctx, command.Operation(ctx))
}

func (r Repository) removeCard(ctx context.Context, boardID string, id string) error {
	command := c.RemoveCommand{BoardID: boardID, ID: id}

	return r.repository.Execute(ctx, command.Operation(ctx))
}

func (r Repository) updateCard(ctx context.Context, boardID string, id string, field string, value interface{}) error {
	command := c.UpdateCommand{BoardID: boardID, ID: id, Field: field, Value: value}

	return r.repository.Execute(ctx, command.Operation(ctx))
}
