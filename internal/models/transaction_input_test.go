package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateTransactionInput(t *testing.T) {
	c := require.New(t)

	input := TransactionInput{}

	c.ErrorIs(input.Validate(), ErrInvalidAmount)

	input.Amount = 2000

	c.ErrorIs(input.Validate(), ErrMissingCurrency)

	input.Currency = "usd"

	c.ErrorIs(input.Validate(), ErrMissingPaymentMethod)

	input.PaymentMethod = "pm_card_visa"

	c.NoError(input.Validate())
}
