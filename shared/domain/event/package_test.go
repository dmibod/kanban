package event_test

import (
	"context"
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

		bus.Listen(event.HandleFunc(func(ctx context.Context, event interface{}) error {
			act = event.(eventType)
			return nil
		}))

		bus.Fire(context.TODO())

		test.AssertExpAct(t, exp, act)
		test.AssertExpAct(t, exp.Name, act.Name)

		return nil
	})
}
