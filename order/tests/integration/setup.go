//go:build integration

package integration

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/docker/go-connections/nat"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"

	"github.com/Medveddo/rocket-science/order/internal/repository/order"
	"github.com/Medveddo/rocket-science/platform/pkg/logger"
	pgMigrator "github.com/Medveddo/rocket-science/platform/pkg/migrator/pg"
	"github.com/Medveddo/rocket-science/platform/pkg/testcontainers"
	"github.com/Medveddo/rocket-science/platform/pkg/testcontainers/network"
	"github.com/Medveddo/rocket-science/platform/pkg/testcontainers/path"
	"github.com/Medveddo/rocket-science/platform/pkg/testcontainers/postgres"
)

const (
	startupTimeout = 2 * time.Minute
)

// setupTestEnvironment prepares the test environment: network, containers and returns structure with resources
func setupTestEnvironment(ctx context.Context) *TestEnvironment {
	logger.Info(ctx, "🚀 Подготовка тестового окружения для order service...")

	// Step 1: Create common Docker network
	generatedNetwork, err := network.NewNetwork(ctx, projectName)
	if err != nil {
		logger.Fatal(ctx, "не удалось создать общую сеть", zap.Error(err))
	}
	logger.Info(ctx, "✅ Сеть успешно создана")

	postgresUser := getEnvWithLogging(ctx, testcontainers.PostgresUserKey)
	postgresPassword := getEnvWithLogging(ctx, testcontainers.PostgresPasswordKey)
	postgresDatabase := getEnvWithLogging(ctx, testcontainers.PostgresDatabaseKey)
	postgresImageName := getEnvWithLogging(ctx, testcontainers.PostgresImageNameKey)

	waitStrategy := wait.ForListeningPort(nat.Port("5432/tcp")).
		WithStartupTimeout(startupTimeout)

	// Step 2: Start PostgreSQL container
	generatedPostgres, err := postgres.NewContainer(ctx,
		postgres.WithNetworkName(generatedNetwork.Name()),
		postgres.WithContainerName(testcontainers.PostgresContainerName),
		postgres.WithImageName(postgresImageName),
		postgres.WithDatabase(postgresDatabase),
		postgres.WithUser(postgresUser),
		postgres.WithPassword(postgresPassword),
		postgres.WithLogger(logger.Logger()),
		postgres.WithWaitStrategy(waitStrategy),
	)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork})
		logger.Fatal(ctx, "не удалось запустить контейнер PostgreSQL", zap.Error(err))
	}
	logger.Info(ctx, "✅ Контейнер PostgreSQL успешно запущен")

	// Step 3: Connect to database and run migrations
	connString, err := generatedPostgres.ConnectionString(ctx)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork, Postgres: generatedPostgres})
		logger.Fatal(ctx, "не удалось получить строку подключения к PostgreSQL", zap.Error(err))
	}

	// Connect to database using pgxpool
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork, Postgres: generatedPostgres})
		logger.Fatal(ctx, "не удалось подключиться к PostgreSQL", zap.Error(err))
	}

	// Wait for database to be ready
	err = pool.Ping(ctx)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork, Postgres: generatedPostgres})
		logger.Fatal(ctx, "не удалось подключиться к PostgreSQL", zap.Error(err))
	}

	// Step 4: Run migrations
	err = runMigrations(ctx, pool)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork, Postgres: generatedPostgres})
		logger.Fatal(ctx, "не удалось выполнить миграции", zap.Error(err))
	}

	orderRepo := order.NewOrderRepository(pool)

	logger.Info(ctx, "🎉 Тестовое окружение готово")
	return &TestEnvironment{
		Network:  generatedNetwork,
		Postgres: generatedPostgres,
		DB:       pool,
		Repo:     orderRepo,
		psql:     sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

// runMigrations runs the database migrations using the platform migrator
func runMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	logger.Info(ctx, "🔄 Выполнение миграций...")

	var db *sql.DB = stdlib.OpenDB(*pool.Config().ConnConfig)
	defer db.Close()

	projectRoot := path.GetProjectRoot()
	migrationsDir := filepath.Join(projectRoot, "order", "migrations")
	migrator := pgMigrator.NewMigrator(db, migrationsDir)

	err := migrator.Up()
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	logger.Info(ctx, "✅ Миграции выполнены успешно")
	return nil
}

// getEnvWithLogging returns environment variable value with logging and default values
func getEnvWithLogging(ctx context.Context, key string) string {
	value := os.Getenv(key)
	if value == "" {
		switch key {
		case "POSTGRES_HOST":
			value = "localhost"
		case "POSTGRES_PORT":
			value = "5432"
		case "POSTGRES_DB":
			value = "test_orders"
		case "POSTGRES_USER":
			value = "postgres"
		case "POSTGRES_PASSWORD":
			value = "postgres"
		case "POSTGRES_IMAGE_NAME":
			value = "postgres:17-alpine"
		default:
			logger.Warn(ctx, "Переменная окружения не установлена", zap.String("key", key))
		}
	}

	return value
}
