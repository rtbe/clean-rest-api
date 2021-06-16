package order

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

// Postgre is an abstraction layer that manages order entities inside PostgreSQL DB.
// It's also embed logger for convenience.
type Postgre struct {
	db *sqlx.DB
	logger.Logger
}

// NewPostgreRepo creates a new PostgreSQL repository for order entity.
func NewPostgreRepo(db *sqlx.DB, l logger.Logger) *Postgre {
	return &Postgre{
		db,
		l,
	}
}

// Create creates an new order in PostgreSQL DB.
func (r *Postgre) Create(ctx context.Context, no entity.NewOrder) (entity.Order, error) {
	const query = `
	INSERT INTO orders 
		(order_id, user_id, status, date_created, date_updated) 
	VALUES
		(:order_id, :user_id, :status, :date_created, :date_updated)`

	order := entity.Order{
		ID:          uuid.NewString(),
		UserID:      no.UserID,
		Status:      no.Status,
		DateCreated: time.Now().UTC(),
		DateUpdated: time.Now().UTC(),
	}

	if _, err := r.db.NamedExecContext(ctx, query, order); err != nil {
		return entity.Order{}, errors.Wrap(err, "inserting an order")
	}

	return order, nil
}

// Query gets orders from PostgreSQL.
// This query uses two provided values to implement pagination: last seen id and limit.
// Results of a query sorted by creation date of selected users.
func (r *Postgre) Query(ctx context.Context, lastSeenID, limit string) ([]entity.Order, error) {
	const query = `
	SELECT 
		* 
	FROM 
		orders 
	WHERE
		order_id <= :last_seen_id
	ORDER BY 
		order_id DESC
	FETCH FIRST :limit ROWS ONLY`

	var orders []entity.Order

	data := struct {
		LastSeenID string `db:"last_seen_id"`
		Limit      string `db:"limit"`
	}{
		LastSeenID: lastSeenID,
		Limit:      limit,
	}

	err := database.QuerySlice(ctx, r.db, query, data, &orders)
	if err != nil {
		return []entity.Order{}, errors.Wrap(err, "selecting orders")
	}

	return orders, nil
}

// Query gets an order from PostgreSQL DB by given id.
func (r *Postgre) QueryByID(ctx context.Context, id string) (entity.Order, error) {
	const query = `
	SELECT 
		* 
	FROM 
		orders 
	WHERE 
		order_id = :order_id`

	data := struct {
		ID string `db:"order_id"`
	}{
		ID: id,
	}

	var order entity.Order

	err := database.QueryStruct(ctx, r.db, query, data, &order)
	if err != nil {
		return entity.Order{}, errors.Wrapf(err, "getting an order with id %s", id)
	}

	return order, nil
}

// QueryByUserID gets orders from PostgreSQL DB by given user id.
func (r *Postgre) QueryByUserID(ctx context.Context, userID string) ([]entity.Order, error) {
	const query = `
	SELECT 
		* 
	FROM 
		orders 
	WHERE 
		user_id = :user_id`

	data := struct {
		UserID string `db:"user_id"`
	}{
		UserID: userID,
	}

	var orders []entity.Order

	err := database.QuerySlice(ctx, r.db, query, data, &orders)
	if err != nil {
		return []entity.Order{}, errors.Wrapf(err, "selecting products for user with id %s", userID)
	}

	return orders, nil
}

// Update updates a specific order inside PostgreSQL.
func (r *Postgre) Update(ctx context.Context, id string, updateOrder entity.UpdateOrder) error {
	order, err := r.QueryByID(ctx, id)
	if err != nil {
		return errors.Wrapf(err, "error updating a order with id %s", id)
	}

	const query = `
	UPDATE 
		orders
	SET	
		"status" = :status,
		"date_updated" = :date_updated
	WHERE 
		order_id = :order_id`

	// You should not update not updated order fields.
	if updateOrder.Status != nil {
		order.Status = *updateOrder.Status
	}
	order.DateUpdated = time.Now().UTC()

	_, err = r.db.NamedExecContext(ctx, query, order)
	if err != nil {
		return errors.Wrapf(err, "updating an order with id %s", id)
	}

	return nil
}

// Delete an order from PostgreSQL DB by given order id.
func (r *Postgre) Delete(ctx context.Context, id string) error {
	const query = `
	DELETE FROM 
		orders 
	WHERE 
		order_id = :order_id`

	data := struct {
		ID string `db:"order_id"`
	}{
		ID: id,
	}

	if _, err := r.db.NamedExecContext(ctx, query, data); err != nil {
		return errors.Wrapf(err, "deleting an order with id %s", id)
	}

	return nil
}

// DeleteByUserID deletes orders from PostgreSQL DB by given user id.
func (r *Postgre) DeleteByUserID(ctx context.Context, userID string) error {
	const query = `
	DELETE FROM 
		orders 
	WHERE 
		user_id = :user_id`

	data := struct {
		UserID string `db:"user_id"`
	}{
		UserID: userID,
	}

	if _, err := r.db.NamedExecContext(ctx, query, data); err != nil {
		return errors.Wrapf(err, "deleting an order with user id %s", userID)
	}

	return nil
}
