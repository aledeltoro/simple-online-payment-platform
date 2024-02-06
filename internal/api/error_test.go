package api

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAPIErr(t *testing.T) {
	c := require.New(t)

	customErr := errors.New("custom error")

	testCases := []struct {
		runFunc    func(err error, resource string) error
		errCode    ErrorCode
		statusCode int
		ErrMessage string
		resource   string
		err        error
	}{
		{
			runFunc: func(err error, resource string) error {
				return NewInternalServerError(err)
			},
			errCode:    ErrCodeInternalServerError,
			statusCode: http.StatusInternalServerError,
			ErrMessage: "(500) Internal server error",
			resource:   "",
			err:        customErr,
		},
		{
			runFunc: func(err error, resource string) error {
				return NewInvalidRequestError(err)
			},
			errCode:    ErrCodeInvalidRequestError,
			statusCode: http.StatusBadRequest,
			ErrMessage: fmt.Sprintf("(400) Invalid request: %s", customErr.Error()),
			resource:   "",
			err:        customErr,
		},
		{
			runFunc: func(err error, resource string) error {
				return NewResourceNotFoundError(err, resource)
			},
			errCode:    ErrCodeResourceNotFound,
			statusCode: http.StatusNotFound,
			ErrMessage: "(404) Resource 'transaction' not found",
			resource:   "transaction",
			err:        customErr,
		},
	}

	for _, testCase := range testCases {
		var apiErr APIError

		err := testCase.runFunc(testCase.err, testCase.resource)

		c.ErrorAs(err, &apiErr)
		c.ErrorIs(err, apiErr)
		c.Equal(testCase.errCode, apiErr.Code())
		c.Equal(testCase.statusCode, apiErr.HTTPStatusCode())
		c.EqualError(apiErr, testCase.ErrMessage)
	}
}
