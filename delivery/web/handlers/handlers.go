package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	mid "github.com/rtbe/clean-rest-api/delivery/web/middlewares"
	"github.com/rtbe/clean-rest-api/internal/logger"
)

// Handler is a handler function so we can return an error from request handling functions.
type Handler struct {
	H handlerFunc
	L logger.Logger
}

// handlerFunc is a custom handler function to enable error handling.
type handlerFunc func(w http.ResponseWriter, r *http.Request) error

// ServeHTTP lets Handler implements http.Handler interface.
// That lets us return and handle errors from our custom handlerFunc.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Centralize error handling from handlerFunc here.
	if err := h.H(w, r); err != nil {
		info, err := mid.GetRequestInfo(r.Context())
		if err != nil {
			return
		}

		h.L.Log("error", fmt.Sprintf(
			"%s: ERROR: %v ",
			info.ID, err.Error(),
		))
		// Customize handling of known and unknown errors.
		switch err := errors.Cause(err).(type) {

		// Known error handling.
		case RequestError:
			info.StatusCode = err.Status
			if err := respond(ctx, w, err, err.Status); err != nil {
				h.L.Log("error", fmt.Sprintf("responding with error: %v", err.Error()))
			}

			// Unknown error handling.
		default:
			e := RequestError{
				ErrorText: http.StatusText(http.StatusInternalServerError),
			}
			info.StatusCode = http.StatusInternalServerError

			if err := respond(ctx, w, e, http.StatusInternalServerError); err != nil {
				h.L.Log("error", fmt.Sprintf("responding with error: %v", err.Error()))
			}
		}
	}
}

// respond is a helper function for handling json responses.
func respond(ctx context.Context, w http.ResponseWriter, data interface{}, statusCode int) error {

	// Set status code of request into it's context for logging it later.
	requestInfo, err := mid.GetRequestInfo(ctx)
	if err != nil {
		return err
	}
	requestInfo.StatusCode = statusCode

	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}

	json, err := json.Marshal(data)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)

	if _, err := w.Write(json); err != nil {
		return err
	}

	return nil
}

// parseURLParamID gets value from URL by it's key and checks if it's validity to UUID format.
func parseURLParamID(r *http.Request, key string) (string, error) {
	v := chi.URLParam(r, key)
	if v == "" {
		e := fmt.Sprintf("there is no value with key: %v inside URL", key)
		return "", errors.New(e)
	}

	if _, err := uuid.Parse(v); err != nil {
		return "", RequestError{
			ErrorText: "id is not in UUID format",
			Status:    http.StatusBadRequest,
		}
	}

	return "", nil
}
