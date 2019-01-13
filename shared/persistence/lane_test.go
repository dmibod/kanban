// +build integration

package persistence_test

import (
	"context"
	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"testing"

	"github.com/dmibod/kanban/shared/tools/logger/console"
	"github.com/dmibod/kanban/shared/tools/test"

	"github.com/dmibod/kanban/shared/persistence"
)

func TestLanes(t *testing.T) {
	l := console.New(console.WithDebug(true))
	s := persistence.CreateSessionFactory(mongo.CreateSessionFactory(mongo.WithLogger(l)), l)
	p := mongo.CreateSessionProvider(s, l)
	e := persistence.CreateOperationExecutor(p, l)
	f := persistence.CreateRepositoryFactory(e, l)
	r := persistence.CreateLaneRepository(f)
	c := context.TODO()

	id, err := r.Create(c, &persistence.LaneEntity{Name: "Sample", Layout: kernel.HLayout, Kind: kernel.LKind, Children: []string{"dummy"}})
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
	test.AssertExpAct(t, kernel.HLayout, entity.Layout)
	test.AssertExpAct(t, kernel.LKind, entity.Kind)
	test.AssertExpAct(t, 1, len(entity.Children))
	test.AssertExpAct(t, "dummy", entity.Children[0])

	test.Ok(t, r.Remove(c, id))

	found, err = r.FindByID(c, id)
	test.Assert(t, found == nil, "Entity must be nil")
	test.Assert(t, err != nil, "Entity must not be found")
}
