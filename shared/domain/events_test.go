package domain_test

import (
	"github.com/dmibod/kanban/shared/tools/test"
	"testing"

	"github.com/dmibod/kanban/shared/domain"
)

func TestShouldFireEvent(t *testing.T) {
	eventManager := domain.CreateEventManager()

	type eventType struct{ Name string }

	exp := eventType{"Test"}

	eventManager.Register(exp)

	var act eventType

	eventManager.Listen(domain.HandleFunc(func(event interface{}) {
		act = event.(eventType)
	}))

	eventManager.Fire()

	test.AssertExpAct(t, exp, act)
	test.AssertExpAct(t, exp.Name, act.Name)
}
