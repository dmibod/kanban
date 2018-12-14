package test

import (
	"testing"

	"github.com/dmibod/kanban/kernel"

	"github.com/dmibod/kanban/query"
	"github.com/dmibod/kanban/tools/db"
)

type repository struct {
	id     string
	entity *query.CardEntity
	count  int
}

func (r *repository) Create(e interface{}) (string, error) {
	return r.id, nil
}

func (r *repository) FindById(id string) (interface{}, error) {
	r.id = id
	return r.entity, nil
}

func (r *repository) Find(f interface{}, v db.Visitor) error {
	v(r.entity)
	return nil
}

func (r *repository) Count(f interface{}) (int, error) {
	return r.count, nil
}

func (r *repository) Update(e interface{}) error {
	return nil
}

func (r *repository) Remove(id string) error {
	r.id = id
	return nil
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

func mockService(r db.Repository) *query.CardService {
	return query.CreateCardService(r)
}
func TestGetCardByID(t *testing.T) {
	r := mockRepository("newid", 10)

	_, err := mockService(r).GetCardByID(kernel.Id(r.id))

	if err != nil {
		t.Fatal(err)
	}

	if r.id != "newid!" {
		t.Fatal("Id does not match")
	}
}
