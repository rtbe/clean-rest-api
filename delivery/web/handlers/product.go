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

type ProductGroup struct {
	ProductService *usecase.ProductService
}

// swagger:route POST /products/ product createProduct
//
// Creates a new product
// .
//
// Consumes:
// - application/json
// Produces:
// - application/json
//
// Responses:
//   201: Product
//   500: errorResponse
func (pg *ProductGroup) CreateProduct(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var newProduct entity.NewProduct
	if err := json.NewDecoder(r.Body).Decode(&newProduct); err != nil {
		return err
	}

	if err := validation.Check(newProduct); err != nil {
		return RequestError{
			ErrorText: "validation error",
			Fields:    err.Error(),
			Status:    http.StatusBadRequest,
		}
	}

	product, err := pg.ProductService.Create(ctx, newProduct)
	if err != nil {
		return err
	}

	return respond(ctx, w, product, http.StatusCreated)
}

// swagger:route GET /products/{lastSeenID}/{id} product listProducts
//
// Gets paginated list of products.
// This request uses two provided values to implement pagination: last seen id and limit.
// Results of a request sorted by creation date of selected users and sended back as JSON.
//
// Produces:
// - application/json
//
// Responses:
//   200: []Product
//   500: errorResponse
func (pg *ProductGroup) ListProducts(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	lastSeenID, err := parseURLParamID(r, "lastSeenID")
	if err != nil {
		return err
	}
	limit := chi.URLParam(r, "limit")

	products, err := pg.ProductService.Query(ctx, lastSeenID, limit)
	if err != nil {
		return err
	}

	return respond(ctx, w, products, http.StatusOK)
}

// swagger:route GET /products/{id} product getProduct
//
// Gets a product by it\`s id
// and returns it\`s JSON representation.
//
// Produces:
// - application/json
//
// Responses:
//   200: Product
//   500: errorResponse
func (pg *ProductGroup) GetProduct(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	id, err := parseURLParamID(r, "id")
	if err != nil {
		return err
	}

	product, err := pg.ProductService.QueryByID(ctx, id)
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

	return respond(ctx, w, product, http.StatusOK)
}

// swagger:route PATCH /products/{id} product updateProduct
//
// Updates a product
// .
//
// Consumes:
// - application/json
//
// Responses:
//   201: emptyResponse
//   500: errorResponse
func (pg *ProductGroup) UpdateProduct(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var updateProduct entity.UpdateProduct
	if err := json.NewDecoder(r.Body).Decode(&updateProduct); err != nil {
		return err
	}

	id, err := parseURLParamID(r, "id")
	if err != nil {
		return err
	}

	if err := pg.ProductService.Update(ctx, id, updateProduct); err != nil {
		return err
	}

	return respond(ctx, w, nil, http.StatusCreated)
}

// swagger:route DELETE /products/{id} product deleteProduct
//
// Deletes a product by it\`s id
// and returns it\`s JSON representation.
//
// Produces:
// - application/json
//
// Responses:
//   204: emptyResponse
//   500: errorResponse
func (pg *ProductGroup) DeleteProduct(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	id, err := parseURLParamID(r, "id")
	if err != nil {
		return err
	}

	if err := pg.ProductService.Delete(ctx, id); err != nil {
		return err
	}

	return respond(ctx, w, nil, http.StatusNoContent)
}
