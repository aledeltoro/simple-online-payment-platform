package events

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVerifyEvent(t *testing.T) {
	c := require.New(t)

	mockStripe := MockStripe{}

	mockStripe.On("VerifyEvent").Return(nil)

	c.NoError(mockStripe.VerifyEvent())
}

func TestProcessEvent(t *testing.T) {
	c := require.New(t)

	mockStripe := MockStripe{}

	mockStripe.On("ProcessEvent", context.Background()).Return(nil)

	err := mockStripe.ProcessEvent(context.Background())
	c.NoError(err)
}
