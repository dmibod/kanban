package bus_test

import (
	"context"
	"testing"
	"time"

	"github.com/dmibod/kanban/shared/tools/logger"

	"github.com/dmibod/kanban/shared/tools/bus/nats"
	"github.com/dmibod/kanban/shared/tools/bus/stan"
	"github.com/dmibod/kanban/shared/tools/logger/console"

	"github.com/dmibod/kanban/shared/tools/bus"
)

var enable bool = false

func TestBusWithNats(t *testing.T) {
	if enable {
		testBus(t, true)
	}
}

func TestBusWithStan(t *testing.T) {
	if enable {
		testBus(t, false)
	}
}

func testBus(t *testing.T, isNats bool) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	ch := make(chan struct{}, 1)

	l := console.New(console.WithDebug(true))

	l.Debugln("Subscribe topics")
	sub := bus.Subscribe("test.bus1", bus.HandleFunc(func(m []byte) {
		act := string(m)
		exp := "Hello"
		assertf(t, act == exp, "Wrong value:\nwant: %v\ngot: %v\n", exp, act)

		ch <- struct{}{}
	}))

	bus.Subscribe("test.bus2", bus.HandleFunc(func(m []byte) {
		act := string(m)
		exp := "Bye"
		assertf(t, act == exp, "Wrong value:\nwant: %v\ngot: %v\n", exp, act)
	}))

	var conn bus.Connection

	if isNats {
		conn = natsConnection(l)
	} else {
		conn = stanConnection(l)
	}

	l.Debugln("Connect and Serve")
	ok(t, bus.ConnectAndServe(ctx, conn))

	l.Debugln("Publish messages")
	ok(t, bus.Publish("test.bus1", []byte("Hello")))
	ok(t, bus.Publish("test.bus2", []byte("Bye")))

	<-ch

	l.Debugln("Unsubscribe")
	ok(t, sub.Unsubscribe())

	l.Debugln("Close connection")
	conn.Disconnect()
}

func natsConnection(l logger.Logger) bus.Connection {
	return nats.CreateConnection(
		nats.WithName("test"),
		nats.WithLogger(l))
}

func stanConnection(l logger.Logger) bus.Connection {
	return stan.CreateConnection(
		stan.WithClientID("test"),
		stan.WithLogger(l))
}

func ok(t *testing.T, e error) {
	if e != nil {
		t.Fatal(e)
	}
}

func assert(t *testing.T, exp bool, msg string) {
	if !exp {
		t.Fatal(msg)
	}
}

func assertf(t *testing.T, exp bool, f string, v ...interface{}) {
	if !exp {
		t.Fatalf(f, v...)
	}
}
