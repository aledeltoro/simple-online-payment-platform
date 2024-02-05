package postgres

import (
	"context"
	"encoding/json"
	"regexp"
	"testing"

	"github.com/aledeltoro/simple-online-payment-platform/internal/models"
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

func TestGetTransaction(t *testing.T) {
	c := require.New(t)

	mock, err := pgxmock.NewPool()
	c.NoError(err)

	defer mock.Close()

	marshalledAdditionalFields, err := json.Marshal(map[string]interface{}{
		"charge_id": "ch_123",
	})
	c.NoError(err)

	rows := mock.NewRows([]string{"transaction_id", "status", "description", "failure_reason", "payment_provider", "amount", "currency", "type", "additional_fields"})
	rows.AddRow("TXN123", models.TransactionStatusSucceeded, "Description", "", models.PaymentProviderStripe, 2000, "usd", models.TransactionTypeCharge, string(marshalledAdditionalFields))

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
	c.NotNil(transaction)
}

func TestUpdateTransaction(t *testing.T) {
	c := require.New(t)

	mock, err := pgxmock.NewPool()
	c.NoError(err)

	defer mock.Close()

	marshalledAdditionalFields, err := json.Marshal(map[string]interface{}{
		"charge_id": "ch_123",
	})
	c.NoError(err)

	rows := mock.NewRows([]string{"transaction_id", "type", "additional_fields"})
	rows.AddRow("TXN123", models.TransactionTypeCharge, string(marshalledAdditionalFields))

	query := `
	UPDATE transactions_history
	SET
		status = COALESCE($1, status),
		type = COALESCE($2, type),
		additional_fields = COALESCE($3, additional_fields)
	WHERE transaction_id = $4`

	transaction := &models.Transaction{
		TransactionID: "TXN123",
		Status:        models.TransactionStatusSucceeded,
		Type:          models.TransactionTypeRefund,
		AdditionalFields: map[string]interface{}{
			"charge_id": "ch_123",
			"refund_id": "rf_123",
		},
	}

	mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(transaction.Status, transaction.Type, transaction.AdditionalFields, transaction.TransactionID).WillReturnResult(pgxmock.NewResult("1", 1))

	service := postgresService{pool: mock}

	err = service.UpdateTransaction(context.Background(), "TXN123", transaction)
	c.NoError(err)
}
