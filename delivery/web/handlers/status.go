package handlers

import (
	"net/http"
	"os"
)

type StatusGroup struct {
}

// swagger:route GET /status status status
//
// Gets application status.
//
// Produces:
// - application/json
//
// Responses:
//   200: statusResponse
func (sg StatusGroup) Status(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	host, err := os.Hostname()
	if err != nil {
		host = "unavailable"
	}

	info := struct {
		Status string `json:"status"`
		Host   string `json:"host"`
	}{
		Status: "up",
		Host:   host,
	}

	return respond(ctx, w, info, http.StatusOK)
}
