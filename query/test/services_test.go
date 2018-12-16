package test

import (
	"testing"

	"github.com/dmibod/kanban/kernel"

	"github.com/dmibod/kanban/query"
)

type repository struct {
	id     string
	entity *query.CardEntity
	count  int
}

func (r *repository) FindById(id string) (interface{}, error) {
	r.id = id
	return r.entity, nil
}

func mockRepository(id string, count int) *repository {
	return &repository{
		id: id,
		entity: &query.CardEntity{
			Name: "newentity",
		},
		count: count,
	}
}
func TestGetCardByID(t *testing.T) {
	r := mockRepository("newid", 10)

	s := &query.CardService{Repository: r}

	_, err := s.GetCardByID(kernel.Id(r.id))

	if err != nil {
		t.Fatal(err)
	}

	if r.id != "newid" {
		t.Fatal("Id does not match")
	}
}
