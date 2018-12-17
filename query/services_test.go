package query_test

import (
	"testing"

	"github.com/dmibod/kanban/shared/persistence"

	"github.com/dmibod/kanban/shared/tools/log/logger"

	"github.com/dmibod/kanban/shared/kernel"

	"github.com/dmibod/kanban/query"
)

type repository struct {
	id     string
	entity *persistence.CardEntity
}

func (r *repository) FindById(id string) (interface{}, error) {
	r.id = id
	return r.entity, nil
}

func mockRepository() *repository {
	return &repository{
		entity: &persistence.CardEntity{
			Name: "newentity",
		},
	}
}
func TestGetCardByID(t *testing.T) {
	repo := mockRepository()

	service := &query.CardService{Logger: logger.New(), Repository: repo}

	_, err := service.GetCardByID(kernel.Id("newid"))
	ok(t, err)

	assert(t, repo.id == "newid", "Id does not match")
}
