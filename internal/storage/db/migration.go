package db

import (
	"embed"
	"fmt"
	"log/slog"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/msmkdenis/yap-infokeeper/pkg/caller"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// Migrations represents the go migrate instance.
type Migrations struct {
	migrations *migrate.Migrate
}

// NewMigrations creates a new Migrations instance.
//
// It takes a connection string and a logger as parameters and returns a
// pointer to Migrations and an error.
func NewMigrations(connection string) (*Migrations, error) {
	dbConfig, err := pgxpool.ParseConfig(connection)
	if err != nil {
		return nil, fmt.Errorf("%s %w", caller.CodeLine(), err)
	}
	slog.Info("Successful connection string parsing", slog.String("database", dbConfig.ConnConfig.Database))

	dbURL := makeDBURL(dbConfig, parseSSLMode(connection))

	driver, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return nil, fmt.Errorf("%s %w", caller.CodeLine(), err)
	}
	slog.Info("Successful connection", slog.String("database", dbConfig.ConnConfig.Database))

	migrations, err := migrate.NewWithSourceInstance("iofs", driver, dbURL)
	if err != nil {
		return nil, fmt.Errorf("%s %w", caller.CodeLine(), err)
	}

	return &Migrations{migrations: migrations}, nil
}

// MigrateUp perform migrations up.
func (m *Migrations) MigrateUp() error {
	err := m.migrations.Up()
	if err != nil && err.Error() != "no change" {
		return fmt.Errorf("%s %w", caller.CodeLine(), err)
	}
	return nil
}

func makeDBURL(config *pgxpool.Config, sslMode string) string {
	var dbURL strings.Builder

	dbURL.WriteString("postgres://")
	dbURL.WriteString(config.ConnConfig.User)
	dbURL.WriteString(":")
	dbURL.WriteString(config.ConnConfig.Password)
	dbURL.WriteString("@")
	dbURL.WriteString(config.ConnConfig.Host)
	dbURL.WriteString(":")
	dbURL.WriteString(fmt.Sprint(config.ConnConfig.Port))
	dbURL.WriteString("/")
	dbURL.WriteString(config.ConnConfig.Database)
	dbURL.WriteString("?sslmode=")
	if config.ConnConfig.TLSConfig == nil {
		dbURL.WriteString("disable")
	} else {
		dbURL.WriteString(sslMode)
	}

	return dbURL.String()
}

func parseSSLMode(connection string) string {
	con := strings.Split(connection, " ")
	sslMode := ""
	for _, v := range con {
		pair := strings.Split(v, "=")
		if pair[0] == "sslmode" {
			sslMode = pair[1]
		}
	}

	return sslMode
}
