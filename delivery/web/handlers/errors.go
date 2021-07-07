package handlers

import (
	"encoding/json"
)

// RequestError is our known error that will be presented to an app user.
type RequestError struct {
	ErrorText string `json:"error"`
	Fields    string `json:"fields,omitempty"`
	Status    int    `json:"status_code"`
}

// Error implements error interface.
func (re RequestError) Error() string {
	data, err := json.Marshal(re)
	if err != nil {
		return err.Error()
	}

	return string(data)
}
