package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/Medveddo/rocket-science/order/internal/app"
	"github.com/Medveddo/rocket-science/order/internal/config"
	"github.com/Medveddo/rocket-science/platform/pkg/closer"
	"github.com/Medveddo/rocket-science/platform/pkg/logger"
)

const configPath = "../deploy/compose/order/.env"

func main() {
	ctx := context.Background()
	err := config.Load(configPath)
	if err != nil {
		logger.Error(ctx, "cannot load config", zap.Error(err))
		return
	}

	appCtx, appCancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer appCancel()
	defer gracefulShutdown()

	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	a, err := app.New(appCtx)
	if err != nil {
		logger.Error(appCtx, "❌ failed to create application", zap.Error(err))
		return
	}

	err = a.Run(appCtx)
	if err != nil {
		logger.Error(appCtx, "❌ error running application", zap.Error(err))
		return
	}
}

func gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := closer.CloseAll(ctx); err != nil {
		logger.Error(ctx, "❌ error shutting down application", zap.Error(err))
	}
}
