package repository

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/msmkdenis/yap-infokeeper/internal/model"
	"github.com/msmkdenis/yap-infokeeper/internal/storage/db"
	"github.com/msmkdenis/yap-infokeeper/internal/text_data/specification"
)

//go:embed queries/insert_text_data.sql
var insertTextData string

//go:embed queries/select_all_text_data.sql
var selectAllTextData string

type PostgresTextDataRepository struct {
	postgresPool *db.PostgresPool
}

func NewPostgresTextDataRepository(postgresPool *db.PostgresPool) *PostgresTextDataRepository {
	return &PostgresTextDataRepository{postgresPool: postgresPool}
}

func (r *PostgresTextDataRepository) Insert(ctx context.Context, textData model.TextData) error {
	_, err := r.postgresPool.DB.Exec(ctx, insertTextData,
		textData.ID,
		textData.OwnerID,
		textData.Data,
		textData.Metadata)

	return err
}

func (r *PostgresTextDataRepository) SelectAll(ctx context.Context, spec *specification.TextDataSpecification) ([]model.TextData, error) {
	query, args := spec.GetQueryArgs(selectAllTextData)
	fmt.Println(query, args)

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
