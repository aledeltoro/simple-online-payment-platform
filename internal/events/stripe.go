package events

import (
	"context"

	"github.com/aledeltoro/simple-online-payment-platform/internal/database"
	"github.com/stripe/stripe-go/v76"
)

type stripeEvents struct {
	database database.Database
	event    stripe.Event
}

func newStripeEvent(database database.Database) Events {
	return &stripeEvents{
		database: database,
	}
}

func (e *stripeEvents) VerifyEvent(payload []byte) error {
	// Verify signature
	// Verify if it's supported event
	return nil
}

func (e *stripeEvents) ProcessEvent(ctx context.Context) error {
	return nil
}
