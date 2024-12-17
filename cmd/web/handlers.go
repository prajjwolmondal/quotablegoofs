package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"quotablegooofs.prajjmon.net/internal/models"
)

func (app *application) randomJoke(w http.ResponseWriter, r *http.Request) {
	var limit int
	if len(r.URL.Query().Get("limit")) == 0 {
		limit = 1
	} else {
		var err error
		limit, err = strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			app.serverError(w, r, err)
			return
		}
	}

	jokes, err := app.jokes.GetRandomJokes(limit)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	jsonJokes, err := json.Marshal(jokes)
	if err != nil {
		app.serverError(w, r, err)
	}

	w.Write(jsonJokes)
}

func (app *application) randomQuote(w http.ResponseWriter, r *http.Request) {
	var limit int
	if len(r.URL.Query().Get("limit")) == 0 {
		limit = 1
	} else {
		var err error
		limit, err = strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			app.serverError(w, r, err)
			return
		}
	}

	quotes, err := app.quotes.GetRandomQuotes(limit)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	jsonQuotes, err := json.Marshal(quotes)
	if err != nil {
		app.serverError(w, r, err)
	}

	w.Write(jsonQuotes)
}

func (app *application) getJoke(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	joke, err := app.jokes.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
			return
		} else {
			app.serverError(w, r, err)
			return
		}
	}

	jsonJoke, err := json.Marshal(joke)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.Write(jsonJoke)
}

func (app *application) insertJoke(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	body = bytes.TrimSpace(body)
	var joke models.Joke
	err = json.Unmarshal([]byte(body), &joke)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	if errs := joke.Validate(); len(errs) > 0 {
		app.clientError(w, r, http.StatusBadRequest, errs)
		return
	}

	if len(joke.Source) == 0 {
		joke.Source = "Unknown"
	}

	now := time.Now()
	joke.CreatedAt = now
	joke.UpdatedAt = now

	joke, err = app.jokes.Insert(joke)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	jokeAsJson, err := json.Marshal(joke)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.Write(jokeAsJson)
}

func (app *application) getQuote(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	quote, err := app.quotes.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
			return
		} else {
			app.serverError(w, r, err)
			return
		}
	}

	jsonQuote, err := json.Marshal(quote)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.Write(jsonQuote)
}

func (app *application) insertQuote(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	body = bytes.TrimSpace(body)
	var quote models.Quote
	err = json.Unmarshal([]byte(body), &quote)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	if errs := quote.Validate(); len(errs) > 0 {
		app.clientError(w, r, http.StatusBadRequest, errs)
		return
	}

	if len(quote.Source) == 0 {
		quote.Source = "Unknown"
	}

	now := time.Now()
	quote.CreatedAt = now
	quote.UpdatedAt = now

	quote, err = app.quotes.Insert(quote)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	quoteJson, err := json.Marshal(quote)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.Write(quoteJson)
}
