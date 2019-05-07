// +build integration

package lane_test

import (
	"context"
	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/services"
	"github.com/dmibod/kanban/shared/services/board"
	"github.com/dmibod/kanban/shared/services/lane"
	setup "github.com/dmibod/kanban/shared/services/test"
	"github.com/dmibod/kanban/shared/tools/test"
	"testing"
)

func Test(t *testing.T) {
	setup.WithLane(t, func(c context.Context, f *services.ServiceFactory, b *board.Model, l *lane.Model) {
		testLanes(t, c, f, b, l)
	})
}

func testLanes(t *testing.T, c context.Context, f *services.ServiceFactory, b *board.Model, l *lane.Model) {
	s := f.CreateLaneService()

	//Update lane
	test.Ok(t, s.Name(c, l.ID.WithSet(b.ID), l.Name+"_updated"))
	test.Ok(t, s.Describe(c, l.ID.WithSet(b.ID), l.Description+"_updated"))
	test.Ok(t, s.Layout(c, l.ID.WithSet(b.ID), kernel.HLayout))

	//Confirm not attached
	list, err := s.GetByBoardID(c, b.ID)
	test.Ok(t, err)
	test.AssertExpAct(t, 0, len(list))

	bs := f.CreateBoardService()

	//Attach lane
	test.Ok(t, bs.AppendLane(c, l.ID.WithSet(b.ID)))

	//Confirm attached
	list, err = s.GetByBoardID(c, b.ID)
	test.Ok(t, err)
	test.AssertExpAct(t, 1, len(list))
	test.AssertExpAct(t, l.ID, list[0].ID)

	//Detach lane
	test.Ok(t, bs.ExcludeLane(c, l.ID.WithSet(b.ID)))

	//Confirm detached
	list, err = s.GetByBoardID(c, b.ID)
	test.Ok(t, err)
	test.AssertExpAct(t, 0, len(list))

	testNestedLanes(t, c, f, b, l)
}

func testNestedLanes(t *testing.T, c context.Context, f *services.ServiceFactory, b *board.Model, l *lane.Model) {
	s := f.CreateLaneService()

	model := &lane.CreateModel{
		Type:        kernel.CKind,
		Name:        "test_nested_name",
		Description: "test_nested_description",
		Layout:      kernel.VLayout}

	//Create lane
	id, err := s.Create(c, b.ID, model)
	test.Ok(t, err)

	//Confirm created
	n, err := s.GetByID(c, id.WithSet(b.ID))
	test.Ok(t, err)
	test.Assert(t, n != nil, "Lane should be found")
	test.AssertExpAct(t, id, n.ID)

	//Confirm not attached
	list, err := s.GetByLaneID(c, l.ID.WithSet(b.ID))
	test.Ok(t, err)
	test.AssertExpAct(t, 0, len(list))

	//Attach lane
	test.Ok(t, s.AppendChild(c, l.ID.WithSet(b.ID), n.ID))

	//Confirm attached
	list, err = s.GetByLaneID(c, l.ID.WithSet(b.ID))
	test.Ok(t, err)
	test.AssertExpAct(t, 1, len(list))
	test.AssertExpAct(t, n.ID, list[0].ID)

	//Detach lane
	test.Ok(t, s.ExcludeChild(c, l.ID.WithSet(b.ID), n.ID))

	//Confirm detached
	list, err = s.GetByLaneID(c, l.ID.WithSet(b.ID))
	test.Ok(t, err)
	test.AssertExpAct(t, 0, len(list))
}
