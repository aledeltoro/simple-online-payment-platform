package stripe

import (
	"github.com/aledeltoro/simple-online-payment-platform/internal/models"
	"github.com/aledeltoro/simple-online-payment-platform/internal/paymentservice"
)

type stripeService struct{}

// New initializes implementation of Stripe service
func New() paymentservice.PaymentService {
	return stripeService{}
}

// PerformTransaction method to perform a transaction
func (s stripeService) PerformTransaction(input *models.TransactionInput) (*models.Transaction, error) {
	return nil, nil
}

func (s stripeService) QueryTransaction(transactionID string) (*models.Transaction, error) {
	return nil, nil
}

func (s stripeService) RefundTransaction(transactionID string) (*models.Transaction, error) {
	return nil, nil
}
