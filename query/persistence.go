package query

import (
	"github.com/dmibod/kanban/tools/db"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

// CardEntity maps card to/from mongo db
type CardEntity struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"name"`
}

// CreateCardRepository creates new cards repository
func CreateCardRepository(repoFactory db.RepoFactory) db.Repository {

	instance := func() interface{} {
		return &CardEntity{}
	}

	return repoFactory.Create("cards", instance)
}
