package persistence

import (
	"gopkg.in/mgo.v2/bson"

	"github.com/dmibod/kanban/shared/tools/db"
)

// BoardEntity entity
type BoardEntity struct {
	ID   bson.ObjectId `bson:"_id,omitempty"`
	Name string        `bson:"name"`
}

// CreateBoardRepository creates new repository
func CreateBoardRepository(f db.RepositoryFactory) db.Repository {
	instance := func() interface{} {
		return &BoardEntity{}
	}
	identity := func(entity interface{}) string {
		return entity.(*BoardEntity).ID.Hex()
	}
	return f.CreateRepository("boards", instance, identity)
}
