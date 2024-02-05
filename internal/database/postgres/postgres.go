package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/aledeltoro/simple-online-payment-platform/internal/database"
	"github.com/aledeltoro/simple-online-payment-platform/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// pgxIface interface required to mock behavior of PGX Pool library
type pgxIface interface {
	Begin(context.Context) (pgx.Tx, error)
	Close()
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
}

type postgresService struct {
	pool pgxIface
}

// Init initializes PostgreSQL implementation
func Init(ctx context.Context) (database.Database, error) {
	host := os.Getenv("DATABASE_HOST")
	port := os.Getenv("DATABASE_PORT")
	user := os.Getenv("DATABASE_USER")
	databaseName := os.Getenv("DATABASE_NAME")

	connectionString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable", host, port, user, databaseName)

	pool, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return postgresService{
		pool: pool,
	}, nil
}

func (p postgresService) Close() {
	p.pool.Close()
}

func (p postgresService) InsertTransaction(ctx context.Context, transaction *models.Transaction) error {
	query := `
	INSERT INTO transactions_history(
		transaction_id,
		status,
		failure_reason,
		payment_provider,
		amount,
		currency,
		type,
		additional_fields
	) VALUES($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := p.pool.Exec(ctx, query, transaction.TransactionID, transaction.Status, transaction.FailureReason, transaction.Provider, transaction.Amount, transaction.Currency, transaction.Type, transaction.AdditionalFields)

	return err
}

func (p postgresService) GetTransaction(ctx context.Context, transactionID string) (*models.Transaction, error) {
	query := `
	SELECT
		transaction_id,
		status,
		failure_reason,
		payment_provider,
		amount,
		currency,
		type,
		additional_fields
	FROM transactions_history
	WHERE transaction_id = $1
	`

	row := p.pool.QueryRow(ctx, query, transactionID)

	var transaction models.Transaction
	var additionalFieldsJSON string

	err := row.Scan(
		&transaction.TransactionID,
		&transaction.Status,
		&transaction.FailureReason,
		&transaction.Provider,
		&transaction.Amount,
		&transaction.Currency,
		&transaction.Type,
		&additionalFieldsJSON,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, database.ErrTransactionNotFound
	}

	if err != nil {
		return nil, err
	}

	if additionalFieldsJSON != "" {
		err = json.Unmarshal([]byte(additionalFieldsJSON), &transaction.AdditionalFields)
		if err != nil {
			return nil, err
		}
	}

	return &transaction, nil
}

func (p postgresService) UpdateTransaction(ctx context.Context, transactionID string, updatedTransaction *models.Transaction) error {
	query := `
	UPDATE transactions_history
	SET
		status = COALESCE($1, status),
		type = COALESCE($2, type),
		additional_fields = COALESCE($3, additional_fields)
	WHERE transaction_id = $4
	`

	_, err := p.pool.Exec(ctx, query, updatedTransaction.Status, updatedTransaction.Type, updatedTransaction.AdditionalFields, transactionID)
	if err != nil {
		return err
	}

	return nil
}
