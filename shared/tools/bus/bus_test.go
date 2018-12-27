package bus_test

import (
	"context"
	"github.com/dmibod/kanban/shared/tools/bus/nats"
	"github.com/dmibod/kanban/shared/tools/logger/console"
	"testing"

	"github.com/dmibod/kanban/shared/tools/bus"
)

var enable bool = false

func TestBus(t *testing.T) {
	if enable {
		testBus(t)
	}
}

func testBus(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	t.Log("Subscribe topic")
	bus.Subscribe("test.bus", bus.HandleFunc(func(m []byte) {
		act := string(m)
		exp := "Hello"
		assertf(t, act == exp, "Wrong value:\nwant: %v\ngot: %v\n", exp, act)
	}))

	conn := nats.CreateConnection(
		nats.WithContext(ctx),
		nats.WithClientID("test"),
		nats.WithLogger(console.New(console.WithDebug(true))))

	t.Log("Connect and Serve")
	err := bus.ConnectAndServe(conn)
	ok(t, err)

	t.Log("Publish message")
	err = bus.Publish("test.bus", []byte("Hello"))
	ok(t, err)

	t.Log("Close connection")
	cancel()
	<-conn.Close()
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
