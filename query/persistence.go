package query

import (
	"github.com/dmibod/kanban/tools/db"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

type DbCard struct {
	Id   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"name"`
}

func CreateCardRepository(repoFactory db.RepoFactory) db.Repository {

	instance := func() interface{} {
		return &DbCard{}
	}

	return repoFactory.Create("cards", instance)
}
