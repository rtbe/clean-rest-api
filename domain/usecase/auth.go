// Package usecase is responsible for managing application business logic.
// This package defines set of interfaces that abstract business logic related to particular entities.
// It also defines set of services (implementations of these interfaces)
// that connects business logic with database abstraction layer.
package usecase

import (
	"context"

	"github.com/pkg/errors"
	"github.com/rtbe/clean-rest-api/domain/entity"
	"github.com/rtbe/clean-rest-api/repository/auth"
	"golang.org/x/crypto/bcrypt"
)

// Auth is an interface that represents authentication business domain use case.
type Auth interface {
	SignUp(ctx context.Context, newUser entity.NewUser) (entity.User, error)
	SignIn(ctx context.Context, user entity.User) (entity.TokenPair, error)
	SignOut(ctx context.Context, user entity.User) error
	Refresh(ctx context.Context, tokenPair *entity.TokenPair) (entity.TokenPair, error)
}

// AuthService is an business domain intermidiate layer
// between auth entity and User DB layer (repository).
type AuthService struct {
	authRepo    auth.Repository
	userService User
}

// NewAuthService creates a new auth entity service.
func NewAuthService(ar auth.Repository, u User) *AuthService {
	return &AuthService{
		authRepo:    ar,
		userService: u,
	}
}

// SignUp creates a new user.
func (s *AuthService) SignUp(ctx context.Context, nu entity.NewUser) (entity.User, error) {
	u, err := s.userService.Create(ctx, nu)
	if err != nil {
		return entity.User{}, err
	}
	return u, nil
}

// SignIn issues pair of access and refresh token for particular user.
func (s *AuthService) SignIn(ctx context.Context, user entity.User) (entity.TokenPair, error) {
	u, err := s.userService.QueryByID(ctx, user.UserName)
	if err != nil {
		return entity.TokenPair{}, err
	}

	if isPasswordValid := bcrypt.CompareHashAndPassword(u.Password, user.Password); isPasswordValid == nil {
		return entity.TokenPair{}, errors.Wrap(isPasswordValid, "user password is not valid")
	}

	JWTTokenPair, err := entity.NewTokenPair(u.ID, u.Roles)
	if err != nil {
		return entity.TokenPair{}, err
	}

	err = s.authRepo.Create(ctx, JWTTokenPair.RefreshToken)
	if err != nil {
		return entity.TokenPair{}, err
	}

	return entity.TokenPair{
		AccessToken:  JWTTokenPair.AccessToken.Token,
		RefreshToken: JWTTokenPair.RefreshToken.Token,
	}, nil
}

// SignOut deletes refresh token for particular user.
func (s *AuthService) SignOut(ctx context.Context, user entity.User) error {

	u, err := s.userService.QueryByID(ctx, user.UserName)
	if err != nil {
		return err
	}

	if isPasswordValid := bcrypt.CompareHashAndPassword(u.Password, user.Password); isPasswordValid == nil {
		return errors.Wrap(isPasswordValid, "user password is not valid")
	}

	err = s.authRepo.Delete(ctx, u.ID)
	if err != nil {
		return err
	}

	return nil
}

// Refresh refreshes access and refresh JWT tokens.
func (s *AuthService) Refresh(ctx context.Context, tokenPair *entity.JWTTokenPair) (entity.TokenPair, error) {

	accessTokenClaims, err := entity.ParseAccessTokenClaims(tokenPair.AccessToken.Token)
	if err != nil {
		return entity.TokenPair{}, err
	}

	JWTTokenPair, err := entity.NewTokenPair(accessTokenClaims.User_id, accessTokenClaims.User_roles)
	if err != nil {
		return entity.TokenPair{}, err
	}

	err = s.authRepo.Refresh(ctx, accessTokenClaims.User_id, JWTTokenPair.RefreshToken)
	if err != nil {
		return entity.TokenPair{}, err
	}

	return entity.TokenPair{
		AccessToken:  JWTTokenPair.AccessToken.Token,
		RefreshToken: JWTTokenPair.RefreshToken.Token,
	}, nil
}
