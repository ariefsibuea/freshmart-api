package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	pkgerr "github.com/ariefsibuea/freshmart-api/internal/pkg/errors"
	"github.com/ariefsibuea/freshmart-api/internal/pkg/logger"

	"github.com/labstack/echo/v4"
)

type Response struct {
	Status     string       `json:"status"`
	Data       any          `json:"data,omitempty"`
	Pagination *Pagination  `json:"pagination,omitempty"`
	Error      *ErrorDetail `json:"error,omitempty"`
}

type Pagination struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

type ErrorDetail struct {
	Message string `json:"message"`
}

const (
	StatusSuccess = "success"
	StatusError   = "error"
)

func Success(c echo.Context, statusCode int, data any) error {
	return c.JSON(statusCode, Response{
		Status: StatusSuccess,
		Data:   data,
	})
}

func SuccessWithPagination(c echo.Context, statusCode int, data any, pagination Pagination) error {
	return c.JSON(statusCode, Response{
		Status:     StatusSuccess,
		Data:       data,
		Pagination: &pagination,
	})
}

func Error(c echo.Context, statusCode int, message string) error {
	return c.JSON(statusCode, Response{
		Status: StatusError,
		Error:  &ErrorDetail{Message: message},
	})
}

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

	c.JSON(code, Response{
		Status: StatusError,
		Error:  &ErrorDetail{Message: message},
	})
}
