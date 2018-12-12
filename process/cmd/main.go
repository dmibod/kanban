package main

import (
	"time"
	"log"
	"os/signal"
	"os"
	"github.com/dmibod/kanban/process"
	"context"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	log.Println("Starting processor...");

	process.Boot(ctx)

	log.Println("Processor started!");

	<-c

	log.Println("Interrupt signal received!");

	cancel()

	time.Sleep(time.Second)
}
