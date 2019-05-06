package board

import (
	"context"

	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"gopkg.in/mgo.v2"
)

// UpdateCommand type
type UpdateCommand struct {
	ID    string
	Field string
	Value interface{}
}

// Operation to update board
func (command UpdateCommand) Operation(ctx context.Context) mongo.Operation {
	return func(col *mgo.Collection) error {
		return mongo.Update(
			ctx,
			col,
			mongo.FromID(command.ID),
			mongo.Set(command.Field, command.Value))
	}
}
