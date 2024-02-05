package service

import (
	"context"
	"errors"

	"github.com/aledeltoro/simple-online-payment-platform/internal/api"
	"github.com/aledeltoro/simple-online-payment-platform/internal/database"
	"github.com/aledeltoro/simple-online-payment-platform/internal/models"
	"github.com/aledeltoro/simple-online-payment-platform/internal/paymentprocessor"
)

var (
	// ErrMissingTransactionID error when transaction ID is missing
	ErrMissingTransactionID = api.NewInvalidRequestError(errors.New("missing transaction id"))
)

// OnlinePaymentService interface to implement business logic for the online payment platform
type OnlinePaymentService interface {
	ProcessPayment(ctx context.Context, amount int64, currency, paymentMethod, description string) (*models.Transaction, error)
	QueryPayment(ctx context.Context, transactionID string) (*models.Transaction, error)
	RefundPayment(ctx context.Context, transactionID string) (*models.Transaction, error)
}

type onlinePaymentService struct {
	database         database.Database
	paymentProcessor paymentprocessor.PaymentProcessor
}

// NewOnlinePaymentService constructor for online payment service
func NewOnlinePaymentService(database database.Database, paymentProcessor paymentprocessor.PaymentProcessor) OnlinePaymentService {
	return onlinePaymentService{
		database:         database,
		paymentProcessor: paymentProcessor,
	}
}

// ProcessPayment handles business logic to process a payment
func (o onlinePaymentService) ProcessPayment(ctx context.Context, amount int64, currency, paymentMethod, description string) (*models.Transaction, error) {
	input := &models.TransactionInput{
		Amount:        amount,
		Currency:      currency,
		PaymentMethod: paymentMethod,
		Description:   description,
	}

	err := input.Validate()
	if err != nil {
		return nil, api.NewInvalidRequestError(err)
	}

	transaction, err := o.paymentProcessor.PerformTransaction(input)
	if err != nil {
		return nil, err
	}

	err = o.database.InsertTransaction(ctx, transaction)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

// QueryPayment handles business logic to query a payment
func (o onlinePaymentService) QueryPayment(ctx context.Context, transactionID string) (*models.Transaction, error) {
	if transactionID == "" {
		return nil, ErrMissingTransactionID
	}

	return o.database.GetTransaction(ctx, transactionID)
}

// RefundPayment handles business logic to refund a payment
func (o onlinePaymentService) RefundPayment(ctx context.Context, transactionID string) (*models.Transaction, error) {
	if transactionID == "" {
		return nil, ErrMissingTransactionID
	}

	transaction, err := o.database.GetTransaction(ctx, transactionID)
	if err != nil {
		return nil, err
	}

	refundedTransaction, err := o.paymentProcessor.RefundTransaction(transaction.AdditionalFields)
	if err != nil {
		return nil, err
	}

	updatedTransaction, err := o.database.UpdateTransaction(ctx, transactionID, refundedTransaction)
	if err != nil {
		return nil, err
	}

	return updatedTransaction, nil
}
