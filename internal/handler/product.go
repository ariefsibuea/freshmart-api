package handler

import (
	"net/http"
	"strconv"

	"github.com/ariefsibuea/freshmart-api/internal/model"
	"github.com/ariefsibuea/freshmart-api/internal/usecase"

	"github.com/labstack/echo/v4"
)

type productHandler struct {
	usecase usecase.ProductUsecase
}

func InitProductHandler(e *echo.Group, usecase usecase.ProductUsecase) {
	handler := &productHandler{
		usecase: usecase,
	}

	e.POST("/products", handler.create)
	e.GET("/products/:id", handler.get)
}

func (h *productHandler) create(c echo.Context) error {
	var req model.CreateProductRequest

	if err := c.Bind(&req); err != nil {
		return Error(c, http.StatusBadRequest, "invalid request body")
	}

	product, err := h.usecase.Create(c.Request().Context(), req)
	if err != nil {
		return err
	}

	return Success(c, http.StatusCreated, product)
}

func (h *productHandler) get(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return Error(c, http.StatusBadRequest, "invalid product id")
	}

	product, err := h.usecase.Get(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return Success(c, http.StatusOK, product)
}
