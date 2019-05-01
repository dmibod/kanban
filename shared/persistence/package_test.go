// +build integration

package persistence_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/dmibod/kanban/shared/tools/logger/noop"
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
	r := f.CreateRepository()
	c := context.Background()

	err := r.FindBoardsByOwner(c, "112099531976694150844", func(board *persistence.BoardListModel) error {
		fmt.Printf("%+v\n", board)
		return nil
	})
	test.Ok(t, err)
}
