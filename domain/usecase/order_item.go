package usecase

import (
	"context"

	"github.com/rtbe/clean-rest-api/domain/entity"
	orderitem "github.com/rtbe/clean-rest-api/repository/order_item"
)

// OrderItem is an interface that represents order item domain use case.
type OrderItem interface {
	Create(ctx context.Context, newOrderItem entity.NewOrderItem) (entity.OrderItem, error)
	Query(ctx context.Context, lastSeenID, limit string) ([]entity.OrderItem, error)
	QueryByID(ctx context.Context, id string) (entity.OrderItem, error)
	QueryByOrderID(ctx context.Context, orderID string) ([]entity.OrderItem, error)
	Update(ctx context.Context, id string, updateOrderItem entity.UpdateOrderItem) error
	Delete(ctx context.Context, id string) error
	DeleteByOrderID(ctx context.Context, orderID string) error
}

// OrderItemService is an business domain intermidiate layer
// between order entity and DB layer (repository).
type OrderItemService struct {
	repo orderitem.Repository
}

// NewOrderItemService creates a new order item entity service.
func NewOrderItemService(r orderitem.Repository) *OrderItemService {
	return &OrderItemService{
		repo: r,
	}
}

// Create creates a new order item.
func (s *OrderItemService) Create(ctx context.Context, no entity.NewOrderItem) (entity.OrderItem, error) {
	return s.repo.Create(ctx, no)
}

// Query gets a paginated list of orders items.
func (s *OrderItemService) Query(ctx context.Context, lastSeenID, limit string) ([]entity.OrderItem, error) {
	return s.repo.Query(ctx, lastSeenID, limit)
}

// QueryByID gets an order item by given id.
func (s *OrderItemService) QueryByID(ctx context.Context, id string) (entity.OrderItem, error) {
	return s.repo.QueryByID(ctx, id)
}

// QueryByOrderID gets all orders items items for particular order.
func (s *OrderItemService) QueryByOrderID(ctx context.Context, orderID string) ([]entity.OrderItem, error) {
	return s.repo.QueryByOrderID(ctx, orderID)
}

// Update updates order item.
func (s *OrderItemService) Update(ctx context.Context, id string, uoi entity.UpdateOrderItem) error {
	return s.repo.Update(ctx, id, uoi)
}

// Delete deletes an order item by given id.
func (s *OrderItemService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// DeleteByOrderID deletes all order items for particular order.
func (s *OrderItemService) DeleteByOrderID(ctx context.Context, orderID string) error {
	return s.repo.DeleteByOrderID(ctx, orderID)
}
