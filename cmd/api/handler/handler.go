package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/aledeltoro/simple-online-payment-platform/internal/api"
	"github.com/aledeltoro/simple-online-payment-platform/internal/models"
	"github.com/aledeltoro/simple-online-payment-platform/internal/service"
	"github.com/go-chi/chi/v5"
)

var (
	errMissingTransactionID = api.NewInvalidRequestError(errors.New("missing transaction id"))
	errInvalidInput         = api.NewInvalidRequestError(errors.New("invalid input"))
)

// Handler interface to handle incoming requests to online payment plataform API
type Handler interface {
	HandleProcessPayment(ctx context.Context) http.HandlerFunc
	HandleQueryPayment(ctx context.Context) http.HandlerFunc
	HandleRefundPayment(ctx context.Context) http.HandlerFunc
}

type handler struct {
	service service.OnlinePaymentService
}

// NewHandler constructor to handle incoming requests to API
func NewHandler(service service.OnlinePaymentService) Handler {
	return handler{
		service: service,
	}
}

// HandleProcessPayments handles requests to create a payment
func (h handler) HandleProcessPayment(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			api.WriteErrorResponse(w, errInvalidInput)
			return
		}

		amount, err := strconv.ParseInt(r.FormValue("amount"), 10, 64)
		if err != nil {
			api.WriteErrorResponse(w, api.NewInvalidRequestError(models.ErrInvalidAmount))
			return
		}

		transaction, err := h.service.ProcessPayment(ctx, amount, r.FormValue("currency"), r.FormValue("payment_method"), r.FormValue("description"))
		if err != nil {
			api.WriteErrorResponse(w, err)
			return
		}

		api.WriteJSONResponse(w, http.StatusOK, transaction)
	}
}

// HandleQueryPayment handles requests to query a specific payment
func (h handler) HandleQueryPayment(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		transactionID := chi.URLParam(r, "id")
		if transactionID == "" {
			api.WriteErrorResponse(w, errMissingTransactionID)
			return
		}

		transaction, err := h.service.QueryPayment(ctx, transactionID)
		if err != nil {
			api.WriteErrorResponse(w, err)
			return
		}

		api.WriteJSONResponse(w, http.StatusOK, transaction)
	}
}

// HandleRefundPayment handles requests to refund a specific payment
func (h handler) HandleRefundPayment(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		transactionID := chi.URLParam(r, "id")
		if transactionID == "" {
			api.WriteErrorResponse(w, errMissingTransactionID)
			return
		}

		transaction, err := h.service.RefundPayment(ctx, transactionID)
		if err != nil {
			api.WriteErrorResponse(w, err)
			return
		}

		api.WriteJSONResponse(w, http.StatusOK, transaction)
	}
}
