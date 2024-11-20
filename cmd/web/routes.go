package main

import "net/http"

func (app *application) routes() http.Handler {

	mux := http.NewServeMux()

	// Jokes relalted routes
	mux.HandleFunc("GET /random-joke", app.randomJoke)
	mux.HandleFunc("GET /joke/{id}", app.getJoke)
	mux.HandleFunc("POST /joke", app.insertJoke)

	// Quotes related routes
	mux.HandleFunc("GET /random-quote", app.randomQuote)

	return mux
}
