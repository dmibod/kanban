package event_test

import (
	"testing"

	"github.com/dmibod/kanban/shared/domain/event"
	"github.com/dmibod/kanban/shared/tools/test"
)

func TestShouldFireEvent(t *testing.T) {
	event.Execute(func(bus event.Bus) error {

		type eventType struct{ Name string }

		exp := eventType{"Test"}

		bus.Register(exp)

		var act eventType

		bus.Listen(event.HandleFunc(func(event interface{}) {
			act = event.(eventType)
		}))

		bus.Fire()

		test.AssertExpAct(t, exp, act)
		test.AssertExpAct(t, exp.Name, act.Name)

		return nil
	})
}
