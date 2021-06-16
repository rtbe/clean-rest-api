// package middlewares contains the set of middlewares.
package middlewares

import "errors"

var (
	errNoContext         = errors.New("there is no request context")
	errNoClaimsInContext = errors.New("there is no JWT claims in request context")
)

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation. This technique
// for defining context keys was copied from Go 1.7's new use of context in net/http.
type contextKey struct {
	name string
}
