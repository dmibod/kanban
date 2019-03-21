// +build integration

package board_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/dmibod/kanban/shared/tools/logger/noop"
	"gopkg.in/mgo.v2/bson"

	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/tools/db/mongo"

	"github.com/dmibod/kanban/shared/tools/test"

	"github.com/dmibod/kanban/shared/persistence"
)

func TestBoards(t *testing.T) {
	l := &noop.Logger{}
	s := persistence.CreateSessionFactory(mongo.CreateSessionFactory(), l)
	p := mongo.CreateSessionProvider(s, l)
	e := persistence.CreateOperationExecutor(p, l)
	f := persistence.CreateRepositoryFactory(e, l)
	r := f.CreateRepository("fullboards")
	c := context.Background()

	entity := createBoard(true, "Test", "Board", "Desc")

	id, err := r.Create(c, entity)
	test.Ok(t, err)

	found, err := r.FindByID(c, id, &persistence.Board{})
	test.Ok(t, err)

	entity, ok := found.(*persistence.Board)
	test.Assert(t, ok, "invalid type")
}

func TestBoards2(t *testing.T) {
	l := &noop.Logger{}
	s := persistence.CreateSessionFactory(mongo.CreateSessionFactory(), l)
	p := mongo.CreateSessionProvider(s, l)
	e := persistence.CreateOperationExecutor(p, l)
	f := persistence.CreateRepositoryFactory(e, l)
	r := f.CreateRepository("fullboards")
	c := context.Background()

	//criteria := bson.M{}
	err := r.Find(c, nil, &persistence.Board{}, func(item interface{}) error {
		fmt.Println(item)
		return nil
	})
	test.Ok(t, err)
	t.Fatal("")
}

func createBoard(shared bool, owner, name, description string) *persistence.Board {
	board := &persistence.Board{
		Owner:       owner,
		Name:        name,
		Description: description,
		Layout:      kernel.VLayout,
		Shared:      shared,
	}

	lanes := []persistence.Lane{
		*createLane(board, true, "name_1", "desc_1", 2),
		*createLane(board, false, "name_2", "desc_2", 0),
		*createLane(board, true, "name_3", "desc_3", 3),
	}

	children := []bson.ObjectId{}

	for _, lane := range lanes {
		children = append(children, lane.ID)
	}

	board.Children = children

	return board
}

func createLane(board *persistence.Board, composite bool, name, description string, levels int) *persistence.Lane {
	kind := kernel.CKind
	layout := ""

	var lane *persistence.Lane
	var card *persistence.Card

	children := []bson.ObjectId{}

	if composite {
		kind = kernel.LKind
		layout = kernel.HLayout
		if levels > 0 {

			lane = createLane(board, true, fmt.Sprintf("CompLane-%v-0", levels), "desc", levels-1)
			children = append(children, lane.ID)

			board.Lanes = append(board.Lanes, *lane)

			lane = createLane(board, false, fmt.Sprintf("CardLane-%v-0", levels), "desc", levels-1)
			children = append(children, lane.ID)

			board.Lanes = append(board.Lanes, *lane)
		}
	} else {
		for i := 0; i < 10; i++ {
			card = createCard(fmt.Sprintf("Card-%v-%v", levels, i), "desc")

			children = append(children, card.ID)

			board.Cards = append(board.Cards, *card)
		}
	}

	return &persistence.Lane{
		ID:          bson.NewObjectId(),
		Kind:        kind,
		Layout:      layout,
		Name:        name,
		Description: description,
		Children:    children,
	}
}

func createCard(name, description string) *persistence.Card {
	return &persistence.Card{
		ID:          bson.NewObjectId(),
		Name:        name,
		Description: description,
	}
}
