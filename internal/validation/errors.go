package validation

import "encoding/json"

// FieldError is a representation of validation error on particular struct field.
type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

// FieldErrors is a collection of field errors.
type FieldErrors []FieldError

// Error implements error interface.
func (fe FieldErrors) Error() string {
	data, err := json.Marshal(fe)
	if err != nil {
		return err.Error()
	}

	return string(data)
}
