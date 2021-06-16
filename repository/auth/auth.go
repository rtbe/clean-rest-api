// Package order is responsible for managing information about JWT tokens in database-agnostic way.
// This package defines repository interface for abstracting interaction with particular database.
package auth

import (
	"context"

	"github.com/rtbe/clean-rest-api/domain/entity"
)

// Repository is an interface that represents persistent storage abstraction.
// This is a port in hexagonal architecture terms,
// so concrete implementation of database should implements the set of these methods.
type Repository interface {
	Create(ctx context.Context, refreshToken entity.RefreshToken) error
	Refresh(ctx context.Context, userID string, refreshToken entity.RefreshToken) error
	Delete(ctx context.Context, userID string) error
}
