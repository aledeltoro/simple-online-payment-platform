package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/aledeltoro/simple-online-payment-platform/internal/models"
	"github.com/aledeltoro/simple-online-payment-platform/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestHandleProcessPayment(t *testing.T) {
	c := require.New(t)

	mockService := service.MockOnlinePaymentService{}

	expectedTransaction := &models.Transaction{
		TransactionID: "TXN_123",
		Status:        models.TransactionStatusSucceeded,
		Description:   "Transaction for payment amount of 2000",
		Provider:      models.PaymentProviderStripe,
		Amount:        2000,
		Currency:      "usd",
		Type:          models.TransactionTypeCharge,
		AdditionalFields: map[string]interface{}{
			"charge_id":         "ch_123",
			"payment_intent_id": "pi_123",
		},
	}

	mockService.On("ProcessPayment", context.Background(), int64(2000), "usd", "card_pm_visa", "Transaction for payment amount of 2000").Return(expectedTransaction, nil)

	form := url.Values{}
	form.Add("amount", "2000")
	form.Add("currency", "usd")
	form.Add("payment_method", "card_pm_visa")
	form.Add("description", "Transaction for payment amount of 2000")

	handler := NewHandler(&mockService)

	router := chi.NewRouter()
	router.Post("/payments", http.HandlerFunc(handler.HandleProcessPayment(context.Background())))

	req := httptest.NewRequest(http.MethodPost, "/payments", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	response := recorder.Result()

	defer response.Body.Close()

	c.Equal(http.StatusOK, response.StatusCode)

	var transaction *models.Transaction

	err := json.NewDecoder(response.Body).Decode(&transaction)
	c.NoError(err)
	c.Equal(expectedTransaction, transaction)
}

func TestHandleGetPayment(t *testing.T) {
	c := require.New(t)

	mockService := service.MockOnlinePaymentService{}

	expectedTransaction := &models.Transaction{
		TransactionID: "TXN_123",
		Status:        models.TransactionStatusSucceeded,
		Description:   "Transaction for payment amount of 2000",
		Provider:      models.PaymentProviderStripe,
		Amount:        2000,
		Currency:      "usd",
		Type:          models.TransactionTypeCharge,
		AdditionalFields: map[string]interface{}{
			"charge_id":         "ch_123",
			"payment_intent_id": "pi_123",
		},
	}

	mockService.On("QueryPayment", context.Background(), "TXN_123").Return(expectedTransaction, nil)

	handler := NewHandler(&mockService)

	router := chi.NewRouter()
	router.Get("/payments/{id}", http.HandlerFunc(handler.HandleQueryPayment(context.Background())))

	req := httptest.NewRequest(http.MethodGet, "/payments/TXN_123", nil)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	response := recorder.Result()

	defer response.Body.Close()

	c.Equal(http.StatusOK, response.StatusCode)

	var transaction *models.Transaction

	err := json.NewDecoder(response.Body).Decode(&transaction)
	c.NoError(err)
	c.Equal(expectedTransaction, transaction)
}

func TestHandleRefundPayment(t *testing.T) {
	c := require.New(t)

	mockService := service.MockOnlinePaymentService{}

	expectedTransaction := &models.Transaction{
		TransactionID: "TXN_123",
		Status:        models.TransactionStatusSucceeded,
		Description:   "Transaction for payment amount of 2000",
		Provider:      models.PaymentProviderStripe,
		Amount:        2000,
		Currency:      "usd",
		Type:          models.TransactionTypeRefund,
		AdditionalFields: map[string]interface{}{
			"refund_id":         "rf_123",
			"charge_id":         "ch_123",
			"payment_intent_id": "pi_123",
		},
	}

	mockService.On("RefundPayment", context.Background(), "TXN_123").Return(expectedTransaction, nil)

	handler := NewHandler(&mockService)

	router := chi.NewRouter()
	router.Post("/payments/{id}/refunds", http.HandlerFunc(handler.HandleRefundPayment(context.Background())))

	req := httptest.NewRequest(http.MethodPost, "/payments/TXN_123/refunds", nil)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	response := recorder.Result()

	defer response.Body.Close()

	c.Equal(http.StatusOK, response.StatusCode)

	var transaction *models.Transaction

	err := json.NewDecoder(response.Body).Decode(&transaction)
	c.NoError(err)
	c.Equal(expectedTransaction, transaction)
}
