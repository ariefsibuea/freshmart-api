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
	e.GET("/products", handler.fetch)
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

func (h *productHandler) fetch(c echo.Context) error {
	ctx := c.Request().Context()

	filter := model.ProductFilter{
		Name:     c.QueryParam("name"),
		SortBy:   c.QueryParam("sort_by"),
		Order:    c.QueryParam("order"),
		Page:     model.DefaultPage,
		PageSize: model.DefaultPageSize,
	}

	if productType := c.QueryParam("product_type"); productType != "" {
		filter.ProductType = model.ProductType(productType)
	}

	if page := c.QueryParam("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil {
			filter.Page = p
		}
	}

	if pageSize := c.QueryParam("page_size"); pageSize != "" {
		if ps, err := strconv.Atoi(pageSize); err == nil {
			filter.PageSize = ps
		}
	}

	if err := filter.Validate(); err != nil {
		return err
	}

	products, total, err := h.usecase.Fetch(ctx, filter)
	if err != nil {
		return err
	}

	return SuccessWithPagination(c, http.StatusOK, products, Pagination{
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		TotalItems: int(total),
		TotalPages: int((total + int64(filter.PageSize) - 1) / int64(filter.PageSize)),
	})
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
