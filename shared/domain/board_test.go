package domain_test

import (
	"github.com/stretchr/testify/mock"
	"github.com/dmibod/kanban/shared/kernel"
	"testing"

	"github.com/dmibod/kanban/shared/domain/mocks"
	"github.com/dmibod/kanban/shared/tools/test"

	"github.com/dmibod/kanban/shared/domain"
)

func TestNewBoard(t *testing.T) {

	type testcase struct {
		arg0 string
		arg1 domain.Repository
		arg2 domain.EventRegistry
		err  error
	}

	tests := []testcase{
		{"", nil, nil, domain.ErrInvalidArgument},
		{"test", nil, nil, domain.ErrInvalidArgument},
		{"test", &mocks.Repository{}, nil, domain.ErrInvalidArgument},
		{"test", &mocks.Repository{}, domain.CreateEventManager(), nil},
	}

	for _, c := range tests {
		_, err := domain.NewBoard(c.arg0, c.arg1, c.arg2)
		test.AssertExpAct(t, c.err, err)
	}
}

func TestLoadBoard(t *testing.T) {

	type testcase struct {
		arg0 kernel.ID
		arg1 domain.Repository
		arg2 domain.EventRegistry
		err  error
	}

	id := kernel.ID("test")

	repository := &mocks.Repository{}
	repository.On("Fetch", mock.Anything).Return(&domain.BoardEntity{ID:id}, nil)

	tests := []testcase{
		{kernel.EmptyID, nil, nil, domain.ErrInvalidID},
		{id, nil, nil, domain.ErrInvalidArgument},
		{id, repository, nil, domain.ErrInvalidArgument},
		{id, repository, domain.CreateEventManager(), nil},
	}

	for _, c := range tests {
		_, err := domain.LoadBoard(c.arg0, c.arg1, c.arg2)
		test.AssertExpAct(t, c.err, err)
	}
}
