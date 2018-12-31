package nats_test

import (
	"testing"
	"time"

	"github.com/dmibod/kanban/shared/tools/test"

	"github.com/dmibod/kanban/shared/tools/bus"

	"github.com/dmibod/kanban/shared/tools/bus/nats"
)

var enable bool = false

func TestConn(t *testing.T) {
	if enable {
		testConn(t)
	}
}

func testConn(t *testing.T) {
	conn := nats.CreateConnection(nats.WithName("test"))
	test.Ok(t, conn.Connect())

	ch := make(chan struct{}, 1)

	sub, err := conn.Subscribe("test.nats", "", bus.HandleFunc(func(m []byte) {
		act := string(m)
		exp := "Hello"
		test.AssertExpAct(t, exp, act)

		ch <- struct{}{}
	}))
	test.Ok(t, err)

	test.Ok(t, conn.Publish("test.nats", []byte("Hello")))

	select {
	case <-ch:
		test.Ok(t, conn.Unsubscribe(sub))
		conn.Disconnect()
	case <-time.After(time.Second * 5):
		t.Fatal("Failed to connect")
	}
}
