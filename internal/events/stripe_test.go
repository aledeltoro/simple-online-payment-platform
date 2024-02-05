package events

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/aledeltoro/simple-online-payment-platform/internal/database/postgres"
	"github.com/aledeltoro/simple-online-payment-platform/internal/models"
	"github.com/stretchr/testify/require"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/webhook"
)

func TestVerifyEvent(t *testing.T) {
	c := require.New(t)

	t.Setenv("STRIPE_WEBHOOK_SECRET_KEY", "stripe_webhook_secret_key")

	stripeEvent := stripe.Event{
		Type:       stripe.EventTypeAccountUpdated,
		APIVersion: "2023-10-16",
	}

	payload, err := json.Marshal(stripeEvent)
	c.NoError(err)

	body := bytes.NewBuffer(payload)

	signedPayload := webhook.GenerateTestSignedPayload(&webhook.UnsignedPayload{Payload: payload, Secret: "stripe_webhook_secret_key"})

	req, err := http.NewRequest(http.MethodPost, "/payments/stripe/events", body)
	c.NoError(err)

	req.Header.Set("Stripe-Signature", signedPayload.Header)

	eventHandler := stripeEvents{
		request: req,
	}

	err = eventHandler.VerifyEvent()
	c.NoError(err)
}

func TestProcessEventUnsupportedEvent(t *testing.T) {
	c := require.New(t)

	stripeEvent := stripe.Event{
		Type: stripe.EventTypeAccountUpdated,
	}

	eventHandler := stripeEvents{
		event: stripeEvent,
	}

	err := eventHandler.ProcessEvent(context.Background())
	c.ErrorIs(err, ErrUnsupportedEvent)
}

func TestProcessEventPaymentIntentEvent(t *testing.T) {
	c := require.New(t)

	paymentIntent := &stripe.PaymentIntent{
		ID:     "pi_123",
		Status: stripe.PaymentIntentStatusSucceeded,
		Metadata: map[string]string{
			"transaction_id": "TXN_123",
		},
	}

	transaction := &models.Transaction{
		TransactionID: paymentIntent.Metadata["transaction_id"],
		Status:        models.TransactionStatus(paymentIntent.Status),
		Type:          models.TransactionTypeCharge,
	}

	rawData, err := json.Marshal(paymentIntent)
	c.NoError(err)

	stripeEvent := stripe.Event{
		Type: stripe.EventTypePaymentIntentSucceeded,
		Data: &stripe.EventData{
			Raw: rawData,
		},
	}

	mockDatabase := postgres.MockPostgres{}

	mockDatabase.On("UpdateTransaction", context.Background(), "TXN_123", transaction).Return(transaction, nil)

	eventHandler := stripeEvents{
		event:    stripeEvent,
		database: &mockDatabase,
	}

	err = eventHandler.ProcessEvent(context.Background())
	c.NoError(err)
}

func TestProcessEventChargeRefundedEvent(t *testing.T) {
	c := require.New(t)

	charge := &stripe.Charge{
		ID:     "ch_123",
		Status: stripe.ChargeStatusSucceeded,
		Metadata: map[string]string{
			"transaction_id": "TXN_123",
		},
	}

	transaction := &models.Transaction{
		TransactionID: charge.Metadata["transaction_id"],
		Status:        models.TransactionStatus(charge.Status),
		Type:          models.TransactionTypeRefund,
	}

	rawData, err := json.Marshal(charge)
	c.NoError(err)

	stripeEvent := stripe.Event{
		Type: stripe.EventTypeChargeRefunded,
		Data: &stripe.EventData{
			Raw: rawData,
		},
	}

	mockDatabase := postgres.MockPostgres{}

	mockDatabase.On("UpdateTransaction", context.Background(), "TXN_123", transaction).Return(transaction, nil)

	eventHandler := stripeEvents{
		event:    stripeEvent,
		database: &mockDatabase,
	}

	err = eventHandler.ProcessEvent(context.Background())
	c.NoError(err)
}
