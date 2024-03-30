package repository

import (
	"context"
	_ "embed"
	"errors"

	"github.com/msmkdenis/yap-infokeeper/internal/model"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/msmkdenis/yap-infokeeper/internal/storage/db"
	apperr "github.com/msmkdenis/yap-infokeeper/pkg/apperror"
)

//go:embed queries/insert_credit_card.sql
var insertCreditCard string

type PostgresCreditCardRepository struct {
	postgresPool *db.PostgresPool
}

func NewPostgresCreditCardRepository(postgresPool *db.PostgresPool) *PostgresCreditCardRepository {
	return &PostgresCreditCardRepository{postgresPool: postgresPool}
}

func (r *PostgresCreditCardRepository) Insert(ctx context.Context, ownerID string, card model.CreditCard) error {
	_, err := r.postgresPool.DB.Exec(ctx, insertCreditCard, card.ID, card.Number, card.Owner,
		card.ExpiresAt, card.CVVCode, card.PinCode, card.Metadata)

	var e *pgconn.PgError
	if errors.As(err, &e) && e.Code == pgerrcode.UniqueViolation {
		return apperr.NewValueError("Credit card with this number already exists", apperr.Caller(), err)
	}

	return err
}
