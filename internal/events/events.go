package events

import (
	"context"
	"errors"
	"net/http"

	"github.com/aledeltoro/simple-online-payment-platform/internal/database"
	"github.com/aledeltoro/simple-online-payment-platform/internal/models"
)

var (
	// ErrUnsupportedProvider error when incoming provider is unsupported
	ErrUnsupportedProvider = errors.New("unsupported provider")
	// ErrUnsupportedEvent error when event is not supported by event handler
	ErrUnsupportedEvent = errors.New("unsupported event")
)

// Events interface to implement business logic to handle incoming events from the payment provider
type Events interface {
	VerifyEvent() error
	ProcessEvent(ctx context.Context) error
}

// NewEvent constructor to return the proper event handler
func NewEvent(provider models.PaymentProvider, database database.Database, request *http.Request) (Events, error) {
	if provider == models.PaymentProviderStripe {
		return newStripeEvent(database, request), nil
	}

	return nil, ErrUnsupportedProvider
}
