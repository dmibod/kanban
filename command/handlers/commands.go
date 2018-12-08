package command

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/dmibod/kanban/tools/msg"
)

type Env struct{
   msg msg.Transport
}

func (env *Env) PostCommands(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
			http.Error(w, http.StatusText(405), 405)
			return
	}

	body, readErr := ioutil.ReadAll(r.Body)

	if readErr != nil {
			http.Error(w, http.StatusText(500), 500)
			return
	}

	commands := []command.Command
	jsonErr := json.Unmarshal(body, commands)

	if jsonErr != nil {
		http.Error(w, http.StatusText(500), 500)
		return
}

	bks, err := env.msg.AllBooks()
	if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
	}

	enc := json.NewEncoder(w)
  d := struct{
		Success bool `json:"success"`
	}{ true }
	enc.Encode(d)
}