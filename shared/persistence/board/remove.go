package board

import (
	"context"

	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"gopkg.in/mgo.v2"
)

// RemoveCommand type
type RemoveCommand struct {
	ID string
}

// Operation to remove card
func (command RemoveCommand) Operation(ctx context.Context) mongo.Operation {
	return func(col *mgo.Collection) error {
		return mongo.Remove(ctx, col, mongo.FromID(command.ID))
	}
}
