package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/aledeltoro/simple-online-payment-platform/internal/service"
	"github.com/go-chi/chi/v5"
)

type Handler interface {
	HandleProcessPayment(ctx context.Context) http.HandlerFunc
	HandleQueryPayment(ctx context.Context) http.HandlerFunc
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
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			errMessage := fmt.Errorf("parse form failed: %w", err)
			http.Error(w, errMessage.Error(), http.StatusInternalServerError)
			return
		}

		amount, err := strconv.ParseInt(r.FormValue("amount"), 10, 64)
		if err != nil {
			errMessage := fmt.Errorf("parsing integer failed: %w", err)
			http.Error(w, errMessage.Error(), http.StatusInternalServerError)
			return
		}

		transaction, err := h.service.ProcessPayment(ctx, amount, r.FormValue("currency"), r.FormValue("payment_method"), r.FormValue("description"))
		if err != nil {
			errMessage := fmt.Errorf("processing payment failed: %w", err)
			http.Error(w, errMessage.Error(), http.StatusInternalServerError)
			return
		}

		_ = json.NewEncoder(w).Encode(transaction)
	}
}

func (h handler) HandleQueryPayment(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		transactionID := chi.URLParam(r, "id")
		if transactionID == "" {
			http.Error(w, "Missing transaction id", http.StatusBadRequest)
			return
		}

		transaction, err := h.service.QueryPayment(ctx, transactionID)
		if err != nil {
			errMessage := fmt.Errorf("querying payment failed: %w", err)
			http.Error(w, errMessage.Error(), http.StatusInternalServerError)
			return
		}

		_ = json.NewEncoder(w).Encode(transaction)
	}
}

func (h handler) HandleRefundPayment(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		transactionID := chi.URLParam(r, "id")
		if transactionID == "" {
			http.Error(w, "Missing transaction id", http.StatusBadRequest)
			return
		}

		transaction, err := h.service.RefundPayment(ctx, transactionID)
		if err != nil {
			errMessage := fmt.Errorf("refunding payment failed: %w", err)
			http.Error(w, errMessage.Error(), http.StatusInternalServerError)
			return
		}

		_ = json.NewEncoder(w).Encode(transaction)
	}
}
