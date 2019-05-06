package board

import (
	"context"
	err "github.com/dmibod/kanban/shared/persistence/error"
	"github.com/dmibod/kanban/shared/persistence/models"
	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// ListQuery type
type ListQuery struct {
	Owner string
}

// Operation to query boards list
func (query ListQuery) Operation(ctx context.Context, visitor func(*models.BoardListModel) error) mongo.Operation {
	return func(col *mgo.Collection) error {
		return mongo.QueryList(ctx, col, query.criteria(), &models.BoardListModel{}, func(entity interface{}) error {
			if board, ok := entity.(*models.BoardListModel); ok {
				return visitor(board)
			}

			return err.ErrInvalidType
		})
	}
}

func (query ListQuery) criteria() bson.M {
	if query.Owner == "" {
		return bson.M{"shared": true}
	}

	return bson.M{"$or": []bson.M{bson.M{"shared": true}, bson.M{"owner": query.Owner}}}
}
