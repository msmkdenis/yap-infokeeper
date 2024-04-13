package db

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/msmkdenis/yap-infokeeper/pkg/caller"
)

type PostgresPool struct {
	DB *pgxpool.Pool
}

func NewPostgresPool(ctx context.Context, connection string) (*PostgresPool, error) {
	dbPool, err := pgxpool.New(ctx, connection)
	if err != nil {
		return nil, fmt.Errorf("%s %w", caller.CodeLine(), err)
	}
	slog.Info("Successful connection", slog.String("database", dbPool.Config().ConnConfig.Database))

	err = dbPool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s %w", caller.CodeLine(), err)
	}
	slog.Info("Successful ping", slog.String("database", dbPool.Config().ConnConfig.Database))

	return &PostgresPool{DB: dbPool}, nil
}
