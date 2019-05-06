package board

import (
	"context"
	err "github.com/dmibod/kanban/shared/persistence/error"
	"github.com/dmibod/kanban/shared/persistence/models"
	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// OneQuery type
type OneQuery struct {
	ID string
}

// Operation to query board
func (query OneQuery) Operation(ctx context.Context, visitor func(*models.Board) error) mongo.Operation {
	return func(col *mgo.Collection) error {
		return mongo.QueryOne(ctx, col, query.criteria(), &models.Board{}, func(entity interface{}) error {
			if board, ok := entity.(*models.Board); ok {
				return visitor(board)
			}

			return err.ErrInvalidType
		})
	}
}

func (query OneQuery) criteria() bson.M {
	return mongo.FromID(query.ID)
}
