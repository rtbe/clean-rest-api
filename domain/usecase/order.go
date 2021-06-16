package usecase

import (
	"context"

	"github.com/rtbe/clean-rest-api/domain/entity"
	"github.com/rtbe/clean-rest-api/repository/order"
)

// Order is an interface that represents order domain use case.
type Order interface {
	Create(ctx context.Context, newOrder entity.NewOrder) (entity.Order, error)
	Query(ctx context.Context, lastSeenID, limit string) ([]entity.Order, error)
	QueryByID(ctx context.Context, id string) (entity.Order, error)
	QueryByUserID(ctx context.Context, userID string) ([]entity.Order, error)
	Update(ctx context.Context, id string, updateOrder entity.UpdateOrder) error
	Delete(ctx context.Context, id string) error
	DeleteByUserID(ctx context.Context, userID string) error
}

// OrderService is an business domain intermidiate layer
// between order entity, order item entity and database layer (repository).
type OrderService struct {
	orderRepo order.Repository
}

// NewOrderService creates a new order entity service.
func NewOrderService(orderRepo order.Repository) *OrderService {
	return &OrderService{
		orderRepo: orderRepo,
	}
}

// Create creates a new order.
func (s *OrderService) Create(ctx context.Context, no entity.NewOrder) (entity.Order, error) {
	return s.orderRepo.Create(ctx, no)
}

// Query gets a paginated list of orders.
func (s *OrderService) Query(ctx context.Context, lastSeenID, limit string) ([]entity.Order, error) {
	return s.orderRepo.Query(ctx, lastSeenID, limit)
}

// QueryByID queries a specific order.
func (s *OrderService) QueryByID(ctx context.Context, id string) (entity.Order, error) {
	return s.orderRepo.QueryByID(ctx, id)
}

// QueryByUserID queries all orders belonging to a specific user.
func (s *OrderService) QueryByUserID(ctx context.Context, userID string) ([]entity.Order, error) {
	return s.orderRepo.QueryByUserID(ctx, userID)
}

// Update updates a specific order.
func (s *OrderService) Update(ctx context.Context, id string, uo entity.UpdateOrder) error {
	return s.orderRepo.Update(ctx, id, uo)
}

// Delete deletes a specific order.
func (s *OrderService) Delete(ctx context.Context, id string) error {
	return s.orderRepo.Delete(ctx, id)
}

// DeleteByUserID deletes orders belonging to specific user.
func (s *OrderService) DeleteByUserID(ctx context.Context, userID string) error {
	return s.orderRepo.DeleteByUserID(ctx, userID)
}
