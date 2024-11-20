package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"runtime/debug"
)

// Writes a log entry at Error level (including the request method and URI as attributes),
// then sends a generic 500 Internal Server Error response to the user.
func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)

	app.logger.Error(err.Error(), slog.String("method", method), slog.String("uri", uri), slog.String("trace", trace))
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// Sends a specific status code and corresponding description to the user. Should be used to send
// responses like 400 "Bad Request" when there's a problem with the request that the user sent.
func (app *application) clientError(w http.ResponseWriter, r *http.Request, status int, errs []error) {

	var errorMessages []string
	for _, err := range errs {
		errorMessages = append(errorMessages, err.Error())
	}

	type ErrorResponse struct {
		Errors []string `json:"errors"`
	}

	response := ErrorResponse{
		Errors: errorMessages,
	}
	jsonData, err := json.Marshal(response)
	if err != nil {
		app.serverError(w, r, err)
	}
	http.Error(w, string(jsonData), status)
}
