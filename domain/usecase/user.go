package usecase

import (
	"context"

	"github.com/rtbe/clean-rest-api/domain/entity"
	"github.com/rtbe/clean-rest-api/repository/user"
)

// User is an interface that represents user business domain use case.
type User interface {
	Create(ctx context.Context, newUser entity.NewUser) (entity.User, error)
	Query(ctx context.Context, lastSeenUserID, limit string) ([]entity.User, error)
	QueryByID(ctx context.Context, id string) (entity.User, error)
	Update(ctx context.Context, id string, updateUser entity.UpdateUser) error
	Delete(ctx context.Context, id string) error
}

// UserService is an business domain intermidiate layer
// between user entity and user DB layer (repository).
type UserService struct {
	repo user.Repository
}

// NewUserService creates a new user entity service.
func NewUserService(r user.Repository) *UserService {
	return &UserService{
		repo: r,
	}
}

// Create creates a new user.
func (s *UserService) Create(ctx context.Context, nu entity.NewUser) (entity.User, error) {
	return s.repo.Create(ctx, nu)
}

// Query gets a paginated list of users.
func (s *UserService) Query(ctx context.Context, lastSeenUserID, limit string) ([]entity.User, error) {
	return s.repo.Query(ctx, lastSeenUserID, limit)
}

// QueryByID queries a user by his id.
func (s *UserService) QueryByID(ctx context.Context, id string) (entity.User, error) {
	return s.repo.QueryByID(ctx, id)
}

// Update updates particular user.
func (s *UserService) Update(ctx context.Context, id string, uu entity.UpdateUser) error {
	return s.repo.Update(ctx, id, uu)
}

// Delete deletes user by his id.
func (s *UserService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
