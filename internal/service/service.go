package handler

import (
	"github.com/aledeltoro/simple-online-payment-platform/internal/database"
	"github.com/aledeltoro/simple-online-payment-platform/internal/models"
	"github.com/aledeltoro/simple-online-payment-platform/internal/paymentservice"
)

// OnlinePaymentService interface to implement business logic for the online payment platform
type OnlinePaymentService interface {
	ProcessPayment(input *models.TransactionInput) (*models.Transaction, error)
	QueryPayment(transactionID string) (*models.Transaction, error)
	RefundPayment(transactionID string) (*models.Transaction, error)
}

type onlinePaymentService struct {
	database         database.Database
	paymentProcessor paymentservice.PaymentService
}

func NewOnlinePaymentService(database database.Database, paymentProcessor paymentservice.PaymentService) OnlinePaymentService {
	return onlinePaymentService{
		database:         database,
		paymentProcessor: paymentProcessor,
	}
}

func (o onlinePaymentService) ProcessPayment(input *models.TransactionInput) (*models.Transaction, error) {
	return nil, nil
}

func (o onlinePaymentService) QueryPayment(transactionID string) (*models.Transaction, error) {
	return nil, nil
}

func (o onlinePaymentService) RefundPayment(transactionID string) (*models.Transaction, error) {
	return nil, nil
}
