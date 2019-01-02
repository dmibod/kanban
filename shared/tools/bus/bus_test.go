package bus_test

import (
	"context"
	"testing"
	"time"

	"github.com/dmibod/kanban/shared/tools/test"

	"github.com/dmibod/kanban/shared/tools/bus/nats"
	"github.com/dmibod/kanban/shared/tools/bus/stan"

	"github.com/dmibod/kanban/shared/tools/bus"
)

var enable bool = false

func TestBus(t *testing.T) {
	if !enable {
		return
	}
	testBus(t, true)
	testBus(t, false)
}

func testBus(t *testing.T, isNats bool) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	ch := make(chan struct{}, 1)

	sub := bus.Subscribe("test.bus1", bus.HandleFunc(func(m []byte) {
		act := string(m)
		exp := "Hello"
		test.AssertExpAct(t, exp, act)
		ch <- struct{}{}
	}))

	bus.Subscribe("test.bus2", bus.HandleFunc(func(m []byte) {
		act := string(m)
		exp := "Bye"
		test.AssertExpAct(t, exp, act)
	}))

	var conn bus.Connection
	var tran bus.Transport

	if isNats {
		conn, tran = natsConnection()
	} else {
		conn, tran = stanConnection()
	}

	test.Ok(t, bus.ConnectAndServe(ctx, conn, tran))
	test.Ok(t, bus.Publish("test.bus1", []byte("Hello")))
	test.Ok(t, bus.Publish("test.bus2", []byte("Bye")))

	<-ch

	test.Ok(t, sub.Unsubscribe())
	bus.Disconnect()
}

func natsConnection() (bus.Connection, bus.Transport) {
	conn := nats.CreateConnection(nats.WithName("test"))
	return conn, conn
}

func stanConnection() (bus.Connection, bus.Transport) {
	conn := stan.CreateConnection(stan.WithClientID("test"))
	return conn, conn
}
