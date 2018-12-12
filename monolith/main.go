package main

import (
	"log"
	"context"
	"time"

	"github.com/dmibod/kanban/command"
	"github.com/dmibod/kanban/notify"
	"github.com/dmibod/kanban/process"
	"github.com/dmibod/kanban/query"
	"github.com/dmibod/kanban/tools/mux/http"
	"github.com/dmibod/kanban/update"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	mux := http.New(http.WithPort(3000))

	command.Boot(mux)
	notify.Boot(mux)
	query.Boot(mux)
	update.Boot(mux)
	process.Boot(ctx)

	log.Println("Starting mux at port 3000... ")

	mux.Start()

	cancel()

	time.Sleep(time.Second)
}
