package handler

import (
	pkgres "github.com/ariefsibuea/freshmart-api/internal/pkg/response"

	"github.com/labstack/echo/v4"
)

type Pagination = pkgres.Pagination

func Success(c echo.Context, statusCode int, data any) error {
	return c.JSON(statusCode, pkgres.Response{
		Status: pkgres.StatusSuccess,
		Data:   data,
	})
}

func SuccessWithPagination(c echo.Context, statusCode int, data any, pagination Pagination) error {
	return c.JSON(statusCode, pkgres.Response{
		Status:     pkgres.StatusSuccess,
		Data:       data,
		Pagination: &pagination,
	})
}
