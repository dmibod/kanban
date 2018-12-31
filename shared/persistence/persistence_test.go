package persistence_test

import (
	"context"
	"testing"

	"github.com/dmibod/kanban/shared/tools/logger/console"
	"github.com/dmibod/kanban/shared/tools/test"

	"github.com/dmibod/kanban/shared/persistence"
)

const enable = false

func TestDB(t *testing.T) {
	if enable {
		testDB(t)
	}
}

func testDB(t *testing.T) {
	l := console.New(console.WithDebug(true))
	f := persistence.CreateFactory(persistence.CreateService(l), l)
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
