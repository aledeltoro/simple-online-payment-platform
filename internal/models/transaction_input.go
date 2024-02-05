package models

import (
	"errors"
	"fmt"
)

// TransactionInput inputs to perform a transaction
type TransactionInput struct {
	Amount        int64  `json:"amount"`
	Currency      string `json:"currency"`
	PaymentMethod string `json:"payment_method"`
	Description   string `json:"description"`
}

var (
	// ErrInvalidAmount error when amount is equal or less than zero
	ErrInvalidAmount = errors.New("invalid amount")
	// ErrMissingCurrency error when currency is missing
	ErrMissingCurrency = errors.New("missing currency")
	// ErrMissingPaymentMethod error when payment method is missing
	ErrMissingPaymentMethod = errors.New("missing payment method")
)

// Validate validate the inputs required for a transaction
func (ti *TransactionInput) Validate() error {
	if ti.Amount <= 0 {
		return ErrInvalidAmount
	}

	if ti.Currency == "" {
		return ErrMissingCurrency
	}

	if ti.PaymentMethod == "" {
		return ErrMissingPaymentMethod
	}

	if ti.Description == "" {
		ti.Description = fmt.Sprintf("Transaction for payment amount of %d", ti.Amount)
	}

	return nil
}
