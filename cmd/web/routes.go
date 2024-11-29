package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	mux := http.NewServeMux()

	chain := alice.New(app.recoverFromPanic, app.logRequests, commonHeaders)

	// Jokes relalted routes
	mux.Handle("GET /random-joke", chain.ThenFunc(app.randomJoke))
	mux.Handle("GET /joke/{id}", chain.ThenFunc(app.getJoke))
	mux.Handle("POST /joke", chain.ThenFunc(app.insertJoke))

	// Quotes related routes
	mux.Handle("GET /random-quote", chain.ThenFunc(app.randomQuote))
	mux.Handle("GET /quote/{id}", chain.ThenFunc(app.getQuote))
	mux.Handle("POST /quote", chain.ThenFunc(app.insertQuote))

	// flow of control (left to right) looks like this:
	// 		recoverFromPanic ↔ logRequest ↔ commonHeaders ↔ servemux ↔ application handler
	return chain.Then(mux)
}
