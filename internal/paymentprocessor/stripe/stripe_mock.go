package stripe

import (
	"bytes"

	"github.com/aledeltoro/simple-online-payment-platform/internal/models"
	"github.com/stretchr/testify/mock"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/form"
)

// MockStripe mock object for Stripe implementation
type MockStripe struct {
	mock.Mock
}

// PerformTransaction mock implementation
func (m *MockStripe) PerformTransaction(input *models.TransactionInput) (*models.Transaction, error) {
	args := m.Called(input)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

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
func (m *MockStripe) RefundTransaction(metadata map[string]interface{}) (*models.Transaction, error) {
	args := m.Called(metadata)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*models.Transaction), args.Error(1)
}

// mockStripeBackend mock for Stripe Backend interface
type mockStripeBackend struct {
	mock.Mock
}

// Call mock for Call method in Stripe Backend interface
func (m *mockStripeBackend) Call(method, path, key string, params stripe.ParamsContainer, v stripe.LastResponseSetter) error {
	args := m.Called(method, path, key, params, v)

	return args.Error(0)
}

// CallStreaming mock for Call method in Stripe Backend interface
func (m *mockStripeBackend) CallStreaming(method, path, key string, params stripe.ParamsContainer, v stripe.StreamingLastResponseSetter) error {
	args := m.Called(method, path, key, params, v)

	return args.Error(0)
}

// CallRaw mock for Call method in Stripe Backend interface
func (m *mockStripeBackend) CallRaw(method, path, key string, body *form.Values, params *stripe.Params, v stripe.LastResponseSetter) error {
	args := m.Called(method, path, key, params, v)

	return args.Error(0)
}

// CallMultipart mock for Call method in Stripe Backend interface
func (m *mockStripeBackend) CallMultipart(method, path, key, boundary string, body *bytes.Buffer, params *stripe.Params, v stripe.LastResponseSetter) error {
	args := m.Called(method, path, key, params, v)

	return args.Error(0)
}

// CallMultipart mock for Call method in Stripe Backend interface
func (m *mockStripeBackend) SetMaxNetworkRetries(maxNetworkRetries int64) {
	m.Called(maxNetworkRetries)
}
