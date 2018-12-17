package update_test

import (
	"testing"

	dbm "github.com/dmibod/kanban/shared/tools/db/mocks"
	logm "github.com/dmibod/kanban/shared/tools/log/mocks"
	"github.com/dmibod/kanban/update"
)

func TestCardService(t *testing.T) {

	id := "000"
	card := update.CardPayload{Name: "Sample"}

	r := &dbm.Repository{}
	r.On("Create", &card).Return(id, nil).Once()

	s := update.CreateCardService(&logm.Logger{}, r)

	s.CreateCard(&card)

	r.AssertExpectations(t)
}
