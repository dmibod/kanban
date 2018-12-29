package stan_test

import (
	"testing"
	"time"

	"github.com/dmibod/kanban/shared/tools/logger/console"

	"github.com/dmibod/kanban/shared/tools/bus"

	"github.com/dmibod/kanban/shared/tools/bus/stan"
)

var enable bool = false

func TestStan(t *testing.T) {
	if enable {
		testStan(t)
	}
}

func testStan(t *testing.T) {
	l := console.New(console.WithDebug(true))

	conn := stan.CreateConnection(
		stan.WithClientID("test"),
		stan.WithLogger(l))

	ok(t, conn.Connect())

	ch := make(chan struct{}, 1)

	l.Debugln("Subscribe topic")
	sub, err := conn.Subscribe("test.stan", "", bus.HandleFunc(func(m []byte) {
		act := string(m)
		exp := "Hello"
		assertf(t, act == exp, "Wrong value:\nwant: %v\ngot: %v\n", exp, act)

		ch <- struct{}{}
	}))
	ok(t, err)

	l.Debugln("Publish message")
	err = conn.Publish("test.stan", []byte("Hello"))
	ok(t, err)

	select {

	case <-ch:
		l.Debugln("Unsubscribe")
		ok(t, conn.Unsubscribe(sub))

		l.Debugln("Dusconnect")
		conn.Disconnect()

	case <-time.After(time.Second * 5):
		t.Fatal("Failed to connect")
	}
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
