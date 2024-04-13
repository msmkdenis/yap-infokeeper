package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/msmkdenis/yap-infokeeper/internal/model"
	"github.com/msmkdenis/yap-infokeeper/internal/text_data/specification"
	"github.com/msmkdenis/yap-infokeeper/pkg/caller"
)

func (r *PostgresTextDataRepository) SelectAll(ctx context.Context, spec *specification.TextDataSpecification) ([]model.TextData, error) {
	query, args := spec.GetQueryArgs(selectAllTextData)

	queryRows, err := r.postgresPool.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s %w", caller.CodeLine(), err)
	}

	textData, err := pgx.CollectRows(queryRows, pgx.RowToStructByPos[model.TextData])
	if err != nil {
		return nil, fmt.Errorf("%s %w", caller.CodeLine(), err)
	}

	return textData, nil
}
