package board

import (
	"context"

	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// DetachCommand type
type DetachCommand struct {
	ID      string
	ChildID string
}

// Operation to detach lane from board
func (command DetachCommand) Operation(ctx context.Context) mongo.Operation {
	return func(col *mgo.Collection) error {
		return mongo.Update(
			ctx,
			col,
			mongo.FromID(command.ID),
			mongo.RemoveFromSet("children", bson.ObjectIdHex(command.ChildID)))
	}
}
