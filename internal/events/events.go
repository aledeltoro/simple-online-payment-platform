package events

import (
	"context"
	"errors"
	"net/http"

	"github.com/aledeltoro/simple-online-payment-platform/internal/database"
	"github.com/aledeltoro/simple-online-payment-platform/internal/models"
)

var (
	ErrUnsupportedProvider = errors.New("unsupported provider")
	ErrUnsupportedEvent    = errors.New("unsupported event")
)

// Events interface to implement business logic to handle incoming events from the payment provider
type Events interface {
	VerifyEvent() error
	ProcessEvent(ctx context.Context) error
}

func NewEvent(provider models.PaymentProvider, database database.Database, request *http.Request) (Events, error) {
	if provider == models.PaymentProviderStripe {
		return newStripeEvent(database, request), nil
	}

	return nil, ErrUnsupportedProvider
}
