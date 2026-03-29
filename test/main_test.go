package test

import (
	"errors"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/DangeL187/erax"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"back/internal/features/subs/handler"
	"back/internal/features/subs/usecase"
	"back/internal/shared/test"
)

var (
	testOnce sync.Once
	testEcho *echo.Echo
	testDB   *test.DB
)

func TestMain(m *testing.M) {
	testOnce.Do(func() {
		logger, _ := zap.NewDevelopment()
		zap.ReplaceGlobals(logger)
		defer func() {
			_ = logger.Sync()
		}()

		var err error
		testDB, err = test.SetupTestDB()
		if err != nil {
			zap.L().Fatal("Failed to setup test database", zap.String("error", "\n"+erax.Format(err)))
		}

		crudlUseCase := usecase.NewCrudlUseCase(testDB.CrudlRepo)
		sumUseCase := usecase.NewSumUseCase(testDB.SumRepo)
		subsHandler := handler.NewSubsHandler(crudlUseCase, sumUseCase)

		testEcho = echo.New()

		subs := testEcho.Group("/subscriptions")

		subs.POST("", subsHandler.Create)
		subs.GET("/:id", subsHandler.Get)
		subs.PUT("/:id", subsHandler.Update)
		subs.DELETE("/:id", subsHandler.Delete)
		subs.GET("", subsHandler.List)
		subs.GET("/costs", subsHandler.GetSubscriptionsCostForPeriod)

		go func() {
			err = testEcho.Start(":8081")
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				zap.L().Fatal("Failed to start HTTP server", zap.String("error", "\n"+erax.Format(err)))
			}
		}()

		time.Sleep(200 * time.Millisecond)
	})

	code := m.Run()
	os.Exit(code)
}

func beforeTest(t *testing.T) {
	// Очищаем таблицу subscriptions перед каждым тестом:
	require.NoError(t, testDB.DB.Exec(`
		TRUNCATE TABLE subscriptions RESTART IDENTITY CASCADE
	`).Error)
}

func setupTestServer(t *testing.T) (*echo.Echo, *test.DB) {
	beforeTest(t)
	return testEcho, testDB
}
