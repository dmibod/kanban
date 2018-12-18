package query_test

import (
	"testing"

	"github.com/dmibod/kanban/shared/persistence"

	"github.com/dmibod/kanban/shared/tools/log/mocks"

	"github.com/dmibod/kanban/shared/kernel"

	"github.com/dmibod/kanban/query"
)

type repository struct {
	fn func(string)
}

func (r *repository) FindById(id string) (interface{}, error) {

	r.fn(id)

	return &persistence.CardEntity{
		Name: "newentity",
	}, nil
}
func TestGetCardByID(t *testing.T) {

	var act string
	var call int

	r := &repository{
		fn: func(id string) {
			act = id
			call++
		},
	}

	exp := "newid"

	_, err := query.CreateCardService(&mocks.Logger{}, r).GetCardByID(kernel.Id(exp))
	ok(t, err)

	assert(t, call == 1, "repository must be called once")
	assert(t, act == exp, "id does not match")
}
