package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Message struct {
	Content string `json:"content"`
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		var m Message
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&m)
		//TODO: Actually log the errors instead of just printing them
		if err != nil {
			fmt.Println("Error with JSON decode")
			m.Content = "Your tale was lost on its way! Try speaking with the Tale-Teller again"
		}

		body, err := json.Marshal(m)

		if err != nil {
			fmt.Println("Error with json Marshal")
			body = []byte(`{
				"content":"Your tale was lost in the Astral Sea! Try again later!"
			}`)
		}
		//TODO:Make this url an environment variable to ingest on startup
		_, err = http.Post("{{DISCORD WEBHOOK URL}}", "application/json", bytes.NewBuffer(body))
		if err != nil {
			fmt.Println(err)
		}
	})
	http.ListenAndServe(":3000", r)
}
