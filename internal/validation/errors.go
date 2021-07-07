package validation

import (
	"fmt"
	"strings"
)

// FieldError is a representation of validation error on particular struct field.
type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

// FieldErrors is a collection of field errors.
type FieldErrors []FieldError

// Error implements error interface.
func (fe FieldErrors) Error() string {
	var buf []string
	var str strings.Builder

	// Default json marshaller will escape all of the fields inside FieldErrors with "\",
	// So here we build json encoded message by yourself.
	for i, err := range fe {
		if i == 0 {
			str.WriteString("[")
		}

		str.WriteString(fmt.Sprintf("{%s: %s}", err.Field, err.Error))

		if i == len(fe)-1 {
			str.WriteString("]")
		}

		buf = append(buf, str.String())
		str.Reset()
	}

	return strings.Join(buf, ",")
}
