package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	pkgerr "github.com/ariefsibuea/freshmart-api/internal/pkg/errors"
	"github.com/ariefsibuea/freshmart-api/internal/pkg/logger"

	"github.com/labstack/echo/v4"
)

func ErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return // response has been already sent to the client by handler or some middleware
	}

	var (
		code    int
		message string
		apiErr  *pkgerr.APIError
	)

	switch {
	case errors.As(err, &apiErr):
		code = apiErr.Code()
		message = apiErr.Error()
	case errors.Is(err, context.DeadlineExceeded):
		code = http.StatusRequestTimeout
		message = "request timeout"
	default:
		if httpErr, ok := err.(*echo.HTTPError); ok {
			code = httpErr.Code
			switch msg := httpErr.Message.(type) {
			case string:
				message = msg
			case error:
				message = msg.Error()
			default:
				message = fmt.Sprintf("%v", msg)
			}
		} else {
			code = http.StatusInternalServerError
			message = "internal server error"

			logger.FromContext(c.Request().Context()).Error("unexpected error",
				logger.FieldError, err.Error(),
			)
		}
	}

	c.JSON(code, errorResponse(message))
}

func errorResponse(message string) map[string]any {
	return map[string]any{
		"status": "error",
		"error":  map[string]string{"message": message},
	}
}
