package handler

import (
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
