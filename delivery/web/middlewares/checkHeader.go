package middlewares

import (
	"fmt"
	"net/http"
)

// CheckHeader is an middleware that`s checks validity of passed request header for set of routes
// And returns an error if it's not valid.
func CheckHeader(header, value string, status int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if str := r.Header.Get(header); str != value {
				http.Error(w, fmt.Sprintf("%s header should have value of %s", header, value), status)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// CheckJSONHeader is an middleware that checks for content-type: "application/json" header
// in request.
// And returns an error if it's not there.
func CheckJSONHeader() func(http.Handler) http.Handler {
	return CheckHeader("content-type", "application/json", http.StatusUnsupportedMediaType)
}
