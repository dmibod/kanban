package persistence

import (
	"context"
	"github.com/dmibod/kanban/shared/tools/db"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

// CardEntity maps card to/from mongo db
type CardEntity struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"name"`
}

// CreateCardRepository creates new cards repository
func CreateCardRepository(c context.Context, f db.Factory) db.Repository {
	instance := func() interface{} {
		return &CardEntity{}
	}

	return f.CreateRepository(c, "cards", instance)
}
