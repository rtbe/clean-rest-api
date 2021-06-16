package orderitem

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

// Postgre is an abstraction layer that manages product item entities inside PostgreSQL DB.
// It's also embed logger for convenience.
type Postgre struct {
	db *sqlx.DB
	logger.Logger
}

// NewPostgreRepo creates a new PostgreSQL repository for order item entity.
func NewPostgreRepo(db *sqlx.DB, l logger.Logger) *Postgre {
	return &Postgre{
		db,
		l,
	}
}

// Create a new order item for particular order in PostgreSQL DB.
func (r *Postgre) Create(ctx context.Context, newOrderItem entity.NewOrderItem) (entity.OrderItem, error) {
	const query = `
	INSERT INTO order_items
		(order_item_id, order_id, product_id, quantity, date_created, date_updated) 
	VALUES
		(:order_item_id, :order_id, :product_id , :quantity, :date_created, :date_updated)`

	orderItem := entity.OrderItem{
		ID:          uuid.NewString(),
		OrderID:     newOrderItem.OrderID,
		ProductID:   newOrderItem.ProductID,
		Quantity:    newOrderItem.Quantity,
		DateCreated: time.Now().UTC(),
		DateUpdated: time.Now().UTC(),
	}

	if _, err := r.db.NamedExecContext(ctx, query, orderItem); err != nil {
		return entity.OrderItem{}, errors.Wrap(err, "inserting an order item")
	}

	return orderItem, nil
}

// Query gets order items from PostgreSQL DB.
// This query uses two provided values to implement pagination: last seen id and limit.
// Results of a query sorted by creation date of selected users.
func (r *Postgre) Query(ctx context.Context, lastSeenID, limit string) ([]entity.OrderItem, error) {
	const query = `
	SELECT 
		* 
	FROM 
		order_items 
	WHERE
		order_item_id <= :last_seen_id
	ORDER BY 
		order_item_id DESC
	FETCH FIRST :limit ROWS ONLY`

	var orderItems []entity.OrderItem

	data := struct {
		LastSeenID string `db:"last_seen_id"`
		Limit      string `db:"limit"`
	}{
		LastSeenID: lastSeenID,
		Limit:      limit,
	}

	err := database.QuerySlice(ctx, r.db, query, data, &orderItems)
	if err != nil {
		return nil, errors.Wrap(err, "selecting order items")
	}

	return orderItems, nil
}

// QueryByID gets an order item from PostgreSQL DB by given id.
func (r *Postgre) QueryByID(ctx context.Context, id string) (entity.OrderItem, error) {
	const query = `
	SELECT 
		* 
	FROM 
		order_items 
	WHERE 
		order_item_id = :order_item_id`

	data := struct {
		ID string `db:"order_item_id"`
	}{
		ID: id,
	}

	var orderItem entity.OrderItem

	err := database.QueryStruct(ctx, r.db, query, data, &orderItem)
	if err != nil {
		return entity.OrderItem{}, errors.Wrapf(err, "getting an order item with id %s", id)
	}

	return orderItem, nil
}

// QueryByOrderID gets all order items for particular order from PostgreSQL DB.
func (r *Postgre) QueryByOrderID(ctx context.Context, orderID string) ([]entity.OrderItem, error) {
	const query = `
	SELECT 
		* 
	FROM 
		order_items 
	WHERE 
		order_id = :order_id 
	ORDER BY 
		date_created`

	data := struct {
		ID string `db:"order_id"`
	}{
		ID: orderID,
	}

	var orderItems []entity.OrderItem

	err := database.QuerySlice(ctx, r.db, query, data, &orderItems)
	if err != nil {
		return nil, errors.Wrap(err, "selecting order items by order id")
	}

	return orderItems, nil
}

// Update an order item in PostgreSQL DB.
func (r *Postgre) Update(ctx context.Context, id string, updateOrderItem entity.UpdateOrderItem) error {
	orderItem, err := r.QueryByID(ctx, id)
	if err != nil {
		return errors.Wrapf(err, "updating an order item with id %s", id)
	}

	const query = `
	UPDATE 
		order_items 
	SET	
		"quantity" = :quantity,		
		"date_updated" = :date_updated 
	WHERE
		order_item_id = :order_item_id`

	if updateOrderItem.Quantity != nil {
		orderItem.Quantity = *updateOrderItem.Quantity
	}
	orderItem.DateUpdated = time.Now().UTC()

	_, err = r.db.NamedExecContext(ctx, query, orderItem)
	if err != nil {
		return errors.Wrapf(err, "updating an order item with id %s", id)
	}

	return nil
}

// Delete an order item from PostgreSQL DB by given order item id.
func (r *Postgre) Delete(ctx context.Context, id string) error {
	const query = `
	DELETE FROM 
		order_items 
	WHERE 
		order_item_id = :order_item_id`

	data := struct {
		ID string `db:"order_item_id"`
	}{
		ID: id,
	}

	if _, err := r.db.NamedExecContext(ctx, query, data); err != nil {
		return errors.Wrapf(err, "deleting an order item with id %s", id)
	}

	return nil
}

// DeleteByOrderID deletes order items from PostgreSQL DB by given order id.
func (r *Postgre) DeleteByOrderID(ctx context.Context, orderID string) error {
	const query = `
	DELETE FROM 
		order_items 
	WHERE 
		order_id = :order_id`

	data := struct {
		OrderID string `db:"order_id"`
	}{
		OrderID: orderID,
	}

	if _, err := r.db.NamedExecContext(ctx, query, data); err != nil {
		return errors.Wrapf(err, "deleting order items for order with id %s", orderID)
	}

	return nil
}
