package command

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Env struct {
	msg chan<- []byte
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

	commands := []Command{}
	jsonErr := json.Unmarshal(body, &commands)

	if jsonErr != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	m, msgErr := json.Marshal(commands)
	if msgErr != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	env.msg <- m

	enc := json.NewEncoder(w)
	d := struct {
		Success bool `json:"success"`
	}{true}
	enc.Encode(d)
}
