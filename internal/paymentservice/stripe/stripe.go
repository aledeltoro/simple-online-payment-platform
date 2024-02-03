package stripe

import (
	"errors"
	"fmt"
	"os"

	"github.com/aledeltoro/simple-online-payment-platform/internal/models"
	"github.com/aledeltoro/simple-online-payment-platform/internal/paymentservice"
	"github.com/oklog/ulid/v2"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/client"
)

// ErrMissingAPIKey error when missing api key
var ErrMissingAPIKey = errors.New("missing api key")

type stripeService struct {
	client *client.API
}

// New initializes implementation of Stripe service
func New() (paymentservice.PaymentService, error) {
	stripeKey := os.Getenv("STRIPE_SECRET_KEY")

	if stripeKey == "" {
		return nil, ErrMissingAPIKey
	}

	return stripeService{
		client: client.New(stripeKey, nil),
	}, nil
}

// PerformTransaction method to perform a transaction
func (s stripeService) PerformTransaction(input *models.TransactionInput) (*models.Transaction, error) {
	params := &stripe.PaymentIntentParams{
		Amount:        stripe.Int64(input.Amount),
		Currency:      stripe.String(input.Currency),
		Description:   stripe.String(input.Description),
		PaymentMethod: stripe.String(input.PaymentMethod),
		Confirm:       stripe.Bool(true),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			AllowRedirects: stripe.String("never"),
			Enabled:        stripe.Bool(true),
		},
	}

	result, err := s.client.PaymentIntents.New(params)
	if err != nil {
		// Handle all errors, except card error
		return nil, fmt.Errorf("performing transaction: %w", err)
	}

	// Card error should be handled here, if found, we should return a normal transaction object
	transaction := &models.Transaction{
		TransactionID: fmt.Sprintf("TXN_%s", ulid.Make().String()),
		Status:        models.TransactionStatusSucceeded,
		Description:   result.Description,
		Provider:      models.PaymentProviderStripe,
		Amount:        int(result.Amount),
		Currency:      string(result.Currency),
		Type:          models.TransactionTypeCharge,
		AdditionalFields: map[string]interface{}{
			"charge_id":         result.LatestCharge.ID,
			"payment_intent_id": result.ID,
		},
	}

	return transaction, nil
}

func (s stripeService) QueryTransaction(id string) (*models.Transaction, error) {
	return nil, nil
}

func (s stripeService) RefundTransaction(chargeID string) (*models.Transaction, error) {
	params := &stripe.RefundParams{
		Charge: stripe.String(chargeID),
	}

	result, err := s.client.Refunds.New(params)
	if err != nil {
		return nil, fmt.Errorf("performing refund: %w", err)
	}

	transaction := &models.Transaction{
		Status: models.TransactionStatusSucceeded,
		Type:   models.TransactionTypeRefund,
		AdditionalFields: map[string]interface{}{
			"charge_id":         result.Charge.ID,
			"payment_intent_id": result.PaymentIntent.ID,
			"refund_id":         result.ID,
		},
	}

	return transaction, nil
}
