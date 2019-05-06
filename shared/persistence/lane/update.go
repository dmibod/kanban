package lane

import (
	"context"

	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// UpdateCommand type
type UpdateCommand struct {
	BoardID string
	ID      string
	Field   string
	Value   interface{}
}

// Operation to update lane
func (command UpdateCommand) Operation(ctx context.Context) mongo.Operation {
	return func(col *mgo.Collection) error {
		return mongo.Update(
			ctx,
			col,
			bson.M{"_id": bson.ObjectIdHex(command.BoardID), "lanes._id": bson.ObjectIdHex(command.ID)},
			mongo.Set("lanes.$."+command.Field, command.Value))
	}
}
