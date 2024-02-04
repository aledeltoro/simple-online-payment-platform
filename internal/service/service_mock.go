package service

import (
	"context"

	"github.com/aledeltoro/simple-online-payment-platform/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockOnlinePaymentService struct {
	mock.Mock
}

func (m *MockOnlinePaymentService) ProcessPayment(ctx context.Context, amount int64, currency, paymentMethod, description string) (*models.Transaction, error) {
	args := m.Called(ctx, amount, currency, paymentMethod, description)

	return args.Get(0).(*models.Transaction), args.Error(1)
}

func (m *MockOnlinePaymentService) QueryPayment(ctx context.Context, transactionID string) (*models.Transaction, error) {
	args := m.Called(ctx, transactionID)

	return args.Get(0).(*models.Transaction), args.Error(1)
}

func (m *MockOnlinePaymentService) RefundPayment(ctx context.Context, transactionID string) (*models.Transaction, error) {
	args := m.Called(ctx, transactionID)

	return args.Get(0).(*models.Transaction), args.Error(1)
}
