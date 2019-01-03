package persistence_test

import (
	"context"
	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"testing"

	"github.com/dmibod/kanban/shared/tools/logger/console"
	"github.com/dmibod/kanban/shared/tools/test"

	"github.com/dmibod/kanban/shared/persistence"
)

const enableCardTest = false

func TestCards(t *testing.T) {
	if enableCardTest {
		testCards(t)
	}
}

func testCards(t *testing.T) {
	l := console.New(console.WithDebug(true))
	s := persistence.CreateSessionFactory(mongo.CreateSessionFactory(mongo.WithLogger(l)), l)
	p := mongo.CreateSessionProvider(s, l)
	e := persistence.CreateOperationExecutor(p, l)
	f := persistence.CreateRepositoryFactory(e, l)
	r := persistence.CreateCardRepository(f)
	c := context.TODO()

	id, err := r.Create(c, &persistence.CardEntity{Name: "Sample"})
	test.Ok(t, err)

	found, err := r.FindByID(c, id)
	test.Ok(t, err)

	entity := found.(*persistence.CardEntity)
	entity.Name = entity.Name + "!"
	test.Ok(t, r.Update(c, entity))

	found, err = r.FindByID(c, id)
	test.Ok(t, err)
	entity = found.(*persistence.CardEntity)

	exp := "Sample!"
	act := entity.Name
	test.AssertExpAct(t, exp, act)

	test.Ok(t, r.Remove(c, id))

	found, err = r.FindByID(c, id)
	test.Assert(t, found == nil, "Entity must be nil")
	test.Assert(t, err != nil, "Entity must not be found")
}
