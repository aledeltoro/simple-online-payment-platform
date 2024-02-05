package stripe

import (
	"errors"
	"fmt"
	"os"

	"github.com/aledeltoro/simple-online-payment-platform/internal/api"
	"github.com/aledeltoro/simple-online-payment-platform/internal/models"
	"github.com/aledeltoro/simple-online-payment-platform/internal/paymentprocessor"
	"github.com/oklog/ulid/v2"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/client"
)

var (
	// ErrMissingAPIKey error when missing api key
	ErrMissingAPIKey = errors.New("missing api key")
	// ErrChargeAlreadyRefunded error when charge has been refunded already
	ErrChargeAlreadyRefunded = errors.New("charge already refunded")
	// ErrMissingChargeID error when missing charge ID
	ErrMissingChargeID = errors.New("missing charge ID")
)

type stripeService struct {
	client *client.API
}

// New initializes implementation of Stripe service
func New() (paymentprocessor.PaymentProcessor, error) {
	stripeKey := os.Getenv("STRIPE_SECRET_KEY")

	if stripeKey == "" {
		return nil, ErrMissingAPIKey
	}

	return stripeService{
		client: client.New(stripeKey, nil),
	}, nil
}

// PerformTransaction performs transaction to payment processor
func (s stripeService) PerformTransaction(input *models.TransactionInput) (*models.Transaction, error) {
	transactionID := fmt.Sprintf("TXN_%s", ulid.Make().String())

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
		Metadata: map[string]string{
			"transaction_id": transactionID,
		},
	}

	var stripeErr *stripe.Error

	result, err := s.client.PaymentIntents.New(params)
	if err != nil {
		if errors.As(err, &stripeErr) && stripeErr.Type != stripe.ErrorTypeCard {
			return nil, api.NewInternalServerError(fmt.Errorf("performing transaction: %s", stripeErr.Code))
		}
	}

	transaction := &models.Transaction{
		TransactionID: transactionID,
		Status:        models.TransactionStatusPending,
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

	if stripeErr != nil && stripeErr.Type == stripe.ErrorTypeCard {
		transaction.Status = models.TransactionStatusFailure
		transaction.FailureReason = string(stripeErr.Code)
	}

	return transaction, nil
}

// RefundTransaction performs refund to payment processor
func (s stripeService) RefundTransaction(metadata map[string]interface{}) (*models.Transaction, error) {
	chargeID, ok := metadata["charge_id"].(string)
	if !ok {
		return nil, ErrMissingChargeID
	}

	params := &stripe.RefundParams{
		Charge: stripe.String(chargeID),
	}

	var stripeErr *stripe.Error

	result, err := s.client.Refunds.New(params)
	if err != nil {
		if errors.As(err, &stripeErr) && stripeErr.Code == stripe.ErrorCodeChargeAlreadyRefunded {
			return nil, api.NewInvalidRequestError(ErrChargeAlreadyRefunded)
		}

		return nil, api.NewInternalServerError(fmt.Errorf("performing refund: %w", err))
	}

	transaction := &models.Transaction{
		Status: models.TransactionStatusPending,
		Type:   models.TransactionTypeRefund,
		AdditionalFields: map[string]interface{}{
			"charge_id":         result.Charge.ID,
			"payment_intent_id": result.PaymentIntent.ID,
			"refund_id":         result.ID,
		},
	}

	return transaction, nil
}
