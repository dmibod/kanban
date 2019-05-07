package test

import (
	"context"
	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/message/mocks"
	"github.com/dmibod/kanban/shared/persistence"
	"github.com/dmibod/kanban/shared/services"
	"github.com/dmibod/kanban/shared/services/board"
	"github.com/dmibod/kanban/shared/services/card"
	"github.com/dmibod/kanban/shared/services/lane"
	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"github.com/dmibod/kanban/shared/tools/logger/noop"
	"github.com/dmibod/kanban/shared/tools/test"
	"github.com/stretchr/testify/mock"
	"testing"
)

func withServices(t *testing.T, h func(context.Context, *services.ServiceFactory)) {
	l := &noop.Logger{}
	s := persistence.CreateSessionFactory(mongo.CreateSessionFactory(), l)
	p := mongo.CreateSessionProvider(s, l)
	e := mongo.CreateExecutor(p, l)
	f := persistence.CreateRepositoryFactory(e, l)

	publisher := &mocks.Publisher{}
	publisher.On("Publish", mock.Anything).Return(nil)

	h(context.Background(), services.CreateServiceFactory(f, publisher, l))
}

// WithBoard test template
func WithBoard(t *testing.T, h func(c context.Context, f *services.ServiceFactory, b *board.Model)) {
	withServices(t, func(c context.Context, f *services.ServiceFactory) {
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

// WithLane test template
func WithLane(t *testing.T, h func(c context.Context, f *services.ServiceFactory, b *board.Model, l *lane.Model)) {
	WithBoard(t, func(c context.Context, f *services.ServiceFactory, b *board.Model) {
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

// WithCard test template
func WithCard(t *testing.T, h func(ctx context.Context, f *services.ServiceFactory, b *board.Model, l *lane.Model, c *card.Model)) {
	WithLane(t, func(ctx context.Context, f *services.ServiceFactory, b *board.Model, l *lane.Model) {
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
