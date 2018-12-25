package message_test

import (
	"context"
	"github.com/dmibod/kanban/shared/message"
	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/logger/console"
	"github.com/dmibod/kanban/shared/tools/msg/nats"
	natz "github.com/nats-io/go-nats"
	"github.com/nats-io/go-nats-streaming"
	"testing"
	"time"
)

const enable = false

func TestMessage(t *testing.T) {
	if enable {
		testMessage(t)
	}
}

func testMessage(t *testing.T) {
	l := console.New(console.WithDebug(true))
	f := nats.CreateTransport(context.TODO(), service(l))

	err := f.CreateReceiver("topic").Receive("", func(msg []byte) {
		act := string(msg)
		exp := "Hello World!"
		assertf(t, act == exp, "Wrong message:\nwant: %v\ngot: %v\n", act, exp)
		l.Debugf("Received message: %v\n", act)
	})
	ok(t, err)

	err = f.CreateSender("topic").Send([]byte("Hello World!"))
	ok(t, err)
}

func wrapped_service(l logger.Logger) nats.OperationExecutor {
	return message.CreateService(l)
}

func service(l logger.Logger) nats.OperationExecutor {
	return nats.CreateExecutor(
		nats.WithLogger(l),
		nats.WithReconnectDelay(time.Second),
		nats.WithName("KANBAN"),
		nats.WithClusterID("test-cluster"),
		nats.WithClientID("KANBAN-CLIENT"),
		nats.WithConnectionLostHandler(func(c stan.Conn, reason error) { l.Debugf("connection lost, reason %v...", reason) }),
		nats.WithReconnectHandler(func(c *natz.Conn) { l.Debugln("reconnect...") }),
		nats.WithDisconnectHandler(func(c *natz.Conn) { l.Debugln("disconnect...") }),
		nats.WithCloseHandler(func(c *natz.Conn) { l.Debugln("close...") }))
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
