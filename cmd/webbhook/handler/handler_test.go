package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

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
