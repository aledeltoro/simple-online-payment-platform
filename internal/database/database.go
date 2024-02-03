package database

import (
	"context"

	"github.com/aledeltoro/simple-online-payment-platform/internal/models"
)

// Database service to handle database integrations
type Database interface {
	InsertTransaction(context.Context, *models.Transaction) error
	GetTransaction(ctx context.Context, transactionID string) (*models.Transaction, error)
}
