package database

import (
	"context"
	"errors"

	"github.com/aledeltoro/simple-online-payment-platform/internal/models"
)

var (
	// ErrTransactionNotFound error when transaction was not found
	ErrTransactionNotFound = errors.New("transaction not found")
	// ErrMultipleRowsAffected error when multiple rows were affeted in an operation
	ErrMultipleRowsAffected = errors.New("multiple rows affected")
)

// Database service to handle database integrations
type Database interface {
	InsertTransaction(context.Context, *models.Transaction) error
	GetTransaction(ctx context.Context, transactionID string) (*models.Transaction, error)
	UpdateTransaction(ctx context.Context, transactionID string, updatedTransaction *models.Transaction) error
	Close()
}
