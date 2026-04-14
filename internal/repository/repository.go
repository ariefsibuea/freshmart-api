package repository

import (
	"context"

	"github.com/ariefsibuea/freshmart-api/internal/model"
)

type ProductRepository interface {
	Create(ctx context.Context, product model.Product) (model.Product, error)
	Fetch(ctx context.Context, filter model.ProductFilter) ([]model.Product, int64, error)
	Get(ctx context.Context, id int64) (model.Product, error)
}
