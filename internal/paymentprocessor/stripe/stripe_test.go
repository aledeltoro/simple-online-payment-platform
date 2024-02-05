package stripe

import (
	"testing"

	"github.com/aledeltoro/simple-online-payment-platform/internal/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/client"
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
		Status:      models.TransactionStatusPending,
		Description: "Testing stripe service",
		Provider:    models.PaymentProviderStripe,
		Amount:      2000,
		Currency:    "usd",
		Type:        models.TransactionTypeCharge,
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

	stripeBackendMock := new(mockStripeBackend)
	stripeTestBackends := &stripe.Backends{
		API:     stripeBackendMock,
		Connect: stripeBackendMock,
		Uploads: stripeBackendMock,
	}

	stripeBackendMock.On("Call", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		mockPaymentIntentResult := args.Get(4).(*stripe.PaymentIntent)

		*mockPaymentIntentResult = stripe.PaymentIntent{
			ID:          "payment_intent_id",
			Description: expectedTransaction.Description,
			Amount:      int64(expectedTransaction.Amount),
			Currency:    stripe.Currency(expectedTransaction.Currency),
			LatestCharge: &stripe.Charge{
				ID: "charge_id",
			},
		}
	}).Return(nil)

	mockStripeClient := client.New("sk_test", stripeTestBackends)

	service := stripeService{
		client: mockStripeClient,
	}

	transaction, err := service.PerformTransaction(expectedInput)
	c.NoError(err)

	transaction.TransactionID = ""

	c.Equal(expectedTransaction, transaction)
}

func TestPerformTransactionCardError(t *testing.T) {
	c := require.New(t)

	expectedTransaction := &models.Transaction{
		Status:        models.TransactionStatusFailure,
		FailureReason: string(stripe.ErrorCodeCardDeclined),
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

	stripeBackendMock := new(mockStripeBackend)
	stripeTestBackends := &stripe.Backends{
		API:     stripeBackendMock,
		Connect: stripeBackendMock,
		Uploads: stripeBackendMock,
	}

	stripeErr := &stripe.Error{
		Type: stripe.ErrorTypeCard,
		Code: stripe.ErrorCodeCardDeclined,
	}

	stripeBackendMock.On("Call", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		mockPaymentIntentResult := args.Get(4).(*stripe.PaymentIntent)

		*mockPaymentIntentResult = stripe.PaymentIntent{
			ID:          "payment_intent_id",
			Description: expectedTransaction.Description,
			Amount:      int64(expectedTransaction.Amount),
			Currency:    stripe.Currency(expectedTransaction.Currency),
			LatestCharge: &stripe.Charge{
				ID: "charge_id",
			},
		}
	}).Return(stripeErr)

	mockStripeClient := client.New("sk_test", stripeTestBackends)

	service := stripeService{
		client: mockStripeClient,
	}

	transaction, err := service.PerformTransaction(expectedInput)
	c.NoError(err)
	c.Equal(expectedTransaction.Status, transaction.Status)
	c.Equal(expectedTransaction.FailureReason, transaction.FailureReason)
}

func TestRefundTransaction(t *testing.T) {
	c := require.New(t)

	expectedTransaction := &models.Transaction{
		Status: models.TransactionStatusPending,
		Type:   models.TransactionTypeRefund,
		AdditionalFields: map[string]interface{}{
			"charge_id":         "charge_id",
			"payment_intent_id": "payment_intent_id",
			"refund_id":         "refund_id",
		},
	}

	stripeBackendMock := new(mockStripeBackend)
	stripeTestBackends := &stripe.Backends{
		API:     stripeBackendMock,
		Connect: stripeBackendMock,
		Uploads: stripeBackendMock,
	}

	stripeBackendMock.On("Call", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		mockRefund := args.Get(4).(*stripe.Refund)

		*mockRefund = stripe.Refund{
			ID: "refund_id",
			Charge: &stripe.Charge{
				ID: "charge_id",
			},
			PaymentIntent: &stripe.PaymentIntent{
				ID: "payment_intent_id",
			},
		}
	}).Return(nil)

	mockStripeClient := client.New("sk_test", stripeTestBackends)

	service := stripeService{
		client: mockStripeClient,
	}

	updatedTransaction, err := service.RefundTransaction(expectedTransaction.AdditionalFields)
	c.NoError(err)
	c.Equal(expectedTransaction, updatedTransaction)
}

func TestRefundTransactionAlreadyRefunded(t *testing.T) {
	c := require.New(t)

	expectedTransaction := &models.Transaction{
		Status: models.TransactionStatusPending,
		Type:   models.TransactionTypeRefund,
		AdditionalFields: map[string]interface{}{
			"charge_id":         "charge_id",
			"payment_intent_id": "payment_intent_id",
			"refund_id":         "refund_id",
		},
	}

	stripeBackendMock := new(mockStripeBackend)
	stripeTestBackends := &stripe.Backends{
		API:     stripeBackendMock,
		Connect: stripeBackendMock,
		Uploads: stripeBackendMock,
	}

	stripeErr := &stripe.Error{
		Type: stripe.ErrorTypeInvalidRequest,
		Code: stripe.ErrorCodeChargeAlreadyRefunded,
	}

	stripeBackendMock.On("Call", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(stripeErr)

	mockStripeClient := client.New("sk_test", stripeTestBackends)

	service := stripeService{
		client: mockStripeClient,
	}

	updatedTransaction, err := service.RefundTransaction(expectedTransaction.AdditionalFields)
	c.Nil(updatedTransaction)
	c.ErrorIs(err, ErrChargeAlreadyRefunded)
}

func TestRefundTransactionMissingChargeID(t *testing.T) {
	c := require.New(t)

	service := stripeService{}

	_, err := service.RefundTransaction(map[string]interface{}{})
	c.ErrorIs(err, ErrMissingChargeID)
}
