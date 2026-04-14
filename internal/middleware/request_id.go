package middleware

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const (
	HeaderRequestID     = "X-Request-ID"
	ContextKeyRequestID = "request_id"
)

func GetRequestID(c echo.Context) string {
	id, _ := c.Get(ContextKeyRequestID).(string)
	return id
}

func RequestID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()

			requestID := req.Header.Get(HeaderRequestID)
			if _, err := uuid.Parse(requestID); err != nil {
				requestID = uuid.New().String()
			}

			c.Set(ContextKeyRequestID, requestID)
			req.Header.Set(HeaderRequestID, requestID)
			res.Header().Set(HeaderRequestID, requestID)

			return next(c)
		}
	}
}
