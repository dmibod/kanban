package message_test

import (
	"context"
	"github.com/dmibod/kanban/shared/message"
	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/logger/console"
	"github.com/dmibod/kanban/shared/tools/msg/nats"
	"testing"
)

const enable = false

func TestMessage(t *testing.T) {
	if enable {
		testMessage(t)
	}
}

func TestSendMessage(t *testing.T) {
	if enable {
		testSendMessage(t)
	}
}

func testMessage(t *testing.T) {
	l := console.New(console.WithDebug(true))
	f := nats.CreateTransport(context.TODO(), service(l), l)

	_, err := f.Subscriber("topic").Subscribe("", func(msg []byte) {
		act := string(msg)
		exp := "Hello World!"
		assertf(t, act == exp, "Wrong message:\nwant: %v\ngot: %v\n", act, exp)
		l.Debugf("Received message: %v\n", act)
	})
	ok(t, err)

	err = f.Publisher("topic").Publish([]byte("Hello World!"))
	ok(t, err)
}

func testSendMessage(t *testing.T) {
	l := console.New(console.WithDebug(true))
	f := nats.CreateTransport(context.TODO(), service(l), l)

	err := f.Publisher("topic").Publish([]byte("Hello World!"))
	assertf(t, err != nil, "Sending message should fail")
}

func service(l logger.Logger) nats.OperationExecutor {
	return message.CreateService("test", l)
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
