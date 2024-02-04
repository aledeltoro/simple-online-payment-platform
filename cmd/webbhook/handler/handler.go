package handler

import (
	"context"
	"net/http"

	"github.com/aledeltoro/simple-online-payment-platform/internal/events"
)

type Handler interface {
	HandlePaymentEvents(ctx context.Context) http.HandlerFunc
}

type handler struct {
	events events.Events
}

func NewHandler(events events.Events) Handler {
	return handler{
		events: events,
	}
}

func (h handler) HandlePaymentEvents(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}
