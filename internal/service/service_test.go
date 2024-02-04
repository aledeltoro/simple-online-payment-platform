package service

import (
	"context"
	"testing"

	"github.com/aledeltoro/simple-online-payment-platform/internal/database/postgres"
	"github.com/aledeltoro/simple-online-payment-platform/internal/models"
	"github.com/aledeltoro/simple-online-payment-platform/internal/paymentservice/stripe"
	"github.com/stretchr/testify/require"
)

func TestProcessPayment(t *testing.T) {
	c := require.New(t)

	mockDatabase := postgres.MockPostgres{}
	mockPaymentProcessor := stripe.MockStripe{}

	input := &models.TransactionInput{
		Amount:        2000,
		Currency:      "usd",
		PaymentMethod: "card_pm_visa",
		Description:   "Transaction for payment maount of 2000",
	}

	expectedTransaction := &models.Transaction{
		TransactionID: "TXN_123",
		Status:        models.TransactionStatusSucceeded,
		Description:   input.Description,
		FailureReason: "",
		Provider:      models.PaymentProviderStripe,
		Amount:        int(input.Amount),
		Currency:      input.Currency,
		Type:          models.TransactionTypeCharge,
		AdditionalFields: map[string]interface{}{
			"charge_id":         "ch_123",
			"payment_intent_id": "pi_123",
		},
	}

	mockPaymentProcessor.On("PerformTransaction", input).Return(expectedTransaction, nil)
	mockDatabase.On("InsertTransaction", context.Background(), expectedTransaction).Return(nil)

	onlinePaymentService := onlinePaymentService{
		database:         &mockDatabase,
		paymentProcessor: &mockPaymentProcessor,
	}

	transaction, err := onlinePaymentService.ProcessPayment(context.Background(), input.Amount, input.Currency, input.PaymentMethod, input.Description)
	c.NoError(err)
	c.Equal(expectedTransaction, transaction)
}

func TestQueryPayment(t *testing.T) {
	c := require.New(t)

	mockDatabase := postgres.MockPostgres{}

	input := &models.TransactionInput{
		Amount:        2000,
		Currency:      "usd",
		PaymentMethod: "card_pm_visa",
		Description:   "Transaction for payment maount of 2000",
	}

	expectedTransaction := &models.Transaction{
		TransactionID: "TXN_123",
		Status:        models.TransactionStatusSucceeded,
		Description:   input.Description,
		FailureReason: "",
		Provider:      models.PaymentProviderStripe,
		Amount:        int(input.Amount),
		Currency:      input.Currency,
		Type:          models.TransactionTypeCharge,
		AdditionalFields: map[string]interface{}{
			"charge_id":         "ch_123",
			"payment_intent_id": "pi_123",
		},
	}

	mockDatabase.On("GetTransaction", context.Background(), "TXN_123").Return(expectedTransaction, nil)

	onlinePaymentService := onlinePaymentService{
		database:         &mockDatabase,
		paymentProcessor: nil,
	}

	transaction, err := onlinePaymentService.QueryPayment(context.Background(), "TXN_123")
	c.NoError(err)
	c.Equal(expectedTransaction, transaction)
}

func TestRefundPayment(t *testing.T) {
	c := require.New(t)

	mockDatabase := postgres.MockPostgres{}
	mockPaymentProcessor := stripe.MockStripe{}

	refundedTransaction := &models.Transaction{
		TransactionID: "TXN_123",
		Status:        models.TransactionStatusSucceeded,
		Type:          models.TransactionTypeRefund,
		AdditionalFields: map[string]interface{}{
			"charge_id":         "ch_123",
			"payment_intent_id": "pi_123",
			"refund_id":         "rf_123",
		},
	}

	expectedTransaction := &models.Transaction{
		TransactionID: "TXN_123",
		Status:        models.TransactionStatusSucceeded,
		Description:   "Transaction for payment amount of 2000",
		FailureReason: "",
		Provider:      models.PaymentProviderStripe,
		Amount:        2000,
		Currency:      "usd",
		Type:          models.TransactionTypeRefund,
		AdditionalFields: map[string]interface{}{
			"charge_id":         "ch_123",
			"payment_intent_id": "pi_123",
			"refund_id":         "rf_123",
		},
	}

	mockPaymentProcessor.On("RefundTransaction", "TXN_123").Return(refundedTransaction, nil)
	mockDatabase.On("UpdateTransaction", context.Background(), "TXN_123", refundedTransaction).Return(nil)
	mockDatabase.On("GetTransaction", context.Background(), "TXN_123").Return(expectedTransaction, nil)

	onlinePaymentService := onlinePaymentService{
		database:         &mockDatabase,
		paymentProcessor: &mockPaymentProcessor,
	}

	transaction, err := onlinePaymentService.RefundPayment(context.Background(), "TXN_123")
	c.NoError(err)
	c.Equal(expectedTransaction, transaction)
}
