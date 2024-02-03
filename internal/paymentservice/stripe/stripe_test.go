package stripe

import (
	"fmt"
	"testing"

	"github.com/aledeltoro/simple-online-payment-platform/internal/models"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	c := require.New(t)

	t.Setenv("STRIPE_SECRET_KEY", "secret_key")

	stripeService, err := New()
	c.NoError(err)
	c.IsType(stripeService, stripeService)
}

func TestNewErrMissingAPIKey(t *testing.T) {
	c := require.New(t)

	stripeService, err := New()
	c.Nil(stripeService)
	c.ErrorIs(err, ErrMissingAPIKey)
}

func TestPerformTransaction(t *testing.T) {
	c := require.New(t)

	expectedTransaction := &models.Transaction{
		TransactionID: fmt.Sprintf("TXN_%s", ulid.Make().String()),
		Status:        models.TransactionStatusSucceeded,
		Description:   "Testing stripe service",
		Provider:      models.PaymentProviderStripe,
		Amount:        2000,
		Currency:      "usd",
		Type:          models.TransactionTypeCharge,
		AdditionalFields: map[string]interface{}{
			"charge_id":         "charge_id",
			"payment_intent_id": "payment_intent_id",
		},
	}

	expectedInput := &models.TransactionInput{
		Amount:        2000,
		Currency:      "usd",
		PaymentMethod: "pm_card_visa",
		Description:   "Testing stripe service",
	}

	mockService := MockStripe{}

	mockService.On("PerformTransaction", expectedInput).Return(expectedTransaction, nil)

	transaction, err := mockService.PerformTransaction(expectedInput)
	c.NoError(err)
	c.Equal(expectedTransaction, transaction)
}

func TestQueryTransaction(t *testing.T) {
	c := require.New(t)

	mockService := MockStripe{}

	mockService.On("QueryTransaction", "id").Return(nil, nil)

	transaction, err := mockService.QueryTransaction("id")
	c.Nil(transaction)
	c.NoError(err)
}

func TestRefundTransaction(t *testing.T) {
	c := require.New(t)

	expectedTransaction := &models.Transaction{
		Status: models.TransactionStatusSucceeded,
		Type:   models.TransactionTypeRefund,
		AdditionalFields: map[string]interface{}{
			"charge_id":         "charge_id",
			"payment_intent_id": "payment_intent_id",
			"refund_id":         "refund_id",
		},
	}

	mockService := MockStripe{}

	mockService.On("RefundTransaction", "charge_id").Return(expectedTransaction, nil)

	updatedTransaction, err := mockService.RefundTransaction("charge_id")
	c.NoError(err)
	c.Equal(expectedTransaction, updatedTransaction)
}
