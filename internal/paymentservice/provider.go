package paymentservice

import "github.com/aledeltoro/simple-online-payment-platform/internal/models"

// PaymentService service to handle interactions with an integrated payment provider
type PaymentService interface {
	PerformTransaction(input *models.TransactionInput) (*models.Transaction, error)
	QueryTransaction(transactionID string) (*models.Transaction, error)
	RefundTransaction(transactionID string) (*models.Transaction, error)
}
