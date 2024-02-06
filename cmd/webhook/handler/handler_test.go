package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aledeltoro/simple-online-payment-platform/internal/api"
	"github.com/aledeltoro/simple-online-payment-platform/internal/database"
	"github.com/aledeltoro/simple-online-payment-platform/internal/database/postgres"
	"github.com/aledeltoro/simple-online-payment-platform/internal/events"
	"github.com/aledeltoro/simple-online-payment-platform/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestHandlePaymentEvents(t *testing.T) {
	c := require.New(t)

	mockEvents := events.MockStripe{}

	mockEvents.On("VerifyEvent").Return(nil)
	mockEvents.On("ProcessEvent").Return(nil)

	copyNewEventHandlerFunc := newEventHandlerFunc
	newEventHandlerFunc = func(provider models.PaymentProvider, database database.Database, request *http.Request) (events.Events, error) {
		return &mockEvents, nil
	}

	t.Cleanup(func() {
		newEventHandlerFunc = copyNewEventHandlerFunc
	})

	mockDatabase := postgres.MockPostgres{}

	handler := NewHandler(&mockDatabase)

	router := chi.NewRouter()
	router.Post("/payments/{provider}/events", http.HandlerFunc(handler.HandlePaymentEvents(context.Background())))

	body := bytes.NewReader([]byte(`{"hello": "world"}`))

	req := httptest.NewRequest(http.MethodPost, "/payments/mock/events", body)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	response := recorder.Result()

	defer response.Body.Close()

	c.Equal(http.StatusOK, response.StatusCode)
}

func TestHandlerPaymentEventsUnsupportedProvider(t *testing.T) {
	c := require.New(t)

	mockDatabase := postgres.MockPostgres{}

	handler := NewHandler(&mockDatabase)

	router := chi.NewRouter()
	router.Post("/payments/{provider}/events", http.HandlerFunc(handler.HandlePaymentEvents(context.Background())))

	req := httptest.NewRequest(http.MethodPost, "/payments/invalid/events", nil)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	response := recorder.Result()

	defer response.Body.Close()

	c.Equal(http.StatusBadRequest, response.StatusCode)

	var apiErr api.APIErr

	err := json.NewDecoder(response.Body).Decode(&apiErr)
	c.NoError(err)
	c.Equal(api.ErrCodeInvalidRequestError, apiErr.Code())
	c.Contains(apiErr.Error(), events.ErrUnsupportedProvider.Error())
}

func TestHandlerPaymentEventsVerificationFailed(t *testing.T) {
	c := require.New(t)

	mockEvents := events.MockStripe{}

	mockEvents.On("VerifyEvent").Return(api.NewInvalidRequestError(events.ErrEventVerificationFailed))

	copyNewEventHandlerFunc := newEventHandlerFunc
	newEventHandlerFunc = func(provider models.PaymentProvider, database database.Database, request *http.Request) (events.Events, error) {
		return &mockEvents, nil
	}

	t.Cleanup(func() {
		newEventHandlerFunc = copyNewEventHandlerFunc
	})

	mockDatabase := postgres.MockPostgres{}

	handler := NewHandler(&mockDatabase)

	router := chi.NewRouter()
	router.Post("/payments/{provider}/events", http.HandlerFunc(handler.HandlePaymentEvents(context.Background())))

	body := bytes.NewReader([]byte(`{"hello": "world"}`))

	req := httptest.NewRequest(http.MethodPost, "/payments/mock/events", body)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	response := recorder.Result()

	defer response.Body.Close()

	c.Equal(http.StatusBadRequest, response.StatusCode)

	var apiErr api.APIErr

	err := json.NewDecoder(response.Body).Decode(&apiErr)
	c.NoError(err)
	c.Equal(api.ErrCodeInvalidRequestError, apiErr.Code())
	c.Contains(apiErr.Error(), events.ErrEventVerificationFailed.Error())
}

func TestHandlerPaymentEventsProcessEventFailure(t *testing.T) {
	c := require.New(t)

	mockEvents := events.MockStripe{}

	unknownErr := errors.New("unknown error")

	mockEvents.On("VerifyEvent").Return(nil)
	mockEvents.On("ProcessEvent").Return(api.NewInternalServerError(unknownErr))

	copyNewEventHandlerFunc := newEventHandlerFunc
	newEventHandlerFunc = func(provider models.PaymentProvider, database database.Database, request *http.Request) (events.Events, error) {
		return &mockEvents, nil
	}

	t.Cleanup(func() {
		newEventHandlerFunc = copyNewEventHandlerFunc
	})

	mockDatabase := postgres.MockPostgres{}

	handler := NewHandler(&mockDatabase)

	router := chi.NewRouter()
	router.Post("/payments/{provider}/events", http.HandlerFunc(handler.HandlePaymentEvents(context.Background())))

	body := bytes.NewReader([]byte(`{"hello": "world"}`))

	req := httptest.NewRequest(http.MethodPost, "/payments/mock/events", body)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	response := recorder.Result()

	defer response.Body.Close()

	c.Equal(http.StatusInternalServerError, response.StatusCode)

	var apiErr api.APIErr

	err := json.NewDecoder(response.Body).Decode(&apiErr)
	c.NoError(err)
	c.Equal(api.ErrCodeInternalServerError, apiErr.Code())
	c.Contains(apiErr.Error(), "Internal server error")
}
