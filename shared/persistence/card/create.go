package card

import (
	"context"
	"github.com/dmibod/kanban/shared/persistence/models"

	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"gopkg.in/mgo.v2"
)

// CreateCommand type
type CreateCommand struct {
	BoardID string
	Card    *models.Card
}

// Operation to create card
func (command CreateCommand) Operation(ctx context.Context) mongo.Operation {
	return func(col *mgo.Collection) error {
		return mongo.Update(
			ctx,
			col,
			mongo.FromID(command.BoardID),
			mongo.AddToSet("cards", command.Card))
	}
}
