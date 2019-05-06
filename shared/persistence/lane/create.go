package lane

import (
	"context"
	"github.com/dmibod/kanban/shared/persistence/models"

	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"gopkg.in/mgo.v2"
)

// CreateCommand type
type CreateCommand struct {
	BoardID string
	Lane    *models.Lane
}

// Operation to create lane
func (command CreateCommand) Operation(ctx context.Context) mongo.Operation {
	return func(col *mgo.Collection) error {
		return mongo.Update(
			ctx,
			col,
			mongo.FromID(command.BoardID),
			mongo.AddToSet("lanes", command.Lane))
	}
}
