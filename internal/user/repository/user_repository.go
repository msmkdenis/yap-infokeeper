package repository

import (
	"context"
	_ "embed"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/msmkdenis/yap-infokeeper/internal/storage/db"
	"github.com/msmkdenis/yap-infokeeper/internal/user/model"
	apperr "github.com/msmkdenis/yap-infokeeper/pkg/apperror"
)

//go:embed queries/insert_user.sql
var insertUser string

type PostgresUserRepository struct {
	postgresPool *db.PostgresPool
}

func NewPostgresUserRepository(postgresPool *db.PostgresPool) *PostgresUserRepository {
	return &PostgresUserRepository{postgresPool: postgresPool}
}

func (r *PostgresUserRepository) Insert(ctx context.Context, user model.User) error {
	_, err := r.postgresPool.DB.Exec(ctx, insertUser, user.ID, user.Login, user.Password)

	var e *pgconn.PgError
	if errors.As(err, &e) && e.Code == pgerrcode.UniqueViolation {
		return apperr.NewValueError("User already exists", apperr.Caller(), err)
	}

	return err
}
