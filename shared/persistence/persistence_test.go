package persistence_test

import (
	"context"
	"testing"

	"github.com/dmibod/kanban/shared/persistence"
	"github.com/dmibod/kanban/shared/tools/logger/noop"

	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

const enable = false

func TestDB(t *testing.T) {
	if enable {
		testDB(t)
	}
}

func testDB(t *testing.T) {
	i := func() interface{} {
		e := struct {
			ID   primitive.ObjectID `bson:"_id,omitempty"`
			Name string             `bson:"name"`
		}{}
		return &e
	}

	l := &noop.Logger{}
	s := persistence.CreateService(l)
	f := mongo.CreateFactory(mongo.WithDatabase("kanban"), mongo.WithExecutor(s))
	r := f.CreateRepository(context.TODO(), "cards", i)

	_, err := r.FindByID("5c16dd24c7ee6e5dcf626266")
	if err != nil {
		t.Fatal(err)
	}
}
