package service

import (
	"context"

	"github.com/aledeltoro/simple-online-payment-platform/internal/models"
	"github.com/stretchr/testify/mock"
)

// MockOnlinePaymentService mock object for online payment service implementation
type MockOnlinePaymentService struct {
	mock.Mock
}

// ProcessPayment mock implementation
func (m *MockOnlinePaymentService) ProcessPayment(ctx context.Context, amount int64, currency, paymentMethod, description string) (*models.Transaction, error) {
	args := m.Called(ctx, amount, currency, paymentMethod, description)

	return args.Get(0).(*models.Transaction), args.Error(1)
}

// QueryPayment mock implementation
func (m *MockOnlinePaymentService) QueryPayment(ctx context.Context, transactionID string) (*models.Transaction, error) {
	args := m.Called(ctx, transactionID)

	return args.Get(0).(*models.Transaction), args.Error(1)
}

// RefundPayment mock implementation
func (m *MockOnlinePaymentService) RefundPayment(ctx context.Context, transactionID string) (*models.Transaction, error) {
	args := m.Called(ctx, transactionID)

	return args.Get(0).(*models.Transaction), args.Error(1)
}
