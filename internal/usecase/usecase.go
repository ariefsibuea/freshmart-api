package usecase

import (
	"context"

	"github.com/ariefsibuea/freshmart-api/internal/model"
	"github.com/ariefsibuea/freshmart-api/internal/repository"
)

type ProductUsecase interface {
	Create(ctx context.Context, req model.CreateProductRequest) (model.Product, error)
	Fetch(ctx context.Context, filter repository.ProductFilter) ([]model.Product, int64, error)
	Get(ctx context.Context, id int64) (model.Product, error)
}
