package update

import (
	"log"
	"net/http"
)

type Env struct {
	Db interface {}
}

func (*Env) CreateCard(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	log.Println("Request received")
}
