package events

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockStripe mock for Stripe events implementation
type MockStripe struct {
	mock.Mock
}

// VerifyEvent mock implementation
func (m *MockStripe) VerifyEvent() error {
	args := m.Called()

	return args.Error(0)
}

// ProcessEvent mock implementation
func (m *MockStripe) ProcessEvent(ctx context.Context) error {
	args := m.Called()

	return args.Error(0)
}
