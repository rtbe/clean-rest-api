// Package product is responsible for managing information about products in database-agnostic way.
// This package defines repository interface for abstracting interaction with particular database.
package product

import (
	"context"

	"github.com/rtbe/clean-rest-api/domain/entity"
)

// Repository is an interface that represents persistent storage abstraction.
// This is a port in hexagonal architecture terms,
// so concrete implementation of database should implements the set of these methods.
type Repository interface {
	Create(ctx context.Context, newProduct entity.NewProduct) (entity.Product, error)
	Query(ctx context.Context, lastSeenID, limit string) ([]entity.Product, error)
	QueryByID(ctx context.Context, id string) (entity.Product, error)
	Update(ctx context.Context, id string, updateProduct entity.UpdateProduct) error
	Delete(ctx context.Context, id string) error
}
