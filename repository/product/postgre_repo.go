package product

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/rtbe/clean-rest-api/domain/entity"
	"github.com/rtbe/clean-rest-api/internal/database"
	"github.com/rtbe/clean-rest-api/internal/logger"
)

// Postgre is an abstraction layer that manages product entities inside PostgreSQL DB.
type Postgre struct {
	db *sqlx.DB
	logger.Logger
}

// NewPostgreRepo creates a new PostgreSQL repository for Product entity.
// It's also embed logger for convenience.
func NewPostgreRepo(db *sqlx.DB, l logger.Logger) *Postgre {
	return &Postgre{
		db,
		l,
	}
}

// Create a new product in PostgreSQL DB.
func (r *Postgre) Create(ctx context.Context, newProduct entity.NewProduct) (entity.Product, error) {
	const query = `
	INSERT INTO products 
		(product_id, title, description, price, stock, date_created, date_updated) 
	VALUES
		(:product_id, :title, :description, :price, :stock, :date_created, :date_updated)`

	product := entity.Product{
		ID:          uuid.NewString(),
		Title:       newProduct.Title,
		Description: newProduct.Description,
		Price:       newProduct.Price,
		Stock:       newProduct.Stock,
		DateCreated: time.Now().UTC(),
		DateUpdated: time.Now().UTC(),
	}

	if _, err := r.db.NamedExecContext(ctx, query, product); err != nil {
		return entity.Product{}, errors.Wrap(err, "inserting a product")
	}

	return product, nil
}

// Query gets products from PostgreSQL DB.
// This query uses two provided values to implement pagination: last seen id and limit.
// Results of a query sorted by creation date of selected users.
func (r *Postgre) Query(ctx context.Context, lastSeenID, limit string) ([]entity.Product, error) {
	const query = `
	SELECT 
		* 
	FROM 
		products 
	WHERE
		product_id <= :last_seen_id
	ORDER BY 
		product_id DESC	
	FETCH FIRST :limit ROWS ONLY`

	var products []entity.Product

	data := struct {
		LastSeenID string `db:"last_seen_id"`
		Limit      string `db:"limit"`
	}{
		LastSeenID: lastSeenID,
		Limit:      limit,
	}

	err := database.QuerySlice(ctx, r.db, query, data, &products)
	if err != nil {
		return []entity.Product{}, errors.Wrap(err, "selecting products")
	}

	return products, nil
}

// QueryByID gets product from PostgreSQL DB by given id.
func (r *Postgre) QueryByID(ctx context.Context, id string) (entity.Product, error) {
	const query = `
	SELECT 
		* 
	FROM 
		products 
	WHERE 
		product_id = :product_id`

	data := struct {
		ID string `db:"product_id"`
	}{
		ID: id,
	}

	var product entity.Product

	err := database.QueryStruct(ctx, r.db, query, data, &product)
	if err != nil {
		return entity.Product{}, errors.Wrapf(err, "getting a product with id %s", id)
	}

	return product, err
}

// Update a product inside PostgreSQL.
func (r *Postgre) Update(ctx context.Context, id string, updateProduct entity.UpdateProduct) error {
	product, err := r.QueryByID(ctx, id)
	if err != nil {
		return errors.Wrapf(err, "error updating a product with id %s", id)
	}

	const query = `
	UPDATE 
		products 
	SET	
		"title" = :title, 
		"description" = :description, 
		"price" = :price, 
		"stock" = :stock,
		"date_updated" = :date_updated
	WHERE 
		"product_id" = :product_id`

	if updateProduct.Title != nil {
		product.Title = *updateProduct.Title
	}
	if updateProduct.Description != nil {
		product.Description = *updateProduct.Description
	}
	if updateProduct.Price != nil {
		product.Price = *updateProduct.Price
	}
	if updateProduct.Stock != nil {
		product.Stock = *updateProduct.Stock
	}
	product.DateUpdated = time.Now().UTC()

	_, err = r.db.NamedExecContext(ctx, query, product)
	if err != nil {
		return errors.Wrapf(err, "updating a product with id %s", id)
	}

	return nil
}

// Delete a product from PostgreSQL DB by given product id.
func (r *Postgre) Delete(ctx context.Context, id string) error {
	const query = `
	DELETE FROM 
		products 
	WHERE 
		"product_id" = :product_id`

	data := struct {
		ID string `db:"product_id"`
	}{
		ID: id,
	}

	if _, err := r.db.NamedExecContext(ctx, query, data); err != nil {
		return errors.Wrapf(err, "deleting a product with id %s", id)
	}

	return nil
}
