package persistence_test

import (
	"context"
	"github.com/dmibod/kanban/shared/tools/logger/console"
	"testing"

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

	id, createErr := r.Create(c, &persistence.CardEntity{Name: "Sample"})
	ok(t, createErr)

	e, getErr := r.FindByID(c, id)
	ok(t, getErr)

	entity := e.(*persistence.CardEntity)
	entity.Name = entity.Name + "!"
	updErr := r.Update(c, entity)
	ok(t, updErr)

	e, getErr = r.FindByID(c, id)
	ok(t, getErr)
	entity = e.(*persistence.CardEntity)

	exp := "Sample!"
	act := entity.Name
	assertf(t, act == exp, "Wrong value:\nwant: %v\ngot: %v\n", exp, act)

	remErr := r.Remove(c, id)
	ok(t, remErr)

	e, getErr = r.FindByID(c, id)
	assert(t, e == nil, "Entity must be nil")
	assert(t, getErr != nil, "Entity must not be found")
}

func ok(t *testing.T, e error) {
	if e != nil {
		t.Fatal(e)
	}
}

func assert(t *testing.T, exp bool, msg string) {
	if !exp {
		t.Fatal(msg)
	}
}

func assertf(t *testing.T, exp bool, f string, v ...interface{}) {
	if !exp {
		t.Fatalf(f, v...)
	}
}
