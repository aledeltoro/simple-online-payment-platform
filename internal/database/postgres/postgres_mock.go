package postgres

import (
	"context"

	"github.com/aledeltoro/simple-online-payment-platform/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockPostgres struct {
	mock.Mock
}

func (m *MockPostgres) InsertTransaction(ctx context.Context, transaction *models.Transaction) error {
	args := m.Called(ctx, transaction)

	return args.Error(0)
}

func (m *MockPostgres) GetTransaction(ctx context.Context, transactionID string) (*models.Transaction, error) {
	args := m.Called(ctx, transactionID)

	return args.Get(0).(*models.Transaction), args.Error(1)
}

func (m *MockPostgres) UpdateTransaction(ctx context.Context, transactionID string, updatedTransaction *models.Transaction) error {
	args := m.Called(ctx, transactionID, updatedTransaction)

	return args.Error(0)
}

func (m *MockPostgres) Close() {
	return
}
