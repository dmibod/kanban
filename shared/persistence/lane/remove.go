package lane

import (
	"context"

	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"gopkg.in/mgo.v2"
)

// RemoveCommand type
type RemoveCommand struct {
	BoardID string
	ID      string
}

// Operation to remove lane
func (command RemoveCommand) Operation(ctx context.Context) mongo.Operation {
	return func(col *mgo.Collection) error {
		return mongo.Update(
			ctx,
			col,
			mongo.FromID(command.BoardID),
			mongo.RemoveFromSet("lanes", mongo.FromID(command.ID)))
	}
}
