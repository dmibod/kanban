package stan_test

import (
	"context"
	"testing"
	"time"

	"github.com/dmibod/kanban/shared/tools/logger/console"

	"github.com/dmibod/kanban/shared/tools/bus"

	"github.com/dmibod/kanban/shared/tools/bus/stan"
)

var enable bool = true

func TestStan(t *testing.T) {
	if enable {
		testStan(t)
	}
}

func testStan(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	conn := stan.CreateConnection(
		stan.WithContext(ctx),
		stan.WithClientID("test"),
		stan.WithLogger(console.New(console.WithDebug(true))))

OuterLoop:
	for {
		select {

		case status := <-conn.Connect():

			assertf(t, status, "Wrong status:\nwant: true\ngot: false\n")

			t.Log("Subscribe topic")

			_, err := conn.Subscribe("test.nats", "", bus.HandleFunc(func(m []byte) {
				act := string(m)
				exp := "Hello"
				assertf(t, act == exp, "Wrong value:\nwant: %v\ngot: %v\n", exp, act)
			}))
			ok(t, err)

			t.Log("Publish message")

			err = conn.Publish("test.nats", []byte("Hello"))
			ok(t, err)

			break OuterLoop

		case <-time.After(time.Second * 5):

			t.Fatal("Failed to connect")

			//break OuterLoop
		}
	}

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
