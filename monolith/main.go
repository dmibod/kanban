package main

import (
  "fmt"
  "github.com/dmibod/kanban/messaging"
  nats "github.com/dmibod/kanban/messaging/nats"
)

func main(){
  var c messaging.Client = nats.New()
  fmt.Println(c)
}