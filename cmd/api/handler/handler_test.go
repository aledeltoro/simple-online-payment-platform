package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/aledeltoro/simple-online-payment-platform/internal/api"
	"github.com/aledeltoro/simple-online-payment-platform/internal/database"
	"github.com/aledeltoro/simple-online-payment-platform/internal/models"
	"github.com/aledeltoro/simple-online-payment-platform/internal/paymentprocessor/stripe"
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

func TestHandleProcessPaymentParseFormFailure(t *testing.T) {
	c := require.New(t)

	handler := handler{}

	router := chi.NewRouter()
	router.Post("/payments", http.HandlerFunc(handler.HandleProcessPayment(context.Background())))

	req := httptest.NewRequest(http.MethodPost, "/payments", strings.NewReader("foo%3z1%26bar%3D2"))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	response := recorder.Result()

	defer response.Body.Close()

	c.Equal(http.StatusBadRequest, response.StatusCode)

	var apiErr api.APIErr

	err := json.NewDecoder(response.Body).Decode(&apiErr)
	c.NoError(err)
	c.Equal(api.ErrCodeInvalidRequestError, apiErr.Code())
	c.Contains(apiErr.Error(), errInvalidInput.Error())
}

func TestHandleProcessPaymentInvalidAmount(t *testing.T) {
	c := require.New(t)

	form := url.Values{}
	form.Add("amount", "invalid")
	form.Add("currency", "usd")
	form.Add("payment_method", "card_pm_visa")
	form.Add("description", "Transaction for payment amount of 2000")

	handler := handler{}

	router := chi.NewRouter()
	router.Post("/payments", http.HandlerFunc(handler.HandleProcessPayment(context.Background())))

	req := httptest.NewRequest(http.MethodPost, "/payments", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	response := recorder.Result()

	defer response.Body.Close()

	c.Equal(http.StatusBadRequest, response.StatusCode)

	var apiErr api.APIErr

	err := json.NewDecoder(response.Body).Decode(&apiErr)
	c.NoError(err)
	c.Equal(api.ErrCodeInvalidRequestError, apiErr.Code())
	c.Contains(apiErr.Error(), models.ErrInvalidAmount.Error())
}

func TestHandleProcessPaymentFailure(t *testing.T) {
	c := require.New(t)

	mockService := service.MockOnlinePaymentService{}

	unknownErr := errors.New("unknown error")

	mockService.On("ProcessPayment", context.Background(), int64(2000), "usd", "card_pm_visa", "Transaction for payment amount of 2000").Return(nil, unknownErr)

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

	c.Equal(http.StatusInternalServerError, response.StatusCode)

	var apiErr api.APIErr

	err := json.NewDecoder(response.Body).Decode(&apiErr)
	c.NoError(err)
	c.Equal(api.ErrCodeInternalServerError, apiErr.Code())
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

func TestHandleQueryPaymentMissingTransactionID(t *testing.T) {
	c := require.New(t)

	handler := handler{}

	router := chi.NewRouter()
	router.Get("/payments/", http.HandlerFunc(handler.HandleQueryPayment(context.Background())))

	req := httptest.NewRequest(http.MethodGet, "/payments/", nil)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	response := recorder.Result()

	defer response.Body.Close()

	c.Equal(http.StatusBadRequest, response.StatusCode)

	var apiErr api.APIErr

	err := json.NewDecoder(response.Body).Decode(&apiErr)
	c.NoError(err)
	c.Equal(api.ErrCodeInvalidRequestError, apiErr.Code())
	c.Contains(apiErr.Error(), errMissingTransactionID.Error())
}

func TestHandleGetPaymentTransactionNotFound(t *testing.T) {
	c := require.New(t)

	mockService := service.MockOnlinePaymentService{}

	errTransactionNotFound := api.NewResourceNotFoundError(database.ErrTransactionNotFound, "transaction")

	mockService.On("QueryPayment", context.Background(), "TXN_123").Return(nil, errTransactionNotFound)

	handler := NewHandler(&mockService)

	router := chi.NewRouter()
	router.Get("/payments/{id}", http.HandlerFunc(handler.HandleQueryPayment(context.Background())))

	req := httptest.NewRequest(http.MethodGet, "/payments/TXN_123", nil)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	response := recorder.Result()

	defer response.Body.Close()

	c.Equal(http.StatusNotFound, response.StatusCode)

	var apiErr api.APIErr

	err := json.NewDecoder(response.Body).Decode(&apiErr)
	c.NoError(err)
	c.Equal(api.ErrCodeResourceNotFound, apiErr.Code())
	c.Contains(apiErr.Error(), errTransactionNotFound.Error())
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

func TestHandleRefundPaymentMissingTransactionID(t *testing.T) {
	c := require.New(t)

	handler := handler{}

	router := chi.NewRouter()
	router.Get("/payments/refunds", http.HandlerFunc(handler.HandleRefundPayment(context.Background())))

	req := httptest.NewRequest(http.MethodGet, "/payments/refunds", nil)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	response := recorder.Result()

	defer response.Body.Close()

	c.Equal(http.StatusBadRequest, response.StatusCode)

	var apiErr api.APIErr

	err := json.NewDecoder(response.Body).Decode(&apiErr)
	c.NoError(err)
	c.Equal(api.ErrCodeInvalidRequestError, apiErr.Code())
	c.Contains(apiErr.Error(), errMissingTransactionID.Error())
}

func TestHandleRefundPaymentChargeAlreadyRefunded(t *testing.T) {
	c := require.New(t)

	mockService := service.MockOnlinePaymentService{}

	errChargeAlreadyRefunded := api.NewInvalidRequestError(stripe.ErrChargeAlreadyRefunded)

	mockService.On("RefundPayment", context.Background(), "TXN_123").Return(nil, errChargeAlreadyRefunded)

	handler := NewHandler(&mockService)

	router := chi.NewRouter()
	router.Get("/payments/{id}/refunds", http.HandlerFunc(handler.HandleRefundPayment(context.Background())))

	req := httptest.NewRequest(http.MethodGet, "/payments/TXN_123/refunds", nil)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	response := recorder.Result()

	defer response.Body.Close()

	c.Equal(http.StatusBadRequest, response.StatusCode)

	var apiErr api.APIErr

	err := json.NewDecoder(response.Body).Decode(&apiErr)
	c.NoError(err)
	c.Equal(api.ErrCodeInvalidRequestError, apiErr.Code())
	c.Contains(apiErr.Error(), errChargeAlreadyRefunded.Error())
}
