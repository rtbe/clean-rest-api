package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/rtbe/clean-rest-api/domain/entity"
	"github.com/rtbe/clean-rest-api/domain/usecase"
	"github.com/rtbe/clean-rest-api/internal/validation"
)

type AuthGroup struct {
	AuthService *usecase.AuthService
}

// swagger:route POST /auth/signup auth signUp
//
// Creates a new user.
//
//
// Consumes:
// - application/json
// Produces:
// - application/json
//
// Responses:
//   201: User
//   500: errorResponse
func (ag *AuthGroup) SignUp(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var newUser entity.NewUser
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		return err
	}

	if err := validation.Check(newUser); err != nil {
		return RequestError{
			ErrorText: "validation error",
			Fields:    err.Error(),
			Status:    http.StatusBadRequest,
		}
	}

	user, err := ag.AuthService.SignUp(ctx, newUser)
	if err != nil {
		// Check if user already registered
		if strings.Contains(err.Error(), "duplicate") {
			return nil
		}
		return err
	}

	return respond(ctx, w, user, http.StatusCreated)
}

// swagger:route POST /auth/signin auth signIn
//
// Issues pair of access/refresh tokens.
//
// Consumes:
// - application/json
// Produces:
// - application/json
//
// Responses:
//   200: TokenPair
//   500: errorResponse
func (ag *AuthGroup) SignIn(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var decodedUser entity.User
	err := json.NewDecoder(r.Body).Decode(&decodedUser)
	if err != nil {
		return err
	}

	tokenPair, err := ag.AuthService.SignIn(ctx, decodedUser)
	if err != nil {
		// Check error if tokens already issued
		if strings.Contains(err.Error(), "duplicate") {
			return nil
		}
		return nil
	}

	return respond(ctx, w, tokenPair, http.StatusOK)
}

// swagger:route POST /auth/signout auth signOut
//
// Deletes a refresh token belonging to specific  user.
//
// Consumes:
// - application/json
// Produces:
// - application/json
//
// Responses:
//   204: emptyResponse
//   500: errorResponse
func (ag *AuthGroup) SignOut(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var decodedUser entity.User
	err := json.NewDecoder(r.Body).Decode(&decodedUser)
	if err != nil {
		return err
	}

	err = ag.AuthService.SignOut(ctx, decodedUser)
	if err != nil {
		return err
	}

	return respond(ctx, w, nil, http.StatusNoContent)
}

// swagger:route POST /refreshTokens auth refreshTokens
//
// Receives pair of access/refresh tokens and returns fresh pair.
//
// Consumes:
// - application/json
// Produces:
// - application/json
//
// Responses:
//   200: TokenPair
//   500: errorResponse
func (ag *AuthGroup) RefreshTokens(w http.ResponseWriter, r *http.Request) error {
	//ctx := r.Context()

	var decodedTokenPair entity.TokenPair
	err := json.NewDecoder(r.Body).Decode(&decodedTokenPair)
	if err != nil {
		return err
	}

	accessToken := decodedTokenPair.AccessToken
	accessTokenClaims, err := entity.ParseAccessTokenClaims(accessToken)
	if err != nil {
		return err
	}

	refreshToken := decodedTokenPair.RefreshToken
	refreshTokenClaims, err := entity.ParseRefreshTokenClaims(refreshToken)
	if err != nil {
		return err
	}
	if accessTokenClaims.Refresh_uuid != refreshTokenClaims.UUID {
		return err
	}

	//respTokenPair, err := ag.AuthService.Refresh(ctx, &entity.JWTTokenPair{AccessToken: refreshTokenClaims})
	//if err != nil {
	//	return err
	//}

	return nil
	//return respond(ctx, w, respTokenPair, http.StatusOK)
}
