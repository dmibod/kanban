package main

func main() {
	db, err := models.NewDB("postgres://user:pass@localhost/bookstore")
	if err != nil {
			log.Panic(err)
	}

	env := &Env{db}

	http.HandleFunc("/commands", env.booksIndex)
	http.ListenAndServe(":3000", nil)
}