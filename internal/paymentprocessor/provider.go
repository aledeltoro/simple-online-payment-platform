package paymentprocessor

import "github.com/aledeltoro/simple-online-payment-platform/internal/models"

// PaymentProcessor service to handle interactions with an integrated payment provider
type PaymentProcessor interface {
	PerformTransaction(input *models.TransactionInput) (*models.Transaction, error)
	QueryTransaction(id string) (*models.Transaction, error)
	RefundTransaction(metadata map[string]interface{}) (*models.Transaction, error)
}
