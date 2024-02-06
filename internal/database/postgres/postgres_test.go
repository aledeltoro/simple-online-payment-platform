package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"regexp"
	"testing"

	"github.com/aledeltoro/simple-online-payment-platform/internal/database"
	"github.com/aledeltoro/simple-online-payment-platform/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/require"
)

func TestInsertTransaction(t *testing.T) {
	c := require.New(t)

	mock, err := pgxmock.NewPool()
	c.NoError(err)

	defer mock.Close()

	transaction := &models.Transaction{
		TransactionID: "TXN123",
		Status:        models.TransactionStatusSucceeded,
		Description:   "Sample description",
		Provider:      models.PaymentProviderStripe,
		Amount:        2000,
		Currency:      "USD",
		Type:          models.TransactionTypeCharge,
		AdditionalFields: map[string]interface{}{
			"charge_id": "ch_123",
		},
	}

	mock.ExpectExec("INSERT INTO transactions_history").WithArgs(
		transaction.TransactionID,
		transaction.Status,
		transaction.Description,
		transaction.FailureReason,
		transaction.Provider,
		transaction.Amount,
		transaction.Currency,
		transaction.Type,
		transaction.AdditionalFields,
	).WillReturnResult(pgxmock.NewResult("INSERT", 1))

	service := postgresService{pool: mock}

	err = service.InsertTransaction(context.Background(), transaction)
	c.NoError(err)
}

func TestInsertTransactionFailure(t *testing.T) {
	c := require.New(t)

	mock, err := pgxmock.NewPool()
	c.NoError(err)

	defer mock.Close()

	transaction := &models.Transaction{
		TransactionID: "TXN123",
		Status:        models.TransactionStatusSucceeded,
		Description:   "Sample description",
		Provider:      models.PaymentProviderStripe,
		Amount:        2000,
		Currency:      "USD",
		Type:          models.TransactionTypeCharge,
		AdditionalFields: map[string]interface{}{
			"charge_id": "ch_123",
		},
	}

	mock.ExpectExec("INSERT INTO transactions_history").WithArgs(
		transaction.TransactionID,
		transaction.Status,
		transaction.Description,
		transaction.FailureReason,
		transaction.Provider,
		transaction.Amount,
		transaction.Currency,
		transaction.Type,
		transaction.AdditionalFields,
	).WillReturnError(sql.ErrConnDone)

	service := postgresService{pool: mock}

	err = service.InsertTransaction(context.Background(), transaction)
	c.ErrorIs(err, sql.ErrConnDone)
}

func TestGetTransaction(t *testing.T) {
	c := require.New(t)

	mock, err := pgxmock.NewPool()
	c.NoError(err)

	defer mock.Close()

	columns := []string{"transaction_id", "status", "description", "failure_reason", "payment_provider", "amount", "currency", "type", "additional_fields"}

	rows := mock.NewRows(columns)

	expectedtransaction := &models.Transaction{
		TransactionID: "TXN123",
		Status:        models.TransactionStatusSucceeded,
		Description:   "Sample description",
		Provider:      models.PaymentProviderStripe,
		Amount:        2000,
		Currency:      "USD",
		Type:          models.TransactionTypeCharge,
		AdditionalFields: map[string]interface{}{
			"charge_id": "ch_123",
		},
	}

	marshalledAdditionalFields, err := json.Marshal(expectedtransaction.AdditionalFields)
	c.NoError(err)

	rows.AddRow(
		expectedtransaction.TransactionID,
		expectedtransaction.Status,
		expectedtransaction.Description,
		expectedtransaction.FailureReason,
		expectedtransaction.Provider,
		expectedtransaction.Amount,
		expectedtransaction.Currency,
		expectedtransaction.Type,
		string(marshalledAdditionalFields),
	)

	query := `
	SELECT
		transaction_id,
		status,
		description,
		failure_reason,
		payment_provider,
		amount,
		currency,
		type,
		additional_fields
	FROM transactions_history
	WHERE transaction_id = $1
	`

	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("TXN123").WillReturnRows(rows)

	service := postgresService{pool: mock}

	transaction, err := service.GetTransaction(context.Background(), "TXN123")
	c.NoError(err)
	c.Equal(expectedtransaction, transaction)
}

