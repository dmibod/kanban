package query

import (
	"log"
	"net/http"
)

type Env struct {
	Db interface {}
}

func (*Env) GetCards(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	log.Println("Request received")
}
