package entity

import (
	"time"
)

// OrderItem type is particular order.
//
// swagger:model
type OrderItem struct {
	// UUID of an order item
	//
	ID string `db:"order_item_id" json:"order_item_id"`

	// UUID of an order that an order item belongs to
	//
	// required: true
	OrderID string `db:"order_id" json:"order_id"`

	// UUID of a product that an order item belongs to
	//
	// required: true
	ProductID string `db:"product_id" json:"product_id"`

	// Quantity of an order item
	//
	// min: 0
	// required : true
	Quantity int `db:"quantity" json:"quantity"`

	// Date of an order item creation
	//
	DateCreated time.Time `db:"date_created" json:"date_created"`

	// Date of an order item last modification
	//
	DateUpdated time.Time `db:"date_updated" json:"date_updated"`
}

// NewOrderItem is an information needed to create a new order item.
//
// swagger:model
type NewOrderItem struct {
	// UUID of an order that an order item belongs to
	//
	// required: true
	OrderID string `json:"order_id" validate:"required"`

	// UUID of a product that an order item belongs to
	//
	// required: true
	ProductID string `json:"product_id" validate:"required"`

	// Quantity of an order item
	//
	// min: 0
	// required : true
	Quantity int `json:"quantity" validate:"gte=1"`
}

// UpdateOrderItem is an information needed to update an existing order item.
//
// swagger:model
type UpdateOrderItem struct {
	// Quantity of an order item
	//
	// min: 0
	// required : true
	Quantity *int `json:"quantity"`
}
