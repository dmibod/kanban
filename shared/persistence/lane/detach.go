package lane

import (
	"context"

	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// DetachCommand type
type DetachCommand struct {
	BoardID string
	ID      string
	ChildID string
}

// Operation to detach child from lane
func (command DetachCommand) Operation(ctx context.Context) mongo.Operation {
	return func(col *mgo.Collection) error {
		return mongo.Update(
			ctx,
			col,
			bson.M{"_id": bson.ObjectIdHex(command.BoardID), "lanes._id": bson.ObjectIdHex(command.ID)},
			mongo.RemoveFromSet("lanes.$.children", bson.ObjectIdHex(command.ChildID)))
	}
}
