// Package orderitem is responsible for managing information about order items in database-agnostic way.
// This package defines repository interface for abstracting interaction with particular database.
package orderitem

import (
	"context"

	"github.com/rtbe/clean-rest-api/domain/entity"
)

// Repository is an interface that represents persistent storage abstraction.
// This is a port in hexagonal architecture terms,
// so concrete implementation of database should implements the set of these methods.
type Repository interface {
	Create(ctx context.Context, newOrderItem entity.NewOrderItem) (entity.OrderItem, error)
	Query(ctx context.Context, lastSeenID, limit string) ([]entity.OrderItem, error)
	QueryByID(ctx context.Context, id string) (entity.OrderItem, error)
	QueryByOrderID(ctx context.Context, orderID string) ([]entity.OrderItem, error)
	Update(ctx context.Context, id string, updateOrderItem entity.UpdateOrderItem) error
	Delete(ctx context.Context, id string) error
	DeleteByOrderID(ctx context.Context, orderID string) error
}
