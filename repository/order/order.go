// Package order is responsible for managing information about orders in database-agnostic way.
// This package defines repository interface for abstracting interaction with particular database.
package order

import (
	"context"

	"github.com/rtbe/clean-rest-api/domain/entity"
)

// Repository is an interface that represents persistent storage abstraction.
// This is a port in hexagonal architecture terms,
// so concrete implementation of database should implements the set of these methods.
type Repository interface {
	Create(ctx context.Context, newOrder entity.NewOrder) (entity.Order, error)
	Query(ctx context.Context, lastSeenID, limit string) ([]entity.Order, error)
	QueryByID(ctx context.Context, id string) (entity.Order, error)
	QueryByUserID(ctx context.Context, userID string) ([]entity.Order, error)
	Update(ctx context.Context, id string, updateOrder entity.UpdateOrder) error
	Delete(ctx context.Context, id string) error
	DeleteByUserID(ctx context.Context, userID string) error
}
