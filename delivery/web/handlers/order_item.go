package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"github.com/rtbe/clean-rest-api/domain/entity"
	"github.com/rtbe/clean-rest-api/domain/usecase"
	"github.com/rtbe/clean-rest-api/internal/database"
	"github.com/rtbe/clean-rest-api/internal/validation"
)

type OrderItemGroup struct {
	OrderItemService *usecase.OrderItemService
}

// swagger:route POST /products/ product createOrderItem
//
// Creates a new order item
// .
//
// Consumes:
// - application/json
// Produces:
// - application/json
//
// Responses:
//   200: OrderItem
//   500: errorResponse
func (oig *OrderItemGroup) CreateOrderItem(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var newOrderItem entity.NewOrderItem
	if err := json.NewDecoder(r.Body).Decode(&newOrderItem); err != nil {
		return err
	}

	if err := validation.Check(newOrderItem); err != nil {
		return RequestError{
			ErrorText: "validation error",
			Fields:    err.Error(),
			Status:    http.StatusBadRequest,
		}
	}

	orderItem, err := oig.OrderItemService.Create(ctx, newOrderItem)
	if err != nil {
		return err
	}

	return respond(ctx, w, orderItem, http.StatusCreated)
}

// swagger:route GET /orderItems/{id} orderItem getOrderItem
//
// Gets an order item by it\`s id
// and returns it\`s JSON representation.
//
// Produces:
// - application/json
//
// Responses:
//   201: OrderItem
//   500: errorResponse
func (oig *OrderItemGroup) GetOrderItem(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	id, err := parseURLParamID(r, "id")
	if err != nil {
		switch err {
		case database.ErrNotFound:
			return RequestError{
				ErrorText: err.Error(),
				Status:    http.StatusNotFound,
			}
		default:
			return errors.Wrapf(err, "ID: %s", id)
		}
	}

	orderItem, err := oig.OrderItemService.QueryByID(ctx, id)
	if err != nil {
		return err
	}

	return respond(ctx, w, orderItem, http.StatusOK)
}

// swagger:route PATCH /orders/user/{id} orderItem updateOrderItem
//
// Updates an order item
// .
//
// Consumes:
// - application/json
// Produces:
// - application/json
//
// Responses:
//   204: emptyResponse
//   500: errorResponse
func (oig *OrderItemGroup) UpdateOrderItem(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var updateOrderItem entity.UpdateOrderItem
	if err := json.NewDecoder(r.Body).Decode(&updateOrderItem); err != nil {
		return err
	}

	id, err := parseURLParamID(r, "id")
	if err != nil {
		return err
	}

	if err := oig.OrderItemService.Update(ctx, id, updateOrderItem); err != nil {
		return err
	}

	return respond(ctx, w, nil, http.StatusNoContent)
}

// swagger:route DELETE /orderItems/{id} orderItem deleteOrderItem
//
// Deletes an order item by it\`s id
// .
//
// Produces:
// - application/json
//
// Responses:
//   204: emptyResponse
//   500: errorResponse
func (oig *OrderItemGroup) DeleteOrderItem(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	id, err := parseURLParamID(r, "orderID")
	if err != nil {
		return err
	}

	if err := oig.OrderItemService.Delete(ctx, id); err != nil {
		return err
	}

	return respond(ctx, w, nil, http.StatusNoContent)
}

// swagger:route GET /orders/{orderID}/orderItems/ orderItem listOrderOrderItems
//
// Gets list of order items for particular order
//
// Produces:
// - application/json
//
// Responses:
//   200: []OrderItem
//   500: errorResponse
func (oig *OrderItemGroup) ListOrderOrderItems(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	id := chi.URLParam(r, "orderID")
	orders, err := oig.OrderItemService.QueryByOrderID(ctx, id)
	if err != nil {
		return err
	}

	return respond(ctx, w, orders, http.StatusOK)
}

// swagger:route DELETE /orders/{orderID}/orderItems/ order deleteOrderOrderItems
//
// Deletes all order items related to an particular order.
//
// Consumes:
// - application/json
//
// Responses:
//   204: emptyResponse
//   500: errorResponse
func (oig *OrderItemGroup) DeleteOrderOrderItems(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	id, err := parseURLParamID(r, "orderID")
	if err != nil {
		return err
	}

	if err := oig.OrderItemService.DeleteByOrderID(ctx, id); err != nil {
		return err
	}

	return respond(ctx, w, nil, http.StatusNoContent)
}
