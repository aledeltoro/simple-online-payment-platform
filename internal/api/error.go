package api

import (
	"fmt"
	"net/http"
	"os"
)

// ErrorCode error code for each type of error in the service
type ErrorCode string

var (
	// ErrCodeInternalServerError error code when service suffered an unexpected error
	ErrCodeInternalServerError ErrorCode = "internal_server_error"
	// ErrCodeInvalidRequestError error code when service received an invalid request
	ErrCodeInvalidRequestError ErrorCode = "invalid_request"
	// ErrCodeResourceNotFound error code when resource was not found
	ErrCodeResourceNotFound ErrorCode = "resource_not_found"
)

// APIError interface to handle API errors in the service
type APIError interface {
	error
	Unwrap() error
	HTTPStatusCode() int
	Code() ErrorCode
}

// APIErr error type to standardize errors in the service
type APIErr struct {
	ErrCode    ErrorCode `json:"code"`
	StatusCode int       `json:"status_code"`
	Message    string    `json:"message"`
	err        error
}

// Unwrap returns the error
func (e APIErr) Unwrap() error {
	return e.err
}

// HTTPStatusCode return HTTP status code assigned to the error
func (e APIErr) HTTPStatusCode() int {
	return e.StatusCode
}

// Code returns the error code assigned to the error
func (e APIErr) Code() ErrorCode {
	return e.ErrCode
}

// Error returns a formatted error message
func (e APIErr) Error() string {
	return fmt.Sprintf("(%d) %s", e.StatusCode, e.Message)
}

// NewInternalServerError API error to handle unexpected behavior in the service
func NewInternalServerError(err error) APIErr {
	apiErr := APIErr{
		ErrCode:    ErrCodeInternalServerError,
		StatusCode: http.StatusInternalServerError,
		Message:    "Internal server error",
		err:        err,
	}

	if os.Getenv("DEBUG_MODE") == "true" {
		apiErr.Message = fmt.Sprintf("Internal server error: %s", err.Error())
	}

	return apiErr
}

// NewInvalidRequestError API error to handle invalid request in the service
func NewInvalidRequestError(err error) APIErr {
	return APIErr{
		ErrCode:    ErrCodeInvalidRequestError,
		StatusCode: http.StatusBadRequest,
		Message:    fmt.Sprintf("Invalid request: %s", err.Error()),
		err:        err,
	}
}

// NewResourceNotFoundError API error when resource was not found in the service
func NewResourceNotFoundError(err error, resource string) APIErr {
	return APIErr{
		ErrCode:    ErrCodeResourceNotFound,
		StatusCode: http.StatusNotFound,
		Message:    fmt.Sprintf("Resource '%s' not found", resource),
		err:        err,
	}
}
