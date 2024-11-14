package main

import "net/http"

func (app *application) routes() http.Handler {

	mux := http.NewServeMux()

	// Jokes relalted routes
	mux.HandleFunc("/random-joke", app.randomJoke)

	// Quotes related routes
	mux.HandleFunc("/random-qupte", app.randomQuote)

	return mux
}
