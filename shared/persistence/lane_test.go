package persistence_test

import (
	"context"
	"testing"

	"github.com/dmibod/kanban/shared/tools/logger/console"
	"github.com/dmibod/kanban/shared/tools/test"

	"github.com/dmibod/kanban/shared/persistence"
)

const enableLaneTest = false

func TestLanes(t *testing.T) {
	if enableLaneTest {
		testLanes(t)
	}
}

func testLanes(t *testing.T) {
	l := console.New(console.WithDebug(true))
	s, _ := persistence.CreateService(l)
	f := persistence.CreateFactory(s, l)
	r := persistence.CreateLaneRepository(f)
	c := context.TODO()

	id, err := r.Create(c, &persistence.LaneEntity{Name: "Sample", Layout: persistence.HLayout, Type: persistence.LType, Children: []string{"dummy"}})
	test.Ok(t, err)

	found, err := r.FindByID(c, id)
	test.Ok(t, err)

	entity := found.(*persistence.LaneEntity)
	entity.Name = entity.Name + "!"
	test.Ok(t, r.Update(c, entity))

	found, err = r.FindByID(c, id)
	test.Ok(t, err)
	entity = found.(*persistence.LaneEntity)

	test.AssertExpAct(t, "Sample!", entity.Name)
	test.AssertExpAct(t, persistence.HLayout, entity.Layout)
	test.AssertExpAct(t, persistence.LType, entity.Type)
	test.AssertExpAct(t, 1, len(entity.Children))
	test.AssertExpAct(t, "dummy", entity.Children[0])

	test.Ok(t, r.Remove(c, id))

	found, err = r.FindByID(c, id)
	test.Assert(t, found == nil, "Entity must be nil")
	test.Assert(t, err != nil, "Entity must not be found")
}