func TestGetTransactionNoRows(t *testing.T) {
	c := require.New(t)

	mock, err := pgxmock.NewPool()
	c.NoError(err)

	defer mock.Close()

	query := `
	SELECT
		transaction_id,
		status,
		description,
		failure_reason,
		payment_provider,
		amount,
		currency,
		type,
		additional_fields
	FROM transactions_history
	WHERE transaction_id = $1
	`

	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("TXN123").WillReturnError(pgx.ErrNoRows)

	service := postgresService{pool: mock}

	transaction, err := service.GetTransaction(context.Background(), "TXN123")
	c.Nil(transaction)
	c.ErrorIs(err, database.ErrTransactionNotFound)
}

func TestGetTransactionFailure(t *testing.T) {
	c := require.New(t)

	mock, err := pgxmock.NewPool()
	c.NoError(err)

	defer mock.Close()

	query := `
	SELECT
		transaction_id,
		status,
		description,
		failure_reason,
		payment_provider,
		amount,
		currency,
		type,
		additional_fields
	FROM transactions_history
	WHERE transaction_id = $1
	`

	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("TXN123").WillReturnError(sql.ErrConnDone)

	service := postgresService{pool: mock}

	transaction, err := service.GetTransaction(context.Background(), "TXN123")
	c.Nil(transaction)
	c.ErrorIs(err, sql.ErrConnDone)
}

func TestUpdateTransaction(t *testing.T) {
	c := require.New(t)

	mock, err := pgxmock.NewPool()
	c.NoError(err)

	defer mock.Close()

	expectedTransaction := &models.Transaction{
		TransactionID: "TXN_123",
		Status:        models.TransactionStatusSucceeded,
		Description:   "Sample description",
		Amount:        2000,
		Provider:      models.PaymentProviderStripe,
		Currency:      "usd",
		Type:          models.TransactionTypeRefund,
		AdditionalFields: map[string]interface{}{
			"charge_id": "ch_123",
			"refund_id": "re_123",
		},
	}

	marshalledAdditionalFields, err := json.Marshal(expectedTransaction.AdditionalFields)
	c.NoError(err)

	columns := []string{"transaction_id", "status", "description", "failure_reason", "provider", "amount", "currency", "type", "additional_fields"}

	rows := mock.NewRows(columns)

	rows.AddRow(
		expectedTransaction.TransactionID,
		expectedTransaction.Status,
		expectedTransaction.Description,
		expectedTransaction.FailureReason,
		expectedTransaction.Provider,
		expectedTransaction.Amount,
		expectedTransaction.Currency,
		expectedTransaction.Type,
		string(marshalledAdditionalFields),
	)

	query := `
	UPDATE transactions_history
	SET
		status = COALESCE($1, status),
		type = COALESCE($2, type),
		additional_fields = COALESCE($3, additional_fields)
	WHERE transaction_id = $4`

	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(expectedTransaction.Status, expectedTransaction.Type, expectedTransaction.AdditionalFields, expectedTransaction.TransactionID).WillReturnRows(rows)

	service := postgresService{pool: mock}

	transaction, err := service.UpdateTransaction(context.Background(), "TXN_123", expectedTransaction)
	c.NoError(err)
	c.Equal(expectedTransaction, transaction)
}

func TestUpdateTransactionFailure(t *testing.T) {
	c := require.New(t)

	mock, err := pgxmock.NewPool()
	c.NoError(err)

	defer mock.Close()

	transaction := &models.Transaction{
		TransactionID: "TXN_123",
		Status:        models.TransactionStatusSucceeded,
		Type:          models.TransactionTypeRefund,
		AdditionalFields: map[string]interface{}{
			"charge_id": "ch_123",
			"refund_id": "re_123",
		},
	}

	marshalledAdditionalFields, err := json.Marshal(transaction.AdditionalFields)
	c.NoError(err)

	columns := []string{"transaction_id", "status", "type", "additional_fields"}

	rows := mock.NewRows(columns)

	rows.AddRow(
		transaction.TransactionID,
		transaction.Status,
		transaction.Type,
		string(marshalledAdditionalFields),
	)

	query := `
	UPDATE transactions_history
	SET
		status = COALESCE($1, status),
		type = COALESCE($2, type),
		additional_fields = COALESCE($3, additional_fields)
	WHERE transaction_id = $4
	RETURNING
		transaction_id,
		status,
		description,
		failure_reason,
		payment_provider,
		amount,
		currency,
		type,
		additional_fields`

	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(transaction.Status, transaction.Type, transaction.AdditionalFields, transaction.TransactionID).WillReturnError(sql.ErrConnDone)

	service := postgresService{pool: mock}

	_, err = service.UpdateTransaction(context.Background(), "TXN_123", transaction)
	c.ErrorIs(err, sql.ErrConnDone)

}
