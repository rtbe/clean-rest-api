package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	mid "github.com/rtbe/clean-rest-api/delivery/web/middlewares"
	"github.com/rtbe/clean-rest-api/domain/entity"
	"github.com/rtbe/clean-rest-api/domain/usecase"
	"github.com/rtbe/clean-rest-api/internal/database"
)

type UserGroup struct {
	UserService *usecase.UserService
}

// swagger:route GET /users/{lastSeenID}/{limit} user listUsers
//
// Gets paginated list of users.
// This request uses two provided values to implement pagination: last seen id and limit.
// Results of a request sorted by creation date of selected users and sended back as JSON.
//
// Produces:
// - application/json
//
// Responses:
//   200: []User
//   500: errorResponse
func (ug *UserGroup) ListUsers(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	lastSeenID, err := parseURLParamID(r, "lastSeenID")
	if err != nil {
		return err
	}

	limit := chi.URLParam(r, "limit")
	users, err := ug.UserService.Query(ctx, lastSeenID, limit)
	if err != nil {
		return err
	}

	return respond(ctx, w, users, http.StatusOK)
}

// swagger:route GET /users/id/{id} user getUser
//
// Gets a user by his id
// and returns it\`s JSON representation.
//
// Consumes:
// - application/json
// Produces:
// - application/json
//
// Responses:
//   201: User
//   500: errorResponse
func (ug *UserGroup) GetUserByID(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	id, err := parseURLParamID(r, "id")
	if err != nil {
		return err
	}

	user, err := ug.UserService.QueryByID(ctx, id)
	if err != nil {
		switch err {
		case database.ErrNotFound:
			return RequestError{
				ErrorText: err.Error(),
				Status:    http.StatusNotFound,
			}
		default:
			return errors.Wrapf(err, "ID: %s", id)
		}
	}

	return respond(ctx, w, user, http.StatusOK)
}

// swagger:route PATCH /users/{id} user updateUser
//
// Updates a user
// .
//
// Consumes:
// - application/json
//
// Responses:
//   204: emptyResponse
//   500: errorResponse
func (ug *UserGroup) UpdateUser(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	id, err := parseURLParamID(r, "id")
	if err != nil {
		return err
	}

	claims, err := mid.GetJWTClaims(ctx)
	if err != nil {
		return err
	}

	// Restrict modification of other user.
	if id != claims.Id {
		return RequestError{
			ErrorText: "modification of other user is restricted",
			Status:    http.StatusBadRequest,
		}
	}

	var user entity.UpdateUser
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return err
	}

	if err := ug.UserService.Update(ctx, id, user); err != nil {
		return err
	}

	return respond(ctx, w, nil, http.StatusNoContent)
}

// swagger:route DELETE /users/{id} user deleteUser
//
// Deletes a user
// .
//
// Produces:
// - application/json
//
// Responses:
//   204: emptyResponse
//   500: errorResponse
func (ug *UserGroup) DeleteUser(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	id, err := parseURLParamID(r, "id")
	if err != nil {
		return err
	}

	if err := ug.UserService.Delete(ctx, id); err != nil {
		return err
	}

	return respond(ctx, w, nil, http.StatusNoContent)
}
