// +build integration

package persistence_test

import (
	"context"
	"github.com/dmibod/kanban/shared/domain/board"
	"github.com/dmibod/kanban/shared/domain/card"
	"github.com/dmibod/kanban/shared/domain/lane"
	"github.com/dmibod/kanban/shared/kernel"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"testing"

	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"github.com/dmibod/kanban/shared/tools/logger/noop"

	"github.com/dmibod/kanban/shared/tools/test"

	"github.com/dmibod/kanban/shared/persistence"
)

func TestBoards(t *testing.T) {
	testBoard(t, func(c context.Context, r persistence.Repository, b board.Entity) {
		testBoards(t, c, r, b)
	})
}

func testBoards(t *testing.T, c context.Context, r persistence.Repository, b board.Entity) {
	//Update board
	test.Ok(t, r.Handle(c, board.NameChangedEvent{ID: b.ID, OldValue: b.Name, NewValue: b.Name + "_updated"}))
	test.Ok(t, r.Handle(c, board.DescriptionChangedEvent{ID: b.ID, OldValue: b.Description, NewValue: b.Description + "_updated"}))
	test.Ok(t, r.Handle(c, board.LayoutChangedEvent{ID: b.ID, OldValue: b.Layout, NewValue: kernel.VLayout}))
	test.Ok(t, r.Handle(c, board.SharedChangedEvent{ID: b.ID, OldValue: b.Shared, NewValue: !b.Shared}))

	count := 0
	test.Ok(t, r.FindBoardsByOwner(c, "test_owner", func(board *persistence.BoardListModel) error {
		id := kernel.ID(board.ID.Hex())
		test.AssertExpAct(t, b.ID, id)
		count++
		return nil
	}))
	test.AssertExpAct(t, 1, count)
}

func TestLanes(t *testing.T) {
	testLane(t, func(c context.Context, r persistence.Repository, b board.Entity, l lane.Entity) {
		testLanes(t, c, r, b, l)
	})
}

func testLanes(t *testing.T, c context.Context, r persistence.Repository, b board.Entity, l lane.Entity) {
	//Update lane
	test.Ok(t, r.Handle(c, lane.NameChangedEvent{ID: l.ID, OldValue: l.Name, NewValue: l.Name + "_updated"}))
	test.Ok(t, r.Handle(c, lane.DescriptionChangedEvent{ID: l.ID, OldValue: l.Description, NewValue: l.Description + "_updated"}))
	test.Ok(t, r.Handle(c, lane.LayoutChangedEvent{ID: l.ID, OldValue: l.Layout, NewValue: kernel.HLayout}))

	//Confirm not attached
	test.Ok(t, r.FindLanesByParent(c, b.ID.WithID(kernel.EmptyID), func(lane *persistence.LaneListModel) error {
		test.Fail(t, "Lanes should not be attached to Board")
		return nil
	}))

	//Attach lane
	test.Ok(t, r.Handle(c, board.ChildAppendedEvent{ID: l.ID}))

	//Confirm attached
	count := 0
	test.Ok(t, r.FindLanesByParent(c, b.ID.WithID(kernel.EmptyID), func(lane *persistence.LaneListModel) error {
		id := kernel.ID(lane.ID.Hex())
		test.AssertExpAct(t, l.ID.ID, id)
		count++
		return nil
	}))
	test.AssertExpAct(t, 1, count)

	//Detach lane
	test.Ok(t, r.Handle(c, board.ChildRemovedEvent{ID: l.ID}))

	//Confirm detached
	test.Ok(t, r.FindLanesByParent(c, b.ID.WithID(kernel.EmptyID), func(lane *persistence.LaneListModel) error {
		test.Fail(t, "Lane should be detached from Board")
		return nil
	}))
}

func TestCards(t *testing.T) {
	testCard(t, func(ctx context.Context, r persistence.Repository, b board.Entity, l lane.Entity, c card.Entity) {
		testCards(t, ctx, r, b, l, c)
	})
}

func testCards(t *testing.T, ctx context.Context, r persistence.Repository, b board.Entity, l lane.Entity, c card.Entity) {
	//Update card
	test.Ok(t, r.Handle(ctx, card.NameChangedEvent{ID: c.ID, OldValue: c.Name, NewValue: c.Name + "_updated"}))
	test.Ok(t, r.Handle(ctx, card.DescriptionChangedEvent{ID: c.ID, OldValue: c.Description, NewValue: c.Description + "_updated"}))

	//Confirm not attached
	test.Ok(t, r.FindCardsByParent(ctx, l.ID, func(entity *persistence.Card) error {
		test.Fail(t, "Cards should not be attached to Lane")
		return nil
	}))

	//Attach card
	test.Ok(t, r.Handle(ctx, lane.ChildAppendedEvent{ID: l.ID, ChildID: c.ID.ID}))

	//Confirm attached
	count := 0
	test.Ok(t, r.FindCardsByParent(ctx, l.ID, func(entity *persistence.Card) error {
		id := kernel.ID(entity.ID.Hex())
		test.AssertExpAct(t, c.ID.ID, id)
		count++
		return nil
	}))
	test.AssertExpAct(t, 1, count)

	//Detach card
	test.Ok(t, r.Handle(ctx, lane.ChildRemovedEvent{ID: l.ID, ChildID: c.ID.ID}))

	//Confirm detached
	test.Ok(t, r.FindCardsByParent(ctx, l.ID, func(entity *persistence.Card) error {
		test.Fail(t, "Card should be detached from Lane")
		return nil
	}))
}

