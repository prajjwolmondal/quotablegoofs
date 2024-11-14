package main

import "net/http"

func (app *application) randomJoke(w http.ResponseWriter, r *http.Request) {
	joke := `{"joke": ["I told them I wanted to be a comedian, and they laughed; then I became a comedian, no one's laughing now"], "source": "Unknown"}`

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(joke))
}

func (app *application) randomQuote(w http.ResponseWriter, r *http.Request) {
	quote := `{"quote": "An ounce of action is worth a ton of theory.", "source": "Friedrich Engels"}`

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(quote))
}
