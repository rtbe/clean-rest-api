package middlewares

import (
	"fmt"
	"net/http"
	"time"

	"github.com/rtbe/clean-rest-api/internal/logger"
)

func Logger(l logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			info, err := GetRequestInfo(r.Context())
			if err != nil {
				return
			}

			// Log info about request before invoking next handler.
			start := fmt.Sprintf("%s: started  : %s %s -> %s",
				info.ID,
				r.Method, r.URL.Path, r.RemoteAddr,
			)
			l.Log("info", start)

			next.ServeHTTP(w, r)

			// Log info about request after invoking next handler.
			finish := fmt.Sprintf("%s: completed: %s %s -> %s (%d) (%s)",
				info.ID,
				r.Method, r.URL.Path, r.RemoteAddr,
				info.StatusCode, time.Since(info.Now),
			)
			l.Log("info", finish)
		})
	}
}
