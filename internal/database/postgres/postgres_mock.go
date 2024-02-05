package postgres

import (
	"context"

	"github.com/aledeltoro/simple-online-payment-platform/internal/models"
	"github.com/stretchr/testify/mock"
)

// MockPostgres mock implementation
type MockPostgres struct {
	mock.Mock
}

// InsertTransaction mocks operation to insert an item to the database
func (m *MockPostgres) InsertTransaction(ctx context.Context, transaction *models.Transaction) error {
	args := m.Called(ctx, transaction)

	return args.Error(0)
}

// GetTransaction mocks operation to fetch an item given its ID
func (m *MockPostgres) GetTransaction(ctx context.Context, transactionID string) (*models.Transaction, error) {
	args := m.Called(ctx, transactionID)

	return args.Get(0).(*models.Transaction), args.Error(1)
}

// UpdateTransaction mocks operation to update an item given its ID
func (m *MockPostgres) UpdateTransaction(ctx context.Context, transactionID string, updatedTransaction *models.Transaction) (*models.Transaction, error) {
	args := m.Called(ctx, transactionID, updatedTransaction)

	return args.Get(0).(*models.Transaction), args.Error(1)
}

// Close mock operation to close a database connection
func (m *MockPostgres) Close() {}
