package repository

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/msmkdenis/yap-infokeeper/internal/storage/db"
)

const (
	testDatabaseName    = "yap-infokeeper_server-test"
	testDatabaseUser    = "postgres"
	testDatabasePass    = "postgres"
	testDatabasePort    = "5432/tcp"
	containerMappedPort = "5432"
)

type CreditCardRepositoryTestSuite struct {
	suite.Suite
	pool                 *db.PostgresPool
	creditCardRepository *PostgresCreditCardRepository
	container            testcontainers.Container
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(CreditCardRepositoryTestSuite))
}

func (s *CreditCardRepositoryTestSuite) SetupTest() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
	slog.SetDefault(logger)
	container, pool, err := setupTestDatabase()
	if err != nil {
		slog.Error("Unable to setup test database", slog.String("error", err.Error()))
	}
	require.NoError(s.T(), err)
	s.pool = pool
	s.container = container
	s.creditCardRepository = NewPostgresCreditCardRepository(s.pool)
}

func (s *CreditCardRepositoryTestSuite) TearDownTest() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := s.container.Terminate(ctx)
	assert.NoError(s.T(), err)
}

func setupTestDatabase() (testcontainers.Container, *db.PostgresPool, error) {
	containerReq := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{testDatabasePort},
		WaitingFor:   wait.ForListeningPort(testDatabasePort),
		Env: map[string]string{
			"POSTGRES_DB":       testDatabaseName,
			"POSTGRES_PASSWORD": testDatabaseUser,
			"POSTGRES_USER":     testDatabasePass,
		},
	}
	dbContainer, err := testcontainers.GenericContainer(
		context.Background(),
		testcontainers.GenericContainerRequest{
			ContainerRequest: containerReq,
			Started:          true,
		})
	if err != nil {
		return nil, nil, err
	}

	port, err := dbContainer.MappedPort(context.Background(), containerMappedPort)
	if err != nil {
		return nil, nil, err
	}

	host, err := dbContainer.Host(context.Background())
	if err != nil {
		return nil, nil, err
	}

	connection := fmt.Sprintf("user=postgres password=postgres host=%s database=yap-infokeeper_server-test sslmode=disable port=%d", host, port.Int())

	pool, err := db.NewPostgresPool(context.Background(), connection)
	if err != nil {
		return nil, nil, err
	}

	migrations, err := db.NewMigrations(connection)
	if err != nil {
		return nil, nil, err
	}

	err = migrations.MigrateUp()
	if err != nil {
		return nil, nil, err
	}

	return dbContainer, pool, err
}
