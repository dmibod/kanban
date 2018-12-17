package query_test

import (
	"github.com/dmibod/kanban/shared/persistence"
	"testing"

	"github.com/dmibod/kanban/shared/tools/log/logger"

	"github.com/dmibod/kanban/shared/kernel"

	"github.com/dmibod/kanban/query"
)

type repository struct {
	id     string
	entity *persistence.CardEntity
	count  int
}

func (r *repository) FindById(id string) (interface{}, error) {
	r.id = id
	return r.entity, nil
}

func mockRepository(id string, count int) *repository {
	return &repository{
		id: id,
		entity: &persistence.CardEntity{
			Name: "newentity",
		},
		count: count,
	}
}
func TestGetCardByID(t *testing.T) {
	r := mockRepository("newid", 10)

	s := &query.CardService{Logger: logger.New(), Repository: r}

	_, err := s.GetCardByID(kernel.Id(r.id))

	if err != nil {
		t.Fatal(err)
	}

	if r.id != "newid" {
		t.Fatal("Id does not match")
	}
}
