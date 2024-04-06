package repository

import (
	"context"

	"github.com/msmkdenis/yap-infokeeper/internal/model"
)

func (r *PostgresTextDataRepository) Insert(ctx context.Context, textData model.TextData) error {
	_, err := r.postgresPool.DB.Exec(ctx, insertTextData,
		textData.ID,
		textData.OwnerID,
		textData.Data,
		textData.Metadata)

	return err
}
