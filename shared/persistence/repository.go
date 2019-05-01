package persistence

import (
	"context"
	"gopkg.in/mgo.v2/bson"
	"github.com/dmibod/kanban/shared/tools/db/mongo"

	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/domain/board"
	"github.com/dmibod/kanban/shared/domain/lane"
	"github.com/dmibod/kanban/shared/domain/card"
)

// Repository type
type Repository struct {
	repository *mongo.Repository
}

// FindBoardByID method
func (r Repository) FindBoardByID(ctx context.Context, id kernel.ID, visitor func(*Board) error) error {
	query := BoardQuery{ID: id.String()}

	return r.repository.Execute(ctx, query.Operation(ctx, visitor))
}

// FindBoardsByOwner method
func (r Repository) FindBoardsByOwner(ctx context.Context, owner string, visitor func(*BoardListModel) error) error {
	query := BoardListQuery{Owner: owner}

	return r.repository.Execute(ctx, query.Operation(ctx, visitor))
}

// FindLaneByID method
func (r Repository) FindLaneByID(ctx context.Context, id kernel.MemberID, visitor func(*Lane) error) error {
	query := LaneQuery{BoardID: id.SetID.String(), ID: id.ID.String()}

	return r.repository.Execute(ctx, query.Operation(ctx, visitor))
}

// FindLanesByParent method
func (r Repository) FindLanesByParent(ctx context.Context, id kernel.MemberID, visitor func(*LaneListModel) error) error {
	query := LaneListQuery{BoardID: id.SetID.String(), ParentID: id.ID.String()}

	return r.repository.Execute(ctx, query.Operation(ctx, visitor))
}

// FindCardByID method
func (r Repository) FindCardByID(ctx context.Context, id kernel.MemberID, visitor func(*Card) error) error {
	query := CardQuery{BoardID: id.SetID.String(), ID: id.ID.String()}

	return r.repository.Execute(ctx, query.Operation(ctx, visitor))
}

// FindCardsByParent method
func (r Repository) FindCardsByParent(ctx context.Context, id kernel.MemberID, visitor func(*Card) error) error {
	query := CardListQuery{BoardID: id.SetID.String(), LaneID: id.ID.String()}

	return r.repository.Execute(ctx, query.Operation(ctx, visitor))
}

// Handle domain event
func (r Repository) Handle(ctx context.Context, event interface{}) {
	if event == nil {
		return
	}

	switch e := event.(type) {
	case board.CreatedEvent:
		r.createBoard(ctx, r.mapBoard(&e.Entity))
	case board.DeletedEvent:
		r.removeBoard(ctx, e.ID.String())
	case board.NameChangedEvent:
		r.updateBoard(ctx, e.ID.String(), "name", e.NewValue)
	case board.DescriptionChangedEvent:
		r.updateBoard(ctx, e.ID.String(), "description", e.NewValue)
	case board.LayoutChangedEvent:
		r.updateBoard(ctx, e.ID.String(), "layout", e.NewValue)
	case board.SharedChangedEvent:
		r.updateBoard(ctx, e.ID.String(), "shared", e.NewValue)
	case board.ChildAppendedEvent:
		r.attachToBoard(ctx, e.ID.SetID.String(), e.ID.ID.String())
	case board.ChildRemovedEvent:
		r.detachFromBoard(ctx, e.ID.SetID.String(), e.ID.ID.String())
	case lane.CreatedEvent:
		r.createLane(ctx, e.ID.SetID.String(), r.mapLane(&e.Entity))
	case lane.DeletedEvent:
		r.removeLane(ctx, e.ID.SetID.String(), e.ID.ID.String())
	case lane.NameChangedEvent:
		r.updateLane(ctx, e.ID.SetID.String(), e.ID.ID.String(), "name", e.NewValue)
	case lane.DescriptionChangedEvent:
		r.updateLane(ctx, e.ID.SetID.String(), e.ID.ID.String(), "description", e.NewValue)
	case lane.LayoutChangedEvent:
		r.updateLane(ctx, e.ID.SetID.String(), e.ID.ID.String(), "layout", e.NewValue)
	case lane.ChildAppendedEvent:
		r.attachToLane(ctx, e.ID.SetID.String(), e.ID.ID.String(), e.ChildID.String())
	case lane.ChildRemovedEvent:
		r.detachFromLane(ctx, e.ID.SetID.String(), e.ID.ID.String(), e.ChildID.String())
	case card.CreatedEvent:
		r.createCard(ctx, e.ID.SetID.String(), r.mapCard(&e.Entity))
	case card.DeletedEvent:
		r.removeCard(ctx, e.ID.SetID.String(), e.ID.ID.String())
	case card.NameChangedEvent:
		r.updateCard(ctx, e.ID.SetID.String(), e.ID.ID.String(), "name", e.NewValue)
	case card.DescriptionChangedEvent:
		r.updateCard(ctx, e.ID.SetID.String(), e.ID.ID.String(), "description", e.NewValue)
	default:
		return
	}
}

