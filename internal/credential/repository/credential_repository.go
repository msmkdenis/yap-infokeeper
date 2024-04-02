package repository

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/msmkdenis/yap-infokeeper/internal/credential/specification"
	"github.com/msmkdenis/yap-infokeeper/internal/model"
	"github.com/msmkdenis/yap-infokeeper/internal/storage/db"
)

//go:embed queries/insert_credential.sql
var insertCredential string

//go:embed queries/select_all_credentials.sql
var selectAllCredentials string

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

func (r *PostgresCredentialsRepository) SelectAll(ctx context.Context, spec *specification.CredentialSpecification) ([]model.Credential, error) {
	query, args := spec.GetQueryArgs(selectAllCredentials)
	fmt.Println(query, args)

	queryRows, err := r.postgresPool.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	credentials, err := pgx.CollectRows(queryRows, pgx.RowToStructByPos[model.Credential])
	if err != nil {
		return nil, err
	}

	fmt.Println(credentials)

	return credentials, nil
}
