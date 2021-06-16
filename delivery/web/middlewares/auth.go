package middlewares

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/rtbe/clean-rest-api/domain/entity"
)

var (
	errAuthWrongHeaderFormat = errors.New("expected 'Authorization' header format: Bearer <token>")
	errAuthHeaderMissing     = errors.New("header 'Authorization' is missing")
)

// ClaimsKey is the context.Context key to store authentication claims.
var ClaimsKey = &contextKey{"claims"}

// Authenticate is an middleware that validates passed access JWT token in `Authorization` header.
// So it serves as gateway to all of app incoming requests.
// Access token claims then passed into request context so you can get them later with ClaimsCtxKey.
func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Expecting: bearer <token>
		bearer := r.Header.Get("Authorization")
		if len(bearer) == 0 {
			http.Error(w, errAuthHeaderMissing.Error(), http.StatusBadRequest)
			return
		}
		if len(bearer) < 6 || strings.ToLower(bearer[0:6]) != "bearer" {
			http.Error(w, errAuthWrongHeaderFormat.Error(), http.StatusBadRequest)
			return
		}

		accessToken := bearer[7:]
		accessTokenClaims, err := entity.ParseAccessTokenClaims(accessToken)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Add claims to the context so we can retrieve them later
		ctx = context.WithValue(ctx, ClaimsKey, accessTokenClaims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Authorize is an middleware that filters requests based on provided role in access token claims.
// So it serves as gateway to particular routes in the app.
func Authorize(requiredRole string) func(http.Handler) http.Handler {
	requiredRole = strings.ToLower(requiredRole)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			claims, err := GetJWTClaims(ctx)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Check role inside of an array of roles in access token claims.
			var role bool
			for _, userRole := range claims.User_roles {
				if strings.EqualFold(requiredRole, userRole) {
					role = true
				}
			}
			if !role {
				s := fmt.Sprintf(
					"you are not authorized for that action; get roles: %v, expected: %s",
					claims.User_roles,
					requiredRole,
				)
				http.Error(w, s, http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// GetJWTClaims returns JWT access token claims from the given context.
func GetJWTClaims(ctx context.Context) (*entity.AccessTokenClaims, error) {
	if ctx == nil {
		return &entity.AccessTokenClaims{}, errNoContext
	}

	claims, ok := ctx.Value(ClaimsKey).(*entity.AccessTokenClaims)
	if !ok {
		return &entity.AccessTokenClaims{}, errNoClaimsInContext
	}

	return claims, nil
}
