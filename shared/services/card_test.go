package services_test

import (
	"context"
	"github.com/dmibod/kanban/shared/tools/db"
	"github.com/dmibod/kanban/shared/tools/test"
	"github.com/stretchr/testify/mock"
	"gopkg.in/mgo.v2/bson"
	"testing"

	"github.com/dmibod/kanban/shared/tools/logger/noop"

	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/persistence"
	"github.com/dmibod/kanban/shared/services"

	"github.com/dmibod/kanban/shared/tools/db/mocks"
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

	repository := &mocks.Repository{}
	repository.On("FindByID", mock.Anything, id).Return(entity, nil).Once()

	act, err := getService(repository).GetByID(context.TODO(), exp.ID)
	test.Ok(t, err)

	repository.AssertExpectations(t)

	test.AssertExpAct(t, *exp, *act)
}

func TestCreateCard(t *testing.T) {

	exp := "5c16dd24c7ee6e5dcf626266"
	payload := &services.CardPayload{Name: "Sample"}

	entity := &persistence.CardEntity{Name: payload.Name}
	repository := &mocks.Repository{}
	repository.On("Create", mock.Anything, entity).Return(exp, nil).Once()

	id, err := getService(repository).Create(context.TODO(), payload)
	test.Ok(t, err)

	repository.AssertExpectations(t)

	act := string(id)

	test.AssertExpAct(t, exp, act)
}

func getService(r db.Repository) services.CardService {
	factory := &mocks.RepositoryFactory{}
	factory.On("CreateRepository", mock.Anything, mock.Anything, mock.Anything).Return(r)

	return services.CreateServiceFactory(factory, &noop.Logger{}).CreateCardService()
}
