package middlewares

import (
	"net/http"
)

// Cors is an middleware that allows to fetch data from
// different origin with request methods and headers restrictions.
func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//Allowed request origins
		w.Header().Set("Access-Control-Allow-Origin", "*")
		//Allowed request methods
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		//Allowed request headers
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}
