package routes

import (
	"context"
	"net/http"

	"github.com/aledeltoro/simple-online-payment-platform/internal/service"
)

type Handler interface {
	HandleProcessPayment(ctx context.Context) http.HandlerFunc
	HandleGetPayment(ctx context.Context) http.HandlerFunc
	HandleRefundPayment(ctx context.Context) http.HandlerFunc
}

type handler struct {
	service service.OnlinePaymentService
}

func NewHandler(service service.OnlinePaymentService) Handler {
	return handler{
		service: service,
	}
}

func (h handler) HandleProcessPayment(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func (h handler) HandleGetPayment(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func (h handler) HandleRefundPayment(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}
