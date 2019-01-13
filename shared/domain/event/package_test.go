package event_test

import (
	"github.com/dmibod/kanban/shared/domain/event"
	"github.com/dmibod/kanban/shared/tools/test"
	"testing"
)

func TestShouldFireEvent(t *testing.T) {
	manager := event.CreateEventManager()

	type eventType struct{ Name string }

	exp := eventType{"Test"}

	manager.Register(exp)

	var act eventType

	manager.Listen(event.HandleFunc(func(event interface{}) {
		act = event.(eventType)
	}))

	manager.Fire()

	test.AssertExpAct(t, exp, act)
	test.AssertExpAct(t, exp.Name, act.Name)
}
