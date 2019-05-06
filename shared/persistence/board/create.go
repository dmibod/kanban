package board

import (
	"context"
	"github.com/dmibod/kanban/shared/persistence/models"

	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"gopkg.in/mgo.v2"
)

// CreateCommand type
type CreateCommand struct {
	Board *models.Board
}

// Operation to create board
func (command CreateCommand) Operation(ctx context.Context) mongo.Operation {
	return func(col *mgo.Collection) error {
		return mongo.Insert(ctx, col, command.Board)
	}
}