func testCard(t *testing.T, h func(ctx context.Context, r persistence.Repository, b board.Entity, l lane.Entity, c card.Entity)) {
	testLane(t, func(ctx context.Context, r persistence.Repository, b board.Entity, l lane.Entity) {
		id := kernel.ID(bson.NewObjectId().Hex())

		entity := card.Entity{
			ID:          id.WithSet(b.ID),
			Name:        "test_name",
			Description: "test_description",
		}

		//Create card
		test.Ok(t, r.Handle(ctx, card.CreatedEvent{Entity: entity}))

		//Find by ID
		count := 0
		test.Ok(t, r.FindCardByID(ctx, entity.ID, func(c *persistence.Card) error {
			test.AssertExpAct(t, entity.ID.ID, kernel.ID(c.ID.Hex()))
			count++
			return nil
		}))
		test.AssertExpAct(t, 1, count)

		h(ctx, r, b, l, entity)

		//Remove card
		test.Ok(t, r.Handle(ctx, card.DeletedEvent{Entity: entity}))

		//Confirm removed
		test.Assert(t, r.FindCardByID(ctx, entity.ID, func(*persistence.Card) error {
			test.Fail(t, "Card should be removed")
			return nil
		}) == mgo.ErrNotFound, "Card should not be found")
	})
}

func testLane(t *testing.T, h func(c context.Context, r persistence.Repository, b board.Entity, l lane.Entity)) {
	testBoard(t, func(c context.Context, r persistence.Repository, b board.Entity) {
		id := kernel.ID(bson.NewObjectId().Hex())

		entity := lane.Entity{
			ID:          id.WithSet(b.ID),
			Kind:        kernel.LKind,
			Name:        "test_name",
			Description: "test_description",
			Layout:      kernel.VLayout,
			Children:    []kernel.ID{}}

		//Create lane
		test.Ok(t, r.Handle(c, lane.CreatedEvent{Entity: entity}))

		//Confirm created
		count := 0
		test.Ok(t, r.FindLaneByID(c, entity.ID, func(l *persistence.Lane) error {
			test.AssertExpAct(t, entity.ID.ID, kernel.ID(l.ID.Hex()))
			count++
			return nil
		}))
		test.AssertExpAct(t, 1, count)

		h(c, r, b, entity)

		//Remove lane
		test.Ok(t, r.Handle(c, lane.DeletedEvent{Entity: entity}))

		//Confirm removed
		test.Assert(t, r.FindLaneByID(c, entity.ID, func(*persistence.Lane) error {
			test.Fail(t, "Lane should be removed")
			return nil
		}) == mgo.ErrNotFound, "Lane should not be found")
	})
}

func testBoard(t *testing.T, h func(c context.Context, r persistence.Repository, b board.Entity)) {
	testRepository(t, func(c context.Context, r persistence.Repository) {
		id := kernel.ID(bson.NewObjectId().Hex())

		entity := board.Entity{
			ID:          id,
			Owner:       "test_owner",
			Name:        "test_name",
			Description: "test_description",
			Layout:      kernel.HLayout,
			Shared:      true,
			Children:    []kernel.ID{}}

		//Create board
		test.Ok(t, r.Handle(c, board.CreatedEvent{Entity: entity}))

		//Confirm created
		count := 0
		test.Ok(t, r.FindBoardByID(c, entity.ID, func(b *persistence.Board) error {
			test.AssertExpAct(t, entity.ID, kernel.ID(b.ID.Hex()))
			count++
			return nil
		}))
		test.AssertExpAct(t, 1, count)

		h(c, r, entity)

		//Remove board
		test.Ok(t, r.Handle(c, board.DeletedEvent{Entity: entity}))

		//Confirm removed
		test.Assert(t, r.FindBoardByID(c, entity.ID, func(*persistence.Board) error {
			test.Fail(t, "Board should be removed")
			return nil
		}) == mgo.ErrNotFound, "Board should not be found")
	})
}

func testRepository(t *testing.T, h func(c context.Context, r persistence.Repository)) {
	l := &noop.Logger{}
	s := persistence.CreateSessionFactory(mongo.CreateSessionFactory(), l)
	p := mongo.CreateSessionProvider(s, l)
	e := mongo.CreateExecutor(p, l)
	f := persistence.CreateRepositoryFactory(e, l)
	r := f.CreateRepository()
	c := context.Background()

	h(c, r)
}
