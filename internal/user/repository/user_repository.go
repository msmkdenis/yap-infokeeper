package repository

import (
	"context"
	_ "embed"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/msmkdenis/yap-infokeeper/internal/model"
	"github.com/msmkdenis/yap-infokeeper/internal/storage/db"
	apperr "github.com/msmkdenis/yap-infokeeper/pkg/apperror"
)

//go:embed queries/insert_user.sql
var insertUser string

//go:embed queries/select_user_by_login.sql
var selectUser string

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

func (r *PostgresUserRepository) SelectByLogin(ctx context.Context, login string) (*model.User, error) {
	var user model.User
	err := r.postgresPool.DB.QueryRow(ctx, selectUser, login).Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = apperr.ErrUserNotFound
		} else {
			err = apperr.NewValueError("query failed", apperr.Caller(), err)
		}
		return nil, err
	}
	return &user, nil
}
