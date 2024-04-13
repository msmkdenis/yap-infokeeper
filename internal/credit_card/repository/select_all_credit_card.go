package repository

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/msmkdenis/yap-infokeeper/internal/credit_card/specification"
	"github.com/msmkdenis/yap-infokeeper/internal/model"
	"github.com/msmkdenis/yap-infokeeper/pkg/caller"
)

//go:embed queries/select_all_credit_cards.sql
var selectAllCreditCards string

func (r *PostgresCreditCardRepository) SelectAll(ctx context.Context, spec *specification.CreditCardSpecification) ([]model.CreditCard, error) {
	query, args := spec.GetQueryArgs(selectAllCreditCards)

	queryRows, err := r.postgresPool.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s %w", caller.CodeLine(), err)
	}

	textData, err := pgx.CollectRows(queryRows, pgx.RowToStructByPos[model.CreditCard])
	if err != nil {
		return nil, fmt.Errorf("%s %w", caller.CodeLine(), err)
	}

	return textData, nil
}
