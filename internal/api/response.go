package api

import (
	"encoding/json"
	"errors"
	"net/http"
)

// WriteJSONResponse writes a new JSON response given a serializable value
func WriteJSONResponse(w http.ResponseWriter, statusCode int, v interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	_ = json.NewEncoder(w).Encode(v)
}

// WriteErrorResponse writes an error JSON response
func WriteErrorResponse(w http.ResponseWriter, err error) {
	var apiErr APIErr

	if errors.As(err, &apiErr) {
		WriteJSONResponse(w, apiErr.HTTPStatusCode(), apiErr)

		return
	}

	WriteJSONResponse(w, http.StatusInternalServerError, NewInternalServerError(apiErr))
}
