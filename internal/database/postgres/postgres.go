package postgres

import (
	"context"
	"os"

	"github.com/aledeltoro/simple-online-payment-platform/internal/database"
	"github.com/aledeltoro/simple-online-payment-platform/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresService struct {
	pool *pgxpool.Pool
}

// Init initializes PostgreSQL implementation
func Init(ctx context.Context) (database.Database, error) {
	pool, err := pgxpool.New(ctx, os.Getenv("CONNECTION_URL"))
	if err != nil {
		return nil, err
	}

	return postgresService{
		pool: pool,
	}, nil
}

func (p postgresService) InsertTransaction(ctx context.Context, transaction *models.Transaction) error {
	return nil
}

func (p postgresService) GetTransaction(ctx context.Context, transactionID string) (*models.Transaction, error) {
	return nil, nil
}
