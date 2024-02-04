package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aledeltoro/simple-online-payment-platform/internal/database"
	"github.com/aledeltoro/simple-online-payment-platform/internal/events"
	"github.com/aledeltoro/simple-online-payment-platform/internal/models"
	"github.com/go-chi/chi/v5"
)

var newEventHandlerFunc = events.NewEvent

type Handler interface {
	HandlePaymentEvents(ctx context.Context) http.HandlerFunc
}

type handler struct {
	database database.Database
}

func NewHandler(database database.Database) Handler {
	return handler{
		database: database,
	}
}

func (h handler) HandlePaymentEvents(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		provider := chi.URLParam(r, "provider")

		eventHandler, err := newEventHandlerFunc(models.PaymentProvider(provider), h.database, r)
		if err != nil {
			errMessage := fmt.Errorf("initializing event handler: %w", err)
			http.Error(w, errMessage.Error(), http.StatusOK)
			return
		}

		err = eventHandler.VerifyEvent()
		if err != nil {
			errMessage := fmt.Errorf("verifying event: %w", err)
			http.Error(w, errMessage.Error(), http.StatusOK)
			return
		}

		err = eventHandler.ProcessEvent(ctx)
		if err != nil {
			errMessage := fmt.Errorf("processing event: %w", err)
			http.Error(w, errMessage.Error(), http.StatusOK)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
