package services_test

import (
	"context"
	"github.com/dmibod/kanban/shared/tools/db"
	"github.com/stretchr/testify/mock"
	"gopkg.in/mgo.v2/bson"
	"testing"

	"github.com/dmibod/kanban/shared/tools/logger/noop"

	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/persistence"
	"github.com/dmibod/kanban/shared/services"

	_db "github.com/dmibod/kanban/shared/tools/db/mocks"
)

func TestGetCardByID(t *testing.T) {

	id := "5c16dd24c7ee6e5dcf626266"
	exp := &services.CardModel{
		ID:   kernel.Id(id),
		Name: "Sample",
	}

	entity := &persistence.CardEntity{
		ID:   bson.ObjectIdHex(id),
		Name: "Sample",
	}

	repository := &_db.Repository{}
	repository.On("FindByID", mock.Anything, id).Return(entity, nil).Once()

	act, err := getService(repository).GetCardByID(context.TODO(), exp.ID)
	ok(t, err)

	repository.AssertExpectations(t)

	assert(t, *act == *exp, "model does not match")
}

func TestCreateCard(t *testing.T) {

	exp := "5c16dd24c7ee6e5dcf626266"
	payload := &services.CardPayload{Name: "Sample"}

	entity := &persistence.CardEntity{Name: payload.Name}
	repository := &_db.Repository{}
	repository.On("Create", mock.Anything, entity).Return(exp, nil).Once()

	id, err := getService(repository).CreateCard(context.TODO(), payload)
	ok(t, err)

	repository.AssertExpectations(t)

	act := string(id)

	assertf(t, act == exp, "Wrong id\nwant: %v\ngot: %v\n", exp, act)
}

func getService(r db.Repository) services.CardService {
	factory := &_db.RepositoryFactory{}
	factory.On("CreateRepository", mock.Anything, mock.Anything, mock.Anything).Return(r)

	return services.CreateFactory(&noop.Logger{}, factory).CreateCardService()
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
