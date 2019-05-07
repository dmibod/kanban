// +build integration

package services_test

import (
	"github.com/dmibod/kanban/shared/services/card"
	"github.com/dmibod/kanban/shared/services/lane"
	"github.com/stretchr/testify/mock"
	"github.com/dmibod/kanban/shared/message/mocks"
	"github.com/dmibod/kanban/shared/persistence"
	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"github.com/dmibod/kanban/shared/tools/logger/noop"
	"github.com/dmibod/kanban/shared/tools/test"
	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/services/board"
	"github.com/dmibod/kanban/shared/services"
	"context"
	"testing"
)

func TestBoards(t *testing.T) {
	testBoard(t, func(c context.Context, f *services.ServiceFactory, b *board.Model) {
		testBoards(t, c, f, b)
	})
}

func testBoards(t *testing.T, c context.Context, f *services.ServiceFactory, b *board.Model) {
	//Update board
	s := f.CreateBoardService()

	test.Ok(t, s.Name(c, b.ID, b.Name + "_updated"))
	test.Ok(t, s.Describe(c, b.ID, b.Description + "_updated"))
	test.Ok(t, s.Layout(c, b.ID, kernel.VLayout))
	test.Ok(t, s.Share(c, b.ID, !b.Shared))

	list, err := s.GetByOwner(c, "test_owner")
	test.Ok(t, err)
	test.AssertExpAct(t, 1, len(list))
	test.AssertExpAct(t, b.ID, list[0].ID)
}

func TestLanes(t *testing.T) {
	testLane(t, func(c context.Context, f *services.ServiceFactory, b *board.Model, l *lane.Model) {
		testLanes(t, c, f, b, l)
	})
}

func testLanes(t *testing.T, c context.Context, f *services.ServiceFactory, b *board.Model, l *lane.Model) {
	s := f.CreateLaneService()

	//Update lane
	test.Ok(t, s.Name(c, l.ID.WithSet(b.ID), l.Name + "_updated"))
	test.Ok(t, s.Describe(c, l.ID.WithSet(b.ID), l.Description + "_updated"))
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

func TestCards(t *testing.T) {
	testCard(t, func(ctx context.Context, f *services.ServiceFactory, b *board.Model, l *lane.Model, c *card.Model) {
		testCards(t, ctx, f, b, l, c)
	})
}

func testCards(t *testing.T, ctx context.Context, f *services.ServiceFactory, b *board.Model, l *lane.Model, c *card.Model) {
	s := f.CreateCardService()

	//Update card
	test.Ok(t, s.Name(ctx, c.ID.WithSet(b.ID), c.Name + "_updated"))
	test.Ok(t, s.Describe(ctx, c.ID.WithSet(b.ID), c.Description + "_updated"))

	//Confirm not attached
	list, err := s.GetByLaneID(ctx, l.ID.WithSet(b.ID))
	test.Ok(t, err)
	test.AssertExpAct(t, 0, len(list))

	ls := f.CreateLaneService()

	//Attach card
	test.Ok(t, ls.AppendChild(ctx, l.ID.WithSet(b.ID), c.ID))

	//Confirm attached
	list, err = s.GetByLaneID(ctx, l.ID.WithSet(b.ID))
	test.Ok(t, err)
	test.AssertExpAct(t, 1, len(list))
	test.AssertExpAct(t, c.ID, list[0].ID)

	//Detach card
	test.Ok(t, ls.ExcludeChild(ctx, l.ID.WithSet(b.ID), c.ID))	

	//Confirm detached
	list, err = s.GetByLaneID(ctx, l.ID.WithSet(b.ID))
	test.Ok(t, err)
	test.AssertExpAct(t, 0, len(list))
}

func testCard(t *testing.T, h func(ctx context.Context, f *services.ServiceFactory, b *board.Model, l *lane.Model, c *card.Model)) {
	testLane(t, func(ctx context.Context, f *services.ServiceFactory, b *board.Model, l *lane.Model) {
		s := f.CreateCardService()

		model := &card.CreateModel{
			Name:        "test_name",
			Description: "test_description"}

		//Create card
		id, err := s.Create(ctx, b.ID, model)
		test.Ok(t, err)

		//Find by ID
		c, err := s.GetByID(ctx, id.WithSet(b.ID))
		test.Ok(t, err)
		test.Assert(t, c != nil, "Card should be found")
		test.AssertExpAct(t, id, c.ID)

		h(ctx, f, b, l, c)

		//Remove card
		test.Ok(t, s.Remove(ctx, id.WithSet(b.ID)))

		//Confirm removed
		c, err = s.GetByID(ctx, id.WithSet(b.ID))
		test.Assert(t, err != nil, "Card should not be found")
		test.Assert(t, c == nil, "Card should not be found")
	})
}

func testLane(t *testing.T, h func(c context.Context, f *services.ServiceFactory, b *board.Model, l *lane.Model)) {
	testBoard(t, func(c context.Context, f *services.ServiceFactory, b *board.Model) {
		s := f.CreateLaneService()

		model := &lane.CreateModel{
			Type:        kernel.LKind,
			Name:        "test_name",
			Description: "test_description",
			Layout:      kernel.VLayout}

		//Create lane
		id, err := s.Create(c, b.ID, model)
		test.Ok(t, err)

		//Confirm created
		l, err := s.GetByID(c, id.WithSet(b.ID))
		test.Ok(t, err)
		test.Assert(t, l != nil, "Lane should be found")
		test.AssertExpAct(t, id, l.ID)

		h(c, f, b, l)

		//Remove lane
		test.Ok(t, s.Remove(c, id.WithSet(b.ID)))

		//Confirm removed
		l, err = s.GetByID(c, id.WithSet(b.ID))
		test.Assert(t, err != nil, "Lane should not be found")
		test.Assert(t, l == nil, "Lane should not be found")
	})
}

func testBoard(t *testing.T, h func(c context.Context, f *services.ServiceFactory, b *board.Model)) {
	testServices(t, func(c context.Context, f *services.ServiceFactory) {
		s := f.CreateBoardService()

		model := &board.CreateModel{
			Owner:       "test_owner",
			Name:        "test_name",
			Description: "test_description",
			Layout:      kernel.HLayout}

		//Create board
		id, err := s.Create(c, model)
		test.Ok(t, err)

		//Confirm created
		b, err := s.GetByID(c, id)
		test.Ok(t, err)
		test.Assert(t, b != nil, "Board should be found")
		test.AssertExpAct(t, id, b.ID)

		h(c, f, b)

		//Remove board
		test.Ok(t, s.Remove(c, id))

		//Confirm removed
		b, err = s.GetByID(c, id)
		test.Assert(t, err != nil, "Board should not be found")
		test.Assert(t, b == nil, "Board should not be found")
	})
}

func testServices(t *testing.T, h func(context.Context, *services.ServiceFactory)) {
	l := &noop.Logger{}
	s := persistence.CreateSessionFactory(mongo.CreateSessionFactory(), l)
	p := mongo.CreateSessionProvider(s, l)
	e := mongo.CreateExecutor(p, l)
	f := persistence.CreateRepositoryFactory(e, l)

	publisher := &mocks.Publisher{}
	publisher.On("Publish", mock.Anything).Return(nil)
	
	h(context.Background(), services.CreateServiceFactory(f, publisher, l))
}
