package entity

import (
	"time"
)

// Order is an particular order.
//
// swagger:model
type Order struct {
	// UUID of order
	//
	ID string `db:"order_id" json:"order_id"`

	// User UUID of an order
	//
	// required: true
	UserID string `db:"user_id" json:"user_id"`

	// Status of an order
	//
	Status string `db:"status" json:"status"`

	// Date of an order creation
	//
	DateCreated time.Time `db:"date_created" json:"date_created"`

	// Date of an order last modification
	//
	DateUpdated time.Time `db:"date_updated" json:"date_updated"`
}

// NewOrder is an information needed to create a new order.
//
// swagger:model
type NewOrder struct {
	// UUID of a user related to particular order
	//
	// required: true
	UserID string `json:"user_id" validate:"required"`

	// Status of an order
	//
	Status string `json:"status" validate:"required"`
}

// UpdateOrder is an information needed to update an existing order.
//
// swagger:model
type UpdateOrder struct {
	// Status of an order
	//
	Status *string `json:"status"`
}
