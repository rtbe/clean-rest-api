package entity

import (
	"time"

	"github.com/lib/pq"
)

// User is a particular user.
//
// swagger:model
type User struct {
	// The UUID of a user
	//
	// required: true
	ID string `db:"user_id" json:"user_id"`

	// Username of a user
	//
	// required: true
	UserName string `db:"user_name" json:"user_name"`

	// First name of a user
	//
	FirstName string `db:"first_name" json:"first_name"`

	// Last name of a user
	//
	LastName string `db:"last_name" json:"last_name"`

	// Password of a user
	//
	// required: true
	Password []byte `db:"password" json:"password"`

	// Email of a user
	//
	// example: user@google.com
	// required: true
	Email string `db:"email" json:"email"`

	// Set of user roles
	//
	Roles pq.StringArray `db:"roles" json:"roles"`

	// Date of a user creation
	//
	DateCreated time.Time `db:"date_created" json:"date_created"`

	// Date of a user last modification
	//
	DateUpdated time.Time `db:"date_updated" json:"date_updated"`
}

// NewUser is an information needed to create a new user.
//
// swagger:model
type NewUser struct {
	// Username of a user
	//
	// required: true
	UserName string `json:"user_name" validate:"required"`

	// First name of a user
	//
	FirstName string `json:"first_name"`

	// Last name of a user
	//
	LastName string `json:"last_name"`

	// Email of a user
	//
	// example: user@google.com
	// required: true
	Email string `json:"email" validate:"required,email"`

	// Password of a user
	//
	// required: true
	Password string `json:"password" validate:"required"`

	// Confirmation password of a user
	//
	// required: true
	PasswordConfirm string `json:"password_confirm" validate:"required,eqfield=Password"`

	// Set of user roles
	//
	// required: true
	Roles []string `json:"roles" validate:"required"`
}

// UpdateUser is an information needed to update a existing user.
//
// swagger:model
type UpdateUser struct {
	// Username of a user
	//
	UserName *string `json:"user_name"`

	// First name of a user
	//
	FirstName *string `json:"first_name"`

	// Last name of a user
	//
	LastName *string `json:"last_name"`

	// Email of a user
	//
	// example: user@google.com
	Email *string `json:"email"`

	// Password of a user
	//
	Password *string `json:"password"`

	// Confirmation password of a user
	//
	PasswordConfirm *string `json:"password_confirm"`

	// Set of user roles
	//
	Roles []string `json:"roles"`
}
