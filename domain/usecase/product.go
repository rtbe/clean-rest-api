package usecase

import (
	"context"

	"github.com/rtbe/clean-rest-api/domain/entity"
	"github.com/rtbe/clean-rest-api/repository/product"
)

// Product is an interface that represents product business domain use case.
type Product interface {
	Create(ctx context.Context, newProduct entity.NewProduct) (entity.Product, error)
	Query(ctx context.Context, lastSeenID, limit string) ([]entity.Product, error)
	QueryByID(ctx context.Context, id string) (entity.Product, error)
	Update(ctx context.Context, id string, updateProduct entity.UpdateProduct) error
	Delete(ctx context.Context, id string) error
}

// ProductService is an business domain intermidiate layer
// between product entity and product DB layer (repository).
type ProductService struct {
	repo product.Repository
}

// NewProductService creates a new product entity service.
func NewProductService(r product.Repository) *ProductService {
	return &ProductService{
		repo: r,
	}
}

// Create creates a new product.
func (s *ProductService) Create(ctx context.Context, np entity.NewProduct) (entity.Product, error) {
	return s.repo.Create(ctx, np)
}

// Query gets a paginated list of products.
func (s *ProductService) Query(ctx context.Context, lastSeenID, limit string) ([]entity.Product, error) {
	return s.repo.Query(ctx, lastSeenID, limit)
}

// QueryByID queries product by given id.
func (s *ProductService) QueryByID(ctx context.Context, id string) (entity.Product, error) {
	return s.repo.QueryByID(ctx, id)
}

// Update updates particular product.
func (s *ProductService) Update(ctx context.Context, id string, up entity.UpdateProduct) error {
	return s.repo.Update(ctx, id, up)
}

// Delete deletes Product by given id.
func (s *ProductService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
