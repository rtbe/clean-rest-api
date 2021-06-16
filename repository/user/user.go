// Package user is responsible for managing information about users in database-agnostic way.
// This package defines repository interface for abstracting interaction with particular database.
package user

import (
	"context"

	"github.com/rtbe/clean-rest-api/domain/entity"
)

// Repository is an interface that represents persistent storage abstraction.
// This is a port in hexagonal architecture terms,
// so concrete implementation of database should implements the set of these methods.
type Repository interface {
	Create(ctx context.Context, newUser entity.NewUser) (entity.User, error)
	Query(ctx context.Context, lastSeenUserID, limit string) ([]entity.User, error)
	QueryByID(ctx context.Context, id string) (entity.User, error)
	Update(ctx context.Context, userID string, user entity.UpdateUser) error
	Delete(ctx context.Context, id string) error
	DeleteByUserName(ctx context.Context, userName string) error
}
