package services_test

import (
	"testing"

	"github.com/dmibod/kanban/shared/tools/logger/noop"

	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/persistence"
	"github.com/dmibod/kanban/shared/services"
	"github.com/mongodb/mongo-go-driver/bson/primitive"

	_db "github.com/dmibod/kanban/shared/tools/db/mocks"
)

func TestGetCardByID(t *testing.T) {

	id := "5c16dd24c7ee6e5dcf626266"
	exp := &services.CardModel{
		ID:   kernel.Id(id),
		Name: "Sample",
	}

	_id, idErr := primitive.ObjectIDFromHex(id)
	ok(t, idErr)

	entity := &persistence.CardEntity{
		ID:   _id,
		Name: "Sample",
	}

	repository := &_db.Repository{}
	repository.On("FindByID", id).Return(entity, nil).Once()

	service := services.CreateCardService(&noop.Logger{}, repository)

	act, err := service.GetCardByID(exp.ID)
	ok(t, err)

	repository.AssertExpectations(t)

	assert(t, *act == *exp, "model does not match")
}
func TestCreateCard(t *testing.T) {

	exp := "5c16dd24c7ee6e5dcf626266"
	payload := &services.CardPayload{Name: "Sample"}

	entity := &persistence.CardEntity{Name: payload.Name}
	repository := &_db.Repository{}
	repository.On("Create", entity).Return(exp, nil).Once()

	service := services.CreateCardService(&noop.Logger{}, repository)

	id, err := service.CreateCard(payload)
	ok(t, err)

	repository.AssertExpectations(t)

	act := string(id)

	assertf(t, act == exp, "Wrong id\nwant: %v\ngot: %v\n", exp, act)
}

func ok(t *testing.T, e error) {
	if e != nil {
		t.Fatal(e)
	}
}

func assert(t *testing.T, exp bool, msg string) {
	if !exp {
		t.Fatal(msg)
	}
}

func assertf(t *testing.T, exp bool, f string, v ...interface{}) {
	if !exp {
		t.Fatalf(f, v...)
	}
}
