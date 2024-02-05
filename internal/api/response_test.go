package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWriteJSONResponse(t *testing.T) {
	c := require.New(t)

	writer := httptest.NewRecorder()

	WriteJSONResponse(writer, http.StatusOK, map[string]string{"hello": "world"})

	response := writer.Result()

	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	c.NoError(err)

	c.Equal(http.StatusOK, response.StatusCode)
	c.Equal("application/json", response.Header.Get("Content-Type"))
	c.JSONEq(`{"hello":"world"}`, string(data))
}

func TestWriteErrorResponse(t *testing.T) {
	c := require.New(t)

	writer := httptest.NewRecorder()

	apiErr := NewInvalidRequestError(errors.New("invalid input"))

	WriteErrorResponse(writer, apiErr)

	response := writer.Result()

	defer response.Body.Close()

	expectedJSONErr, err := json.Marshal(apiErr)
	c.NoError(err)

	responseData, err := io.ReadAll(response.Body)
	c.NoError(err)

	c.Equal(http.StatusBadRequest, response.StatusCode)
	c.Equal("application/json", response.Header.Get("Content-Type"))
	c.JSONEq(string(expectedJSONErr), string(responseData))
}

func TestWriteErrorResponseUnknownError(t *testing.T) {
	c := require.New(t)

	unknownErr := errors.New("unknown error")

	writer := httptest.NewRecorder()

	WriteErrorResponse(writer, unknownErr)

	response := writer.Result()

	defer response.Body.Close()

	expectedJSONErr, err := json.Marshal(NewInternalServerError(unknownErr))
	c.NoError(err)

	responseData, err := io.ReadAll(response.Body)
	c.NoError(err)

	c.Equal(http.StatusInternalServerError, response.StatusCode)
	c.Equal("application/json", response.Header.Get("Content-Type"))
	c.JSONEq(string(expectedJSONErr), string(responseData))
}
