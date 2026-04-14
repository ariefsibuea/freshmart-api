package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/ariefsibuea/freshmart-api/internal/cache"
	"github.com/ariefsibuea/freshmart-api/internal/model"
	"github.com/ariefsibuea/freshmart-api/internal/repository"
)

type productUsecase struct {
	repository repository.ProductRepository
	cache      cache.Cache
	cacheTTL   time.Duration
}

func NewProductUsecase(
	repository repository.ProductRepository,
	cache cache.Cache,
	cacheTTL time.Duration) ProductUsecase {

	return &productUsecase{
		repository: repository,
		cache:      cache,
		cacheTTL:   cacheTTL,
	}
}

func (u *productUsecase) Create(ctx context.Context, req model.CreateProductRequest) (model.Product, error) {
	if err := req.Validate(); err != nil {
		return model.Product{}, err
	}

	product := model.Product{
		Name:        req.Name,
		Price:       req.Price,
		ProductType: req.ProductType,
		Quantity:    req.Quantity,
	}

	if req.Description != nil {
		product.Description = *req.Description
	}

	createdProduct, err := u.repository.Create(ctx, product)
	if err != nil {
		return model.Product{}, fmt.Errorf("create product failed: %w", err)
	}

	return createdProduct, nil
}

func (u *productUsecase) Get(ctx context.Context, id int64) (model.Product, error) {
	cacheKey := fmt.Sprintf("products:%d", id)

	var cachedProduct model.Product
	err := u.cache.Get(ctx, cacheKey, &cachedProduct)
	if err == nil {
		return cachedProduct, nil
	}

	if err != cache.ErrCacheKeyNotFound {
		slog.Warn("cache get error", "cacheKey", cacheKey, "error", err)
	}

	product, err := u.repository.Get(ctx, id)
	if err != nil {
		return model.Product{}, fmt.Errorf("get product '%d' failed: %w", id, err)
	}

	if cacheErr := u.cache.Set(ctx, cacheKey, product, u.cacheTTL); cacheErr != nil {
		slog.Warn("write product to cache failed", "id", id, "cacheKey", cacheKey)
	}

	return product, nil
}

func (u *productUsecase) Fetch(ctx context.Context, filter model.ProductFilter) ([]model.Product, int64, error) {
	products, total, err := u.repository.Fetch(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("fetch products failed: %w", err)
	}

	return products, total, nil
}
