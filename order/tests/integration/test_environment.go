//go:build integration

package integration

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Medveddo/rocket-science/order/internal/repository"
	"github.com/Medveddo/rocket-science/platform/pkg/testcontainers/postgres"
)

// TestEnvironment holds the test environment resources
type TestEnvironment struct {
	Network  interface{} // Network interface from testcontainers
	Postgres *postgres.Container
	DB       *pgxpool.Pool
	Repo     repository.OrderRepository
	psql     squirrel.StatementBuilderType
}

// ClearOrdersTable clears all orders from the orders table
func (env *TestEnvironment) ClearOrdersTable(ctx context.Context) error {
	query, _, err := env.psql.Delete(ordersTable).ToSql()
	if err != nil {
		return fmt.Errorf("failed to delete orders table: %w", err)
	}

	_, err = env.DB.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to clear orders table: %w", err)
	}
	return nil
}

// WaitForDatabase waits for the database to be ready
func (env *TestEnvironment) WaitForDatabase(ctx context.Context) error {
	timeout := 30 * time.Second
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		err := env.DB.Ping(ctx)
		if err == nil {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}

	return fmt.Errorf("database not ready after %v", timeout)
}
