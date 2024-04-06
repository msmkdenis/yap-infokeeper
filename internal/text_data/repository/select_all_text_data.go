package repository

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/msmkdenis/yap-infokeeper/internal/model"
	"github.com/msmkdenis/yap-infokeeper/internal/text_data/specification"
)

func (r *PostgresTextDataRepository) SelectAll(ctx context.Context, spec *specification.TextDataSpecification) ([]model.TextData, error) {
	query, args := spec.GetQueryArgs(selectAllTextData)

	queryRows, err := r.postgresPool.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	textData, err := pgx.CollectRows(queryRows, pgx.RowToStructByPos[model.TextData])
	if err != nil {
		return nil, err
	}

	return textData, nil
}
