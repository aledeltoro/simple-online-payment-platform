package service

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/aledeltoro/simple-online-payment-platform/internal/database/postgres"
	"github.com/aledeltoro/simple-online-payment-platform/internal/models"
	"github.com/aledeltoro/simple-online-payment-platform/internal/paymentprocessor/stripe"
	"github.com/stretchr/testify/require"
)

func TestProcessPayment(t *testing.T) {
	c := require.New(t)

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

	mockDatabase := postgres.MockPostgres{}
	mockPaymentProcessor := stripe.MockStripe{}

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

func TestProcessPaymentInvalidInput(t *testing.T) {
	c := require.New(t)

	onlinePaymentService := onlinePaymentService{}

	_, err := onlinePaymentService.ProcessPayment(context.Background(), -12, "", "", "")
	c.ErrorIs(err, models.ErrInvalidAmount)
}

func TestProcessPaymentPerformTransactionFailure(t *testing.T) {
	c := require.New(t)

	input := &models.TransactionInput{
		Amount:        2000,
		Currency:      "usd",
		PaymentMethod: "card_pm_visa",
		Description:   "Transaction for payment maount of 2000",
	}

	mockPaymentProcessor := stripe.MockStripe{}

	customErr := fmt.Errorf("performing transaction: card_declined")

	mockPaymentProcessor.On("PerformTransaction", input).Return(nil, customErr)

	onlinePaymentService := onlinePaymentService{
		paymentProcessor: &mockPaymentProcessor,
	}

	_, err := onlinePaymentService.ProcessPayment(context.Background(), input.Amount, input.Currency, input.PaymentMethod, input.Description)
	c.ErrorIs(err, customErr)
}

func TestProcessPaymentInsertTransactionFailure(t *testing.T) {
	c := require.New(t)

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

	mockDatabase := postgres.MockPostgres{}
	mockPaymentProcessor := stripe.MockStripe{}

	customErr := fmt.Errorf("inserting transaction: operation failed")

	mockPaymentProcessor.On("PerformTransaction", input).Return(expectedTransaction, nil)
	mockDatabase.On("InsertTransaction", context.Background(), expectedTransaction).Return(customErr)

	onlinePaymentService := onlinePaymentService{
		database:         &mockDatabase,
		paymentProcessor: &mockPaymentProcessor,
	}

	_, err := onlinePaymentService.ProcessPayment(context.Background(), input.Amount, input.Currency, input.PaymentMethod, input.Description)
	c.ErrorIs(err, customErr)
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

func TestQueryPaymentMissingTransactionID(t *testing.T) {
	c := require.New(t)

	onlinePaymentService := onlinePaymentService{}

	_, err := onlinePaymentService.QueryPayment(context.Background(), "")
	c.ErrorIs(err, ErrMissingTransactionID)
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

	mockDatabase.On("GetTransaction", context.Background(), "TXN_123").Return(expectedTransaction, nil)
	mockPaymentProcessor.On("RefundTransaction", expectedTransaction.AdditionalFields).Return(refundedTransaction, nil)
	mockDatabase.On("UpdateTransaction", context.Background(), "TXN_123", refundedTransaction).Return(expectedTransaction, nil)

	onlinePaymentService := onlinePaymentService{
		database:         &mockDatabase,
		paymentProcessor: &mockPaymentProcessor,
	}

	transaction, err := onlinePaymentService.RefundPayment(context.Background(), "TXN_123")
	c.NoError(err)
	c.Equal(expectedTransaction, transaction)
}

func TestRefundPaymentMissingTransactionID(t *testing.T) {
	c := require.New(t)

	onlinePaymentService := onlinePaymentService{}

	_, err := onlinePaymentService.RefundPayment(context.Background(), "")
	c.ErrorIs(err, ErrMissingTransactionID)
}

func TestRefunPaymentTransactionRefundFailure(t *testing.T) {
	c := require.New(t)

	mockDatabase := postgres.MockPostgres{}
	mockPaymentProcessor := stripe.MockStripe{}

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

	customErr := errors.New("refunding transaction: charge already refunded")

	mockDatabase.On("GetTransaction", context.Background(), "TXN_123").Return(expectedTransaction, nil)
	mockPaymentProcessor.On("RefundTransaction", expectedTransaction.AdditionalFields).Return(nil, customErr)

	onlinePaymentService := onlinePaymentService{
		database:         &mockDatabase,
		paymentProcessor: &mockPaymentProcessor,
	}

	_, err := onlinePaymentService.RefundPayment(context.Background(), "TXN_123")
	c.ErrorIs(err, customErr)
}
