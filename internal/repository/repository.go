package repository

import (
	"context"

	"github.com/ariefsibuea/freshmart-api/internal/model"
)

const (
	DefaultPage     = 1
	DefaultPageSize = 10
)

type ProductRepository interface {
	Create(ctx context.Context, product model.Product) (model.Product, error)
	Fetch(ctx context.Context, filter ProductFilter) ([]model.Product, int64, error)
	Get(ctx context.Context, id int64) (model.Product, error)
}

type ProductFilter struct {
	Name        string
	ProductType model.ProductType
	SortBy      string
	Order       string
	Page        int
	PageSize    int
}
