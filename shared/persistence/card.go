package persistence

import (
	"gopkg.in/mgo.v2/bson"

	"github.com/dmibod/kanban/shared/tools/db"
)

// CardEntity maps card to/from mongo db
type CardEntity struct {
	ID   bson.ObjectId `bson:"_id,omitempty"`
	Name string        `bson:"name"`
}

// CreateCardRepository creates new cards repository
func CreateCardRepository(f db.RepositoryFactory) db.Repository {
	instance := func() interface{} {
		return &CardEntity{}
	}
	identity := func(entity interface{}) string {
		return entity.(*CardEntity).ID.Hex()
	}
	return f.CreateRepository("cards", instance, identity)
}
