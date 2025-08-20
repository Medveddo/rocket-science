//go:build integration

package integration

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/Medveddo/rocket-science/platform/pkg/logger"
)

const testsTimeout = 5 * time.Minute

var (
	env *TestEnvironment

	suiteCtx    context.Context
	suiteCancel context.CancelFunc
)

func TestOrderIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Order Service Integration Tests")
}

var _ = BeforeSuite(func() {
	err := logger.Init("debug", true)
	if err != nil {
		panic(fmt.Sprintf("не удалось инициализировать логгер: %v", err))
	}

	suiteCtx, suiteCancel = context.WithTimeout(context.Background(), testsTimeout)

	// Load .env file and set environment variables
	envVars, err := godotenv.Read(filepath.Join("..", "..", "..", "deploy", "compose", "order", ".env"))
	if err != nil {
		// If .env file doesn't exist, continue with default values
		logger.Info(suiteCtx, "No .env file found, using default values")
	} else {
		// Set environment variables in the process
		for key, value := range envVars {
			_ = os.Setenv(key, value)
		}
	}

	logger.Info(suiteCtx, "Запуск тестового окружения...")
	env = setupTestEnvironment(suiteCtx)
})

var _ = AfterSuite(func() {
	logger.Info(context.Background(), "Завершение набора тестов")
	if env != nil {
		cleanupTestEnvironment(suiteCtx, env)
	}
	suiteCancel()
})
