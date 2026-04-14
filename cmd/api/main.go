package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ariefsibuea/freshmart-api/config"
	"github.com/ariefsibuea/freshmart-api/internal/cache"
	"github.com/ariefsibuea/freshmart-api/internal/database"
	"github.com/ariefsibuea/freshmart-api/internal/handler"
	mw "github.com/ariefsibuea/freshmart-api/internal/middleware"
	"github.com/ariefsibuea/freshmart-api/internal/pkg/logger"
	"github.com/ariefsibuea/freshmart-api/internal/repository"
	"github.com/ariefsibuea/freshmart-api/internal/usecase"

	"github.com/labstack/echo/v4"
)

func main() {
	conf := config.Load()

	logger.Init(logger.ToSlogLevel(conf.LogLevel))

	e := echo.New()
	e.Server.ReadTimeout = conf.ServerReadTimeout
	e.Server.WriteTimeout = conf.ServerWriteTimeout
	e.Server.IdleTimeout = conf.ServerIdleTimeout

	e.HTTPErrorHandler = mw.ErrorHandler

	e.Use(mw.RequestID())
	e.Use(mw.Log())
	e.Use(mw.CORS(conf.AllowOrigins))
	e.Use(mw.Timeout(conf.ServerRequestTimeout))

	redisAddr := fmt.Sprintf("%s:%d", conf.Cache.RedisHost, conf.Cache.RedisPort)
	redisCache, err := cache.NewRedisConnection(redisAddr, conf.Cache.RedisPingTimeout)
	if err != nil {
		logger.FromContext(context.Background()).Warn("redis unavailable, continuing without cache",
			"addr", redisAddr,
			logger.FieldError, err.Error(),
		)
	}

	mysqlDB, err := database.NewMySQLConnection(conf.Database)
	if err != nil {
		logger.FromContext(context.Background()).Error("mysql db unavailable",
			logger.FieldError, err.Error(),
		)
		os.Exit(1)
	}
	defer func() {
		if err := mysqlDB.Close(); err != nil {
			logger.FromContext(context.Background()).Warn("close mysql connection failed",
				logger.FieldError, err.Error(),
			)
		}
	}()

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]any{
			"status": "ok",
		})
	})

	apiGroup := e.Group("/api/v1")

	productRepository := repository.NewProductRepository(mysqlDB)
	productUsecase := usecase.NewProductUsecase(productRepository, redisCache, conf.Cache.DefaultCacheTTL)
	handler.InitProductHandler(apiGroup, productUsecase)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		address := fmt.Sprintf(":%d", conf.ServerPort)
		if err := e.Start(address); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.FromContext(context.Background()).Error("shutting down the server",
				logger.FieldError, err.Error(),
			)
			os.Exit(1)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), conf.ServerShutdownTimeout)
	defer cancel()

	if err := e.Shutdown(shutdownCtx); err != nil {
		logger.FromContext(context.Background()).Error("shutdown error",
			logger.FieldError, err.Error(),
		)
	}
}
