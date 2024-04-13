package repository

import (
	_ "embed"

	"github.com/msmkdenis/yap-infokeeper/internal/storage/db"
)

//go:embed queries/insert_credit_card.sql
var insertCreditCard string

type PostgresCreditCardRepository struct {
	postgresPool *db.PostgresPool
}

func NewPostgresCreditCardRepository(postgresPool *db.PostgresPool) *PostgresCreditCardRepository {
	return &PostgresCreditCardRepository{postgresPool: postgresPool}
}
