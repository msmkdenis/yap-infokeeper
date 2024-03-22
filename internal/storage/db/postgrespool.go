package db

import (
	"context"
	"fmt"
	"log/slog"

	apperr "github.com/msmkdenis/yap-infokeeper/pkg/apperror"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresPool struct {
	DB *pgxpool.Pool
}

func NewPostgresPool(ctx context.Context, connection string) (*PostgresPool, error) {
	dbPool, err := pgxpool.New(ctx, connection)
	if err != nil {
		return nil, apperr.NewValueError(fmt.Sprintf("Unable to connect to database with connection %s", connection), apperr.Caller(), err)
	}
	slog.Info("Successful connection", slog.String("database", dbPool.Config().ConnConfig.Database))

	err = dbPool.Ping(ctx)
	if err != nil {
		return nil, apperr.NewValueError("Unable to ping database", apperr.Caller(), err)
	}
	slog.Info("Successful ping", slog.String("database", dbPool.Config().ConnConfig.Database))

	return &PostgresPool{DB: dbPool}, nil
}
