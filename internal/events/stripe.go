package events

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/aledeltoro/simple-online-payment-platform/internal/database"
	"github.com/aledeltoro/simple-online-payment-platform/internal/models"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/webhook"
)

var supportedStripeEvents = map[stripe.EventType]bool{
	stripe.EventTypePaymentIntentSucceeded:     true,
	stripe.EventTypePaymentIntentPaymentFailed: true,
	stripe.EventTypeChargeRefunded:             true,
}

type stripeEvents struct {
	database database.Database
	request  *http.Request
	event    stripe.Event
}

func newStripeEvent(database database.Database, request *http.Request) Events {
	return &stripeEvents{
		database: database,
		request:  request,
	}
}

func (e *stripeEvents) VerifyEvent() error {
	webhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET_KEY")
	stripeSignature := e.request.Header.Get("Stripe-Signature")

	payload, err := io.ReadAll(e.request.Body)
	if err != nil {
		return err
	}

	event, err := webhook.ConstructEvent(payload, stripeSignature, webhookSecret)
	if err != nil {
		return err
	}

	e.event = event

	return nil
}

func (e *stripeEvents) ProcessEvent(ctx context.Context) error {
	if _, ok := supportedStripeEvents[e.event.Type]; !ok {
		return ErrUnsupportedEvent
	}

	transaction := &models.Transaction{}

	switch e.event.Type {
	case stripe.EventTypePaymentIntentSucceeded, stripe.EventTypePaymentIntentPaymentFailed:
		var paymentIntent *stripe.PaymentIntent

		err := json.Unmarshal(e.event.Data.Raw, &paymentIntent)
		if err != nil {
			return err
		}

		transaction.TransactionID = paymentIntent.Metadata["transaction_id"]
		transaction.Type = models.TransactionType(paymentIntent.Status)

	case stripe.EventTypeChargeRefunded:
		var charge *stripe.Charge

		err := json.Unmarshal(e.event.Data.Raw, &charge)
		if err != nil {
			return err
		}

		transaction.TransactionID = charge.Metadata["transaction_id"]
		transaction.Type = models.TransactionType(charge.Status)
	}

	err := e.database.UpdateTransaction(ctx, transaction.TransactionID, transaction)
	if err != nil {
		return err
	}

	return nil
}
