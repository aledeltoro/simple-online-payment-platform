package events

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockStripe struct {
	mock.Mock
}

func (m *MockStripe) VerifyEvent() error {
	args := m.Called()

	return args.Error(0)
}

func (m *MockStripe) ProcessEvent(ctx context.Context) error {
	args := m.Called()

	return args.Error(0)
}
