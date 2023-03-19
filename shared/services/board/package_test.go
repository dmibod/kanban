//go:build integration
// +build integration

package board_test

import (
	"context"
	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/services"
	"github.com/dmibod/kanban/shared/services/board"
	setup "github.com/dmibod/kanban/shared/services/test"
	"github.com/dmibod/kanban/shared/tools/test"
	"testing"
)

func Test(t *testing.T) {
	setup.WithBoard(t, func(c context.Context, f *services.ServiceFactory, b *board.Model) {
		testBoards(t, c, f, b)
	})
}

func testBoards(t *testing.T, c context.Context, f *services.ServiceFactory, b *board.Model) {
	//Update board
	s := f.CreateBoardService()

	test.Ok(t, s.Name(c, b.ID, b.Name+"_updated"))
	test.Ok(t, s.Describe(c, b.ID, b.Description+"_updated"))
	test.Ok(t, s.Layout(c, b.ID, kernel.VLayout))
	test.Ok(t, s.Share(c, b.ID, !b.Shared))

	list, err := s.GetByOwner(c, "test_owner")
	test.Ok(t, err)
	test.AssertExpAct(t, true, len(list) > 0)
	//test.AssertExpAct(t, b.ID, list[0].ID)
}
