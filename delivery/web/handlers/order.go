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

type OrderGroup struct {
	OrderService *usecase.OrderService
}

// swagger:route POST /orders/ order createOrder
//
// Creates a new order
// .
//
// Consumes:
// - application/json
// Produces:
// - application/json
//
// Responses:
//   201: Order
//   500: errorResponse
func (og *OrderGroup) CreateOrder(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var newOrder entity.NewOrder
	if err := json.NewDecoder(r.Body).Decode(&newOrder); err != nil {
		return err
	}

	if err := validation.Check(newOrder); err != nil {
		return RequestError{
			ErrorText: "validation error",
			Fields:    err.Error(),
			Status:    http.StatusBadRequest,
		}
	}

	order, err := og.OrderService.Create(ctx, newOrder)
	if err != nil {
		return err
	}

	return respond(ctx, w, order, http.StatusCreated)
}

// swagger:route GET /orders/{lastSeenID}/{id} order listOrders
//
// Gets paginated list of orders.
// This request uses two provided values to implement pagination: last seen id and limit.
// Results of a request sorted by creation date of selected users and sended back as JSON.
//
// Produces:
// - application/json
//
// Responses:
//   200: []Order
//   500: errorResponse
func (og *OrderGroup) ListOrders(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	lastSeenID, err := parseURLParamID(r, "lastSeenID")
	if err != nil {
		return err
	}
	limit := chi.URLParam(r, "limit")

	orders, err := og.OrderService.Query(ctx, lastSeenID, limit)
	if err != nil {
		return err
	}

	return respond(ctx, w, orders, http.StatusOK)
}

// swagger:route GET /orders/{id} order getOrder
//
// Gets an order by it\`s id
// and returns it\`s JSON representation.
//
// Produces:
// - application/json
//
// Responses:
//   200: Order
//   500: errorResponse
func (og *OrderGroup) GetOrder(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	id, err := parseURLParamID(r, "id")
	if err != nil {
		return err
	}

	order, err := og.OrderService.QueryByID(ctx, id)
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

	return respond(ctx, w, order, http.StatusOK)
}

// swagger:route GET /orders/user/{userID} order listUserOrders
//
// Gets orders by user id
// and returns their`s JSON representation.
//
// Produces:
// - application/json
//
// Responses:
//   200: []Order
//   500: errorResponse
func (og *OrderGroup) ListUserOrders(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	id, err := parseURLParamID(r, "userID")
	if err != nil {
		return err
	}

	orders, err := og.OrderService.QueryByUserID(ctx, id)
	if err != nil {
		return err
	}

	return respond(ctx, w, orders, http.StatusOK)
}

// swagger:route PATCH /orders/{id} order updateOrder
//
// Updates a specific order
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
func (og *OrderGroup) UpdateOrder(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var updateOrder entity.UpdateOrder
	if err := json.NewDecoder(r.Body).Decode(&updateOrder); err != nil {
		return err
	}

	id, err := parseURLParamID(r, "id")
	if err != nil {
		return err
	}

	if err := og.OrderService.Update(ctx, id, updateOrder); err != nil {
		return err
	}

	return respond(ctx, w, nil, http.StatusNoContent)
}

// swagger:route DELETE /orders/{id} order deleteOrder
//
// Deletes an order by it\`s id
// and returns it\`s JSON representation.
//
// Produces:
// - application/json
//
// Responses:
//   204: emptyResponse
//   500: errorResponse
func (og *OrderGroup) DeleteOrder(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	id, err := parseURLParamID(r, "id")
	if err != nil {
		return err
	}

	if err := og.OrderService.Delete(ctx, id); err != nil {
		return err
	}

	return respond(ctx, w, nil, http.StatusNoContent)
}

// swagger:route DELETE /orders/user/{id} order deleteUserOrder
//
// Deletes all orders belonging to a specific user.
//
// Consumes:
// - application/json
//
// Responses:
//   204:
//   500: errorResponse
func (og *OrderGroup) DeleteUserOrders(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	id, err := parseURLParamID(r, "userID")
	if err != nil {
		return err
	}

	if err := og.OrderService.DeleteByUserID(ctx, id); err != nil {
		return err
	}

	return respond(ctx, w, nil, http.StatusNoContent)
}
