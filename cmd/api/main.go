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

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func main() {
	conf := config.Load()

	e := echo.New()
	e.Logger.SetLevel(log.Lvl(conf.LogLevel))

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]any{
			"status": "ok",
		})
	})

	redisAddr := fmt.Sprintf("%s:%d", conf.Cache.RedisHost, conf.Cache.RedisPort)
	// WARN: ignore returned redis instance since we don't use it yet
	_, err := cache.NewRedis(redisAddr, conf.Cache.RedisPingTimeout)
	if err != nil {
		e.Logger.Warnf("redis at %q unavailable, continuing without cache: %v", redisAddr, err)
	}

	e.Server.ReadTimeout = conf.ServerReadTimeout
	e.Server.WriteTimeout = conf.ServerWriteTimeout
	e.Server.IdleTimeout = conf.ServerIdleTimeout

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		address := fmt.Sprintf(":%d", conf.ServerPort)
		if err := e.Start(address); err != nil && !errors.Is(err, http.ErrServerClosed) {
			e.Logger.Fatalf("shutting down the server: %v", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), conf.ServerShutdownTimeout)
	defer cancel()

	if err := e.Shutdown(shutdownCtx); err != nil {
		e.Logger.Fatal(err)
	}
}
