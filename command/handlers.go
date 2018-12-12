package command

import (
	"log"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Env struct {
	CommandQueue chan<- []byte
}

func (env *Env) PostCommands(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), 405)
		log.Println("Wrong HTTP method")
		return
	}

	body, readErr := ioutil.ReadAll(r.Body)

	if readErr != nil {
		http.Error(w, http.StatusText(500), 500)
		log.Println("Error reading body", readErr)
		return
	}

	commands := []Command{}
	jsonErr := json.Unmarshal(body, &commands)

	if jsonErr != nil {
		http.Error(w, http.StatusText(500), 500)
		log.Println("Error parsing json", jsonErr)
		return
	}

	log.Printf("Commands received: %+v\n", commands);

	m, msgErr := json.Marshal(commands)
	if msgErr != nil {
		http.Error(w, http.StatusText(500), 500)
		log.Println("Error marshalling commands", msgErr)
		return
	}

	env.CommandQueue <- m

	enc := json.NewEncoder(w)
	d := struct {
		Count int `json:"count"`
		Success bool `json:"success"`
	}{len(commands),true}
	
	enc.Encode(d)

	log.Printf("Commands sent: %+v\n", len(commands))
}
