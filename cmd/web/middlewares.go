package main

import (
	"fmt"
	"log/slog"
	"net/http"
)

// Security related headers have been set as per recommendations in
// https://cheatsheetseries.owasp.org/cheatsheets/HTTP_Headers_Cheat_Sheet.html
func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")                              // Prevents the browser from interpreting files as a different MIME type than specified in the Content-Type header.
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")             // Send the origin, path, and query string when performing a same-origin request. For cross-origin requests send the origin (only) when the protocol security level stays same (HTTPS→HTTPS).
		w.Header().Set("Access-Control-Allow-Origin", "*")                               // Since we're providing a publicly available API, we're not setting any specific origins
		w.Header().Set("Access-Control-Allow-Headers", "GET, POST")                      // We're only allowing clients to use these HTTP methods
		w.Header().Set("Permissions-Policy", "geolocation=(), camera=(), microphone=()") // Disabling theese features as the site doesn't need them

		next.ServeHTTP(w, r)
	})
}

func commonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Server", "Go")

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			uri    = r.URL.RequestURI()
		)

		app.logger.Info(
			"received request >>>",
			slog.String("ip", ip),
			slog.String("proto", proto),
			slog.String("method", method),
			slog.String("uri", uri))

		next.ServeHTTP(w, r)
	})
}

// Note: our middleware will only recover panics that happen in the same goroutine that
// executed the recoverPanic() middleware.
func (app *application) recoverFromPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Create a deferred function (which will always be run in the event of a panic as
		// Go unwinds the stack).
		defer func() {
			// recover() retrieves the error value passed to the call of panic
			if err := recover(); err != nil {
				// This header acts as a trigger to make Go’s HTTP server automatically close
				// the current connection after a response has been sent. It also informs the
				// user that the connection will be closed
				w.Header().Set("Connection", "close")

				// Since the recover() function has the type "any", and its underlying type
				// could be string, error, or something else — we normalize this into an error
				// by using the fmt.Errorf() function to create a new error object containing
				// the default textual representation of the "any" value
				app.serverError(w, r, fmt.Errorf("%s", err))
			}
		}() // The "()" here can be used to pass arguments to the anon function (func) - if it had params

		next.ServeHTTP(w, r)
	})
}
