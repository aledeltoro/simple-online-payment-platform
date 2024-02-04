package events

import (
	"context"
	"errors"

	"github.com/aledeltoro/simple-online-payment-platform/internal/database"
	"github.com/aledeltoro/simple-online-payment-platform/internal/models"
)

var (
	ErrUnsupportedProvider = errors.New("unsupported provider")
)

// Events interface to implement business logic to handle incoming events from the payment provider
type Events interface {
	VerifyEvent(payload []byte) error
	ProcessEvent(ctx context.Context) error
}

func NewEvent(provider models.PaymentProvider, database database.Database) (Events, error) {
	if provider == models.PaymentProviderStripe {
		return newStripeEvent(database), nil
	}

	return nil, ErrUnsupportedProvider
}
