// +build integration

package stan_test

import (
	"testing"
	"time"

	"github.com/dmibod/kanban/shared/tools/test"

	"github.com/dmibod/kanban/shared/tools/bus"

	"github.com/dmibod/kanban/shared/tools/bus/stan"
)

func TestConn(t *testing.T) {
	conn := stan.CreateConnection(stan.WithClientID("test"))
	test.Ok(t, conn.Connect())

	ch := make(chan struct{}, 1)

	sub, err := conn.Subscribe("test.stan", "", bus.HandleFunc(func(m []byte) {
		act := string(m)
		exp := "Hello"
		test.AssertExpAct(t, exp, act)

		ch <- struct{}{}
	}))
	test.Ok(t, err)

	test.Ok(t, conn.Publish("test.stan", []byte("Hello")))

	select {
	case <-ch:
		test.Ok(t, conn.Unsubscribe(sub))
		conn.Disconnect()
	case <-time.After(time.Second * 5):
		t.Fatal("Failed to connect")
	}
}
