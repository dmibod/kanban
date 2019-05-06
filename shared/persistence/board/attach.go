package board

import (
	"context"

	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// AttachCommand type
type AttachCommand struct {
	ID      string
	ChildID string
}

// Operation to attach lane to board
func (command AttachCommand) Operation(ctx context.Context) mongo.Operation {
	return func(col *mgo.Collection) error {
		return mongo.Update(
			ctx,
			col,
			mongo.FromID(command.ID),
			mongo.AddToSet("children", bson.ObjectIdHex(command.ChildID)))
	}
}