func (r Repository) mapBoard(entity *board.Entity) *Board {
	return &Board{
		ID:          bson.ObjectIdHex(entity.ID.String()),
		Owner:       entity.Owner,
		Name:        entity.Name,
		Description: entity.Description,
		Layout:      entity.Layout,
		Shared:      entity.Shared,
		Children:    []bson.ObjectId{},
		Lanes:       []Lane{},
		Cards:       []Card{},
	}
}

func (r Repository) createBoard(ctx context.Context, board *Board) error {
	command := CreateBoardCommand{Board: board}

	return r.repository.Execute(ctx, command.Operation(ctx))
}

func (r Repository) removeBoard(ctx context.Context, id string) error {
	command := RemoveBoardCommand{ID: id}

	return r.repository.Execute(ctx, command.Operation(ctx))
}

func (r Repository) updateBoard(ctx context.Context, id string, field string, value interface{}) error {
	command := UpdateBoardCommand{ID: id, Field: field, Value: value}

	return r.repository.Execute(ctx, command.Operation(ctx))
}

func (r Repository) attachToBoard(ctx context.Context, id string, childID string) error {
	command := AttachToBoardCommand{ID: id, ChildID: childID}

	return r.repository.Execute(ctx, command.Operation(ctx))
}

func (r Repository) detachFromBoard(ctx context.Context, id string, childID string) error {
	command := DetachFromBoardCommand{ID: id, ChildID: childID}

	return r.repository.Execute(ctx, command.Operation(ctx))
}

func (r Repository) mapLane(entity *lane.Entity) *Lane {
	return &Lane{
		ID:          bson.ObjectIdHex(entity.ID.ID.String()),
		Kind:        entity.Kind,
		Name:        entity.Name,
		Description: entity.Description,
		Layout:      entity.Layout,
		Children:    []bson.ObjectId{},
	}
}

func (r Repository) createLane(ctx context.Context, boardID string, lane *Lane) error {
	command := CreateLaneCommand{Lane: lane}

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

func (r Repository) mapCard(entity *card.Entity) *Card {
	return &Card{
		ID:          bson.ObjectIdHex(entity.ID.ID.String()),
		Name:        entity.Name,
		Description: entity.Description,
	}
}

func (r Repository) createCard(ctx context.Context, boardID string, card *Card) error {
	command := CreateCardCommand{Card: card}

	return r.repository.Execute(ctx, command.Operation(ctx))
}

func (r Repository) removeCard(ctx context.Context, boardID string, id string) error {
	command := RemoveCardCommand{BoardID: boardID, ID: id}

	return r.repository.Execute(ctx, command.Operation(ctx))
}

func (r Repository) updateCard(ctx context.Context, boardID string, id string, field string, value interface{}) error {
	command := UpdateCardCommand{BoardID: boardID, ID: id, Field: field, Value: value}

	return r.repository.Execute(ctx, command.Operation(ctx))
}
