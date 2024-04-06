package repository

import (
	"context"
	_ "embed"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/msmkdenis/yap-infokeeper/internal/model"
	"github.com/msmkdenis/yap-infokeeper/internal/storage/db"
	apperr "github.com/msmkdenis/yap-infokeeper/pkg/apperror"
)

//go:embed queries/insert_credit_card.sql
var insertCreditCard string

//go:embed queries/select_credit_card_by_owner_id_number.sql
var selectCreditCardByOwnerIDNumber string

//go:embed queries/select_all_credit_cards_by_owner_id.sql
var selectAllByOwnerID string

type PostgresCreditCardRepository struct {
	postgresPool *db.PostgresPool
}

func NewPostgresCreditCardRepository(postgresPool *db.PostgresPool) *PostgresCreditCardRepository {
	return &PostgresCreditCardRepository{postgresPool: postgresPool}
}

func (r *PostgresCreditCardRepository) SelectByOwnerIDCardNumber(ctx context.Context, ownerID, number string) (*model.CreditCard, error) {
	var card model.CreditCard
	err := r.postgresPool.DB.QueryRow(ctx, selectCreditCardByOwnerIDNumber, ownerID, number).
		Scan(&card.ID,
			&card.Number,
			&card.OwnerID,
			&card.OwnerName,
			&card.ExpiresAt,
			&card.CVVCode,
			&card.PinCode,
			&card.CreatedAt,
			&card.Metadata)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.ErrCardNotFound
	}

	if err != nil {
		return nil, err
	}

	return &card, nil
}

func (r *PostgresCreditCardRepository) SelectAllByOwnerID(ctx context.Context, ownerID string) ([]model.CreditCard, error) {
	queryRows, err := r.postgresPool.DB.Query(ctx, selectAllByOwnerID, ownerID)
	if err != nil {
		return nil, apperr.NewValueError("query failed", apperr.Caller(), err)
	}
	defer queryRows.Close()

	cards, err := pgx.CollectRows(queryRows, pgx.RowToStructByPos[model.CreditCard])
	if err != nil {
		return nil, apperr.NewValueError("unable to collect rows", apperr.Caller(), err)
	}

	return cards, nil
}
