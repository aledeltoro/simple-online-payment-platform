package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/aledeltoro/simple-online-payment-platform/internal/api"
	"github.com/aledeltoro/simple-online-payment-platform/internal/database"
	"github.com/aledeltoro/simple-online-payment-platform/internal/events"
	"github.com/aledeltoro/simple-online-payment-platform/internal/models"
	"github.com/go-chi/chi/v5"
)

var newEventHandlerFunc = events.NewEvent

// Handler interface to handle incoming events from payment processor providers
type Handler interface {
	HandlePaymentEvents(ctx context.Context) http.HandlerFunc
}

type handler struct {
	database database.Database
}

// NewHandler constructor to handle incoming requests to API
func NewHandler(database database.Database) Handler {
	return handler{
		database: database,
	}
}

// HandlePaymentsEvents validates and processes events from payment processor providers
func (h handler) HandlePaymentEvents(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		provider := chi.URLParam(r, "provider")

		eventHandler, err := newEventHandlerFunc(models.PaymentProvider(provider), h.database, r)
		if err != nil {
			log.Println(err)
			api.WriteErrorResponse(w, err)
			return
		}

		err = eventHandler.VerifyEvent()
		if err != nil {
			log.Println(err)
			api.WriteErrorResponse(w, err)
			return
		}

		err = eventHandler.ProcessEvent(ctx)
		if err != nil {
			log.Println(err)
			api.WriteErrorResponse(w, err)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
