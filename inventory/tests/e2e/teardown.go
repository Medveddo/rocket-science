//go:build integration
package integration

import (
	"context"

	"go.uber.org/zap"

	"github.com/Medveddo/rocket-science/platform/pkg/logger"
)

// cleanupTestEnvironment — очищает ресурсы тестового окружения
func cleanupTestEnvironment(ctx context.Context, env *TestEnvironment) {
	logger.Info(ctx, "🧹 Очистка тестового окружения...")

	if env.App != nil {
		if err := env.App.Terminate(ctx); err != nil {
			logger.Error(ctx, "ошибка при остановке контейнера приложения", zap.Error(err))
		} else {
			logger.Info(ctx, "✅ Контейнер приложения остановлен")
		}
	}

	if env.Mongo != nil {
		if err := env.Mongo.Terminate(ctx); err != nil {
			logger.Error(ctx, "ошибка при остановке контейнера MongoDB", zap.Error(err))
		} else {
			logger.Info(ctx, "✅ Контейнер MongoDB остановлен")
		}
	}

	if env.Network != nil {
		if err := env.Network.Remove(ctx); err != nil {
			logger.Error(ctx, "ошибка при удалении сети", zap.Error(err))
		} else {
			logger.Info(ctx, "✅ Сеть удалена")
		}
	}

	logger.Info(ctx, "🎉 Тестовое окружение очищено")
}
