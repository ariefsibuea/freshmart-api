package middleware

import (
	"errors"
	"net/http"
	"time"

	pkgerr "github.com/ariefsibuea/freshmart-api/internal/pkg/errors"
	"github.com/ariefsibuea/freshmart-api/internal/pkg/logger"

	"github.com/labstack/echo/v4"
)

func Log() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			start := time.Now()

			// NOTE: Pre-load the request ID into the logger so every layer that
			// calls logger.FromContext(ctx) inherits it automatically.
			log := logger.FromContext(req.Context()).With(logger.FieldRequestID, GetRequestID(c))
			c.SetRequest(req.WithContext(logger.IntoContext(req.Context(), log)))

			err := next(c)

			status := c.Response().Status
			if err != nil {
				status = errorStatusCode(err)
			}

			latencyMS := time.Since(start).Milliseconds()

			fields := []any{
				logger.FieldMethod, req.Method,
				logger.FieldPath, req.URL.Path,
				logger.FieldStatus, status,
				logger.FieldLatencyMS, latencyMS,
			}

			if err != nil {
				fields = append(fields, logger.FieldError, err.Error())
			}

			switch {
			case status >= http.StatusInternalServerError:
				log.Error("server error", fields...)
			case status >= http.StatusBadRequest:
				log.Warn("client error", fields...)
			default:
				log.Info("request completed", fields...)
			}

			return err
		}
	}
}

func errorStatusCode(err error) int {
	var httpErr *echo.HTTPError
	if errors.As(err, &httpErr) {
		return httpErr.Code
	}
	return pkgerr.GetErrorCode(err)
}
