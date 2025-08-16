//go:build integration

package integration

import (
	"context"

	"go.uber.org/zap"

	"github.com/Medveddo/rocket-science/platform/pkg/logger"
)

// cleanupTestEnvironment cleans up test environment resources
func cleanupTestEnvironment(ctx context.Context, env *TestEnvironment) {
	logger.Info(ctx, "🧹 Очистка тестового окружения...")

	// Close database connection
	if env.DB != nil {
		env.DB.Close()
	}

	// Stop PostgreSQL container
	if env.Postgres != nil {
		if err := env.Postgres.Stop(ctx); err != nil {
			logger.Error(ctx, "не удалось остановить контейнер PostgreSQL", zap.Error(err))
		}
	}

	// Clean up network
	if env.Network != nil {
		// Assuming Network has a cleanup method
		if network, ok := env.Network.(interface{ Cleanup(context.Context) error }); ok {
			if err := network.Cleanup(ctx); err != nil {
				logger.Error(ctx, "не удалось очистить сеть", zap.Error(err))
			}
		}
	}

	logger.Info(ctx, "✅ Тестовое окружение очищено")
}
