package entity

import (
	"time"
)

// Product is an particular product.
//
// swagger:model
type Product struct {
	// UUID of a product
	//
	ID string `db:"product_id" json:"product_id"`

	// Title of a product
	//
	// required: true
	Title string `db:"title" json:"title"`

	// Description of a product
	//
	// required: true
	Description string `db:"description" json:"description"`

	// Price of a product
	//
	// gte:0.00
	// required: true
	Price float32 `db:"price" json:"price"`

	// Stock of a product
	//
	// gte:0
	// required: true
	Stock int `db:"stock" json:"stock"`

	// Date of a product creation
	//
	DateCreated time.Time `db:"date_created" json:"date_created"`

	// Date of a product last modification
	//
	DateUpdated time.Time `db:"date_updated" json:"date_updated"`
}

// NewProduct is an information needed to create a new product.
//
// swagger:model
type NewProduct struct {
	// Title of a product
	//
	// required: true
	Title string `json:"title,omitempty" validate:"required"`

	// Description of a product
	//
	// required: true
	Description string `json:"description,omitempty" validate:"required"`

	// Price of a product
	//
	// gte:0.00
	// required: true
	Price float32 `json:"price,omitempty" validate:"gte=0"`

	// Stock of a product
	//
	// gte:0
	// required: true
	Stock int `json:"stock,omitempty" validate:"gte=0"`
}

// UpdateProduct is an information needed to update an existing product.
//
// swagger:model
type UpdateProduct struct {
	// Title of a product
	//
	// required: true
	Title *string `json:"title"`

	// Description of a product
	//
	// required: true
	Description *string `json:"description"`

	// Price of a product
	//
	// gte:0
	// required: true
	Price *float32 `json:"price"`

	// Stock of a product
	//
	// gte:0
	// required: true
	Stock *int `json:"stock"`
}
