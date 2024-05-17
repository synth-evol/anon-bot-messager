package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

type Message struct {
	Content string `json:"content"`
}

func main() {
	hookUrl := os.Getenv("DISCORD_URL")

	//Zerolog setup
	logFile, err := os.OpenFile(
		"messager.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	if err != nil {
		//Don't run without logging
		panic(err)
	}

	//Clean up after ourselves
	defer logFile.Close()

	logger := zerolog.New(logFile).With().Timestamp().Logger()

	//Handler setup with HTTP logging from chi middleware
	//Separates error logs from usage logs which don't necessarily need storage in a file
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		var m Message
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&m)

		if err != nil {
			logger.Error().Msg("Error with JSON decode")
			m.Content = "Your tale was lost on its way! Try speaking with the Tale-Teller again"
		}

		body, err := json.Marshal(m)

		if err != nil {
			logger.Error().Msg("Error with json Marshal")
			body = []byte(`{
				"content":"Your tale was lost in the Astral Sea! Try again later!"
			}`)
		}

		_, err = http.Post(hookUrl, "application/json", bytes.NewBuffer(body))
		if err != nil {
			logger.Error().Msg(err.Error())
		}
	})
	http.ListenAndServe(":8000", r)
}
