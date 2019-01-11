// +build integration

package mongo_test

import (
	"context"
	"github.com/dmibod/kanban/shared/tools/logger/console"
	"testing"

	"github.com/dmibod/kanban/shared/tools/test"

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
	r := f.CreateRepository("test", &testEntityRepository{})

	// Find and remove all
	err := r.Find(c, nil, func(e interface{}) error {
		remove, ok := e.(*TestEntity)
		test.Assert(t, ok, "Wrong type")
		test.Ok(t, r.Remove(c, remove.ID.Hex()))
		return nil
	})
	// Check count=0
	count, err := r.Count(c, nil)
	test.Ok(t, err)
	test.AssertExpAct(t, 0, count)

	// Create
	id, err := r.Create(c, &TestEntity{Name: "Test"})
	test.Ok(t, err)
	// Check created
	found, err := r.FindByID(c, id)
	test.Ok(t, err)
	entity, ok := found.(*TestEntity)
	test.Assert(t, ok, "Wrong type")

	// Update
	entity.Name = "Test!"
	test.Ok(t, r.Update(c, entity))
	// Check updated
	found, err = r.FindByID(c, entity.ID.Hex())
	test.Ok(t, err)
	entity, ok = found.(*TestEntity)
	test.Assert(t, ok, "Wrong type")
	test.AssertExpAct(t, "Test!", entity.Name)

	// Remove
	test.Ok(t, r.Remove(c, entity.ID.Hex()))
	// Check removed
	count, err = r.Count(c, nil)
	test.Ok(t, err)
	test.AssertExpAct(t, 0, count)

	// Create 2 entities
	_, err = r.Create(c, &TestEntity{Name: "Test1"})
	test.Ok(t, err)
	_, err = r.Create(c, &TestEntity{Name: "Test2"})
	test.Ok(t, err)

	// Check count=2
	count, err = r.Count(c, nil)
	test.Ok(t, err)
	test.AssertExpAct(t, 2, count)

	// Count by criteria
	count, err = r.Count(c, bson.M{"name": "Test1"})
	test.Ok(t, err)
	test.AssertExpAct(t, 1, count)
}
