package repository

import (
	"context"
	"fmt"

	"github.com/msmkdenis/yap-infokeeper/internal/model"
	"github.com/msmkdenis/yap-infokeeper/pkg/caller"
)

func (r *PostgresTextDataRepository) Insert(ctx context.Context, textData model.TextData) error {
	_, err := r.postgresPool.DB.Exec(ctx, insertTextData,
		textData.ID,
		textData.OwnerID,
		textData.Data,
		textData.Metadata)
	if err != nil {
		return fmt.Errorf("%s %w", caller.CodeLine(), err)
	}

	return nil
}
