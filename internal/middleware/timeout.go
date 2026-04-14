package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func Timeout(timeout time.Duration) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx, cancel := context.WithTimeout(c.Request().Context(), timeout)
			defer cancel()

			c.SetRequest(c.Request().WithContext(ctx))

			done := make(chan error, 1)
			_panic := make(chan any, 1)

			go func() {
				defer func() {
					if p := recover(); p != nil {
						_panic <- p
					}
				}()
				done <- next(c)
			}()

			select {
			case err := <-done:
				return err
			case p := <-_panic:
				panic(p)
			case <-ctx.Done():
				return echo.NewHTTPError(http.StatusRequestTimeout, "request timeout")
			}
		}
	}
}
