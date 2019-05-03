// +build integration

package mongo_test

import (
	"context"
	"github.com/dmibod/kanban/shared/tools/logger/console"
	"testing"

	"github.com/dmibod/kanban/shared/tools/test"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/dmibod/kanban/shared/tools/db/mongo"
)

type TestEntity struct {
	ID   bson.ObjectId `bson:"_id,omitempty"`
	Name string        `bson:"name"`
}

type testEntityRepository struct {
}

func (*testEntityRepository) CreateInstance() interface{} {
	return &TestEntity{}
}

func (*testEntityRepository) GetID(entity interface{}) string {
	return entity.(*TestEntity).ID.Hex()
}

func TestRepository(t *testing.T) {

	c := context.TODO()
	l := console.New(console.WithDebug(true))
	s := mongo.CreateSessionFactory(mongo.WithLogger(l))
	p := mongo.CreateSessionProvider(s, l)
	e := mongo.CreateExecutor(p, l)
	f := mongo.CreateRepositoryFactory("test", e, l)
	r := f.CreateRepository("test")

	// Find and remove all
	test.Ok(t, r.Execute(c, func(col *mgo.Collection) error {
		return mongo.QueryList(c, col, nil, &TestEntity{}, func(e interface{}) error {
			remove, ok := e.(*TestEntity)
			test.Assert(t, ok, "Wrong type")
			test.Ok(t, r.Remove(c, remove.ID.Hex()))
			return nil
		})
	}))

	// Check count=0
	test.Ok(t, r.Execute(c, func(col *mgo.Collection) error {
		return mongo.QueryCount(c, col, nil, func(count int) error {
			test.AssertExpAct(t, 0, count)
			return nil
		})
	}))

	// Create
	id, err := r.Create(c, &TestEntity{Name: "Test"})
	test.Ok(t, err)
	// Check created
	test.Ok(t, r.Execute(c, func(col *mgo.Collection) error {
		return mongo.QueryOne(c, col, mongo.FromID(id), &TestEntity{}, func(e interface{}) error {
			_, ok := e.(*TestEntity)
			test.Assert(t, ok, "Wrong type")
			return nil
		})
	}))

	// Update
	test.Ok(t, r.Execute(c, func(col *mgo.Collection) error {
		return mongo.QueryOne(c, col, mongo.FromID(id), &TestEntity{}, func(e interface{}) error {
			found, ok := e.(*TestEntity)
			test.Assert(t, ok, "Wrong type")
			test.Ok(t, r.Update(c, found.ID.Hex(), mongo.Set("name", "Test!")))
			return nil
		})
	}))

	// Check updated
	test.Ok(t, r.Execute(c, func(col *mgo.Collection) error {
		return mongo.QueryOne(c, col, mongo.FromID(id), &TestEntity{}, func(e interface{}) error {
			found, ok := e.(*TestEntity)
			test.Assert(t, ok, "Wrong type")
			test.AssertExpAct(t, "Test!", found.Name)
			return nil
		})
	}))

	// Remove
	test.Ok(t, r.Remove(c, id))
	// Check removed
	test.Ok(t, r.Execute(c, func(col *mgo.Collection) error {
		return mongo.QueryCount(c, col, nil, func(count int) error {
			test.AssertExpAct(t, 0, count)
			return nil
		})
	}))

	// Create 2 entities
	_, err = r.Create(c, &TestEntity{Name: "Test1"})
	test.Ok(t, err)
	_, err = r.Create(c, &TestEntity{Name: "Test2"})
	test.Ok(t, err)

	// Check count=2
	test.Ok(t, r.Execute(c, func(col *mgo.Collection) error {
		return mongo.QueryCount(c, col, nil, func(count int) error {
			test.AssertExpAct(t, 2, count)
			return nil
		})
	}))

	// Count by criteria
	test.Ok(t, r.Execute(c, func(col *mgo.Collection) error {
		return mongo.QueryCount(c, col, bson.M{"name": "Test1"}, func(count int) error {
			test.AssertExpAct(t, 1, count)
			return nil
		})
	}))
}
