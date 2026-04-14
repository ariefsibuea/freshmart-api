package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/ariefsibuea/freshmart-api/internal/cache"
	"github.com/ariefsibuea/freshmart-api/internal/model"
	"github.com/ariefsibuea/freshmart-api/internal/pkg/logger"
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
		Description: req.Description,
		Quantity:    req.Quantity,
	}

	createdProduct, err := u.repository.Create(ctx, product)
	if err != nil {
		return model.Product{}, fmt.Errorf("create product failed: %w", err)
	}

	return createdProduct, nil
}

func (u *productUsecase) Get(ctx context.Context, id int64) (model.Product, error) {
	cacheKey := fmt.Sprintf("products:%d", id)

	if u.cache != nil {
		var cachedProduct model.Product
		err := u.cache.Get(ctx, cacheKey, &cachedProduct)
		if err == nil {
			return cachedProduct, nil
		}

		if err != cache.ErrCacheKeyNotFound {
			logger.FromContext(ctx).Warn("cache get error",
				"cacheKey", cacheKey,
				logger.FieldError, err.Error(),
			)
		}
	}

	product, err := u.repository.Get(ctx, id)
	if err != nil {
		return model.Product{}, fmt.Errorf("get product '%d' failed: %w", id, err)
	}

	if u.cache != nil {
		if cacheErr := u.cache.Set(ctx, cacheKey, product, u.cacheTTL); cacheErr != nil {
			logger.FromContext(ctx).Warn("write product to cache failed",
				"id", id,
				"cacheKey", cacheKey,
				logger.FieldError, cacheErr.Error(),
			)
		}
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
