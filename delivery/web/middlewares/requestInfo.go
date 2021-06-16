package middlewares

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
)

var (
	errNoRequestInfoInContext = errors.New("error: request info: there is no request info in request context")
)

// Request contains information about each request
type Request struct {
	ID         string
	Now        time.Time
	StatusCode int
}

var RequestKey = &contextKey{"requestInfo"}

// RequestInfo is an middleware that injects information about each passing request into it's context,
// so it can be used later for logging purposes.
func RequestInfo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		info := Request{
			ID:  uuid.NewString(),
			Now: time.Now().UTC(),
		}

		ctx := context.WithValue(r.Context(), RequestKey, &info)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetRequestInfo returns information about request from the given context.
func GetRequestInfo(ctx context.Context) (*Request, error) {
	if ctx == nil {
		return nil, errNoContext
	}

	req, ok := ctx.Value(RequestKey).(*Request)
	if !ok {
		return nil, errNoRequestInfoInContext
	}

	return req, nil
}
