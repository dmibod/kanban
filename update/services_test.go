package update_test

import (
	"testing"

	_db "github.com/dmibod/kanban/shared/tools/db/mocks"
	_log "github.com/dmibod/kanban/shared/tools/log/mocks"
	"github.com/dmibod/kanban/update"
)

func TestCardService(t *testing.T) {

	exp := "000"
	card := update.CardPayload{Name: "Sample"}

	r := &_db.Repository{}
	r.On("Create", &card).Return(exp, nil).Once()

	id, err := update.CreateCardService(&_log.Logger{}, r).CreateCard(&card)
	ok(t, err)
	
	r.AssertExpectations(t)

	act := string(id)

	assertf(t, act == exp, "Wrong id\nwant: %v\ngot: %v\n", exp, act)
}
