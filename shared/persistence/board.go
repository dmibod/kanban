package persistence

import (
	"context"
	"github.com/dmibod/kanban/shared/persistence/models"

	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// BoardListQuery type
type BoardListQuery struct {
	Owner string
}

// Operation to query board list
func (query BoardListQuery) Operation(ctx context.Context, visitor func(*models.BoardListModel) error) mongo.Operation {
	return func(col *mgo.Collection) error {
		return mongo.QueryList(ctx, col, query.criteria(), &models.BoardListModel{}, func(entity interface{}) error {
			if board, ok := entity.(*models.BoardListModel); ok {
				return visitor(board)
			}

			return ErrInvalidType
		})
	}
}

func (query BoardListQuery) criteria() bson.M {
	if query.Owner == "" {
		return bson.M{"shared": true}
	}

	return bson.M{"$or": []bson.M{bson.M{"shared": true}, bson.M{"owner": query.Owner}}}
}

// BoardQuery type
type BoardQuery struct {
	ID string
}

// Operation to query board
func (query BoardQuery) Operation(ctx context.Context, visitor func(*models.Board) error) mongo.Operation {
	return func(col *mgo.Collection) error {
		return mongo.QueryOne(ctx, col, query.criteria(), &models.Board{}, func(entity interface{}) error {
			if board, ok := entity.(*models.Board); ok {
				return visitor(board)
			}

			return ErrInvalidType
		})
	}
}

func (query BoardQuery) criteria() bson.M {
	return mongo.FromID(query.ID)
}

// CreateBoardCommand type
type CreateBoardCommand struct {
	Board *models.Board
}

// Operation to create board
func (command CreateBoardCommand) Operation(ctx context.Context) mongo.Operation {
	return func(col *mgo.Collection) error {
		return mongo.Insert(ctx, col, command.Board)
	}
}

// RemoveBoardCommand type
type RemoveBoardCommand struct {
	ID string
}

// Operation to remove card
func (command RemoveBoardCommand) Operation(ctx context.Context) mongo.Operation {
	return func(col *mgo.Collection) error {
		return mongo.Remove(ctx, col, mongo.FromID(command.ID))
	}
}

// UpdateBoardCommand type
type UpdateBoardCommand struct {
	ID    string
	Field string
	Value interface{}
}

// Operation to update board
func (command UpdateBoardCommand) Operation(ctx context.Context) mongo.Operation {
	return func(col *mgo.Collection) error {
		return mongo.Update(ctx, col, mongo.FromID(command.ID), mongo.Set(command.Field, command.Value))
	}
}

// AttachToBoardCommand type
type AttachToBoardCommand struct {
	ID      string
	ChildID string
}

// Operation to attach to board
func (command AttachToBoardCommand) Operation(ctx context.Context) mongo.Operation {
	return func(col *mgo.Collection) error {
		return mongo.Update(ctx, col, mongo.FromID(command.ID), mongo.AddToSet("children", bson.ObjectIdHex(command.ChildID)))
	}
}

// DetachFromBoardCommand type
type DetachFromBoardCommand struct {
	ID      string
	ChildID string
}

// Operation to detach from board
func (command DetachFromBoardCommand) Operation(ctx context.Context) mongo.Operation {
	return func(col *mgo.Collection) error {
		return mongo.Update(ctx, col, mongo.FromID(command.ID), mongo.RemoveFromSet("children", bson.ObjectIdHex(command.ChildID)))
	}
}
