package repository

import (
	"context"
	_ "embed"

	"github.com/msmkdenis/yap-infokeeper/internal/model"
	"github.com/msmkdenis/yap-infokeeper/internal/storage/db"
)

//go:embed queries/insert_credential.sql
var insertCredential string

type PostgresCredentialsRepository struct {
	postgresPool *db.PostgresPool
}

func NewPostgresCredentialsRepository(postgresPool *db.PostgresPool) *PostgresCredentialsRepository {
	return &PostgresCredentialsRepository{postgresPool: postgresPool}
}

func (r *PostgresCredentialsRepository) Insert(ctx context.Context, credential model.Credential) error {
	_, err := r.postgresPool.DB.Exec(ctx, insertCredential,
		credential.ID,
		credential.OwnerID,
		credential.Login,
		credential.Password,
		credential.Metadata)

	return err
}
