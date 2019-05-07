// +build integration

package card_test

import (
	"context"
	"github.com/dmibod/kanban/shared/services"
	"github.com/dmibod/kanban/shared/services/board"
	"github.com/dmibod/kanban/shared/services/card"
	"github.com/dmibod/kanban/shared/services/lane"
	setup "github.com/dmibod/kanban/shared/services/test"
	"github.com/dmibod/kanban/shared/tools/test"
	"testing"
)

func Test(t *testing.T) {
	setup.WithCard(t, func(ctx context.Context, f *services.ServiceFactory, b *board.Model, l *lane.Model, c *card.Model) {
		testCards(t, ctx, f, b, l, c)
	})
}

func testCards(t *testing.T, ctx context.Context, f *services.ServiceFactory, b *board.Model, l *lane.Model, c *card.Model) {
	s := f.CreateCardService()

	//Update card
	test.Ok(t, s.Name(ctx, c.ID.WithSet(b.ID), c.Name+"_updated"))
	test.Ok(t, s.Describe(ctx, c.ID.WithSet(b.ID), c.Description+"_updated"))

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
