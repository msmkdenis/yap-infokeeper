package repository

import (
	_ "embed"

	"github.com/msmkdenis/yap-infokeeper/internal/storage/db"
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
