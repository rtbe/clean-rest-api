// Package database contains support for different databases.
package database

import (
	"errors"
)

// Set of errors for database related CRUD operations.
var (
	ErrNotFound = errors.New("not found")
)
