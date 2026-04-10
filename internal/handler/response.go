package handler

import (
	"errors"

	pkgerr "github.com/ariefsibuea/freshmart-api/internal/pkg/errors"

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

func Error(c echo.Context, statusCode int, message string) error {
	return c.JSON(statusCode, Response{
		Status: StatusError,
		Error:  &ErrorDetail{Message: message},
	})
}

func ErrorHandler(err error, c echo.Context) {
	var (
		code    int
		message string
	)

	var apiErr *pkgerr.APIError
	switch {
	case errors.As(err, &apiErr):
		code = apiErr.Code()
		message = apiErr.Error()
	default:
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			message = he.Message.(string)
		} else {
			code = pkgerr.GetErrorCode(err)
			message = err.Error()
		}
	}

	c.JSON(code, Response{
		Status: StatusError,
		Error:  &ErrorDetail{Message: message},
	})
}
