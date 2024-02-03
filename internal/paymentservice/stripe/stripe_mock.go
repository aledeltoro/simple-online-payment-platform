package stripe

import (
	"github.com/aledeltoro/simple-online-payment-platform/internal/models"
	"github.com/stretchr/testify/mock"
)

// MockStripe mock object for Stripe implementation
type MockStripe struct {
	mock.Mock
}

// PerformTransaction mock implementation
func (m *MockStripe) PerformTransaction(input *models.TransactionInput) (*models.Transaction, error) {
	args := m.Called(input)

	return args.Get(0).(*models.Transaction), args.Error(1)
}

// QueryTransaction mock implementation
func (m *MockStripe) QueryTransaction(id string) (*models.Transaction, error) {
	args := m.Called(id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*models.Transaction), args.Error(1)
}

// RefundTransaction mock implementation
func (m *MockStripe) RefundTransaction(id string) (*models.Transaction, error) {
	args := m.Called(id)

	return args.Get(0).(*models.Transaction), args.Error(1)
}
