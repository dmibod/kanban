package bus_test

import (
	"context"
	"testing"
	"time"

	"github.com/dmibod/kanban/shared/tools/bus/nats"
	"github.com/dmibod/kanban/shared/tools/bus/stan"
	"github.com/dmibod/kanban/shared/tools/logger/console"

	"github.com/dmibod/kanban/shared/tools/bus"
)

var enable bool = true

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
	ctx, cancel := context.WithCancel(context.Background())

	t.Log("Subscribe topic")
	bus.Subscribe("test.bus", bus.HandleFunc(func(m []byte) {
		act := string(m)
		exp := "Hello"
		assertf(t, act == exp, "Wrong value:\nwant: %v\ngot: %v\n", exp, act)
	}))

	var conn bus.Connection

	if isNats {
		conn = natsConnection(ctx)
	} else {
		conn = stanConnection(ctx)
	}

	go func() {
		<-time.After(time.Second * 5)
		cancel()
		t.Fatal("Failed to connect")
	}()

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

func natsConnection(ctx context.Context) bus.Connection {
	return nats.CreateConnection(
		nats.WithContext(ctx),
		nats.WithName("test"),
		nats.WithLogger(console.New(console.WithDebug(true))))
}

func stanConnection(ctx context.Context) bus.Connection {
	return stan.CreateConnection(
		stan.WithContext(ctx),
		stan.WithClientID("test"),
		stan.WithLogger(console.New(console.WithDebug(true))))
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
