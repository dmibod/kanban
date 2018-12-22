package persistence_test

import (
	"context"
	"testing"

	"github.com/dmibod/kanban/shared/persistence"
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
	f := mongo.CreateFactory(
		mongo.WithDatabase("kanban"), 
		mongo.WithExecutor(persistence.CreateService(&noop.Logger{})))

	r := persistence.CreateCardRepository(context.TODO(), f)

	id, createErr := r.Create(&persistence.CardEntity{Name: "Sample"})
	ok(t, createErr)

	e, getErr := r.FindByID(id)
	ok(t, getErr)

	entity := e.(*persistence.CardEntity)
	entity.Name = entity.Name + "!"
	updErr := r.Update(entity)
	ok(t, updErr)

	e, getErr = r.FindByID(id)
	ok(t, getErr)
	entity = e.(*persistence.CardEntity)

	exp := "Sample!"
	act := entity.Name
	assertf(t, act == exp, "Wrong value:\nwant: %v\ngot: %v\n", act, exp)
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
