// Package tests contains set of global variables and helper functions for testing.
package tests

const (
	PgUser     = "test"
	PgPassword = "test123"
	PgDB       = "test"
	PgPort     = "5432"

	MongoUser     = "test"
	MongoPassword = "test123"
	MongoDB       = "test"
	MongoPort     = "20017"
)

// Success and failure markers.
const (
	Success = "\u2713"
	Failed  = "\u2717"
)

// StrPtr is a helper functions for tests
// that helps with updating string pointer fields on entities.
func StrPtr(s string) *string {
	return &s
}

// StrPtr is a helper functions for tests
// that helps with updating integer pointer fields on entities.
func IntPtr(i int) *int {
	return &i
}

// StrPtr is a helper functions for tests
// that helps with updating float32 pointer fields on entities.
func Float32Ptr(f float32) *float32 {
	return &f
}
