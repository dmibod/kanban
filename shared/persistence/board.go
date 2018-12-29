package persistence

import (
	"context"

	"gopkg.in/mgo.v2/bson"

	"github.com/dmibod/kanban/shared/tools/db"
)

// BoardEntity entity
type BoardEntity struct {
	ID   bson.ObjectId `bson:"_id,omitempty"`
	Name string        `bson:"name"`
}

// CreateBoardRepository creates new repository
func CreateBoardRepository(c context.Context, f db.RepositoryFactory) db.Repository {
	instance := func() interface{} {
		return &BoardEntity{}
	}
	identity := func(entity interface{}) string {
		return entity.(*BoardEntity).ID.Hex()
	}
	return f.CreateRepository(c, "boards", instance, identity)
}
