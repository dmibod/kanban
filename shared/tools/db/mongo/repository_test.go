package mongo_test

import (
	"context"
	"testing"

	"github.com/dmibod/kanban/shared/tools/test"

	"gopkg.in/mgo.v2/bson"

	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"github.com/dmibod/kanban/shared/tools/logger/noop"
)

const enable = false

func TestDB(t *testing.T) {
	if enable {
		testDB(t)
	}
}

func testDB(t *testing.T) {
	instance := func() interface{} {
		entity := struct {
			ID   bson.ObjectId `bson:"_id,omitempty"`
			Name string        `bson:"name"`
		}{}
		return &entity
	}

	identity := func(entity interface{}) string {
		return "5c16dd24c7ee6e5dcf626266"
	}

	r := mongo.CreateFactory(
		"kanban",
		mongo.CreateExecutor(),
		&noop.Logger{}).CreateRepository("cards", instance, identity)

	_, err := r.FindByID(context.TODO(), "5c16dd24c7ee6e5dcf626266")
	test.Ok(t, err)
}
