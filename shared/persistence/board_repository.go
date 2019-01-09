package persistence

import (
	"gopkg.in/mgo.v2/bson"

	"github.com/dmibod/kanban/shared/tools/db"
)

// BoardEntity entity
type BoardEntity struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	Layout   string        `bson:"layout"`
	Name     string        `bson:"name"`
	Children []string      `bson:"children"`
	Owner    string        `bson:"owner"`
	Shared   bool          `bson:"shared"`
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
