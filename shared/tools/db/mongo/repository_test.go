package mongo_test

import (
	"context"
	"testing"

	"gopkg.in/mgo.v2/bson"

	"github.com/dmibod/kanban/shared/tools/logger/noop"
	"github.com/dmibod/kanban/shared/tools/db/mongo"
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
		mongo.WithDatabase("kanban"),
		mongo.WithExecutor(mongo.CreateService(&noop.Logger{}))).CreateRepository(context.TODO(), "cards", instance, identity)

	_, err := r.FindByID("5c16dd24c7ee6e5dcf626266")
	ok(t, err)
}

func ok(t *testing.T, e error) {
	if e != nil {
		t.Fatal(e)
	}
}
