package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/msmkdenis/yap-infokeeper/internal/model"
	apperr "github.com/msmkdenis/yap-infokeeper/pkg/apperror"
)

func (r *PostgresCreditCardRepository) Insert(ctx context.Context, card model.CreditCard) error {
	_, err := r.postgresPool.DB.Exec(ctx, insertCreditCard,
		card.ID,
		card.Number,
		card.OwnerID,
		card.OwnerName,
		card.ExpiresAt,
		card.CVVCode,
		card.PinCode,
		card.Metadata)

	var e *pgconn.PgError
	if errors.As(err, &e) && e.Code == pgerrcode.CheckViolation {
		if e.ConstraintName == "unique_number" {
			return apperr.ErrCardAlreadyExists
		}
	}

	return err
}
